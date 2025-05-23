package interchain

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"hydro/test/interchain/chainsuite"

	"cosmossdk.io/math"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/suite"
)

const (
	HydroWasm               = "hydro.wasm"
	TributeWasm             = "tribute.wasm"
	DAOVotingAdapterWasm    = "dao_voting_adapter.wasm"
	STTokenInfoProviderWasm = "st_token_info_provider.wasm"
	GatekeeperWasm          = "gatekeeper.wasm"
)

var (
	hydroCodeId               string
	tributeCodeId             string
	daoVotingAdapterCodeId    string
	stTokenInfoProviderCodeId string
	gatekeeperCodeId          string
)

type HydroSuite struct {
	suite.Suite
	HubChain     *chainsuite.Chain
	NeutronChain *chainsuite.Chain
	Relayer      *chainsuite.Relayer
	ctx          context.Context
}

func (s *HydroSuite) SetupSuite() {
	ctx, err := chainsuite.NewSuiteContext(&s.Suite)
	s.Require().NoError(err)
	s.ctx = ctx

	// create and start hub chain
	s.HubChain, err = chainsuite.CreateChain(s.GetContext(), s.T(), chainsuite.GetHubSpec())
	s.Require().NoError(err)

	// setup hermes relayer
	relayer, err := chainsuite.NewRelayer(s.GetContext(), s.T())
	s.Require().NoError(err)
	s.Relayer = relayer
	err = relayer.SetupChainKeys(s.GetContext(), s.HubChain)
	s.Require().NoError(err)

	// create and start neutron chain
	s.NeutronChain, err = s.HubChain.AddConsumerChain(s.GetContext(), relayer, chainsuite.NeutronChainID, chainsuite.GetNeutronSpec)
	s.Require().NoError(err)
	s.Require().NoError(s.HubChain.UpdateAndVerifyStakeChange(s.GetContext(), s.NeutronChain, relayer, 1_000_000, 0))

	// copy contract codes to neutron validator
	s.CopyWasmFiles(HydroWasm, TributeWasm, DAOVotingAdapterWasm, STTokenInfoProviderWasm, GatekeeperWasm)

	// store hydro contract code
	hydroCodeId = s.StoreCode(s.GetContractCodePath(HydroWasm))
	// store tribute contract code
	tributeCodeId = s.StoreCode(s.GetContractCodePath(TributeWasm))
	// store dao_voting_adapter contract code
	daoVotingAdapterCodeId = s.StoreCode(s.GetContractCodePath(DAOVotingAdapterWasm))
	// store st_token_info_provider contract code
	stTokenInfoProviderCodeId = s.StoreCode(s.GetContractCodePath(STTokenInfoProviderWasm))
	// store gatekeeper contract code
	gatekeeperCodeId = s.StoreCode(s.GetContractCodePath(GatekeeperWasm))

	// start icq relayer
	sidecarConfig := chainsuite.GetIcqSidecarConfig(s.HubChain, s.NeutronChain)
	dockerClient, dockerNetwork := chainsuite.GetDockerContext(ctx)
	err = s.NeutronChain.NewSidecarProcess(
		s.ctx,
		sidecarConfig.PreStart,
		sidecarConfig.ProcessName,
		s.T().Name(),
		dockerClient,
		dockerNetwork,
		sidecarConfig.Image,
		sidecarConfig.HomeDir,
		0,
		sidecarConfig.Ports,
		sidecarConfig.StartCmd,
		sidecarConfig.Env,
	)
	s.Require().NoError(err)
	err = chainsuite.CopyIcqRelayerKey(s.GetContext(), s.NeutronChain.Sidecars[0])
	s.Require().NoError(err)
	err = s.NeutronChain.StartAllSidecars(s.ctx)
	s.Require().NoError(err)
}

func (s *HydroSuite) GetContext() context.Context {
	s.Require().NotNil(s.ctx, "Tried to GetContext before it was set. SetupSuite must run first")
	return s.ctx
}

func txAmountUatom(txAmount uint64) string {
	return fmt.Sprintf("%d%s", txAmount, chainsuite.Uatom)
}

func (s *HydroSuite) DelegateTokens(node *cosmos.ChainNode, keyMoniker string, valoperAddr string, amount string) {
	_, err := node.ExecTx(
		s.GetContext(),
		keyMoniker,
		"staking", "delegate", valoperAddr, amount,
	)
	s.Require().NoError(err)
}

func (s *HydroSuite) LiquidStakeTokens(node *cosmos.ChainNode, keyMoniker string, valoperAddr string, delegatorAddr string, amount string) string {
	txHash, err := node.ExecTx(
		s.GetContext(),
		keyMoniker,
		"staking", "tokenize-share",
		valoperAddr,
		amount,
		delegatorAddr,
	)
	s.Require().NoError(err)

	response, err := node.TxHashToResponse(s.GetContext(), txHash)
	s.Require().NoError(err)

	recordId, found := getEvtAttribute(response.Events, stakingtypes.EventTypeTokenizeShares, stakingtypes.AttributeKeyShareRecordID)
	s.Require().True(found)

	return recordId
}

func (s *HydroSuite) HubToNeutronShareTokenTransfer(
	validatorIndex int,
	amount math.Int,
	sourceIbcDenom string,
	dstAddress string,
) string {
	keyMoniker := s.HubChain.ValidatorWallets[validatorIndex].Moniker

	hubTransferChannel, err := s.Relayer.GetTransferChannel(s.GetContext(), s.HubChain, s.NeutronChain)
	s.Require().NoError(err)

	dstIbcDenom := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom("transfer", hubTransferChannel.Counterparty.ChannelID, sourceIbcDenom)).IBCDenom()

	// check balance before sending the tokens
	balanceBefore, err := s.NeutronChain.GetBalance(s.GetContext(), dstAddress, dstIbcDenom)
	if err != nil {
		balanceBefore = math.ZeroInt()
	}

	_, err = s.HubChain.Validators[validatorIndex].SendIBCTransfer(s.GetContext(), hubTransferChannel.ChannelID, keyMoniker, ibc.WalletAmount{
		Denom:   sourceIbcDenom,
		Amount:  amount,
		Address: dstAddress,
	}, ibc.TransferOptions{})
	s.Require().NoError(err)

	tCtx, tCancel := context.WithTimeout(s.GetContext(), 30*chainsuite.CommitTimeout)
	defer tCancel()

	// check that tokens are sent
	ibcTokensReceived := false
	for tCtx.Err() == nil {
		time.Sleep(chainsuite.CommitTimeout)
		balanceAfter, err := s.NeutronChain.GetBalance(s.GetContext(), dstAddress, dstIbcDenom)
		if err != nil {
			continue
		}

		if balanceAfter.Sub(balanceBefore).Equal(amount) {
			ibcTokensReceived = true
			break
		}
	}
	s.Require().True(ibcTokensReceived)

	return dstIbcDenom
}

func (s *HydroSuite) CopyWasmFiles(fileNames ...string) {
	for _, fileName := range fileNames {
		binary, err := os.ReadFile(fmt.Sprintf("../../artifacts/%s", fileName))
		s.Require().NoError(err)
		s.Require().NoError(s.NeutronChain.GetNode().WriteFile(s.GetContext(), binary, fileName))
	}
}

func (s *HydroSuite) StoreCode(contractPath string) string {
	node := s.NeutronChain.Validators[0]
	keyMoniker := s.NeutronChain.ValidatorWallets[0].Moniker
	txHash, err := node.ExecTx(s.GetContext(), keyMoniker, "wasm", "store", contractPath, "--gas", "auto")
	s.Require().NoError(err)

	response, err := node.TxHashToResponse(s.GetContext(), txHash)
	s.Require().NoError(err)

	codeId, found := getEvtAttribute(response.Events, wasmtypes.EventTypeStoreCode, wasmtypes.AttributeKeyCodeID)
	s.Require().True(found)

	return codeId
}

func (s *HydroSuite) GetLSMTokenInfoProviderInitMsg(maxValParticipating int) map[string]interface{} {
	neutronTransferChannel, err := s.Relayer.GetTransferChannel(s.GetContext(), s.NeutronChain, s.HubChain)
	s.Require().NoError(err)

	return map[string]interface{}{
		"max_validator_shares_participating": maxValParticipating,
		"hub_connection_id":                  neutronTransferChannel.ConnectionHops[0],
		"hub_transfer_channel_id":            neutronTransferChannel.ChannelID,
		"icq_update_period":                  10,
	}
}

func (s *HydroSuite) GetInstantiateContractMsg(codeId string, initMsg map[string]any, label string, admin *string) map[string]any {
	codeIdInt, err := strconv.Atoi(codeId)
	s.Require().NoError(err)

	initMsgJson, err := json.Marshal(initMsg)
	s.Require().NoError(err)

	intMsgBase64 := base64.StdEncoding.EncodeToString(initMsgJson)

	return map[string]any{
		"code_id": codeIdInt,
		"msg":     intMsgBase64,
		"label":   label,
		"admin":   admin,
	}
}

func (s *HydroSuite) InstantiateHydroContract(
	keyMoniker string,
	codeId string,
	adminAddr string,
	roundLength int,
	tokenInfoProviderLSMInitMsg any,
	tokenInfoProviderContractsInitMsgs []any,
	gatekeeperInitMsg map[string]any,
) string {
	s.Suite.T().Log("Instantiating Hydro contract")

	firstRoundStartTime := time.Now().UnixNano()

	var tokenInfoProvidersInitMsgs []map[string]any

	if tokenInfoProviderLSMInitMsg != nil {
		lsmTokenInfoProviderInitMsg := map[string]any{
			"lsm": tokenInfoProviderLSMInitMsg,
		}

		tokenInfoProvidersInitMsgs = append(tokenInfoProvidersInitMsgs, lsmTokenInfoProviderInitMsg)
	}

	for _, contractInitMsg := range tokenInfoProviderContractsInitMsgs {
		tokenInfoProviderContractInitMsg := map[string]any{
			"token_info_provider_contract": contractInitMsg,
		}

		tokenInfoProvidersInitMsgs = append(tokenInfoProvidersInitMsgs, tokenInfoProviderContractInitMsg)
	}

	initHydro := map[string]any{
		"round_length":      roundLength,
		"lock_epoch_length": roundLength,
		"tranches": []map[string]string{
			{
				"name":     "General tranche",
				"metadata": "General tranche metadata",
			},
			{
				"name":     "Consumer chains tranche",
				"metadata": "Consumer chains tranche metadata",
			},
		},
		"first_round_start":         strconv.FormatInt(firstRoundStartTime, 10),
		"max_locked_tokens":         "1000000000",
		"whitelist_admins":          []string{adminAddr},
		"initial_whitelist":         []string{adminAddr},
		"icq_managers":              []string{adminAddr},
		"max_deployment_duration":   12,
		"round_lock_power_schedule": [][]any{{1, "1"}, {2, "1.25"}, {3, "1.5"}, {6, "2"}, {12, "4"}},
		"token_info_providers":      tokenInfoProvidersInitMsgs,
		"gatekeeper":                gatekeeperInitMsg,
	}

	return s.InstantiateContract(codeId, initHydro, adminAddr, "Hydro Smart Contract")
}

func (s *HydroSuite) InstantiateTributeContract(codeId, hydroContractAddress, adminAddr string) string {
	initTribute := map[string]interface{}{
		"hydro_contract": hydroContractAddress,
	}

	return s.InstantiateContract(codeId, initTribute, adminAddr, "Tribute Smart Contract")
}

func (s *HydroSuite) InstantiateDaoVotingAdapterContract(codeId, hydroContractAddress, adminAddr string) string {
	initDaoVotingAdapter := map[string]interface{}{
		"hydro_contract": hydroContractAddress,
	}

	return s.InstantiateContract(codeId, initDaoVotingAdapter, adminAddr, "DAO Voting Adapter Smart Contract")
}

func (s *HydroSuite) InstantiateContract(
	codeId string,
	initMsg map[string]interface{},
	adminAddr string,
	label string,
) string {
	initJson, err := json.Marshal(initMsg)
	s.Require().NoError(err)

	txHash, err := s.NeutronChain.Validators[0].ExecTx(
		s.GetContext(),
		s.NeutronChain.ValidatorWallets[0].Moniker,
		"wasm", "instantiate", codeId, string(initJson), "--admin", adminAddr, "--label", label, "--gas", "auto",
	)
	s.Require().NoError(err)
	response, err := s.NeutronChain.Validators[0].TxHashToResponse(s.GetContext(), txHash)
	s.Require().NoError(err)
	contractAddr, found := getEvtAttribute(response.Events, wasmtypes.EventTypeInstantiate, wasmtypes.AttributeKeyContractAddr)
	s.Require().True(found)

	return contractAddr
}

func (s *HydroSuite) RegisterInterchainQueries(
	validators []string,
	contractAddr string,
	keyMoniker string,
) {
	icqs := map[string]interface{}{
		"create_icqs_for_validators": map[string]interface{}{
			"validators": validators,
		},
	}

	queryDeposit := len(validators) * chainsuite.NeutronMinQueryDeposit
	_, err := s.WasmExecuteTx(0, keyMoniker, icqs, contractAddr, []string{"--amount", strconv.Itoa(queryDeposit) + "untrn"})
	s.Require().NoError(err)
	// Wait for the relayer to retrieve the initial query data before proceeding with locking
	tCtx, cancelFn := context.WithTimeout(s.GetContext(), 30*chainsuite.CommitTimeout)
	defer cancelFn()

	dataSubmitted := false
	for tCtx.Err() == nil {
		time.Sleep(chainsuite.CommitTimeout)
		queryRes, _, err := s.NeutronChain.Validators[0].ExecQuery(
			s.GetContext(),
			"interchainqueries", "registered-queries", "--owners", contractAddr,
		)
		if err != nil {
			continue
		}

		var queryResponse chainsuite.QueryResponse
		err = json.Unmarshal([]byte(queryRes), &queryResponse)
		s.Require().NoError(err)
		s.Require().NotNil(queryResponse)

		dataSubmitted = true
		for _, query := range queryResponse.RegisteredQueries {
			if query.LastSubmittedResultLocalHeight == "0" {
				dataSubmitted = false
				break
			}
		}

		if dataSubmitted {
			break
		}

	}
	s.Require().True(dataSubmitted)
}

func (s *HydroSuite) WaitForQueryUpdate(contractAddr string, remoteHeight int64) {
	tCtx, cancelFn := context.WithTimeout(s.GetContext(), 30*chainsuite.CommitTimeout)
	defer cancelFn()

	dataSubmitted := false
	for tCtx.Err() == nil {
		time.Sleep(chainsuite.CommitTimeout)
		queryRes, _, err := s.NeutronChain.Validators[0].ExecQuery(
			s.GetContext(),
			"interchainqueries", "registered-queries", "--owners", contractAddr,
		)
		if err != nil {
			continue
		}

		var queryResponse chainsuite.QueryResponse
		err = json.Unmarshal([]byte(queryRes), &queryResponse)
		s.Require().NoError(err)
		s.Require().NotNil(queryResponse)

		dataSubmitted = true
		for _, query := range queryResponse.RegisteredQueries {
			submittedRemoteHeight, err := strconv.ParseInt(query.LastSubmittedResultRemoteHeight.RevisionHeight, 10, 64)
			s.Require().NoError(err)

			if submittedRemoteHeight < remoteHeight {
				dataSubmitted = false
				break
			}
		}

		if dataSubmitted {
			break
		}

	}
	s.Require().True(dataSubmitted)
}

func (s *HydroSuite) LockTokens(validatorIndex int, lockDuration int, lockAmount string, lockDenom string, proof map[string]interface{}, contractAddr string) error {
	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker

	return s.LockTokensWithWallet(validatorIndex, keyMoniker, lockDuration, lockAmount, lockDenom, proof, contractAddr)
}

// Locks tokens with the given keyMoniker that exists on the validator node with the given index.
// Differs from LockTokens() in previous one always using the validator key on the given node.
func (s *HydroSuite) LockTokensWithWallet(validatorIndex int, keyMoniker string, lockDuration int, lockAmount string, lockDenom string, proof map[string]any, contractAddr string) error {
	lockTxData := map[string]any{
		"lock_tokens": map[string]any{
			"lock_duration": lockDuration,
			"proof":         proof,
		},
	}

	response, err := s.WasmExecuteTx(validatorIndex, keyMoniker, lockTxData, contractAddr, []string{"--amount", lockAmount + lockDenom})
	if err != nil {
		return err
	}

	_, found := getEvtAttribute(response.Events, "wasm", "lock_id")
	if !found {
		return errors.New("error locking tokens")
	}

	return nil
}

func (s *HydroSuite) SubmitHydroProposal(
	validatorIndex int,
	contractAddr,
	proposalTitle string,
	trancheId int64,
	deploymentDuration int64,
	minimumAtomLiquidityRequest int,
) error {
	proposalTxData := map[string]interface{}{
		"create_proposal": map[string]interface{}{
			"tranche_id":                     trancheId,
			"title":                          proposalTitle,
			"description":                    "Proposal Description",
			"deployment_duration":            deploymentDuration,
			"minimum_atom_liquidity_request": strconv.Itoa(minimumAtomLiquidityRequest),
		},
	}

	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, proposalTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

type ProposalToLockups struct {
	ProposalId int   `json:"proposal_id"`
	LockIds    []int `json:"lock_ids"`
}

func (s *HydroSuite) VoteForHydroProposal(validatorIndex int, contractAddr string, trancheId int64, proposalsVotes []ProposalToLockups) error {
	voteTxData := map[string]interface{}{
		"vote": map[string]interface{}{
			"tranche_id":      trancheId,
			"proposals_votes": proposalsVotes,
		},
	}

	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, voteTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) GetProposalByTitle(contractAddr string, proposalTitle string, trancheId int64) (chainsuite.Proposal, error) {
	roundId := s.GetCurrentRound(contractAddr)
	// query proposal to get id
	queryData := map[string]interface{}{
		"round_proposals": map[string]interface{}{
			"round_id":   roundId,
			"tranche_id": trancheId,
			"start_from": 0,
			"limit":      100,
		},
	}

	response, err := s.QueryContractState(queryData, contractAddr)
	s.Require().NoError(err)

	var proposals chainsuite.ProposalData
	err = json.Unmarshal([]byte(response), &proposals)
	s.Require().NoError(err)

	for _, proposal := range proposals.Data.Proposals {
		if proposal.Title == proposalTitle {
			return proposal, nil
		}
	}

	return chainsuite.Proposal{}, fmt.Errorf("proposal is not found")
}

func (s *HydroSuite) GetCurrentRound(contractAddr string) int {
	queryData := map[string]interface{}{
		"current_round": map[string]interface{}{},
	}

	response, err := s.QueryContractState(queryData, contractAddr)
	s.Require().NoError(err)

	var roundData chainsuite.RoundData
	err = json.Unmarshal([]byte(response), &roundData)
	s.Require().NoError(err)

	return roundData.Data.RoundID
}

func (s *HydroSuite) WaitForRound(contractAddress string, roundId int) {
	tCtx, cancelFn := context.WithTimeout(s.GetContext(), 1000*chainsuite.CommitTimeout)
	defer cancelFn()

	roundReached := false
	for tCtx.Err() == nil {
		time.Sleep(chainsuite.CommitTimeout)

		if s.GetCurrentRound(contractAddress) >= roundId {
			roundReached = true
			break
		}
	}

	s.Require().True(roundReached)
}

func (s *HydroSuite) GetUserVotingPower(contractAddr string, address string) int64 {
	queryData := map[string]interface{}{
		"user_voting_power": map[string]interface{}{
			"address": address,
		},
	}

	response, err := s.QueryContractState(queryData, contractAddr)
	s.Require().NoError(err)

	var userData chainsuite.UserVotingPower
	err = json.Unmarshal([]byte(response), &userData)
	s.Require().NoError(err)

	return userData.Data.VotingPower
}

func (s *HydroSuite) GetRoundVotingPower(contractAddr string, roundId int64) int64 {
	queryData := map[string]interface{}{
		"round_total_voting_power": map[string]interface{}{
			"round_id": roundId,
		},
	}

	response, err := s.QueryContractState(queryData, contractAddr)
	s.Require().NoError(err)
	var roundData chainsuite.RoundVotingPower
	err = json.Unmarshal([]byte(response), &roundData)
	s.Require().NoError(err)

	roundPower, err := strconv.ParseInt(roundData.Data.VotingPower, 10, 64)
	s.Require().NoError(err)

	return roundPower
}

func (s *HydroSuite) GetGatekeeper(hydroContractAddr string) string {
	queryData := map[string]any{
		"gatekeeper": map[string]any{},
	}

	response, err := s.QueryContractState(queryData, hydroContractAddr)
	s.Require().NoError(err)

	var gatekeeperResponse chainsuite.Gatekeeper
	err = json.Unmarshal([]byte(response), &gatekeeperResponse)
	s.Require().NoError(err)

	return gatekeeperResponse.Data.Gatekeeper
}

func (s *HydroSuite) PauseTheHydroContract(keyMoniker string, contractAddr string) {
	pauseTxData := map[string]interface{}{
		"pause": map[string]interface{}{},
	}

	_, err := s.WasmExecuteTx(0, keyMoniker, pauseTxData, contractAddr, []string{})
	s.Require().NoError(err)
}

func (s *HydroSuite) UnlockTokens(contractAddr string) error {
	unlockTxData := map[string]interface{}{
		"unlock_tokens": map[string]interface{}{},
	}

	validatorIndex := 0
	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, unlockTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) RefreshLock(contractAddr string, new_lock_duration, lock_id int64) error {
	refreshTxData := map[string]interface{}{
		"refresh_lock_duration": map[string]interface{}{
			"lock_ids":      []int64{lock_id},
			"lock_duration": new_lock_duration,
		},
	}

	validatorIndex := 0
	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, refreshTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) UpdateMaxLockedTokens(contractAddr string, newMaxLockedTokens, activate_at int64) error {
	updateMaxLockedTokensTxData := map[string]interface{}{
		"update_config": map[string]interface{}{
			"activate_at":       strconv.FormatInt(activate_at, 10),
			"max_locked_tokens": newMaxLockedTokens,
		},
	}

	validatorIndex := 0
	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, updateMaxLockedTokensTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) WhitelistAccount(contractAddr string, address string) error {
	whiteListTxData := map[string]interface{}{
		"add_account_to_whitelist": map[string]interface{}{
			"address": address,
		},
	}

	validatorIndex := 0
	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, whiteListTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) RemoveFromWhitelist(contractAddr string, address string) error {
	whiteListTxData := map[string]interface{}{
		"remove_account_from_whitelist": map[string]interface{}{
			"address": address,
		},
	}

	validatorIndex := 0
	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, whiteListTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) AddTranche(contractAddr, name, metadata string) error {
	trancheTxData := map[string]interface{}{
		"add_tranche": map[string]interface{}{
			"tranche": map[string]interface{}{
				"name":     name,
				"metadata": metadata,
			},
		},
	}

	validatorIndex := 0
	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, trancheTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) EditTranche(contractAddr, name, metadata string, trancheId int) error {
	trancheTxData := map[string]interface{}{
		"edit_tranche": map[string]interface{}{
			"tranche_id": trancheId,
			"name":       name,
			"metadata":   metadata,
		},
	}

	validatorIndex := 0
	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, trancheTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) AddICQManager(contractAddr, address string) error {
	icqTxData := map[string]interface{}{
		"add_i_c_q_manager": map[string]interface{}{
			"address": address,
		},
	}

	validatorIndex := 0
	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, icqTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) RemoveICQManager(contractAddr, address string) error {
	icqTxData := map[string]interface{}{
		"remove_i_c_q_manager": map[string]interface{}{
			"address": address,
		},
	}

	validatorIndex := 0
	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, icqTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) WithdrawICQFunds(contractAddr string, amount int64) error {
	icqTxData := map[string]interface{}{
		"withdraw_i_cq_funds": map[string]interface{}{
			"amount": amount,
		},
	}

	validatorIndex := 0
	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, icqTxData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) AddLiquidityDeployment(validatorIndex int, contractAddr string, roundId, trancheId, proposalId int, deployedFunds types.Coin) error {
	addLiquidityData := map[string]interface{}{
		"add_liquidity_deployment": map[string]interface{}{
			"round_id":                roundId,
			"tranche_id":              trancheId,
			"proposal_id":             proposalId,
			"destinations":            []string{},
			"deployed_funds":          []types.Coin{deployedFunds},
			"funds_before_deployment": []string{},
			"total_rounds":            1,
			"remaining_rounds":        0,
		},
	}

	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, addLiquidityData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) SubmitTribute(validatorIndex, amount, round_id, trancheId, proposalId int, contractAddr string) (int, error) {
	txData := map[string]interface{}{
		"add_tribute": map[string]interface{}{
			"tranche_id":  trancheId,
			"round_id":    round_id,
			"proposal_id": proposalId,
		},
	}

	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	response, err := s.WasmExecuteTx(validatorIndex, keyMoniker, txData, contractAddr, []string{"--amount", strconv.Itoa(amount) + "untrn"})
	if err != nil {
		return 0, err
	}

	tributeIdStr, found := getEvtAttribute(response.Events, "wasm", "tribute_id")
	s.Require().True(found)

	tributeId, err := strconv.Atoi(tributeIdStr)
	s.Require().NoError(err)

	return tributeId, nil
}

func (s *HydroSuite) ClaimTribute(validatorIndex int, contractAddr, voterAddress string, roundId, trancheId, tributeId int) error {
	txData := map[string]interface{}{
		"claim_tribute": map[string]interface{}{
			"round_id":      roundId,
			"tranche_id":    trancheId,
			"tribute_id":    tributeId,
			"voter_address": voterAddress,
		},
	}

	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, txData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) RefundTribute(validatorIndex int, contractAddr string, roundId, trancheId, tributeId, proposalId int) error {
	txData := map[string]interface{}{
		"refund_tribute": map[string]interface{}{
			"round_id":    roundId,
			"tranche_id":  trancheId,
			"tribute_id":  tributeId,
			"proposal_id": proposalId,
		},
	}

	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, txData, contractAddr, []string{})
	if err != nil {
		return err
	}

	return nil
}

func (s *HydroSuite) RegisterGatekeeperStage(
	validatorIndex int, contractAddr string, activateAt int64,
	merkleRoot string, startNewEpoch bool, hrp *string,
) {
	txData := map[string]any{
		"register_stage": map[string]any{
			"activate_at":     strconv.Itoa(int(activateAt)),
			"merkle_root":     merkleRoot,
			"start_new_epoch": startNewEpoch,
			"hrp":             hrp,
		},
	}

	keyMoniker := s.NeutronChain.ValidatorWallets[validatorIndex].Moniker
	_, err := s.WasmExecuteTx(validatorIndex, keyMoniker, txData, contractAddr, []string{})
	s.Require().NoError(err)
}

func (s *HydroSuite) WaitForGatekeeperStage(contractAddress string, stageId int, timeout time.Duration) {
	tCtx, cancelFn := context.WithTimeout(s.GetContext(), timeout)
	defer cancelFn()

	queryData := map[string]any{
		"current_stage": map[string]any{},
	}

	stageReached := false
	for tCtx.Err() == nil {
		time.Sleep(chainsuite.CommitTimeout)

		response, err := s.QueryContractState(queryData, contractAddress)
		if err != nil {
			if strings.Contains(err.Error(), "Failed to get current stage.") {
				continue
			} else {
				s.Require().NoError(err)
			}
		}

		var stageData chainsuite.StageData
		err = json.Unmarshal(response, &stageData)
		s.Require().NoError(err)

		if stageData.Data.Stage.StageID == stageId {
			stageReached = true
			break
		}
	}

	s.Require().True(stageReached)
}

func (s *HydroSuite) QueryContractState(queryData map[string]any, contractAddr string) ([]byte, error) {
	queryJson, err := json.Marshal(queryData)
	s.Require().NoError(err)

	response, _, err := s.NeutronChain.Validators[0].ExecQuery(
		s.GetContext(),
		"wasm", "contract-state", "smart", contractAddr, string(queryJson),
	)

	// log
	s.Suite.T().Logf("Querying contract state: %s", string(queryJson))
	s.Suite.T().Logf("Response: %s", response)

	return response, err
}

func (s *HydroSuite) WasmExecuteTx(validatorIndex int, keyMoniker string, txData map[string]interface{}, contractAddr string, flags []string) (*types.TxResponse, error) {
	txJson, err := json.Marshal(txData)
	if err != nil {
		return nil, err
	}

	cmdFlags := []string{"wasm", "execute", contractAddr, string(txJson), "--gas", "auto"}
	cmdFlags = append(cmdFlags, flags...)

	txHash, err := s.NeutronChain.Validators[validatorIndex].ExecTx(
		s.GetContext(),
		keyMoniker,
		cmdFlags...,
	)
	if err != nil {
		return nil, err
	}

	response, err := s.NeutronChain.Validators[validatorIndex].TxHashToResponse(s.GetContext(), txHash)
	if err != nil {
		return nil, err
	}
	if response.Code != 0 {
		return nil, fmt.Errorf("ExecuteNeutronTx failed. Error code:%d", response.Code)
	}

	return response, nil
}

// Runs the transaction in generate-only mode to build the resulting TX JSON file, which is then written to the
// container at the specified location. File is written bacause the initial use case is to generate a TX and then
// sign it in second step, for which the TX JSON file must be present in the filesystem.
func (s *HydroSuite) TxGenerateOnly(node *cosmos.ChainNode, keyName string, outputFile string, command ...string) error {
	command = node.BinCommand(append(command,
		"--generate-only",
		"--offline",
		"--account-number", "0",
		"--sequence", "0",
		"--from", keyName,
		"--keyring-backend", keyring.BackendTest,
	)...)

	stdout, _, err := node.Exec(s.GetContext(), command, node.Chain.Config().Env)
	if err != nil {
		return err
	}

	// Redirecting the output with "> out_file.json" doesn't work with the framework, so we need to write it like this
	err = node.WriteFile(s.GetContext(), stdout, outputFile)
	if err != nil {
		return err
	}

	return nil
}

// Given the signed TX JSON file prepares SignatureInfo to be attached to LockTokens transaction.
// Only used when the Gatekeeper current stage refers to external chain addresses.
func (s *HydroSuite) PrepareSignatureInfo(signedTxJson []byte) (*chainsuite.SignatureInfo, error) {
	var tx chainsuite.CosmosTx
	if err := json.Unmarshal(signedTxJson, &tx); err != nil {
		return nil, err
	}

	claimMsg, err := json.Marshal(tx.Body)
	if err != nil {
		return nil, err
	}

	signatureInput := chainsuite.CosmosSignature{
		PubKey:    tx.AuthInfo.SignerInfos[0].PublicKey.Key,
		Signature: tx.Signatures[0],
	}
	signature, err := json.Marshal(signatureInput)
	if err != nil {
		return nil, err
	}

	signatureInfo := chainsuite.SignatureInfo{
		ClaimMsg:  base64.StdEncoding.EncodeToString(claimMsg),
		Signature: base64.StdEncoding.EncodeToString(signature),
	}

	return &signatureInfo, nil
}

func (s *HydroSuite) GetContractCodePath(wasmCodeName string) string {
	return path.Join(s.NeutronChain.GetNode().HomeDir(), wasmCodeName)
}

func (s *HydroSuite) VerifyHistoricalVotingPowers(
	daoVotingAdapterContractAddr string,
	historicalHeight int64,
	expectedTotalPower int64,
	expectedUserPowers map[string]int64,
) {
	queryInput := map[string]interface{}{
		"total_power_at_height": map[string]interface{}{
			"height": historicalHeight,
		},
	}

	response, err := s.QueryContractState(queryInput, daoVotingAdapterContractAddr)
	s.Require().NoError(err)
	var totalPowerAtHeight chainsuite.TotalPowerAtHeight
	err = json.Unmarshal([]byte(response), &totalPowerAtHeight)
	s.Require().NoError(err)

	totalPower, err := strconv.ParseInt(totalPowerAtHeight.Data.Power, 10, 64)
	s.Require().NoError(err)
	s.Require().Equal(expectedTotalPower, totalPower)
	s.Require().Equal(historicalHeight, totalPowerAtHeight.Data.Height)

	for userAddress, expectedVotingPower := range expectedUserPowers {
		queryInput := map[string]interface{}{
			"voting_power_at_height": map[string]interface{}{
				"address": userAddress,
				"height":  historicalHeight,
			},
		}

		response, err := s.QueryContractState(queryInput, daoVotingAdapterContractAddr)
		s.Require().NoError(err)
		var votingPowerAtHeight chainsuite.VotingPowerAtHeight
		err = json.Unmarshal([]byte(response), &votingPowerAtHeight)
		s.Require().NoError(err)

		votingPower, err := strconv.ParseInt(votingPowerAtHeight.Data.Power, 10, 64)
		s.Require().NoError(err)
		s.Require().Equal(expectedVotingPower, votingPower)
		s.Require().Equal(historicalHeight, votingPowerAtHeight.Data.Height)
	}
}

func getEvtAttribute(events []abci.Event, evtType string, key string) (string, bool) {
	for _, evt := range events {
		if evt.GetType() == evtType {
			for _, attr := range evt.Attributes {
				if attr.Key == key {
					return attr.Value, true
				}
			}
		}
	}

	return "", false
}
