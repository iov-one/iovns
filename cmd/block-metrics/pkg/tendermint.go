package pkg

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/app"
	"github.com/iov-one/iovns/x/starname"
	starnametypes "github.com/iov-one/iovns/x/starname/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type TendermintClient struct {
	idCnt uint64

	conn *websocket.Conn

	stop chan struct{}
	mu   sync.Mutex
	resp map[string]chan<- *jsonrpcResponse
}

// DialTendermint returns a client that is maintains a websocket connection to
// tendermint API. The websocket is used instead of standard HTTP connection to
// lower the latency, bypass throttling and to allow subscription requests.
func DialTendermint(websocketURL string) (*TendermintClient, error) {
	c, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "dial")
	}
	cli := &TendermintClient{
		conn: c,
		stop: make(chan struct{}),
		resp: make(map[string]chan<- *jsonrpcResponse),
	}
	go cli.readLoop()
	return cli, nil
}

func (c *TendermintClient) Close() error {
	close(c.stop)
	return c.conn.Close()
}

func (c *TendermintClient) readLoop() {
	for {
		select {
		case <-c.stop:
			return
		default:
		}

		var resp jsonrpcResponse
		if err := c.conn.ReadJSON(&resp); err != nil {
			log.Printf("cannot unmarshal JSONRPC message: %s", err)
			continue
		}

		c.mu.Lock()
		respc, ok := c.resp[resp.CorrelationID]
		delete(c.resp, resp.CorrelationID)
		c.mu.Unlock()

		if ok {
			// repc is expected to be a buffered channel so this
			// operation must never block.
			respc <- &resp
		}
	}
}

// Do makes a jsonrpc call. This method is safe for concurrent calls.
//
// Use API as described in https://tendermint.com/rpc/
func (c *TendermintClient) Do(method string, dest interface{}, args ...interface{}) error {
	params := make([]string, len(args))
	for i, v := range args {
		params[i] = fmt.Sprint(v)
	}
	req := jsonrpcRequest{
		ProtocolVersion: "2.0",
		CorrelationID:   fmt.Sprint(atomic.AddUint64(&c.idCnt, 1)),
		Method:          method,
		Params:          params,
	}

	respc := make(chan *jsonrpcResponse, 1)
	c.mu.Lock()
	c.resp[req.CorrelationID] = respc
	c.mu.Unlock()

	if err := c.conn.WriteJSON(req); err != nil {
		return errors.Wrap(err, "write JSON")
	}

	resp := <-respc

	if resp.Error != nil {
		return errors.Wrapf(ErrFailedResponse,
			"%d: %s",
			resp.Error.Code, resp.Error.Message)
	}
	if err := ModuleCdc.UnmarshalJSON(resp.Result, dest); err != nil {
		return errors.Wrap(err, "cannot unmarshal result")
	}
	return nil
}

type jsonrpcRequest struct {
	ProtocolVersion string   `json:"jsonrpc"`
	CorrelationID   string   `json:"id"`
	Method          string   `json:"method"`
	Params          []string `json:"params,omitempty"`
}

type jsonrpcResponse struct {
	ProtocolVersion string `json:"jsonrpc"`
	CorrelationID   string `json:"id"`
	Result          json.RawMessage
	Error           *struct {
		Code    int64
		Message string
	}
}

// AbciInfo returns abci_info.
func AbciInfo(c *TendermintClient) (*ABCIInfo, error) {
	var payload struct {
		Response struct {
			LastBlockHeight sint64 `json:"last_block_height"`
		} `json:"response"`
	}

	if err := c.Do("abci_info", &payload); err != nil {
		return nil, errors.Wrap(err, "query tendermint")
	}

	return &ABCIInfo{LastBlockHeight: int64(payload.Response.LastBlockHeight)}, nil
}

type ABCIInfo struct {
	LastBlockHeight int64 `json:"last_block_height"`
}

func Commit(ctx context.Context, c *TendermintClient, height int64) (*TendermintCommit, error) {
	var payload struct {
		SignedHeader struct {
			Header struct {
				Height          sint64    `json:"height"`
				Time            time.Time `json:"time"`
				ProposerAddress hexstring `json:"proposer_address"`
				ValidatorsHash  hexstring `json:"validators_hash"`
			} `json:"header"`
			Commit struct {
				BlockID struct {
					Hash hexstring `json:"hash"`
				} `json:"block_id"`
				Precommits []*struct {
					ValidatorAddress hexstring `json:"validator_address"`
				} `json:"precommits"`
			} `json:"commit"`
		} `json:"signed_header"`
	}

	if err := c.Do("commit", &payload, height); err != nil {
		return nil, errors.Wrap(err, "query tendermint")
	}

	commit := TendermintCommit{
		Height:          payload.SignedHeader.Header.Height.Int64(),
		Hash:            payload.SignedHeader.Commit.BlockID.Hash,
		Time:            payload.SignedHeader.Header.Time.UTC(),
		ProposerAddress: payload.SignedHeader.Header.ProposerAddress,
		ValidatorsHash:  payload.SignedHeader.Header.ValidatorsHash,
	}

	for _, pc := range payload.SignedHeader.Commit.Precommits {
		if pc == nil {
			continue
		}
		commit.ParticipantAddresses = append(commit.ParticipantAddresses, pc.ValidatorAddress)
	}

	return &commit, nil
}

type TendermintCommit struct {
	Height               int64
	Hash                 []byte
	Time                 time.Time
	ProposerAddress      []byte
	ValidatorsHash       []byte
	ParticipantAddresses [][]byte
}

func FetchBlock(ctx context.Context, c *TendermintClient, height int64) (*TendermintBlock, error) {
	var payload struct {
		Block struct {
			Header struct {
				Height sint64    `json:"height"`
				Time   time.Time `json:"time"`
			} `json:"header"`
			Data struct {
				Txs [][]byte `json:"txs"`
			} `json:"data"`
		} `json:"block"`
	}

	if err := c.Do("block", &payload, height); err != nil {
		return nil, errors.Wrap(err, "query tendermint")
	}

	block := TendermintBlock{
		Height: payload.Block.Header.Height.Int64(),
		Time:   payload.Block.Header.Time.UTC(),
	}

	for _, rawTx := range payload.Block.Data.Txs {
		var tx auth.StdTx
		if err := ModuleCdc.UnmarshalBinaryLengthPrefixed(rawTx, &tx); err != nil {
			return nil, errors.Wrap(err, "cannot unmarshal transaction")
		}
		block.Transactions = append(block.Transactions, &tx)
		block.TransactionHashes = append(block.TransactionHashes, sha256.Sum256(rawTx))
	}

	return &block, nil
}

type TendermintBlock struct {
	Height            int64
	Time              time.Time
	Transactions      []*auth.StdTx
	TransactionHashes [][32]byte
}

func FetchGenesis(ctx context.Context, c *TendermintClient) (*starname.GenesisState, error) {
	var result coretypes.ResultGenesis
	if err := c.Do("genesis", &result); err != nil {
		return nil, errors.Wrap(err, "query tendermint")
	}
	appState := result.Genesis.AppState
	var st map[string]json.RawMessage
	if err := json.Unmarshal(appState, &st); err != nil {
		return nil, errors.Wrapf(err, "genesis parsing error")
	}
	domainModuleGen, ok := st[starname.ModuleName]
	if !ok {
		return nil, errors.New("cannot get domain module genesis data")
	}
	var genState starname.GenesisState
	if err := ModuleCdc.UnmarshalJSON(domainModuleGen, &genState); err != nil {
		return nil, errors.Wrapf(err, "genesis parsing error")
	}
	return &genState, nil
}

// LcdRequestData is the data required to populate the product_fees table
type LcdRequestData struct {
	// accounts.id
	AccountID int64
	// URL encoded parameters to narrow the query on /txs as much as possible
	Params map[string]string
}

type StarnameResponse struct {
	Height string `json:"height"`
	Result struct {
		Account starnametypes.Account `json:"account"`
	} `json:"result"`
}

type LcdResponseData struct {
	TxResponse       *types.TxResponse
	TxError          *error
	StarnameResponse *StarnameResponse
	StarnameError    *error
}

func fetchTx(ctx context.Context, urlLCD string, messageParams string, height int64, page int64, cdc *codec.Codec) (*types.TxResponse, error) {
	url := fmt.Sprintf("%s/txs?message.module=starname&%s&tx.minheight=%d&tx.maxheight=%d&page=%d", urlLCD, messageParams, height, height, page)
	response, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s", url)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}
	var result types.SearchTxsResult
	err = cdc.UnmarshalJSON(body, &result)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal %s", body)
	}
	if len(result.Txs) == 0 {
		if result.PageNumber < result.PageTotal {
			return fetchTx(ctx, urlLCD, messageParams, height, page+1, cdc)
		} else {
			return nil, errors.New(fmt.Sprintf("failed to find message with params %s and height %d", messageParams, height))
		}
	} else if len(result.Txs) > 1 {
		return nil, errors.New(fmt.Sprintf("expected 1 tx but got %d for query %s and height %d", len(result.Txs), messageParams, height))
	}

	return &result.Txs[0], nil
}

func fetchStarname(ctx context.Context, urlLCD string, starname string, cdc *codec.Codec) (*StarnameResponse, error) {
	payload := strings.NewReader(fmt.Sprintf(`{"starname":"%s"}`, starname))
	url := fmt.Sprintf("%s/starname/query/resolve", urlLCD)
	response, err := http.Post(url, "application/json", payload)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s", url)
	}
	var result StarnameResponse
	// use json decoder instead of cdc because amino fails on int64 valid_until not being a string
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, errors.Wrapf(err, "failed to marshal %s", starname)
	}
	if result.Height == "0" {
		return nil, errors.New(fmt.Sprintf("failed to resolve %s", starname))
	}

	return &result, err
}

func FetchLcdData(ctx context.Context, urlLCD string, queries *[]*LcdRequestData, height int64) (*[]*LcdResponseData, error) {
	responses := make([]*LcdResponseData, len(*queries))
	cdc := app.MakeCodec()
	var wg sync.WaitGroup
	for i, query := range *queries {
		wg.Add(1)
		go func(i int, q *LcdRequestData) {
			defer wg.Done()
			responses[i] = &LcdResponseData{}
			// fetch event related data
			params := url.Values{}
			for k, v := range q.Params {
				params.Add(fmt.Sprintf("message.%s", k), v)
			}
			tx, err := fetchTx(ctx, urlLCD, params.Encode(), height, 1, cdc)
			responses[i].TxResponse = tx
			responses[i].TxError = &err
			// possibly fetch valid_until
			switch q.Params["action"] {
			case "register_account", "register_domain", "renew_account", "renew_domain":
				account := ""
				domain := ""
				switch q.Params["action"] {
				case "register_account", "renew_account":
					account = q.Params["account_name"]
					domain = q.Params["domain_name"]
				case "register_domain", "renew_domain":
					domain = q.Params["domain_name"]
				}
				resolved, err := fetchStarname(ctx, urlLCD, fmt.Sprintf("%s*%s", account, domain), cdc)
				responses[i].StarnameResponse = resolved
				responses[i].StarnameError = &err
			}
		}(i, query)
	}
	wg.Wait()

	return &responses, nil
}
