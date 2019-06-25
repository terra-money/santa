"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const core = require("@terra-money/core");
const rest_1 = require("./utils/rest");
const big_js_1 = require("big.js");
const fs = require("fs");
process.argv = process.argv.slice(2);
const argv = {};
for (let i in process.argv) {
    const tmp = process.argv[i].split('=');
    argv[tmp[0]] = tmp[1];
}
const lcdURL = argv['lcd'] || "http://127.0.0.1:1317";
const outputPath = argv['output'] || "./unsignedTx.json";
const logLevel = argv['log'] || 'debug';
console.log(`lcd path: ${lcdURL}`);
console.log(`output file path: ${outputPath}`);
console.log(`\n`);
const rest = new rest_1.default(lcdURL);
const foundationAddress = "terra1dp0taj85ruc299rkdvzp4z5pfg6z6swaed74e6";
const goliathValAddress = "terravaloper163phlen6dn7sp9khhjar2gqqx6kga0ly8d7h9g";
const marineValAddress = "terravaloper1d3hatwcsvkktgwp3elglw9glca0h42yg6xy4lp";
const ghostValAddress = "terravaloper1rgu3qmm6rllfxlrfk94pgxa0jm37902dynqehm";
const wraithValAddress = "terravaloper1eutun6vh83lmyq0wmyf9vgghvurze2xanl9sq6";
const filterValAddresses = [
    goliathValAddress,
    marineValAddress,
    ghostValAddress,
    wraithValAddress,
];
const filterAddresses = [
    foundationAddress,
    core.convertValAddressToAccAddress(goliathValAddress),
    core.convertValAddressToAccAddress(marineValAddress),
    core.convertValAddressToAccAddress(ghostValAddress),
    core.convertValAddressToAccAddress(wraithValAddress),
];
async function loadFoundationRewards() {
    const promises = [];
    promises.push(rest.loadDelegatorRewards(foundationAddress));
    promises.push(rest.loadValidatorRewards(goliathValAddress));
    promises.push(rest.loadValidatorRewards(marineValAddress));
    promises.push(rest.loadValidatorRewards(ghostValAddress));
    promises.push(rest.loadValidatorRewards(wraithValAddress));
    const rewardMap = {};
    await Promise.all(promises)
        .then(res => {
        for (let i in res) {
            const rewards = res[i];
            if (rewards && rewards.length > 0) {
                for (let j in rewards) {
                    const denom = rewards[j].denom;
                    const amount = rewards[j].amount.split('.')[0];
                    if (rewardMap[denom]) {
                        rewardMap[denom] = big_js_1.default(rewardMap[denom]).plus(amount).toString();
                    }
                    else {
                        rewardMap[denom] = amount;
                    }
                }
            }
        }
    });
    const totalRewards = [];
    for (let denom in rewardMap) {
        totalRewards.push({
            denom: denom,
            amount: rewardMap[denom]
        });
    }
    return totalRewards;
}
const validatorBonusRate = 0.2;
function computeValidatorsRewardRatio(rewardRatioMap, validators) {
    let totalBondedToken = big_js_1.default(0);
    for (let i in validators) {
        if (filterValAddresses.includes(validators[i].operator_address))
            continue;
        totalBondedToken = totalBondedToken.plus(validators[i].tokens);
    }
    for (let i in validators) {
        if (filterValAddresses.includes(validators[i].operator_address))
            continue;
        const validator = validators[i];
        const address = core.convertValAddressToAccAddress(validator.operator_address);
        rewardRatioMap[address] = big_js_1.default(validator.tokens).div(totalBondedToken).mul(validatorBonusRate).toPrecision(10);
    }
    return;
}
async function computeDelegatorRewardRatio(rewardRatioMap, validators) {
    const validatorDelegationMap = {};
    let totalBondedToken = big_js_1.default(0);
    for (let i in validators) {
        const validator = validators[i];
        const delegations = await rest.loadDelegations(validator.operator_address);
        if (!delegations)
            continue;
        validatorDelegationMap[validator.operator_address] = {
            tokens: validator.tokens,
            delegatorShares: validator.delegator_shares,
            delegations: delegations
        };
        for (let j in delegations) {
            const delegation = delegations[j];
            if (filterAddresses.includes(delegation.delegator_address))
                continue;
            const tokens = big_js_1.default(validator.tokens).mul(delegation.shares).div(validator.delegator_shares);
            totalBondedToken = totalBondedToken.plus(tokens);
        }
    }
    for (let v in validatorDelegationMap) {
        const info = validatorDelegationMap[v];
        for (let i in info.delegations) {
            const delegation = info.delegations[i];
            if (filterAddresses.includes(delegation.delegator_address))
                continue;
            const tokens = big_js_1.default(info.tokens).mul(delegation.shares).div(info.delegatorShares);
            const ratio = big_js_1.default(tokens).div(totalBondedToken).mul(1 - validatorBonusRate).toPrecision(10);
            if (rewardRatioMap[delegation.delegator_address]) {
                rewardRatioMap[delegation.delegator_address]
                    = big_js_1.default(rewardRatioMap[delegation.delegator_address])
                        .plus(ratio).toPrecision(10);
            }
            else {
                rewardRatioMap[delegation.delegator_address] = ratio;
            }
        }
    }
    return;
}
async function main() {
    const foundationRewards = await loadFoundationRewards();
    if (logLevel == 'debug') {
        console.debug(`Foundation Rewards:`, foundationRewards);
        console.debug(`\n`);
    }
    const validators = await rest.loadValidators();
    if (!validators) {
        console.error("no validator found");
        return process.exit(-1);
    }
    const rewardRatioMap = {};
    computeValidatorsRewardRatio(rewardRatioMap, validators);
    if (logLevel == 'debug') {
        console.debug(`Validator Bonus Rewards Map:`, rewardRatioMap);
        console.debug(`\n`);
    }
    await computeDelegatorRewardRatio(rewardRatioMap, validators);
    if (logLevel == 'debug') {
        console.debug(`Total Rewards Map:`, rewardRatioMap);
        console.debug(`\n`);
    }
    // Rotate reward ratio and build msg input
    let totalRatio = big_js_1.default(0);
    const outputs = [];
    for (let addr in rewardRatioMap) {
        const ratio = rewardRatioMap[addr];
        totalRatio = totalRatio.plus(ratio);
        const coins = [];
        for (let i in foundationRewards) {
            coins.push({
                denom: foundationRewards[i].denom,
                amount: big_js_1.default(foundationRewards[i].amount).mul(ratio).toFixed(0)
            });
        }
        outputs.push({
            address: addr,
            coins: coins
        });
    }
    if (totalRatio.gt(1)) {
        console.error(`Total Reward Ratio(${totalRatio}) is bigger than 1`);
        return process.exit(-1);
    }
    if (logLevel == 'debug') {
        console.debug(`Total Reward Ratio:${totalRatio}`);
        console.debug(`\n`);
    }
    const inputs = [];
    const coins = [];
    for (let i in foundationRewards) {
        coins.push({
            denom: foundationRewards[i].denom,
            amount: big_js_1.default(foundationRewards[i].amount).toFixed(0)
        });
    }
    inputs.push({
        address: foundationAddress,
        coins: coins
    });
    const multiSendMsg = core.buildMultiSend(inputs, outputs);
    const unSingedTx = core.buildStdTx([multiSendMsg], { gas: "1000000", amount: [{
                denom: "ukrw",
                amount: "1000000"
            }] }, "reward distribution");
    fs.writeFile(outputPath, JSON.stringify(unSingedTx, null, 4), function (err) {
        if (err) {
            console.error("Writing Failed", err);
            return process.exit(-1);
        }
        else {
            console.info("Writing Succeed", `Please check ${outputPath}`);
        }
    });
    return;
}
main();
//# sourceMappingURL=index.js.map