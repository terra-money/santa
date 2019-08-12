package utils

import (
	"bytes"
	"log"
	"os"
	"strings"
	"time"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestWebsocketListen(t *testing.T) {
	app := setupWithPlentyBalanceAccount(t)

	// DO NOT TOUCH - wait next block
	time.Sleep(time.Second * 10)

	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stdout)

	app.TriggerInterval = "1"
	app.ListenNewBLock(true)

	result := logBuf.String()
	require.True(t, strings.Contains(result, "[Success]"))
}

func TestSendSuccess(t *testing.T) {
	app := setupWithPlentyBalanceAccount(t)
	err := app.sendSuccessMessage(10)
	require.NoError(t, err)

	app.SuccessWebHookURL = "https://google.co.kr"
	app.SuccessWebHookDataKey = "text"
	err = app.sendSuccessMessage(10)
	require.NoError(t, err)

	// invalid fee coin
	app.FeeAmount = "u1lu1n1a"
	err = app.sendSuccessMessage(10)
	require.Error(t, err)

	// invalid url
	app.FeeAmount = "1uluna"
	app.SuccessWebHookURL = "http://nonon"
	err = app.sendSuccessMessage(10)
	require.Error(t, err)
}

func TestSendFail(t *testing.T) {
	app := setupWithPlentyBalanceAccount(t)
	err := app.sendFailMessage("msg")
	require.NoError(t, err)

	app.FailWebHookURL = "https://google.co.kr"
	app.FailWebHookDataKey = "text"
	err = app.sendFailMessage("msg")
	require.NoError(t, err)

	// invalid url
	app.FeeAmount = "1uluna"
	app.FailWebHookURL = "http://nonon"
	err = app.sendFailMessage("msg")
	require.Error(t, err)
}