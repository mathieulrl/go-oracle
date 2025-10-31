package providers

import (
    "sync"
    "testing"
    "time"
)

type benchProvider struct {
    price float64
    delay time.Duration
    err   error
}

func (b benchProvider) GetPrice() (float64, error) {
    if b.delay > 0 {
        time.Sleep(b.delay)
    }
    return b.price, b.err
}

func makeMulti(numProviders int, delay time.Duration, concurrency int) MultiProvider {
    providers := make([]PriceProvider, 0, numProviders)
    for i := 0; i < numProviders; i++ {
        providers = append(providers, benchProvider{price: 10000 + float64(i), delay: delay})
    }
    return MultiProvider{Providers: providers, Concurrency: concurrency}
}

// Benchmark: récupère la moyenne sur P providers avec un pool de taille K
func BenchmarkMultiProvider_GetPrice_Cold10_Pool10(b *testing.B) {
    m := makeMulti(10, 2*time.Millisecond, 10)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        if _, err := m.GetPrice(); err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkMultiProvider_GetPrice_Cold50_Pool10(b *testing.B) {
    m := makeMulti(50, 2*time.Millisecond, 10)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        if _, err := m.GetPrice(); err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkMultiProvider_GetPrice_Cold50_Pool50(b *testing.B) {
    m := makeMulti(50, 2*time.Millisecond, 50)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        if _, err := m.GetPrice(); err != nil {
            b.Fatal(err)
        }
    }
}

// Stress test: 1000 requêtes concurrentes sur le même MultiProvider
func TestStress_1000ConcurrentRequests(t *testing.T) {
    t.Parallel()
    m := makeMulti(10, 1*time.Millisecond, 10)

    var wg sync.WaitGroup
    const N = 1000
    wg.Add(N)
    errs := make(chan error, N)
    for i := 0; i < N; i++ {
        go func() {
            defer wg.Done()
            _, err := m.GetPrice()
            if err != nil {
                errs <- err
            }
        }()
    }
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        close(errs)
        for err := range errs {
            if err != nil {
                t.Fatalf("erreur concurrente: %v", err)
            }
        }
    case <-time.After(5 * time.Second):
        t.Fatalf("stress test timeout")
    }
}

// Benchmark parallélisé: simule plusieurs clients en parallèle appelant GetPrice
func BenchmarkMultiProvider_GetPrice_ParallelClients(b *testing.B) {
    m := makeMulti(20, 1*time.Millisecond, 10)
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            if _, err := m.GetPrice(); err != nil {
                b.Fatal(err)
            }
        }
    })
}


