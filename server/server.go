package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CotacaoUsdbrl struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type ApiResultados struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", buscaCotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func buscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cotacao, error := buscaCotacao()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var apiResultados ApiResultados

	apiResultados.Bid = cotacao.Usdbrl.Bid

	json.NewEncoder(w).Encode(apiResultados)
}

func buscaCotacao() (*CotacaoUsdbrl, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)

	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	defer res.Body.Close()
	body, error := ioutil.ReadAll(res.Body)

	if error != nil {
		return nil, error
	}

	var c CotacaoUsdbrl

	error = json.Unmarshal(body, &c)

	if error != nil {
		return nil, error
	}

	error = saveCotacaoDatabase(c)

	if error != nil {
		return nil, error
	}

	return &c, nil
}

func saveCotacaoDatabase(cotacao CotacaoUsdbrl) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	select {
	case <-ctx.Done():
		return errors.New("excedeu tempo limite para salvar cotacao")
	default:
		db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db.AutoMigrate(ApiResultados{})
		err = db.Create(cotacao.Usdbrl).Error
		if err != nil {
			return err
		}
		return nil
	}

}
