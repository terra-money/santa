import axios, { AxiosResponse } from "axios"

class RestInterface {
  lcdAddress: string
  
  constructor(lcdAddress: string) {
    this.lcdAddress = lcdAddress
  }

  async loadValidators(): Promise<void | Array<Validator>> {
    const queryValidatorsURL = `${this.lcdAddress}/staking/validators`
    try {
      const res = await axios.get(queryValidatorsURL)
      return res.data ? res.data : []
    } catch (error) {
      console.error("[LoadValidators]", error);
    }
  }

  async loadValidatorRewards(valAddress: string): Promise<void | Array<Coin>> {
    const queryValidatorRewardsInfoURL = `${this.lcdAddress}/distribution/validators/${valAddress}`
    try {
      const res = await axios.get(queryValidatorRewardsInfoURL)
      // foundation validators have 100% commission, so ignore self_bond_rewards
      return res.data ? res.data.val_commission : []
    } catch (error) {
      console.error("[LoadValidatorRewardsInfoURL]", error);
    }
  }

  async loadDelegations(valAddress: string): Promise<void | Array<Delegation>> {
    const queryDelegatorsURL = `${this.lcdAddress}/staking/validators/${valAddress}/delegations`
    try {
      const res = await axios.get(queryDelegatorsURL)
      return res.data ? res.data : []
    } catch (error) {
      console.error("[LoadDelegations]", error);
    }
  }

  async loadDelegatorRewards(delAddress: string): Promise<void | Array<Coin>> {
    const queryDelegatorRewardsURL = `${this.lcdAddress}/distribution/delegators/${delAddress}/rewards`
    try {
      const res = await axios.get(queryDelegatorRewardsURL)
      return res.data ? res.data : []
    } catch (error) {
      console.error("[LoadDelegatorRewards]", error);
    }
  }
}

export interface Validator {
  operator_address: string
  consensus_pubkey: string
  jailed: boolean
  status: number
  tokens: string
  delegator_shares: string
  description: [Object]
  unbonding_height: string
  unbonding_time: string
  commission: [Object]
  min_self_delegation: string
}

export interface Delegation {
  delegator_address: string
  validator_address: string
  shares: string
}

export interface Coin {
  denom: string
  amount: string
}

export interface ValidatorRewardsInfo {
  operator_address:string
  self_bond_rewards:[Coin]
  val_commission:[Coin]
}

export default RestInterface
