package utils

import "time"

// PriceRecord représente un prix mesuré à un instant donné.
type PriceRecord struct {
    Timestamp time.Time `json:"timestamp"`
    Price     float64   `json:"price"`
    SourceCnt int       `json:"sourceCnt"`
}

// TxRecord représente une transaction envoyée on-chain.
type TxRecord struct {
    Timestamp time.Time `json:"timestamp"`
    TxHash    string    `json:"txHash"`
    Price     float64   `json:"price"`
}

// DBState est la forme persistée sur disque.
type DBState struct {
    Prices []PriceRecord `json:"prices"`
    Txs    []TxRecord    `json:"txs"`
}


