package main

// import "log"
import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	listenAddress := ":8080"
	optPort := os.Getenv("LISTEN_PORT")
	if optPort != "" {
		listenAddress = ":" + os.Getenv("LISTEN_PORT")
	}

	configFile := os.Getenv("RULES_CONFIG")

	if configFile == "" {
		configFile = "rules.hcl"
	}

	LoadConfigFromFile(configFile)

	logged := handlers.CombinedLoggingHandler(os.Stderr, Handlers())

	if err := http.ListenAndServe(listenAddress, logged); err != nil {
		log.Fatal(err)
	}
}
