package utils

import (
	"os"
	"bytes"
	"log"
	"strings"

	"testing"
	"github.com/stretchr/testify/require"
)
func TestWebsocketListen(t *testing.T) {
	app := setupWithPlentyBalanceAccount(t)

	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	app.TriggerInterval = "1"
	app.ListenNewBLock(true)
	
	result := logBuf.String()
	require.True(t, strings.Contains(result, "[Success]"))
}