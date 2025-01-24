package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Estrutura para receber o JSON da AwesomeAPI
type USDBRLResponse struct {
	USDBRL struct {
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

const (
	dbFile               = "exchange.db"
	externalAPI          = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	externalAPITimeout   = 200 * time.Millisecond // Timeout para chamada da API
	databaseWriteTimeout = 10 * time.Millisecond  // Timeout para escrita no DB
)

// Função para inicializar o banco de dados e retornar a conexão
func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Cria tabela se não existir
	createTableSQL := `CREATE TABLE IF NOT EXISTS exchange (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        rate TEXT NOT NULL,
        date TEXT NOT NULL
    );`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

	return db, nil
}

// Função que consome a API externa para obter a cotação do dólar
func getUSDBRLRate(ctx context.Context) (string, error) {
	// Faz a requisição HTTP dentro de um contexto com timeout de 200ms
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, externalAPI, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status da resposta inesperado: %d", resp.StatusCode)
	}

	var apiResponse USDBRLResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", err
	}

	// Retorna o campo "bid"
	return apiResponse.USDBRL.Bid, nil
}

// Função para inserir a cotação no banco de dados usando um contexto com timeout de 10ms
func saveRateToDB(ctx context.Context, db *sql.DB, rate string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Tenta inserir no banco
		insertSQL := `INSERT INTO exchange (rate, date) VALUES (?, ?)`
		_, err := db.ExecContext(ctx, insertSQL, rate, time.Now().Format(time.RFC3339))
		return err
	}
}

// Handler que processa a requisição em /cotacao
func quoteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// 1. Cria um contexto com timeout de 200ms para obter a cotação
	ctxAPI, cancelAPI := context.WithTimeout(context.Background(), externalAPITimeout)
	defer cancelAPI()

	rate, err := getUSDBRLRate(ctxAPI)
	if err != nil {
		log.Println("Erro ao obter cotação do dólar:", err)
		http.Error(w, "Erro ao obter cotação do dólar", http.StatusInternalServerError)
		return
	}

	// 2. Cria um contexto com timeout de 10ms para salvar no banco de dados
	ctxDB, cancelDB := context.WithTimeout(context.Background(), databaseWriteTimeout)
	defer cancelDB()

	if err := saveRateToDB(ctxDB, db, rate); err != nil {
		log.Println("Erro ao salvar cotação no banco de dados:", err)
		http.Error(w, "Erro ao salvar cotação no banco de dados", http.StatusInternalServerError)
		return
	}

	// 3. Retorna o valor de "bid" em JSON para o cliente
	w.Header().Set("Content-Type", "application/json")
	responseBody := map[string]string{"bid": rate}
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		log.Println("Erro ao codificar resposta em JSON:", err)
		http.Error(w, "Erro interno ao gerar resposta", http.StatusInternalServerError)
		return
	}
}

func Client() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		quoteHandler(w, r, db)
	})

	log.Println("Servidor iniciado na porta 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
