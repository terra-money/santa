# Santa Claus

![banner](docs/terra-core.png)

Santa Claus subsidizes block rewards for Terra's Columbus mainnet. It is intended to subsidize the security of the network until Terra transaction volume (and attendent staking returns from transaction fees) has had sufficient time to mature. A detailed explanation on the [motivation](./MOTIVATION.md) for starting Santa Claus is explained.

## Build & Install

First, check out the rep

```
$ git clone https://github.com/terra-project/feegiver.git
$ git checkout master
```

### Install
```
$ make install
```
### Make config
```
$ feegiver config
```
### Change config
```
$ vim ~/.feegiver/config.yaml
```

### Add / Recover Key
```
$ feegiver keys add yun        
```


## Recover Key
```
$ feegiver keys add yun --recover
```



## Start Server
```
$ feegiver start yun              
Enter the passphrase:
```


## Status



## Mechanism


