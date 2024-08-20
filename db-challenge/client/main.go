package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}
	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", body))
	if err != nil {
		panic(err)
	}
	fmt.Println("Arquivo com a cotação do dólar criado com sucesso.")
}
