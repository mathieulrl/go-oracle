package blockchain

import (
	"context"
	"fmt"
    "log/slog"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
    "time"
    "go-oracle/internal/observability"
)

var (
	parsedABI     abi.ABI
	parseABIOnce  sync.Once
)

const priceScale = 100

func getABI() abi.ABI {
	parseABIOnce.Do(func() {
		const abiJSON = `[{"inputs":[{"internalType":"uint256","name":"_price","type":"uint256"}],"name":"updatePrice","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
		var err error
		parsedABI, err = abi.JSON(strings.NewReader(abiJSON))
		if err != nil {
            slog.Error("Erreur parsing ABI", "err", err)
		}
	})
	return parsedABI
}

// NewOracleClient initialise la connexion blockchain
func NewOracleClient(rpcURL, privateKeyHex, contractAddress string, dryRun bool) (*OracleClient, error) {
	contractAddr := common.HexToAddress(contractAddress)

	if dryRun {
		return newDryRunClient(contractAddr), nil
	}
	return newRealClient(rpcURL, privateKeyHex, contractAddr)
}

func newDryRunClient(contractAddr common.Address) *OracleClient {
	return &OracleClient{
		Client:       nil,
		Auth:         nil,
		ContractAddr: contractAddr,
		DryRun:       true,
	}
}

func newRealClient(rpcURL, privateKeyHex string, contractAddr common.Address) (*OracleClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("erreur connexion RPC: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("erreur clé privée: %w", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("erreur chainID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("erreur auth: %w", err)
	}

	return &OracleClient{
		Client:       client,
		Auth:         auth,
		ContractAddr: contractAddr,
		DryRun:       false,
	}, nil
}

// UpdatePrice envoie le prix sur le smart contract
func (o *OracleClient) UpdatePrice(ctx context.Context, price float64) (string, error) {
	priceInt := scalePrice(price)

	if o.DryRun {
		return o.handleDryRun(priceInt)
	}
    return o.sendTransaction(ctx, priceInt)
}

func scalePrice(price float64) *big.Int {
    return new(big.Int).SetInt64(int64(price * priceScale))
}

func (o *OracleClient) handleDryRun(priceInt *big.Int) (string, error) {
	data, err := getABI().Pack("updatePrice", priceInt)
	if err != nil {
		return "", fmt.Errorf("erreur pack ABI: %w", err)
	}
    slog.Info("DRY-RUN updatePrice", "to", o.ContractAddr.Hex(), "data", "0x"+common.Bytes2Hex(data), "price_scaled", priceInt.String())
	return "", nil
}

func (o *OracleClient) sendTransaction(ctx context.Context, priceInt *big.Int) (string, error) {
    start := time.Now()
    contract := bind.NewBoundContract(o.ContractAddr, getABI(), o.Client, o.Client, o.Client)
    opts := *o.Auth
    opts.Context = ctx
    tx, err := contract.Transact(&opts, "updatePrice", priceInt)
    observability.ObserveBlockchainTx(time.Since(start), err)
    if err != nil {
        return "", fmt.Errorf("erreur transaction: %w", err)
    }

    slog.Info("tx sent", "hash", tx.Hash().Hex(), "price_scaled", priceInt.String())
	return tx.Hash().Hex(), nil
}
