package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/microlib/simple"
	"lmzsoftware.com/lzuccarelli/golang-blockchain-interface/pkg/connectors"
	"lmzsoftware.com/lzuccarelli/golang-blockchain-interface/pkg/handlers"
	"lmzsoftware.com/lzuccarelli/golang-blockchain-interface/pkg/validator"
)

func startHttpServer(conn connectors.Clients) {
	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/blockchain/list", func(w http.ResponseWriter, req *http.Request) {
		handlers.GetBlockChainList(w, req, conn)
	}).Methods("OPTIONS", "GET")

	r.HandleFunc("/api/v1/blockchain/{index}", func(w http.ResponseWriter, req *http.Request) {
		handlers.GetBlockChain(w, req, conn)
	}).Methods("OPTIONS", "GET")

	r.HandleFunc("/api/v1/blockchain", func(w http.ResponseWriter, req *http.Request) {
		handlers.WriteBlockChain(w, req, conn)
	}).Methods("OPTIONS", "POST")

	r.HandleFunc("/api/v1/genesis", handlers.Init).Methods("POST")

	r.HandleFunc("/api/v2/sys/info/isalive", handlers.IsAlive).Methods("GET")
	http.Handle("/", r)

	if err := srv.ListenAndServe(); err != nil {
		conn.Error("Httpserver: ListenAndServe() error: %v", err)
	}
}

func main() {

	var logger *simple.Logger

	if os.Getenv("LOG_LEVEL") == "" {
		logger = &simple.Logger{Level: "info"}
	} else {
		logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}
	}

	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}
	conn := connectors.NewClientConnections(logger)
	startHttpServer(conn)
}
