package utils


import (
	"io/ioutil"

	"testing"
	"github.com/stretchr/testify/require"

)
const (
	testChainID = "vodka"

	testInvalidName = ""
	testName = "keyname"
	testName2 = "keyname2"
	testName3= "keyname3"

	testInvalidPassword = ""
	testPassword = "12345678"
	testPassword2 = "12345678910"
	testPassword3 = "1234567891011"
	testUpdatePassword = "123456789"

	testInvalidMnemonic = "crop"
	testMnemonic = "crop vivid nature drastic duck submit night innocent inflict when know divorce fan concert damp stand depart gauge area vanish legend clarify warfare discover"

	// no balance in vodka network
	testMnemonic2 = "baby reduce notice dice eight remember room avoid gravity patient cement unhappy consider exit beyond uncle oblige lamp fault open save rifle airport craft"
	
	// not exists on vodka network
	testMnemonic3 = "rally clarify museum turn slender one fruit october wedding wrong web spike slight domain double connect want flock sport powder gloom yard emerge album"
	

	testNode = "http://15.164.0.235:26657" // vodka network
	testFeeAmount = "1uluna"
)


func setup(t *testing.T) SantaApp {
	dir, err := ioutil.TempDir("/tmp", ".santa")
	require.NoError(t, err)

	app := SantaApp{
		KeyDir: dir,
	}

	return app
}

func setupWithPlentyBalanceAccount(t *testing.T) SantaApp {
	dir, err := ioutil.TempDir("/tmp", ".santa")
	require.NoError(t, err)

	app := SantaApp{
		KeyDir: dir,
		Node: testNode, // test vodka net
		
		KeyName: testName,
		KeyPassword: testPassword,

		FeeAmount: testFeeAmount,
	}

	_, err = app.AddNewKey(testName, testPassword, testMnemonic, false)
	require.NoError(t, err)

	return app
}

func setupWithNoBalanceAccount(t *testing.T) SantaApp {
	dir, err := ioutil.TempDir("/tmp", ".santa")
	require.NoError(t, err)

	app := SantaApp{
		KeyDir: dir,
		Node: testNode, // test vodka net
		
		KeyName: testName2,
		KeyPassword: testPassword2,

		FeeAmount: testFeeAmount,
	}

	_, err = app.AddNewKey(testName2, testPassword2, testMnemonic2, false)
	require.NoError(t, err)

	return app
}


func setupWithNotExistsAccount(t *testing.T) SantaApp {
	dir, err := ioutil.TempDir("/tmp", ".santa")
	require.NoError(t, err)

	app := SantaApp{
		KeyDir: dir,
		Node: testNode, // test vodka net
		
		KeyName: testName3,
		KeyPassword: testPassword3,

		FeeAmount: testFeeAmount,
	}

	_, err = app.AddNewKey(testName3, testPassword3, testMnemonic3, false)
	require.NoError(t, err)

	return app
}
