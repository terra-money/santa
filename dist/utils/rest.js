"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const axios_1 = require("axios");
class RestInterface {
    constructor(lcdAddress) {
        this.lcdAddress = lcdAddress;
    }
    async loadValidators() {
        const queryValidatorsURL = `${this.lcdAddress}/staking/validators`;
        try {
            const res = await axios_1.default.get(queryValidatorsURL);
            return res.data ? res.data : [];
        }
        catch (error) {
            console.error("[LoadValidators]", error);
        }
    }
    async loadValidatorRewards(valAddress) {
        const queryValidatorRewardsInfoURL = `${this.lcdAddress}/distribution/validators/${valAddress}`;
        try {
            const res = await axios_1.default.get(queryValidatorRewardsInfoURL);
            // foundation validators have 100% commission, so ignore self_bond_rewards
            return res.data ? res.data.val_commission : [];
        }
        catch (error) {
            console.error("[LoadValidatorRewardsInfoURL]", error);
        }
    }
    async loadDelegations(valAddress) {
        const queryDelegatorsURL = `${this.lcdAddress}/staking/validators/${valAddress}/delegations`;
        try {
            const res = await axios_1.default.get(queryDelegatorsURL);
            return res.data ? res.data : [];
        }
        catch (error) {
            console.error("[LoadDelegations]", error);
        }
    }
    async loadDelegatorRewards(delAddress) {
        const queryDelegatorRewardsURL = `${this.lcdAddress}/distribution/delegators/${delAddress}/rewards`;
        try {
            const res = await axios_1.default.get(queryDelegatorRewardsURL);
            return res.data ? res.data : [];
        }
        catch (error) {
            console.error("[LoadDelegatorRewards]", error);
        }
    }
}
exports.default = RestInterface;
//# sourceMappingURL=rest.js.map