/**
* This file was automatically generated by @cosmwasm/ts-codegen@1.12.1.
* DO NOT MODIFY IT BY HAND. Instead, modify the source JSONSchema file,
* and run the @cosmwasm/ts-codegen generate command to regenerate this file.
*/

import { CosmWasmClient, SigningCosmWasmClient, ExecuteResult } from "@cosmjs/cosmwasm-stargate";
import { StdFee } from "@cosmjs/amino";
import { Timestamp, Uint64, Binary, Uint128, Decimal, TokenInfoProviderInstantiateMsg, InstantiateMsg, InstantiateContractMsg, TrancheInfo, ExecuteMsg, LockTokensProof, SignatureInfo, ProposalToLockups, Coin, QueryMsg, Addr, AllUserLockupsResponse, LockEntryWithPower, LockEntryV2, AllUserLockupsWithTrancheInfosResponse, LockupWithPerTrancheInfo, PerTrancheLockupInfo, AllVotesResponse, VoteEntry, Vote, AllVotesRoundTrancheResponse, CanLockDenomResponse, ConstantsResponse, Constants, RoundLockPowerSchedule, LockPowerEntry, CurrentRoundResponse, ExpiredUserLockupsResponse, GatekeeperResponse, ICQManagersResponse, LiquidityDeploymentResponse, LiquidityDeployment, ProposalResponse, Proposal, RegisteredValidatorQueriesResponse, RoundEndResponse, RoundProposalsResponse, RoundTotalVotingPowerResponse, RoundTrancheLiquidityDeploymentsResponse, SpecificUserLockupsResponse, SpecificUserLockupsWithTrancheInfosResponse, TokenInfoProvider, TokenInfoProvidersResponse, TokenInfoProviderLSM, TokenInfoProviderDerivative, TopNProposalsResponse, TotalLockedTokensResponse, TotalPowerAtHeightResponse, TranchesResponse, Tranche, UserVotesResponse, VoteWithPower, UserVotingPowerResponse, VotingPowerAtHeightResponse, WhitelistResponse, WhitelistAdminsResponse } from "./HydroBase.types";
export interface HydroBaseReadOnlyInterface {
  contractAddress: string;
  constants: () => Promise<ConstantsResponse>;
  tokenInfoProviders: () => Promise<TokenInfoProvidersResponse>;
  gatekeeper: () => Promise<GatekeeperResponse>;
  tranches: () => Promise<TranchesResponse>;
  allUserLockups: ({
    address,
    limit,
    startFrom
  }: {
    address: string;
    limit: number;
    startFrom: number;
  }) => Promise<AllUserLockupsResponse>;
  specificUserLockups: ({
    address,
    lockIds
  }: {
    address: string;
    lockIds: number[];
  }) => Promise<SpecificUserLockupsResponse>;
  allUserLockupsWithTrancheInfos: ({
    address,
    limit,
    startFrom
  }: {
    address: string;
    limit: number;
    startFrom: number;
  }) => Promise<AllUserLockupsWithTrancheInfosResponse>;
  specificUserLockupsWithTrancheInfos: ({
    address,
    lockIds
  }: {
    address: string;
    lockIds: number[];
  }) => Promise<SpecificUserLockupsWithTrancheInfosResponse>;
  expiredUserLockups: ({
    address,
    limit,
    startFrom
  }: {
    address: string;
    limit: number;
    startFrom: number;
  }) => Promise<ExpiredUserLockupsResponse>;
  userVotingPower: ({
    address
  }: {
    address: string;
  }) => Promise<UserVotingPowerResponse>;
  userVotes: ({
    address,
    roundId,
    trancheId
  }: {
    address: string;
    roundId: number;
    trancheId: number;
  }) => Promise<UserVotesResponse>;
  allVotes: ({
    limit,
    startFrom
  }: {
    limit: number;
    startFrom: number;
  }) => Promise<AllVotesResponse>;
  allVotesRoundTranche: ({
    limit,
    roundId,
    startFrom,
    trancheId
  }: {
    limit: number;
    roundId: number;
    startFrom: number;
    trancheId: number;
  }) => Promise<AllVotesRoundTrancheResponse>;
  currentRound: () => Promise<CurrentRoundResponse>;
  roundEnd: ({
    roundId
  }: {
    roundId: number;
  }) => Promise<RoundEndResponse>;
  roundTotalVotingPower: ({
    roundId
  }: {
    roundId: number;
  }) => Promise<RoundTotalVotingPowerResponse>;
  roundProposals: ({
    limit,
    roundId,
    startFrom,
    trancheId
  }: {
    limit: number;
    roundId: number;
    startFrom: number;
    trancheId: number;
  }) => Promise<RoundProposalsResponse>;
  proposal: ({
    proposalId,
    roundId,
    trancheId
  }: {
    proposalId: number;
    roundId: number;
    trancheId: number;
  }) => Promise<ProposalResponse>;
  topNProposals: ({
    numberOfProposals,
    roundId,
    trancheId
  }: {
    numberOfProposals: number;
    roundId: number;
    trancheId: number;
  }) => Promise<TopNProposalsResponse>;
  whitelist: () => Promise<WhitelistResponse>;
  whitelistAdmins: () => Promise<WhitelistAdminsResponse>;
  iCQManagers: () => Promise<ICQManagersResponse>;
  totalLockedTokens: () => Promise<TotalLockedTokensResponse>;
  registeredValidatorQueries: () => Promise<RegisteredValidatorQueriesResponse>;
  canLockDenom: ({
    tokenDenom
  }: {
    tokenDenom: string;
  }) => Promise<CanLockDenomResponse>;
  liquidityDeployment: ({
    proposalId,
    roundId,
    trancheId
  }: {
    proposalId: number;
    roundId: number;
    trancheId: number;
  }) => Promise<LiquidityDeploymentResponse>;
  roundTrancheLiquidityDeployments: ({
    limit,
    roundId,
    startFrom,
    trancheId
  }: {
    limit: number;
    roundId: number;
    startFrom: number;
    trancheId: number;
  }) => Promise<RoundTrancheLiquidityDeploymentsResponse>;
  totalPowerAtHeight: ({
    height
  }: {
    height?: number;
  }) => Promise<TotalPowerAtHeightResponse>;
  votingPowerAtHeight: ({
    address,
    height
  }: {
    address: string;
    height?: number;
  }) => Promise<VotingPowerAtHeightResponse>;
}
export class HydroBaseQueryClient implements HydroBaseReadOnlyInterface {
  client: CosmWasmClient;
  contractAddress: string;
  constructor(client: CosmWasmClient, contractAddress: string) {
    this.client = client;
    this.contractAddress = contractAddress;
    this.constants = this.constants.bind(this);
    this.tokenInfoProviders = this.tokenInfoProviders.bind(this);
    this.gatekeeper = this.gatekeeper.bind(this);
    this.tranches = this.tranches.bind(this);
    this.allUserLockups = this.allUserLockups.bind(this);
    this.specificUserLockups = this.specificUserLockups.bind(this);
    this.allUserLockupsWithTrancheInfos = this.allUserLockupsWithTrancheInfos.bind(this);
    this.specificUserLockupsWithTrancheInfos = this.specificUserLockupsWithTrancheInfos.bind(this);
    this.expiredUserLockups = this.expiredUserLockups.bind(this);
    this.userVotingPower = this.userVotingPower.bind(this);
    this.userVotes = this.userVotes.bind(this);
    this.allVotes = this.allVotes.bind(this);
    this.allVotesRoundTranche = this.allVotesRoundTranche.bind(this);
    this.currentRound = this.currentRound.bind(this);
    this.roundEnd = this.roundEnd.bind(this);
    this.roundTotalVotingPower = this.roundTotalVotingPower.bind(this);
    this.roundProposals = this.roundProposals.bind(this);
    this.proposal = this.proposal.bind(this);
    this.topNProposals = this.topNProposals.bind(this);
    this.whitelist = this.whitelist.bind(this);
    this.whitelistAdmins = this.whitelistAdmins.bind(this);
    this.iCQManagers = this.iCQManagers.bind(this);
    this.totalLockedTokens = this.totalLockedTokens.bind(this);
    this.registeredValidatorQueries = this.registeredValidatorQueries.bind(this);
    this.canLockDenom = this.canLockDenom.bind(this);
    this.liquidityDeployment = this.liquidityDeployment.bind(this);
    this.roundTrancheLiquidityDeployments = this.roundTrancheLiquidityDeployments.bind(this);
    this.totalPowerAtHeight = this.totalPowerAtHeight.bind(this);
    this.votingPowerAtHeight = this.votingPowerAtHeight.bind(this);
  }
  constants = async (): Promise<ConstantsResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      constants: {}
    });
  };
  tokenInfoProviders = async (): Promise<TokenInfoProvidersResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      token_info_providers: {}
    });
  };
  gatekeeper = async (): Promise<GatekeeperResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      gatekeeper: {}
    });
  };
  tranches = async (): Promise<TranchesResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      tranches: {}
    });
  };
  allUserLockups = async ({
    address,
    limit,
    startFrom
  }: {
    address: string;
    limit: number;
    startFrom: number;
  }): Promise<AllUserLockupsResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      all_user_lockups: {
        address,
        limit,
        start_from: startFrom
      }
    });
  };
  specificUserLockups = async ({
    address,
    lockIds
  }: {
    address: string;
    lockIds: number[];
  }): Promise<SpecificUserLockupsResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      specific_user_lockups: {
        address,
        lock_ids: lockIds
      }
    });
  };
  allUserLockupsWithTrancheInfos = async ({
    address,
    limit,
    startFrom
  }: {
    address: string;
    limit: number;
    startFrom: number;
  }): Promise<AllUserLockupsWithTrancheInfosResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      all_user_lockups_with_tranche_infos: {
        address,
        limit,
        start_from: startFrom
      }
    });
  };
  specificUserLockupsWithTrancheInfos = async ({
    address,
    lockIds
  }: {
    address: string;
    lockIds: number[];
  }): Promise<SpecificUserLockupsWithTrancheInfosResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      specific_user_lockups_with_tranche_infos: {
        address,
        lock_ids: lockIds
      }
    });
  };
  expiredUserLockups = async ({
    address,
    limit,
    startFrom
  }: {
    address: string;
    limit: number;
    startFrom: number;
  }): Promise<ExpiredUserLockupsResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      expired_user_lockups: {
        address,
        limit,
        start_from: startFrom
      }
    });
  };
  userVotingPower = async ({
    address
  }: {
    address: string;
  }): Promise<UserVotingPowerResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      user_voting_power: {
        address
      }
    });
  };
  userVotes = async ({
    address,
    roundId,
    trancheId
  }: {
    address: string;
    roundId: number;
    trancheId: number;
  }): Promise<UserVotesResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      user_votes: {
        address,
        round_id: roundId,
        tranche_id: trancheId
      }
    });
  };
  allVotes = async ({
    limit,
    startFrom
  }: {
    limit: number;
    startFrom: number;
  }): Promise<AllVotesResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      all_votes: {
        limit,
        start_from: startFrom
      }
    });
  };
  allVotesRoundTranche = async ({
    limit,
    roundId,
    startFrom,
    trancheId
  }: {
    limit: number;
    roundId: number;
    startFrom: number;
    trancheId: number;
  }): Promise<AllVotesRoundTrancheResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      all_votes_round_tranche: {
        limit,
        round_id: roundId,
        start_from: startFrom,
        tranche_id: trancheId
      }
    });
  };
  currentRound = async (): Promise<CurrentRoundResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      current_round: {}
    });
  };
  roundEnd = async ({
    roundId
  }: {
    roundId: number;
  }): Promise<RoundEndResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      round_end: {
        round_id: roundId
      }
    });
  };
  roundTotalVotingPower = async ({
    roundId
  }: {
    roundId: number;
  }): Promise<RoundTotalVotingPowerResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      round_total_voting_power: {
        round_id: roundId
      }
    });
  };
  roundProposals = async ({
    limit,
    roundId,
    startFrom,
    trancheId
  }: {
    limit: number;
    roundId: number;
    startFrom: number;
    trancheId: number;
  }): Promise<RoundProposalsResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      round_proposals: {
        limit,
        round_id: roundId,
        start_from: startFrom,
        tranche_id: trancheId
      }
    });
  };
  proposal = async ({
    proposalId,
    roundId,
    trancheId
  }: {
    proposalId: number;
    roundId: number;
    trancheId: number;
  }): Promise<ProposalResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      proposal: {
        proposal_id: proposalId,
        round_id: roundId,
        tranche_id: trancheId
      }
    });
  };
  topNProposals = async ({
    numberOfProposals,
    roundId,
    trancheId
  }: {
    numberOfProposals: number;
    roundId: number;
    trancheId: number;
  }): Promise<TopNProposalsResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      top_n_proposals: {
        number_of_proposals: numberOfProposals,
        round_id: roundId,
        tranche_id: trancheId
      }
    });
  };
  whitelist = async (): Promise<WhitelistResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      whitelist: {}
    });
  };
  whitelistAdmins = async (): Promise<WhitelistAdminsResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      whitelist_admins: {}
    });
  };
  iCQManagers = async (): Promise<ICQManagersResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      i_c_q_managers: {}
    });
  };
  totalLockedTokens = async (): Promise<TotalLockedTokensResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      total_locked_tokens: {}
    });
  };
  registeredValidatorQueries = async (): Promise<RegisteredValidatorQueriesResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      registered_validator_queries: {}
    });
  };
  canLockDenom = async ({
    tokenDenom
  }: {
    tokenDenom: string;
  }): Promise<CanLockDenomResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      can_lock_denom: {
        token_denom: tokenDenom
      }
    });
  };
  liquidityDeployment = async ({
    proposalId,
    roundId,
    trancheId
  }: {
    proposalId: number;
    roundId: number;
    trancheId: number;
  }): Promise<LiquidityDeploymentResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      liquidity_deployment: {
        proposal_id: proposalId,
        round_id: roundId,
        tranche_id: trancheId
      }
    });
  };
  roundTrancheLiquidityDeployments = async ({
    limit,
    roundId,
    startFrom,
    trancheId
  }: {
    limit: number;
    roundId: number;
    startFrom: number;
    trancheId: number;
  }): Promise<RoundTrancheLiquidityDeploymentsResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      round_tranche_liquidity_deployments: {
        limit,
        round_id: roundId,
        start_from: startFrom,
        tranche_id: trancheId
      }
    });
  };
  totalPowerAtHeight = async ({
    height
  }: {
    height?: number;
  }): Promise<TotalPowerAtHeightResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      total_power_at_height: {
        height
      }
    });
  };
  votingPowerAtHeight = async ({
    address,
    height
  }: {
    address: string;
    height?: number;
  }): Promise<VotingPowerAtHeightResponse> => {
    return this.client.queryContractSmart(this.contractAddress, {
      voting_power_at_height: {
        address,
        height
      }
    });
  };
}
export interface HydroBaseInterface extends HydroBaseReadOnlyInterface {
  contractAddress: string;
  sender: string;
  lockTokens: ({
    lockDuration,
    proof
  }: {
    lockDuration: number;
    proof?: LockTokensProof;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  refreshLockDuration: ({
    lockDuration,
    lockIds
  }: {
    lockDuration: number;
    lockIds: number[];
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  unlockTokens: ({
    lockIds
  }: {
    lockIds?: number[];
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  createProposal: ({
    deploymentDuration,
    description,
    minimumAtomLiquidityRequest,
    roundId,
    title,
    trancheId
  }: {
    deploymentDuration: number;
    description: string;
    minimumAtomLiquidityRequest: Uint128;
    roundId?: number;
    title: string;
    trancheId: number;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  vote: ({
    proposalsVotes,
    trancheId
  }: {
    proposalsVotes: ProposalToLockups[];
    trancheId: number;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  unvote: ({
    lockIds,
    trancheId
  }: {
    lockIds: number[];
    trancheId: number;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  addAccountToWhitelist: ({
    address
  }: {
    address: string;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  removeAccountFromWhitelist: ({
    address
  }: {
    address: string;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  updateConfig: ({
    activateAt,
    knownUsersCap,
    maxDeploymentDuration,
    maxLockedTokens
  }: {
    activateAt: Timestamp;
    knownUsersCap?: number;
    maxDeploymentDuration?: number;
    maxLockedTokens?: number;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  deleteConfigs: ({
    timestamps
  }: {
    timestamps: Timestamp[];
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  pause: (fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  addTranche: ({
    tranche
  }: {
    tranche: TrancheInfo;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  editTranche: ({
    trancheId,
    trancheMetadata,
    trancheName
  }: {
    trancheId: number;
    trancheMetadata?: string;
    trancheName?: string;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  createIcqsForValidators: ({
    validators
  }: {
    validators: string[];
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  addICQManager: ({
    address
  }: {
    address: string;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  removeICQManager: ({
    address
  }: {
    address: string;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  withdrawICQFunds: ({
    amount
  }: {
    amount: Uint128;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  addLiquidityDeployment: ({
    deployedFunds,
    destinations,
    fundsBeforeDeployment,
    proposalId,
    remainingRounds,
    roundId,
    totalRounds,
    trancheId
  }: {
    deployedFunds: Coin[];
    destinations: string[];
    fundsBeforeDeployment: Coin[];
    proposalId: number;
    remainingRounds: number;
    roundId: number;
    totalRounds: number;
    trancheId: number;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  removeLiquidityDeployment: ({
    proposalId,
    roundId,
    trancheId
  }: {
    proposalId: number;
    roundId: number;
    trancheId: number;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  updateTokenGroupRatio: ({
    newRatio,
    oldRatio,
    tokenGroupId
  }: {
    newRatio: Decimal;
    oldRatio: Decimal;
    tokenGroupId: string;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  addTokenInfoProvider: ({
    tokenInfoProvider
  }: {
    tokenInfoProvider: TokenInfoProviderInstantiateMsg;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  removeTokenInfoProvider: ({
    providerId
  }: {
    providerId: string;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
  setGatekeeper: ({
    gatekeeperAddr
  }: {
    gatekeeperAddr?: string;
  }, fee_?: number | StdFee | "auto", memo_?: string, funds_?: Coin[]) => Promise<ExecuteResult>;
}
export class HydroBaseClient extends HydroBaseQueryClient implements HydroBaseInterface {
  client: SigningCosmWasmClient;
  sender: string;
  contractAddress: string;
  constructor(client: SigningCosmWasmClient, sender: string, contractAddress: string) {
    super(client, contractAddress);
    this.client = client;
    this.sender = sender;
    this.contractAddress = contractAddress;
    this.lockTokens = this.lockTokens.bind(this);
    this.refreshLockDuration = this.refreshLockDuration.bind(this);
    this.unlockTokens = this.unlockTokens.bind(this);
    this.createProposal = this.createProposal.bind(this);
    this.vote = this.vote.bind(this);
    this.unvote = this.unvote.bind(this);
    this.addAccountToWhitelist = this.addAccountToWhitelist.bind(this);
    this.removeAccountFromWhitelist = this.removeAccountFromWhitelist.bind(this);
    this.updateConfig = this.updateConfig.bind(this);
    this.deleteConfigs = this.deleteConfigs.bind(this);
    this.pause = this.pause.bind(this);
    this.addTranche = this.addTranche.bind(this);
    this.editTranche = this.editTranche.bind(this);
    this.createIcqsForValidators = this.createIcqsForValidators.bind(this);
    this.addICQManager = this.addICQManager.bind(this);
    this.removeICQManager = this.removeICQManager.bind(this);
    this.withdrawICQFunds = this.withdrawICQFunds.bind(this);
    this.addLiquidityDeployment = this.addLiquidityDeployment.bind(this);
    this.removeLiquidityDeployment = this.removeLiquidityDeployment.bind(this);
    this.updateTokenGroupRatio = this.updateTokenGroupRatio.bind(this);
    this.addTokenInfoProvider = this.addTokenInfoProvider.bind(this);
    this.removeTokenInfoProvider = this.removeTokenInfoProvider.bind(this);
    this.setGatekeeper = this.setGatekeeper.bind(this);
  }
  lockTokens = async ({
    lockDuration,
    proof
  }: {
    lockDuration: number;
    proof?: LockTokensProof;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      lock_tokens: {
        lock_duration: lockDuration,
        proof
      }
    }, fee_, memo_, funds_);
  };
  refreshLockDuration = async ({
    lockDuration,
    lockIds
  }: {
    lockDuration: number;
    lockIds: number[];
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      refresh_lock_duration: {
        lock_duration: lockDuration,
        lock_ids: lockIds
      }
    }, fee_, memo_, funds_);
  };
  unlockTokens = async ({
    lockIds
  }: {
    lockIds?: number[];
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      unlock_tokens: {
        lock_ids: lockIds
      }
    }, fee_, memo_, funds_);
  };
  createProposal = async ({
    deploymentDuration,
    description,
    minimumAtomLiquidityRequest,
    roundId,
    title,
    trancheId
  }: {
    deploymentDuration: number;
    description: string;
    minimumAtomLiquidityRequest: Uint128;
    roundId?: number;
    title: string;
    trancheId: number;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      create_proposal: {
        deployment_duration: deploymentDuration,
        description,
        minimum_atom_liquidity_request: minimumAtomLiquidityRequest,
        round_id: roundId,
        title,
        tranche_id: trancheId
      }
    }, fee_, memo_, funds_);
  };
  vote = async ({
    proposalsVotes,
    trancheId
  }: {
    proposalsVotes: ProposalToLockups[];
    trancheId: number;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      vote: {
        proposals_votes: proposalsVotes,
        tranche_id: trancheId
      }
    }, fee_, memo_, funds_);
  };
  unvote = async ({
    lockIds,
    trancheId
  }: {
    lockIds: number[];
    trancheId: number;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      unvote: {
        lock_ids: lockIds,
        tranche_id: trancheId
      }
    }, fee_, memo_, funds_);
  };
  addAccountToWhitelist = async ({
    address
  }: {
    address: string;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      add_account_to_whitelist: {
        address
      }
    }, fee_, memo_, funds_);
  };
  removeAccountFromWhitelist = async ({
    address
  }: {
    address: string;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      remove_account_from_whitelist: {
        address
      }
    }, fee_, memo_, funds_);
  };
  updateConfig = async ({
    activateAt,
    knownUsersCap,
    maxDeploymentDuration,
    maxLockedTokens
  }: {
    activateAt: Timestamp;
    knownUsersCap?: number;
    maxDeploymentDuration?: number;
    maxLockedTokens?: number;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      update_config: {
        activate_at: activateAt,
        known_users_cap: knownUsersCap,
        max_deployment_duration: maxDeploymentDuration,
        max_locked_tokens: maxLockedTokens
      }
    }, fee_, memo_, funds_);
  };
  deleteConfigs = async ({
    timestamps
  }: {
    timestamps: Timestamp[];
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      delete_configs: {
        timestamps
      }
    }, fee_, memo_, funds_);
  };
  pause = async (fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      pause: {}
    }, fee_, memo_, funds_);
  };
  addTranche = async ({
    tranche
  }: {
    tranche: TrancheInfo;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      add_tranche: {
        tranche
      }
    }, fee_, memo_, funds_);
  };
  editTranche = async ({
    trancheId,
    trancheMetadata,
    trancheName
  }: {
    trancheId: number;
    trancheMetadata?: string;
    trancheName?: string;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      edit_tranche: {
        tranche_id: trancheId,
        tranche_metadata: trancheMetadata,
        tranche_name: trancheName
      }
    }, fee_, memo_, funds_);
  };
  createIcqsForValidators = async ({
    validators
  }: {
    validators: string[];
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      create_icqs_for_validators: {
        validators
      }
    }, fee_, memo_, funds_);
  };
  addICQManager = async ({
    address
  }: {
    address: string;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      add_i_c_q_manager: {
        address
      }
    }, fee_, memo_, funds_);
  };
  removeICQManager = async ({
    address
  }: {
    address: string;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      remove_i_c_q_manager: {
        address
      }
    }, fee_, memo_, funds_);
  };
  withdrawICQFunds = async ({
    amount
  }: {
    amount: Uint128;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      withdraw_i_c_q_funds: {
        amount
      }
    }, fee_, memo_, funds_);
  };
  addLiquidityDeployment = async ({
    deployedFunds,
    destinations,
    fundsBeforeDeployment,
    proposalId,
    remainingRounds,
    roundId,
    totalRounds,
    trancheId
  }: {
    deployedFunds: Coin[];
    destinations: string[];
    fundsBeforeDeployment: Coin[];
    proposalId: number;
    remainingRounds: number;
    roundId: number;
    totalRounds: number;
    trancheId: number;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      add_liquidity_deployment: {
        deployed_funds: deployedFunds,
        destinations,
        funds_before_deployment: fundsBeforeDeployment,
        proposal_id: proposalId,
        remaining_rounds: remainingRounds,
        round_id: roundId,
        total_rounds: totalRounds,
        tranche_id: trancheId
      }
    }, fee_, memo_, funds_);
  };
  removeLiquidityDeployment = async ({
    proposalId,
    roundId,
    trancheId
  }: {
    proposalId: number;
    roundId: number;
    trancheId: number;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      remove_liquidity_deployment: {
        proposal_id: proposalId,
        round_id: roundId,
        tranche_id: trancheId
      }
    }, fee_, memo_, funds_);
  };
  updateTokenGroupRatio = async ({
    newRatio,
    oldRatio,
    tokenGroupId
  }: {
    newRatio: Decimal;
    oldRatio: Decimal;
    tokenGroupId: string;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      update_token_group_ratio: {
        new_ratio: newRatio,
        old_ratio: oldRatio,
        token_group_id: tokenGroupId
      }
    }, fee_, memo_, funds_);
  };
  addTokenInfoProvider = async ({
    tokenInfoProvider
  }: {
    tokenInfoProvider: TokenInfoProviderInstantiateMsg;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      add_token_info_provider: {
        token_info_provider: tokenInfoProvider
      }
    }, fee_, memo_, funds_);
  };
  removeTokenInfoProvider = async ({
    providerId
  }: {
    providerId: string;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      remove_token_info_provider: {
        provider_id: providerId
      }
    }, fee_, memo_, funds_);
  };
  setGatekeeper = async ({
    gatekeeperAddr
  }: {
    gatekeeperAddr?: string;
  }, fee_: number | StdFee | "auto" = "auto", memo_?: string, funds_?: Coin[]): Promise<ExecuteResult> => {
    return await this.client.execute(this.sender, this.contractAddress, {
      set_gatekeeper: {
        gatekeeper_addr: gatekeeperAddr
      }
    }, fee_, memo_, funds_);
  };
}