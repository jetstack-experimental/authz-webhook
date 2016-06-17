package main

// import "log"
import (
 "net/http"
 "os" 
 "github.com/gorilla/handlers"
 "log"
)

func main() {
  listenAddress := ":8080"
  optPort := os.Getenv("LISTEN_PORT")
  if optPort != "" {
    listenAddress  = ":" + os.Getenv("LISTEN_PORT")
  }

  logged := handlers.CombinedLoggingHandler(os.Stderr, Handlers() )

  if err := http.ListenAndServe(listenAddress, logged); err != nil {
    log.Fatal(err)
  }
}
