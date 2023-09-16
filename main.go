package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type ResponseData struct {
	APIURL      string
	ElapsedTime time.Duration
	Data        string
	Error       string
}

func fetchFromAPI(apiURL string, ch chan<- ResponseData) {
	startTime := time.Now()
	client := http.Client{Timeout: time.Second}
	response, err := client.Get(apiURL)
	if err != nil {
		ch <- ResponseData{APIURL: apiURL, Error: fmt.Sprintf("Erro ao acessar %s: %v", apiURL, err)}
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ch <- ResponseData{APIURL: apiURL, Error: fmt.Sprintf("Erro ao fechar %s: %v", apiURL, err)}
		}
	}(response.Body)

	elapsedTime := time.Since(startTime)
	if elapsedTime.Seconds() < 1.0 {
		body, _ := io.ReadAll(response.Body)
		ch <- ResponseData{APIURL: apiURL, Data: string(body), ElapsedTime: elapsedTime}
	}
}

func main() {
	cep := "22451-041" // Substitua pelo CEP desejado

	api1URL := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	api2URL := "http://viacep.com.br/ws/" + cep + "/json/"

	ch := make(chan ResponseData, 2)

	go fetchFromAPI(api1URL, ch)
	go fetchFromAPI(api2URL, ch)

	select {
	case response := <-ch:
		if response.Error != "" {
			fmt.Println(response.Error)
		} else {
			fmt.Printf("Resposta da API %s (Tempo: %s):\n%s\n", response.APIURL, response.ElapsedTime, response.Data)
		}
	case <-time.After(time.Second):
		fmt.Println("Timeout - Nenhuma API retornou sucesso dentro do tempo limite.")
	}
}
