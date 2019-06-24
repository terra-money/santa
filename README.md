# Reward Distributor
TERRA Foundation Reward Distributor

This repository is built to distribute all foundation rewards to validator and delegators on terra network.

## Usage
```
npm run build
npm start
```
or
```
npm start -- lcd=https://lcd.terra.dev log=prod output=./unsigned.json
```

It requires active lcd url to get reward information.

It will make unsigned transaction output file (default `./unsignedTx.json`)


## Instructions
### For Both
terracli tx distr set-withdraw-addr --withdraw-to terra1437zllxmq9gag8acyt56rk7dkyrd2zvk9ts02p --from {both} --chain-id columbus-2 --gas-prices 0.015uluna --gas 18000

### For Foundation Validator
terracli tx distr withdraw-rewards --validator {validator} --from {validator} --commission

### For Foundation Delegator
terracli tx distr terracli tx distr withdraw-all-rewards --from {validator} 
