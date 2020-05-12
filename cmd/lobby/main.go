package main

import (
	"context"
	_ "expvar" //
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const version = "0.0.0-alpha"

var port = 8080

func init() {
	flag.IntVar(&port, "port", port, "The port to use for the server")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	os.Exit(execute(context.Background()))
}

func execute(ctx context.Context) int {
	if err := checkFlags(); err != nil {
		log.Println(err)
		return 1
	}
	if err := run(ctx); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func checkFlags() error {
	if port <= 0 {
		return fmt.Errorf("-port must be positive")
	}
	return nil
}

func run(ctx context.Context) error {
	minPort := port + 1
	maxPort := port + 20
	server := NewLobbyServer(version, minPort, maxPort)

	router := mux.NewRouter()
	server.InstallHandlers(router)
	router.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	))
	http.Handle("/", router)

	address := fmt.Sprintf(":%d", port)
	log.Printf("listening at %s", address)
	return http.ListenAndServe(address, nil)
}
