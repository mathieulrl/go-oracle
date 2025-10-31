package config

import (
    "flag"
    "os"
)

// Load centralise flags et variables d'environnement.
// Flags: --dry-run, --db-path
// Env: RPC_URL, PRIVATE_KEY, CONTRACT_ADDRESS
func Load() Config {
    dryRun := flag.Bool("dry-run", false, "n'exécute pas la transaction, affiche seulement les données de la tx")
    dbPath := flag.String("db-path", "data/db.json", "chemin du fichier JSON pour la persistance locale")
    metricsAddr := flag.String("metrics-addr", ":8080", "adresse d'exposition des métriques expvar (ex: :8080)")
    // ne déclare pas ici d'autres flags pour garder la simplicité
    flag.Parse()

    return Config{
        DryRun:          *dryRun,
        DBPath:          *dbPath,
        RPCURL:          os.Getenv("RPC_URL"),
        PrivateKeyHex:   os.Getenv("PRIVATE_KEY"),
        ContractAddress: os.Getenv("CONTRACT_ADDRESS"),
        MetricsAddr:     *metricsAddr,
    }
}


