package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.RemoteAddr, r.Method, r.URL.Path)

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Erro ao ler o body: %v", err)
			} else {
				log.Printf("Payload recebido: %s", string(bodyBytes))
				// Recoloca o body para que o handler possa ler
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL não foi definida")
	}

	serverPORT := os.Getenv("SERVER_PORT")
	if serverPORT == "" {
		serverPORT = "8080"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	defer pool.Close()

	r := mux.NewRouter()
	r.Use(logRequestMiddleware)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API está saudável!")
	}).Methods("GET")

	r.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Lista todas as notas criadas!")
	}).Methods("GET")

	r.HandleFunc("/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Lista uma nota específica!")
	}).Methods("GET")

	r.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Criar uma nova nota!")
	}).Methods("POST")

	r.HandleFunc("/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Atualizar uma nota existente!")
	}).Methods("PATCH")

	r.HandleFunc("/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Deletar uma nota existente!")
	}).Methods("DELETE")

	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Criar um novo usuário!")
	}).Methods("POST")

	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Listar todos os usuários!")
	}).Methods("GET")

	r.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Listar um usuário específico!")
	}).Methods("GET")

	r.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Atualizar um usuário existente!")
	}).Methods("PATCH")

	r.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Deletar um usuário existente!")
	}).Methods("DELETE")

	fmt.Println("Servidor rodando na porta", serverPORT)
	log.Fatal(http.ListenAndServe(":"+serverPORT, r))
}
