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

func (e *EthereumCL) Query() (*babelSDK.SyncStatus, error) {
	resp, err := http.Get(e.URL + "/eth/v1/node/syncing")
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var output struct {
		Data *response `json:"data,omitempty"`
	}
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, err
	}

	syncStatus := &babelSDK.SyncStatus{
		IsSynced:     !output.Data.IsSyncing,
		CurrentBlock: uint64(output.Data.HeadSlot),
		HighestBlock: uint64(output.Data.HeadSlot + output.Data.SyncDistance),
	}
	return syncStatus, nil
}

type response struct {
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
