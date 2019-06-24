declare class RestInterface {
    lcdAddress: string;
    constructor(lcdAddress: string);
    loadValidators(): Promise<void | Array<Validator>>;
    loadValidatorRewards(valAddress: string): Promise<void | Array<Coin>>;
    loadDelegations(valAddress: string): Promise<void | Array<Delegation>>;
    loadDelegatorRewards(delAddress: string): Promise<void | Array<Coin>>;
}
export interface Validator {
    operator_address: string;
    consensus_pubkey: string;
    jailed: boolean;
    status: number;
    tokens: string;
    delegator_shares: string;
    description: [Object];
    unbonding_height: string;
    unbonding_time: string;
    commission: [Object];
    min_self_delegation: string;
}
export interface Delegation {
    delegator_address: string;
    validator_address: string;
    shares: string;
}
export interface Coin {
    denom: string;
    amount: string;
}
export interface ValidatorRewardsInfo {
    operator_address: string;
    self_bond_rewards: [Coin];
    val_commission: [Coin];
}
export default RestInterface;
