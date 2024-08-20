package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Price struct {
	ID         int    `gorm:"primaryKey"`
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8081", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-r.Context().Done():
		log.Println("Requisição cancelada.")
	default:
		price, err := FindPrice("USD-BRL")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in find price: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = Save(price)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in save: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(price.Bid))
	}
}

func FindPrice(coin string) (*Price, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/"+coin, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}
	var result map[string]json.RawMessage
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	var p Price
	err = json.Unmarshal(result[strings.ReplaceAll(coin, "-", "")], &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func Save(price *Price) error {
	db, err := gorm.Open(sqlite.Open("db-challenge.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Price{})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	res := db.WithContext(ctx).Create(&price)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
