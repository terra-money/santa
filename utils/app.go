package utils

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	txbldr "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/x/bank"

	core "github.com/terra-project/core/app"
)

var cdc *codec.Codec

func init() {
	cdc = core.MakeCodec()

	config := sdk.GetConfig()
	config.SetCoinType(330)
	config.SetFullFundraiserPath("44'/330'/0'/0/0")
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.Seal()
}

// SantaApp tx app
type SantaApp struct {
	KeyDir string `json:"key_dir" yaml:"key_dir"`
	Node   string `json:"node" yaml:"node"`

	KeyName     string `json:"key_name" yaml:"key_name,omitempty"`
	KeyPassword string `json:"key_password" yaml:"key_password,omitempty"`

	TriggerInterval string `json:"trigger_interval" yaml:"trigger_interval"`
	FeeAmount       string `json:"fee_amount" yaml:"fee_amount"`

	Version string `yaml:"version,omitempty"`
	Commit  string `yaml:"commit,omitempty"`
	Branch  string `yaml:"branch,omitempty"`
}

// Send BankSend Tx
func (app SantaApp) SendTx(chainID string) (txHash string, err error) {
	kb, err := keys.NewKeyBaseFromDir(app.KeyDir)
	if err != nil {
		return
	}

	info, err := kb.Get(app.KeyName)
	if err != nil {
		return
	}

	acc, err := app.QueryAccount(cdc, info.GetAddress())
	if err != nil {
		return
	}

	targetFeeCoin, err := sdk.ParseCoin(app.FeeAmount)
	if err != nil {
		return
	}

	targetFeeDenom := targetFeeCoin.Denom
	spendableCoins := acc.SpendableCoins(time.Now())
	spendableAmount := spendableCoins.AmountOf(targetFeeDenom)

	if spendableAmount.LT(targetFeeCoin.Amount) {
		err = fmt.Errorf("not enough balance to distribute fee")
		return
	}

	// NOTE - no tax will be charged
	sendAmount := sdk.NewInt(1)
	sendCoins := sdk.NewCoins(sdk.NewCoin(targetFeeDenom, sendAmount))
	feeCoins := sdk.NewCoins(targetFeeCoin)
	stdTx := auth.NewStdTx(
		[]sdk.Msg{bank.NewMsgSend(acc.GetAddress(), acc.GetAddress(), sendCoins)},
		auth.NewStdFee(100000, feeCoins),
		[]auth.StdSignature{},
		"",
	)

	signedTx, err := app.signTx(stdTx, acc, chainID)

	txHash, err = app.BroadcastTx(signedTx)
	return
}

func (app SantaApp) signTx(stdTx auth.StdTx, acc auth.Account, chainID string) (signedTx auth.StdTx, err error) {

	kb, err := keys.NewKeyBaseFromDir(app.KeyDir)
	if err != nil {
		return
	}

	stdSign := txbldr.StdSignMsg{
		Memo:          stdTx.Memo,
		Msgs:          stdTx.Msgs,
		ChainID:       chainID,
		AccountNumber: uint64(acc.GetAccountNumber()),
		Sequence:      uint64(acc.GetSequence()),
		Fee: auth.StdFee{
			Amount: stdTx.Fee.Amount,
			Gas:    uint64(stdTx.Fee.Gas),
		},
	}

	sigBytes, pubkey, err := kb.Sign(app.KeyName, app.KeyPassword, sdk.MustSortJSON(cdc.MustMarshalJSON(stdSign)))
	if err != nil {
		return
	}

	sigs := append(stdTx.GetSignatures(), auth.StdSignature{
		PubKey:    pubkey,
		Signature: sigBytes,
	})

	signedTx = auth.NewStdTx(stdTx.GetMsgs(), stdTx.Fee, sigs, stdTx.GetMemo())
	return
}
