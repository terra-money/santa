package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	amino "github.com/tendermint/go-amino"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpc "github.com/tendermint/tendermint/rpc/lib/types"
	abci "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var aminoCdc = amino.NewCodec()

func init() {
	log.SetOutput(os.Stdout)
	ctypes.RegisterAmino(aminoCdc)
}

// ListenNewBlock listen new block and trigger sendtx
func (app SantaApp) ListenNewBLock(isTest bool) {
	triggerInterval, err := strconv.ParseInt(app.TriggerInterval, 10, 64)
	if err != nil {
		log.Fatal("Trigger interval should be number", err)
	}

	var scheme string
	var host string
	if strings.HasPrefix(app.Node, "https") {
		scheme = "wss"
		host = strings.TrimPrefix(app.Node, "https://")
	} else {
		scheme = "ws"
		host = strings.TrimPrefix(app.Node, "http://")
	}

	u := url.URL{Scheme: scheme, Host: host, Path: "/websocket"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	bz, err := json.Marshal(JSONRPC{
		JSONRPC: "2.0",
		Method:  "subscribe",
		ID:      "0",
		Params: Query{
			Query: "tm.event='NewBlock'",
		},
	})

	c.WriteMessage(websocket.TextMessage, bz)

	for {
		_, bz, _ := c.ReadMessage()

		var rpcResponse rpc.RPCResponse
		var resultEvent ctypes.ResultEvent
		// var newBlockEvent abci.EventDataNewBlock
		err := aminoCdc.UnmarshalJSON(bz, &rpcResponse)
		if err != nil {
			continue
		}

		err = aminoCdc.UnmarshalJSON(rpcResponse.Result, &resultEvent)
		blockEvent, ok := resultEvent.Data.(abci.EventDataNewBlock)
		if !ok {
			continue
		}

		if blockEvent.Block.Height%triggerInterval == 0 {
			txHash, err := app.SendTx(blockEvent.Block.ChainID)
			if err != nil {
				log.Printf("[Fail] to send tx: %s", err.Error())
				app.sendFailMessage(err.Error())
				continue
			}

			log.Printf("[Success] Height: %d,\tTxHash: %s\n", blockEvent.Block.Height, txHash)
			app.sendSuccessMessage(blockEvent.Block.Height)
		}

		if isTest {
			break
		}
	}
}

func (app SantaApp) sendFailMessage(msg string) (err error) {
	if app.FailWebHookURL != "" && app.FailWebHookDataKey != "" {
		notiBody, _ := json.Marshal(map[string]string{
			app.FailWebHookDataKey: fmt.Sprintf("[Fail] sending tx: %s", msg),
		})

		// send notification to slack
		_, err = http.Post(app.FailWebHookURL, "application/json", bytes.NewBuffer(notiBody))
		if err != nil {
			return err
		}
	}

	return nil
}

func (app SantaApp) sendSuccessMessage(height int64) (err error) {
	if app.SuccessWebHookURL != "" && app.SuccessWebHookDataKey != "" {
		feeCoin, err := sdk.ParseCoin(app.FeeAmount)
		if err != nil {
			return err
		}

		amount := sdk.NewDecFromIntWithPrec(feeCoin.Amount, 6).String()
		denom := strings.ToUpper(strings.TrimLeft(feeCoin.Denom, "u"))
		notiBody, _ := json.Marshal(map[string]string{
			app.SuccessWebHookDataKey: fmt.Sprintf("Hohoho! Santa just gave out %s %s to stakeholders in presents in block %d", amount, denom, height),
		})

		// send notification to slack
		_, err = http.Post(app.SuccessWebHookURL, "application/json", bytes.NewBuffer(notiBody))
		if err != nil {
			return err
		}
	}

	return nil
}

// no-lint
type Query struct {
	Query string `json:"query"`
}

// no-lint
type JSONRPC struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      string `json:"id"`
	Params  Query  `json:"params"`
}
