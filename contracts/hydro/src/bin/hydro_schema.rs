use cosmwasm_schema::write_api;
use hydro::{
    msg::{ExecuteMsg, InstantiateMsg},
    query::QueryMsg,
};

fn main() {
    write_api! {
        instantiate: InstantiateMsg,
        execute: ExecuteMsg,
        query: QueryMsg,
    };
}
