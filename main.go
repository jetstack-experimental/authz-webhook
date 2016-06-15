package main

// import "log"
import "net/http"
import "os"
import "github.com/gorilla/handlers"


func main() {
    listenAddress := ":8080"
    optPort := os.Getenv("LISTEN_PORT")
    if optPort != "" {
      listenAddress  = ":" + os.Getenv("LISTEN_PORT")
    }

    logged := handlers.CombinedLoggingHandler(os.Stderr, Handlers() )

    http.ListenAndServe(listenAddress, logged)
}
