package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Response struct {
	Url      string
	Duration time.Duration
	Dados    interface{}
}

// API para buscar os dados da URL e envia o resultado pelo canal
func buscarDados(url string, ch chan<- Response, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erro ao acessar %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()
	duration := time.Since(start)
	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result.Dados); err != nil {
		fmt.Printf("Erro ao decodificar JSON de %s: %v\n", url, err)
		return
	}
	result.Duration = duration
	result.Url = url
	// Envia o resultado pelo canal
	ch <- result
}

func main() {
	urls := []string{
		"https://viacep.com.br/ws/70070140/json/",
		"https://brasilapi.com.br/api/cep/v1/70070140",
	}

	results := make(chan Response, len(urls))
	var wg sync.WaitGroup

	// Inicia as goroutines para buscar os dados
	for _, url := range urls {
		wg.Add(1)
		go buscarDados(url, results, &wg)
	}

	// Fecha o canal quando todas as goroutines terminarem
	go func() {
		wg.Wait()
		close(results)
	}()

	// Processa os resultados recebidos
	tempo := 1 * time.Second
	resposta := Response{}
	for result := range results {
		//fmt.Printf("Recebido: %+v\n", result)
		if result.Duration < tempo {
			tempo = result.Duration
			resposta = result
		}
	}
	fmt.Printf("Resultado\n\n")
	fmt.Printf("Tempo: %v\t\n", tempo)
	fmt.Printf("API: %s\t\n", resposta.Url)
	fmt.Printf("Dados: %+v\t\n", resposta.Dados)
}
