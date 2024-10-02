package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	ApiBrasil = "https://brasilapi.com.br/api/cep/v1/%s"
	ApiViaCEP = "http://viacep.com.br/ws/%s/json/"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	channel := make(chan interface{})

	cep := "95800000"

	go func() {
		res, err := Get(ctx, fmt.Sprintf(ApiBrasil, cep))
		if err != nil {
			fmt.Printf("Error when looking for CEP: %v", err)
			return
		}
		channel <- map[string]interface{}{"api": ApiBrasil, "response": res}
	}()

	go func() {
		res, err := Get(ctx, fmt.Sprintf(ApiViaCEP, cep))
		if err != nil {
			fmt.Printf("Error when looking for CEP: %v", err)
			return
		}
		channel <- map[string]interface{}{"api": ApiViaCEP, "response": res}
	}()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		if err == context.DeadlineExceeded {
			fmt.Println("Timeout")
		} else if err == context.Canceled {
			fmt.Println("Canceled")
		} else {
			fmt.Println("Unknown context error:", err)
		}
	case msg := <-channel:
		fmt.Println(msg)
	}
}

func Get(ctx context.Context, url string) (interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", res.StatusCode, string(body))
	}
	return string(body), nil
}
