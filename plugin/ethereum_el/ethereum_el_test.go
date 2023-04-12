package ethereumel

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	babelSDK "github.com/umbracle/babel/sdk"

	"github.com/stretchr/testify/require"
	"github.com/umbracle/babel"
)

func TestEthereumEL(t *testing.T) {
	srv := babel.NewMockHttpServer(8101)
	defer srv.Stop()

	cb := map[string]string{}

	srv.Mux().HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		var req jsonRPCRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		if req.JsonRPC != "2.0" {
			w.Write([]byte(`invalid request`))
			return
		}

		resp, ok := cb[req.Method]
		if ok {
			w.Write([]byte(resp))
		}
	})

	cb["net_peerCount"] = `{
		"jsonrpc": "2.0",
		"id": 1,
		"result": "0x2"
	}`
	cb["eth_syncing"] = `{
		"jsonrpc": "2.0",
        "id": 1,
        "result": false
    }`
	cb["eth_blockNumber"] = `{
		"jsonrpc": "2.0",
		"id": 1,
		"result": "0x4b7"
	}`

	el := &EthereumEL{URL: srv.Http()}

	sync, err := el.Query()
	require.NoError(t, err)
	require.Equal(t, sync, &babelSDK.SyncStatus{IsSynced: true, CurrentBlock: 1207, NumPeers: 2})

	cb["eth_syncing"] = `{
        "jsonrpc": "2.0",
        "id": 1,
        "result": {
            "currentBlock": "0x540",
            "highestBlock": "0xbad427",
            "knownStates": "0x49d5",
            "pulledStates": "0x6a7",
            "startingBlock": "0x0"
        }
    }`

	sync, err = el.Query()
	require.NoError(t, err)
	require.Equal(t, sync, &babelSDK.SyncStatus{CurrentBlock: 1344, HighestBlock: 12244007, NumPeers: 2})
}
