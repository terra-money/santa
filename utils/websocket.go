package utils

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	amino "github.com/tendermint/go-amino"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpc "github.com/tendermint/tendermint/rpc/lib/types"
	abci "github.com/tendermint/tendermint/types"
)

var aminoCdc = amino.NewCodec()

func init() {
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
				log.Fatal("[Fail] to send tx", err)
			}

			log.Printf("[Success] Height: %d,\tTxHash: %s\n", blockEvent.Block.Height, txHash)
		}

		if isTest {
			break
		}
	}
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
