package utils

import (
	"os"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestBuildSendTx(t *testing.T) {
	//----------------------------------------
	app := setupWithNotExistsAccount(t)
	defer os.RemoveAll(app.KeyDir)

	// not exists
	_, err := app.SendTx(testChainID)
	require.Error(t, err)

	//----------------------------------------
	app2 := setup(t)
	defer os.RemoveAll(app2.KeyDir)

	// not exist key
	_, err = app2.SendTx(testChainID)
	require.Error(t, err)

	//----------------------------------------
	app3 := setupWithNoBalanceAccount(t)
	defer os.RemoveAll(app3.KeyDir)

	// not enough balance to pay fee
	_, err = app3.SendTx(testChainID)
	require.Error(t, err)

	//----------------------------------------
	app4 := setupWithPlentyBalanceAccount(t)
	defer os.RemoveAll(app4.KeyDir)

	// invalid fee amount
	app4.FeeAmount = "1"
	_, err = app4.SendTx(testChainID)
	require.Error(t, err)

	//----------------------------------------
	app5 := setupWithPlentyBalanceAccount(t)
	defer os.RemoveAll(app5.KeyDir)
	_, err = app5.SendTx(testChainID)
	require.NoError(t, err)
}