package utils

import (
	"os"
	"encoding/json"

	"testing"
	"github.com/stretchr/testify/require"

	ckeys "github.com/cosmos/cosmos-sdk/crypto/keys"
)

func TestAddKey(t *testing.T) {
	app := setup(t)
	defer os.RemoveAll(app.KeyDir)

	out, err := app.GetKeys()
	require.NoError(t, err)
	require.Equal(t, 0, len(out))

	// invalid name
	_, err = app.AddNewKey(testInvalidName, testPassword, testMnemonic, false)
	require.Error(t, err)

	// invalid password
	_, err = app.AddNewKey(testName, testInvalidPassword, testMnemonic, false)
	require.Error(t, err)

	// invalid mnemonic
	_, err = app.AddNewKey(testName, testPassword, testInvalidMnemonic, false)
	require.Error(t, err)

	// valid add
	out, err = app.AddNewKey(testName, testPassword, testMnemonic, false)
	require.NoError(t, err)

	var output ckeys.KeyOutput
	err = json.Unmarshal(out, &output)
	require.NoError(t, err)

	require.Equal(t, testName, output.Name)
	require.Equal(t, "terra1ch5ezwqftx8z8969l30j634wzs8772xfp5wur4", output.Address)
	require.Equal(t, "terrapub1addwnpepqvgqm0lpmjn4dga903dugvmw8qtzcush2agl8lx3xz2mxcm8vvwf2adw7e3", output.PubKey)
	require.Equal(t, testMnemonic, output.Mnemonic)
}

func TestGetKey(t *testing.T) {
	app := setup(t)
	defer os.RemoveAll(app.KeyDir)

	out, err := app.GetKey(testName, "acc")
	require.Error(t, err)
	require.Equal(t, 0, len(out))

	out, err = app.AddNewKey(testName, testPassword, testMnemonic, false)
	require.NoError(t, err)

	_, err = app.GetKey(testName, "acc")
	require.NoError(t, err)
}

func TestDeleteKey(t *testing.T) {
	app := setup(t)
	defer os.RemoveAll(app.KeyDir)

	out, err := app.AddNewKey(testName, testPassword, testMnemonic, false)
	require.NoError(t, err)

	out, err = app.GetKeys()
	require.NoError(t, err)

	var outputs []ckeys.KeyOutput
	err = json.Unmarshal(out, &outputs)
	require.Equal(t, 1, len(outputs))

	// not exist name
	err = app.DeleteKey(testInvalidName, testPassword)
	require.Error(t, err)

	// invalid password
	err = app.DeleteKey(testName, testInvalidPassword)
	require.Error(t, err)

	// valid
	err = app.DeleteKey(testName, testPassword)
	require.NoError(t, err)

	out, err = app.GetKeys()
	require.NoError(t, err)

	var outputs2 []ckeys.KeyOutput
	err = json.Unmarshal(out, &outputs2)
	require.Equal(t, 0, len(outputs2))
}

func TestUpdateKey(t *testing.T) {
	app := setup(t)
	defer os.RemoveAll(app.KeyDir)

	out, err := app.AddNewKey(testName, testPassword, testMnemonic, false)
	require.NoError(t, err)

	out, err = app.GetKeys()
	require.NoError(t, err)

	var outputs []ckeys.KeyOutput
	err = json.Unmarshal(out, &outputs)

	// not exist name
	err = app.UpdateKey(testInvalidName, testPassword, testUpdatePassword)
	require.Error(t, err)

	// invalid password
	err = app.UpdateKey(testName, testInvalidPassword, testUpdatePassword)
	require.Error(t, err)

	// valid
	err = app.UpdateKey(testName, testPassword, testUpdatePassword)
	require.NoError(t, err)

	// with old password
	err = app.DeleteKey(testName, testPassword)
	require.Error(t, err)

	// valid
	err = app.DeleteKey(testName, testUpdatePassword)
	require.NoError(t, err)
}
