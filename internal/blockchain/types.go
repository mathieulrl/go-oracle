package blockchain

import (
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
)

// OracleClient contient le client Ethereum et l'authentification
type OracleClient struct {
    Client       *ethclient.Client
    Auth         *bind.TransactOpts
    ContractAddr common.Address
    DryRun       bool
}


