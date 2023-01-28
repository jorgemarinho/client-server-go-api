package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type ApiResultados struct {
	Bid string `json:"bid"`
}

func main() {

	//300ms
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

	file, err := os.Create("cotacao.txt")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %v\n", err)
	}

	defer file.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var data ApiResultados
	err = json.Unmarshal(body, &data)

	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(fmt.Sprintf("Dólar: {%s}", data.Bid))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao escrever no arquivo: %v\n", err)
	}
}
