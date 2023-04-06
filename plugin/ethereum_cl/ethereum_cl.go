package ethereumcl

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	babelSDK "github.com/umbracle/babel/sdk"
)

type EthereumCL struct {
	URL string `mapstructure:"url"`
}

type output struct {
	Data interface{} `json:"data,omitempty"`
}

func (e *EthereumCL) queryImpl(path string, out interface{}) error {
	resp, err := http.Get(e.URL + path)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respObj := &output{
		Data: out,
	}
	if err := json.Unmarshal(data, &respObj); err != nil {
		return err
	}
	return nil
}

func (e *EthereumCL) Query() (*babelSDK.SyncStatus, error) {
	var syncedResp syncResponse
	if err := e.queryImpl("/eth/v1/node/syncing", &syncedResp); err != nil {
		return nil, err
	}

	var numPeers numPeersResponse
	if err := e.queryImpl("/eth/v1/node/peer_count", &numPeers); err != nil {
		return nil, err
	}

	syncStatus := &babelSDK.SyncStatus{
		IsSynced:     !syncedResp.IsSyncing,
		CurrentBlock: uint64(syncedResp.HeadSlot),
		HighestBlock: uint64(syncedResp.HeadSlot + syncedResp.SyncDistance),
		NumPeers:     uint64(numPeers.Connected),
	}
	return syncStatus, nil
}

type numPeersResponse struct {
	Connected argStrUint64 `json:"connected"`
}

type syncResponse struct {
	HeadSlot     argStrUint64 `json:"head_slot"`
	SyncDistance argStrUint64 `json:"sync_distance"`
	IsSyncing    bool         `json:"is_syncing"`
	IsOptimistic bool         `json:"is_optimistic"`
}

// argStrUint64 is a uint represented as a string
type argStrUint64 uint64

func (u *argStrUint64) UnmarshalText(input []byte) error {
	num, err := strconv.ParseUint(string(input), 10, 64)
	if err != nil {
		return err
	}

	*u = argStrUint64(num)
	return nil
}
