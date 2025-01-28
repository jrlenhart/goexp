package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	statusCode int
	err        error
}

func worker(id int, url string, requests int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < requests; i++ {
		log.Printf("[Worker %d] Sending request number %d\n", id, i+1)
		resp, err := http.Get(url)
		if err != nil {
			results <- Result{statusCode: 0, err: err}
			continue
		}
		defer resp.Body.Close()
		results <- Result{statusCode: resp.StatusCode, err: nil}
		log.Printf("[Worker %d] Received response from request number %d with status %d\n", id, i+1, resp.StatusCode)
	}
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado.")
	requests := flag.Int("requests", 100, "Número total de requests.")
	concurrency := flag.Int("concurrency", 10, "Número de chamadas simultâneas.")
	flag.Parse()

	if *url == "" {
		fmt.Println("A URL é obrigatória!")
		return
	}

	requestsPerWorker := *requests / *concurrency
	results := make(chan Result, *requests)
	var wg sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(i, *url, requestsPerWorker, results, &wg)
	}

	wg.Wait()
	close(results)

	elapsedTime := time.Since(startTime)

	var totalRequests, totalRequestsSuccess int
	statusCodes := make(map[int]int)

	for result := range results {
		totalRequests++
		if result.err == nil {
			statusCodes[result.statusCode]++
			if result.statusCode == http.StatusOK {
				totalRequestsSuccess++
			}
		}
	}

	fmt.Println("Relatório de Teste de Carga:")
	fmt.Printf("Tempo total gasto na execução: %v\n", elapsedTime)
	fmt.Printf("Quantidade total de requests realizados: %d\n", totalRequests)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", totalRequestsSuccess)
	for statusCode, count := range statusCodes {
		if statusCode != http.StatusOK {
			fmt.Printf("Quantidade de requests com status HTTP %d: %d\n", statusCode, count)
		}
	}
}
