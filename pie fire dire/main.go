package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode"
)

var url string = "https://baconipsum.com/api/?type=meat-and-filler&paras=99&format=text"
var filePath string = "meat_data.txt"
var client = &http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		MaxIdleConns: 100,
		IdleConnTimeout: 90 * time.Second,
	},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/beef/summary", meatSummaryHandler)

	server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  15 * time.Second,
    }

	// start server in goroutine, so that it will not block the graceful shutdown
	go func() {
        log.Println("Server starting on port 8080...")
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("HTTP server ListenAndServe: %v", err)
        }
    }()

	// graceful shutdown
	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server){
	quit := make(chan os.Signal, 1)
	// use signal.Notify to register the given channel to reeive notis of the signal above
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// block until receive the signal
	<-quit

	// create deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// shutdown server w/o interrupting active connections
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")

}

func meatSummaryHandler(w http.ResponseWriter, r *http.Request) {
    // fetch data from API
    resp, err := client.Get(url)
    if err != nil {
        http.Error(w, "Failed to retrieve meat data", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // read response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Failed to read response body", http.StatusInternalServerError)
        return
    }

	// save data to file
	err = ioutil.WriteFile(filePath, body, 0644)
	if err != nil {
		http.Error(w, "Failed to save meat data to file", http.StatusInternalServerError)
		return 
	}

    // convert body to string
    meatData := string(body)

    // count each words
    counts := countMeats(meatData)

	// nest under key "beef"
	nestCounts := map[string]interface{}{
		"beef": counts,
	}

	// encode to json and send as response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(nestCounts)
}

func countMeats(text string) map[string]int {
	text =  strings.ToLower(text)
	// split words -> below count word that has '-' as one 
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-')
	})

	counts := make(map[string]int)
	for _, word := range words {
		word := strings.ToLower(word)
		trimmedWord := strings.TrimSpace(word)
		if trimmedWord != "" {
			counts[word]++
		}		
	}
	return counts
}
