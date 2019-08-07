package utils

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/terra-project/core/x/treasury"
)

// SimulateGas simulates gas for a transaction
func (g Generator) SimulateGas(cdc *codec.Codec, txbytes []byte) (res uint64, err error) {
	result, err := rpcclient.NewHTTP(g.Node, "/websocket").ABCIQueryWithOptions(
		"/app/simulate",
		cmn.HexBytes(txbytes),
		rpcclient.ABCIQueryOptions{},
	)

	if err != nil {
		return
	}

	if !result.Response.IsOK() {
		return 0, errors.New(result.Response.Log)
	}

	var simulationResult sdk.Result
	if err := cdc.UnmarshalBinaryLengthPrefixed(result.Response.Value, &simulationResult); err != nil {
		return 0, err
	}

	return simulationResult.GasUsed, nil
}

// QueryAccount query account info
func (g Generator) QueryAccount(cdc *codec.Codec, accAddr sdk.AccAddress) (acc auth.Account, err error) {
	bz, err := cdc.MarshalJSON(auth.NewQueryAccountParams(accAddr))
	if err != nil {
		return
	}

	result, err := rpcclient.NewHTTP(g.Node, "/websocket").ABCIQueryWithOptions(
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

// QueryTaxRate query tax rate
func (g Generator) QueryTaxRate(cdc *codec.Codec, epoch int64) (taxRate sdk.Dec, err error) {
	result, err := rpcclient.NewHTTP(g.Node, "/websocket").ABCIQueryWithOptions(
		fmt.Sprintf("custom/treasury/tax-rate/%d", epoch),
		[]byte{},
		rpcclient.ABCIQueryOptions{},
	)

	if err != nil {
		return
	}

	if !result.Response.IsOK() {
		err = errors.New(result.Response.Log)
		return
	}

	var resp treasury.QueryTaxRateResponse
	if err = cdc.UnmarshalJSON(result.Response.Value, &resp); err != nil {
		return
	}

	taxRate = resp.TaxRate
	return
}

// QueryTaxCap query tax cap
func (g Generator) QueryTaxCap(cdc *codec.Codec, denom string) (taxCap sdk.Int, err error) {
	result, err := rpcclient.NewHTTP(g.Node, "/websocket").ABCIQueryWithOptions(
		fmt.Sprintf("custom/treasury/tax-cap/%s", denom),
		[]byte{},
		rpcclient.ABCIQueryOptions{},
	)

	if err != nil {
		return
	}

	if !result.Response.IsOK() {
		err = errors.New(result.Response.Log)
		return
	}

	var resp treasury.QueryTaxCapResponse
	if err = cdc.UnmarshalJSON(result.Response.Value, &resp); err != nil {
		return
	}

	taxCap = resp.TaxCap
	return
}

// BroadcastTx broadcast tx
func (g Generator) BroadcastTx(tx auth.StdTx) (txHash string, err error) {
	txBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return
	}

	result, err := rpcclient.NewHTTP(g.Node, "/websocket").BroadcastTxSync(txBytes)
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
