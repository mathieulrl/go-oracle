package utils

import (
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
    "sync"
    "time"
)

const (
    dirPerm  = 0o755
    filePerm = 0o644
)

// JSONDB est une base très simple basée sur un fichier JSON.
type JSONDB struct {
    mu   sync.Mutex
    path string
    data DBState
}

func NewJSONDB(path string) (*JSONDB, error) {
    db := &JSONDB{path: path, data: emptyState()}
    if err := db.load(); err != nil {
        return nil, err
    }
    return db, nil
}

func (d *JSONDB) load() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.ensureDir(); err != nil {
		return err
	}
	
	state, err := d.readStateFromFile()
	if err != nil {
		return err
	}
	d.data = state
	return nil
}

func (d *JSONDB) ensureDir() error { return os.MkdirAll(filepath.Dir(d.path), dirPerm) }

func (d *JSONDB) readStateFromFile() (DBState, error) {
    if _, err := os.Stat(d.path); errors.Is(err, os.ErrNotExist) {
        return emptyState(), d.flushLocked()
	}
	
	b, err := os.ReadFile(d.path)
	if err != nil {
		return DBState{}, err
	}
	
    if len(b) == 0 {
        return emptyState(), d.flushLocked()
	}
	
	var state DBState
	if err := json.Unmarshal(b, &state); err != nil {
		return DBState{}, err
	}
	return state, nil
}

func (d *JSONDB) flushLocked() error {
    b, err := json.MarshalIndent(d.data, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(d.path, b, filePerm)
}

func emptyState() DBState { return DBState{Prices: []PriceRecord{}, Txs: []TxRecord{}} }

func (d *JSONDB) AddPrice(price float64, sourceCnt int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data.Prices = append(d.data.Prices, PriceRecord{
		Timestamp: time.Now().UTC(),
		Price:      price,
		SourceCnt:  sourceCnt,
	})
	return d.flushLocked()
}

func (d *JSONDB) AddTx(txHash string, price float64) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data.Txs = append(d.data.Txs, TxRecord{
		Timestamp: time.Now().UTC(),
		TxHash:    txHash,
		Price:     price,
	})
	return d.flushLocked()
}

// GetState retourne une copie du state actuel (lecture). 
func (d *JSONDB) GetState() DBState {
    d.mu.Lock()
    defer d.mu.Unlock()
    // copie superficielle
    out := DBState{Prices: make([]PriceRecord, len(d.data.Prices)), Txs: make([]TxRecord, len(d.data.Txs))}
    copy(out.Prices, d.data.Prices)
    copy(out.Txs, d.data.Txs)
    return out
}


