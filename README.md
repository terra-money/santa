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
$ feegiver keys add tmp 12345678 'flash until glimpse chase cradle adjust brick view uncover analyst test pact sponsor example item victory memory attract visual hover pink meadow mosquito torch'
```
## Add Key with Old HD Path
```
$ feegiver keys add tmp 12345678 'flash until glimpse chase cradle adjust brick view uncover analyst test pact sponsor example item victory memory attract visual hover pink meadow mosquito torch' --old-hd-path
```

## Start Server
```
$ feegiver start tmp 12345678
```
