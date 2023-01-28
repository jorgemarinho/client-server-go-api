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

type Usdbrl struct {
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

	apiResultados.Bid = cotacao.Bid

	json.NewEncoder(w).Encode(apiResultados)
}

func buscaCotacao() (*Usdbrl, error) {

	//200ms
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

	var c Usdbrl
	error = json.Unmarshal(body, &c)

	if error != nil {
		return nil, error
	}

	/*err = saveCotacaoDatabase(c)

	if error != nil {
		return nil, error
	}*/

	return &c, nil
}

func saveCotacaoDatabase(cotacao Usdbrl) error {

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
		db.AutoMigrate(&Usdbrl{})
		err = db.Create(cotacao).Error
		if err != nil {
			return err
		}
		return nil
	}

}

/*
Os requisitos para cumprir este desafio são:

O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar. ok

O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente. ok

Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.

O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go. ok

O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor} ok

O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080. ok

Ao finalizar, envie o link do repositório para correção.
*/
