# Go Oracle

A blockchain oracle service that aggregates cryptocurrency prices from multiple providers (Binance, Coinbase) and publishes them to an Ethereum smart contract. The service uses concurrent workers for fast price fetching, persists data to a JSON database, and supports dry-run mode for testing without gas costs.

## Build

```bash
go build -o bin/oracle ./cmd/oracle
```

## Run

### Dry-run mode (no blockchain interaction)

```bash
./bin/oracle --dry-run
```

### Production mode

Set environment variables and run:

```bash
export RPC_URL="https://your-rpc-endpoint"
export PRIVATE_KEY="0x..."
export CONTRACT_ADDRESS="0x..."

./bin/oracle --db-path data/db.json --metrics-addr :8080
```

## Configuration

- `--dry-run`: Run without sending transactions to the blockchain
- `--db-path`: Path to JSON database file (default: `data/db.json`)
- `--metrics-addr`: Address to expose metrics endpoint (default: `:8080`)

View metrics at: `http://localhost:8080/debug/vars`

