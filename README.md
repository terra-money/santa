# Build & Install
```
$ git clone https://github.com/terra-project/feegiver.git
$ git checkout master
```

## Install
```
$ make install
```
## make config
```
$ feegiver config
```
## change config
```
$ vim ~/.feegiver/config.yaml
```

## Add Key
```
$ feegiver keys add yun                        
Enter a passphrase to encrypt your key to disk:
Repeat the passphrase:
{"name":"yun","type":"local","address":"terra1a26sc2vqs20hfx239kejhd88v6cl87yfswvk0t","pubkey":"terrapub1addwnpepqwspkmsl724h9azfvgqgs8jkyuyyr3d6eme432afvlvulk3al0mwwnxwlxv","mnemonic":"decade urge pond sustain unit film milk sunny wash accuse profit staff what black problem treat velvet metal leg review math history juice soccer"}
```

## Recover Key
```
$ feegiver keys add yun --recover
Enter a passphrase to encrypt your key to disk:
Repeat the passphrase:
> Enter your bip39 mnemonic
theory fat merge under hungry utility toss much trend turkey degree glare bread connect trend grain silk toe pupil crouch innocent pause zero shove

{"name":"yun","type":"local","address":"terra1gn37dh0jl4zu4fp48d8y4c0hqs9cel83x7st7v","pubkey":"terrapub1addwnpepqfukqlgu8chxwwqns6adgjxvfny6y6tcvqmqkkn2xk6e6kefdaggvh4j7f0","mnemonic":"spatial fantasy weekend romance entire million celery final moon solid route theory way hockey north trigger advice balcony melody fabric alter bullet twice push"}
```

## Add Key with Old HD Path
```
$ feegiver keys add yun --recover --old-hd-path 
Enter a passphrase to encrypt your key to disk:
Repeat the passphrase:
> Enter your bip39 mnemonic 
flash until glimpse chase cradle adjust brick view uncover analyst test pact sponsor example item victory memory attract visual hover pink meadow mosquito torch

{"name":"yun","type":"local","address":"terra1gn37dh0jl4zu4fp48d8y4c0hqs9cel83x7st7v","pubkey":"terrapub1addwnpepqfukqlgu8chxwwqns6adgjxvfny6y6tcvqmqkkn2xk6e6kefdaggvh4j7f0","mnemonic":"spatial fantasy weekend romance entire million celery final moon solid route theory way hockey north trigger advice balcony melody fabric alter bullet twice push"}
```

## Start Server
```
$ feegiver start yun              
Enter the passphrase:
```
