package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	eth             *ethclient.Client
	currentBlock    *types.Block
	lastBlockUpdate time.Time
	sugGasPrice     *big.Int
	pendingTx       uint
	networkId       *big.Int
	gethInfo *GethInfo
)

func init() {
	gethInfo = new(GethInfo)
	gethInfo.TotalEth = big.NewInt(0)
}

type GethInfo struct {
	ContractsCreated int64
	TokenTransfers   int64
	ContractCalls    int64
	EthTransfers     int64
	BlockSize float64
	LoadTime int64
	TotalEth         *big.Int
}

func main() {
	var err error
	defer eth.Close()
	gethServer := os.Getenv("GETH")
	eth, err = ethclient.Dial(gethServer)
	if err != nil {
		panic(err)
	}

	go Routine()

	fmt.Printf("Geth Exporter running on http://0.0.0.0:9090/metrics\n")

	http.HandleFunc("/metrics", MetricsHttp)
	err = http.ListenAndServe("0.0.0.0:9090", nil)
	if err != nil {
		panic(err)
	}
}

func CalculateTotals(block *types.Block) {
	gethInfo.TotalEth = big.NewInt(0)
	gethInfo.ContractsCreated = 0
	gethInfo.TokenTransfers = 0
	gethInfo.EthTransfers = 0
	for _, b := range block.Transactions() {

		if b.To() == nil {
			gethInfo.ContractsCreated++
		}

		if len(b.Data()) >= 4 {
			method := hexutil.Encode(b.Data()[:4])
			if method == "0xa9059cbb" {
				gethInfo.TokenTransfers++
			}
		}

		if b.Value().Sign() == 1 {
			gethInfo.EthTransfers++
		}

		gethInfo.TotalEth.Add(gethInfo.TotalEth, b.Value())
	}

	size := strings.Split(currentBlock.Size().String(), " ")
	gethInfo.BlockSize = stringToFloat(size[0]) * 1000
}

func Routine() {
	for {
		t1 := time.Now()
		sugGasPrice, _ = eth.SuggestGasPrice(context.TODO())
		pendingTx, _ = eth.PendingTransactionCount(context.TODO())
		newBlock, _ := eth.BlockByNumber(context.TODO(), nil)
		networkId, _ = eth.NetworkID(context.TODO())

		if currentBlock == nil {
			lastBlockUpdate = time.Now()
			currentBlock = newBlock
			fmt.Printf("Received a new block #%v\n", newBlock.NumberU64())
			return
		}
		if newBlock.NumberU64() != currentBlock.NumberU64() {
			fmt.Printf("Received a new block #%v\n", newBlock.NumberU64())
			currentBlock = newBlock
			lastBlockUpdate = time.Now()
			diff := lastBlockUpdate.Sub(t1)
			gethInfo.LoadTime = diff.Nanoseconds()
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func stringToFloat(s string) float64 {
	amount, _ := strconv.ParseFloat(s, 10)
	return amount
}

//
// HTTP response handler for /metrics
func MetricsHttp(w http.ResponseWriter, r *http.Request) {
	var allOut []string

	now := time.Now()

	CalculateTotals(currentBlock)

	allOut = append(allOut, fmt.Sprintf("geth_current_block %v", currentBlock.NumberU64()))
	allOut = append(allOut, fmt.Sprintf("geth_seconds_last_block %0.2f", now.Sub(lastBlockUpdate).Seconds()))
	allOut = append(allOut, fmt.Sprintf("geth_block_transactions %v", len(currentBlock.Transactions())))
	allOut = append(allOut, fmt.Sprintf("geth_block_value %v", ToEther(gethInfo.TotalEth)))
	allOut = append(allOut, fmt.Sprintf("geth_block_gas_used %v", currentBlock.GasUsed()))
	allOut = append(allOut, fmt.Sprintf("geth_block_gas_limit %v", currentBlock.GasLimit()))
	allOut = append(allOut, fmt.Sprintf("geth_block_nonce %v", currentBlock.Nonce()))
	allOut = append(allOut, fmt.Sprintf("geth_block_difficulty %v", currentBlock.Difficulty()))
	allOut = append(allOut, fmt.Sprintf("geth_block_uncles %v", len(currentBlock.Uncles())))
	allOut = append(allOut, fmt.Sprintf("geth_block_size_bytes %v", gethInfo.BlockSize))
	allOut = append(allOut, fmt.Sprintf("geth_gas_price %v", sugGasPrice))
	allOut = append(allOut, fmt.Sprintf("geth_pending_transactions %v", pendingTx))
	allOut = append(allOut, fmt.Sprintf("geth_network_id %v", networkId))
	allOut = append(allOut, fmt.Sprintf("geth_contracts_created %v", gethInfo.ContractsCreated))
	allOut = append(allOut, fmt.Sprintf("geth_token_transfers %v", gethInfo.TokenTransfers))
	allOut = append(allOut, fmt.Sprintf("geth_eth_transfers %v", gethInfo.EthTransfers))
	allOut = append(allOut, fmt.Sprintf("geth_load_time %v", gethInfo.LoadTime))

	fmt.Fprintln(w, strings.Join(allOut, "\n"))
}

//
// CONVERTS WEI TO ETH
func ToEther(o *big.Int) *big.Float {
	pul, int := big.NewFloat(0), big.NewFloat(0)
	int.SetInt(o)
	pul.Mul(big.NewFloat(0.000000000000000001), int)
	return pul
}
