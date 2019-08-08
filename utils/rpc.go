package utils

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// QueryAccount query account info
func (app SantaApp) QueryAccount(cdc *codec.Codec, accAddr sdk.AccAddress) (acc auth.Account, err error) {
	bz, err := cdc.MarshalJSON(auth.NewQueryAccountParams(accAddr))
	if err != nil {
		return
	}

	result, err := rpcclient.NewHTTP(app.Node, "/websocket").ABCIQueryWithOptions(
		"custom/acc/account",
		bz,
		rpcclient.ABCIQueryOptions{},
	)

	if err != nil {
		return
	}

	if !result.Response.IsOK() {
		err = errors.New(result.Response.Log)
		return
	}

	if err = cdc.UnmarshalJSON(result.Response.Value, &acc); err != nil {
		return
	}

	return
}

// BroadcastTx broadcast tx
func (app SantaApp) BroadcastTx(tx auth.StdTx) (txHash string, err error) {
	txBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return
	}

	result, err := rpcclient.NewHTTP(app.Node, "/websocket").BroadcastTxSync(txBytes)
	if err != nil {
		return
	}

	if result.Code != 0 {
		err = errors.New(result.Log)
		return
	}

	txHash = result.Hash.String()
	return
}
