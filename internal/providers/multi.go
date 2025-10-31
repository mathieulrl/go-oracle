package providers

import (
    "fmt"
    "context"
    "sync"
)

// MultiProvider combine plusieurs sources et retourne la moyenne des prix
type MultiProvider struct {
    Providers []PriceProvider
    // Concurrency limite le nombre d'appels simultanés (0 = len(Providers))
    Concurrency int
}

type result struct {
    price float64
    err   error
}

func (m MultiProvider) GetPrice(ctx context.Context) (float64, error) {
    if len(m.Providers) == 0 {
        return 0, fmt.Errorf("aucun provider défini")
    }
    jobs := enqueueJobs(m.Providers)
    results := startWorkers(ctx, m.effectiveConcurrency(len(m.Providers)), jobs)
    avg, success, firstErr := aggregate(ctx, results)
    if success == 0 {
        if firstErr != nil {
            return 0, firstErr
        }
        return 0, fmt.Errorf("échec de récupération des prix")
    }
    return avg, nil
}

func (m MultiProvider) effectiveConcurrency(total int) int {
    if m.Concurrency <= 0 || m.Concurrency > total {
        return total
    }
    return m.Concurrency
}

func enqueueJobs(providers []PriceProvider) <-chan PriceProvider {
    jobs := make(chan PriceProvider)
    go func() {
        for _, p := range providers {
            jobs <- p
        }
        close(jobs)
    }()
    return jobs
}

func startWorkers(ctx context.Context, n int, jobs <-chan PriceProvider) <-chan result {
    results := make(chan result)
    var wg sync.WaitGroup
    worker := func() {
        defer wg.Done()
        for p := range jobs {
            price, err := p.GetPrice(ctx)
            results <- result{price: price, err: err}
        }
    }
    wg.Add(n)
    for i := 0; i < n; i++ {
        go worker()
    }
    go func() {
        wg.Wait()
        close(results)
    }()
    return results
}

func aggregate(ctx context.Context, results <-chan result) (avg float64, success int, firstErr error) {
    var total float64
    for {
        select {
        case <-ctx.Done():
            if success > 0 {
                avg = total / float64(success)
            }
            if firstErr == nil {
                firstErr = ctx.Err()
            }
            return
        case r, ok := <-results:
            if !ok {
                if success > 0 {
                    avg = total / float64(success)
                }
                return
            }
        if r.err != nil {
            if firstErr == nil {
                firstErr = r.err
            }
            continue
        }
        total += r.price
        success++
        }
    }
}
