package ethereumcl

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/umbracle/babel"
	babelSDK "github.com/umbracle/babel/sdk"
)

func TestEthereumCL(t *testing.T) {
	srv := babel.NewMockHttpServer(8100)
	defer srv.Stop()

	srv.Mux().HandleFunc("/eth/v1/node/syncing", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"data": {
				"head_slot": "1",
				"sync_distance": "1",
				"is_syncing": true,
				"is_optimistic": true
			}
		}`))
	})

	cl := &EthereumCL{URL: srv.Http()}

	sync, err := cl.Query()
	require.NoError(t, err)
	require.Equal(t, sync, &babelSDK.SyncStatus{IsSynced: false, CurrentBlock: 1, HighestBlock: 2})
}
