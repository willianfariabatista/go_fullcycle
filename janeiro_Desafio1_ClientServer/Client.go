package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverURL     = "http://localhost:8080/cotacao"
	clientTimeout = 300 * time.Millisecond // Timeout para a chamada ao server
)

// Estrutura que receberá o JSON retornado pelo servidor
type QuoteResponse struct {
	Bid string `json:"bid"`
}

func main() {
	// Cria contexto com timeout de 300ms para fazer a requisição ao servidor
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL, nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Erro ao fazer requisição para %s: %v", serverURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Fatalf("Resposta não OK (status %d). Detalhes: %s", resp.StatusCode, string(bodyBytes))
	}

	var quote QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		log.Fatalf("Erro ao decodificar resposta JSON: %v", err)
	}

	// Escreve no arquivo "cotacao.txt" o valor de bid
	f, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatalf("Erro ao criar arquivo cotacao.txt: %v", err)
	}
	defer f.Close()

	line := fmt.Sprintf("Dólar: %s\n", quote.Bid)
	_, err = f.WriteString(line)
	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %v", err)
	}

	log.Println("Cotação salva em cotacao.txt com sucesso!")
}
