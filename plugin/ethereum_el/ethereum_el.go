package ethereumel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	babelSDK "github.com/umbracle/babel/sdk"
)

var defaultURL = "http://localhost:8545"

type EthereumEL struct {
	URL string `mapstructure:"url"`
}

func (e *EthereumEL) doRequest(method string, obj interface{}) error {
	jsonRPCSyncReq := &jsonRPCRequest{
		JsonRPC: "2.0",
		ID:      1,
		Method:  method,
	}
	postData, err := json.Marshal(jsonRPCSyncReq)
	if err != nil {
		return err
	}

	resp, err := http.Post(e.URL, "application/json", bytes.NewBuffer(postData))
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respObj := &jsonRPCResponse{
		Result: obj,
	}
	if err := json.Unmarshal(data, &respObj); err != nil {
		return err
	}

	if respObj.Error != nil {
		return fmt.Errorf("jsonrpc request failed: %v", respObj.Error.Message)
	}
	return nil
}

func (e *EthereumEL) Query() (*babelSDK.SyncStatus, error) {
	if e.URL == "" {
		e.URL = defaultURL
	}

	// query the number of peers
	var numPeers argHexUint64
	if err := e.doRequest("net_peerCount", &numPeers); err != nil {
		return nil, fmt.Errorf("failed to 'net_peerCount' request: %v", err)
	}

	syncRes := json.RawMessage{}
	if err := e.doRequest("eth_syncing", &syncRes); err != nil {
		return nil, fmt.Errorf("failed to 'eth_syncing' request: %v", err)
	}

	if string(syncRes) == "false" {
		// query to get current block number
		var latestBlock argHexUint64
		if err := e.doRequest("eth_blockNumber", &latestBlock); err != nil {
			return nil, fmt.Errorf("failed to 'eth_blockNumber' request: %v", err)
		}
		syncStatus := &babelSDK.SyncStatus{
			IsSynced:     true,
			CurrentBlock: uint64(latestBlock),
			NumPeers:     uint64(numPeers),
		}
		return syncStatus, nil
	}

	// decode the synced response
	var synced *syncedResponse
	if err := json.Unmarshal(syncRes, &synced); err != nil {
		return nil, err
	}

	syncStatus := &babelSDK.SyncStatus{
		CurrentBlock: uint64(synced.CurrentBlock),
		HighestBlock: uint64(synced.HighestBlock),
		NumPeers:     uint64(numPeers),
	}
	return syncStatus, nil
}

type syncedResponse struct {
	CurrentBlock  argHexUint64 `json:"currentBlock"`
	HighestBlock  argHexUint64 `json:"highestBlock"`
	KnownStates   argHexUint64 `json:"knownStates"`
	PulledStates  argHexUint64 `json:"pulledStates"`
	StartingBlock argHexUint64 `json:"startingBlock"`
}

type jsonRPCRequest struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      uint64          `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type jsonRPCResponse struct {
	ID     uint64        `json:"id"`
	Result interface{}   `json:"result"`
	Error  *jsonRPCError `json:"error,omitempty"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// argHexUint64 is a uint represented as an hexadecimal value
type argHexUint64 uint64

func (u *argHexUint64) UnmarshalText(input []byte) error {
	str := strings.TrimPrefix(string(input), "0x")
	num, err := strconv.ParseUint(str, 16, 64)
	if err != nil {
		return err
	}

	*u = argHexUint64(num)
	return nil
}
