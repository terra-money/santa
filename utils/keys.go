package utils

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/keys"
	ckeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	bip39 "github.com/cosmos/go-bip39"
)

// GetKeys returns key list
func (g Generator) GetKeys() (out []byte, err error) {

	kb, err := keys.NewKeyBaseFromDir(g.KeyDir)
	if err != nil {
		return nil, err
	}

	infos, err := kb.List()
	if err != nil {
		return nil, err
	}

	if len(infos) == 0 {
		return []byte{}, nil
	}

	keysOutput, err := ckeys.Bech32KeysOutput(infos)
	if err != nil {
		return nil, err
	}

	return json.Marshal(keysOutput)
}

// AddNewKey appends new key
func (g Generator) AddNewKey(name, password, mnemonic string, oldHdPath bool) (out []byte, err error) {
	kb, err := keys.NewKeyBaseFromDir(g.KeyDir)
	if err != nil {
		return
	}

	if name == "" || password == "" {
		err = fmt.Errorf("must include both password and name with request")
		return
	}

	// if mnemonic is empty, generate one
	if mnemonic == "" {
		_, mnemonic, _ = ckeys.NewInMemory().CreateMnemonic("inmemorykey", ckeys.English, "123456789", ckeys.Secp256k1)
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		err = fmt.Errorf("invalid mnemonic")
		return
	}

	_, err = kb.Get(name)
	if err == nil {
		err = fmt.Errorf("key %s already exists", name)
		return
	}

	account := uint32(0)
	index := uint32(0)
	coinType := uint32(330)
	if oldHdPath {
		coinType = uint32(118)
	}

	hdParams := hd.NewFundraiserParams(account, coinType, index)
	info, err := kb.Derive(name, mnemonic, "", password, *hdParams)

	if err != nil {
		return
	}

	keyOutput, err := ckeys.Bech32KeyOutput(info)
	if err != nil {
		return
	}

	keyOutput.Mnemonic = mnemonic

	out, err = json.Marshal(keyOutput)
	if err != nil {
		return
	}

	return
}

// GetKey is the handler for the GET /keys/{name}
func (g Generator) GetKey(name, bechPrefix string) (out []byte, err error) {
	kb, err := keys.NewKeyBaseFromDir(g.KeyDir)
	if err != nil {
		return
	}

	if bechPrefix == "" {
		bechPrefix = "acc"
	}

	bechKeyOut, err := getBechKeyOut(bechPrefix)
	if err != nil {
		return
	}

	info, err := kb.Get(name)
	if err != nil {
		return
	}

	keyOutput, err := bechKeyOut(info)
	if err != nil {
		return
	}

	out, err = json.Marshal(keyOutput)
	if err != nil {
		return
	}

	return
}

type bechKeyOutFn func(keyInfo ckeys.Info) (ckeys.KeyOutput, error)

func getBechKeyOut(bechPrefix string) (bechKeyOutFn, error) {
	switch bechPrefix {
	case "acc":
		return ckeys.Bech32KeyOutput, nil
	case "val":
		return ckeys.Bech32ValKeyOutput, nil
	case "cons":
		return ckeys.Bech32ConsKeyOutput, nil
	}

	return nil, fmt.Errorf("invalid Bech32 prefix encoding provided: %s", bechPrefix)
}

// UpdateKey update key password
func (g Generator) UpdateKey(name, newPassword, oldPassword string) (err error) {

	kb, err := keys.NewKeyBaseFromDir(g.KeyDir)
	if err != nil {
		return
	}

	err = kb.Update(name, oldPassword, func() (string, error) { return newPassword, nil })
	if err != nil {
		return
	}

	return
}

// DeleteKey is the handler for the DELETE /keys/{name}
func (g Generator) DeleteKey(name, password string) (err error) {
	kb, err := keys.NewKeyBaseFromDir(g.KeyDir)
	if err != nil {
		return
	}

	err = kb.Delete(name, password, false)
	if err != nil {
		return
	}

	return
}
