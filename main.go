package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	x "github.com/bizmonx/bizmonx-gateway/xymon"
)

var h *x.XymonHost

func main() {
	host := os.Getenv("XYMON_HOST")
	if host == "" {
		log.Println("XYMON_HOST not set, using localhost")
		host = "localhost"
	}

	port := os.Getenv("XYMON_PORT")
	if port == "" {
		log.Println("XYMON_PORT not set, using 1984")
		port = "1984"
	}

	serverPort := os.Getenv("XYMON_GATEWAY_SERVER_PORT")
	if serverPort == "" {
		log.Println("XYMON_GATEWAY_SERVER_PORT not set, using 1976")
		serverPort = "1976"
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	h = &x.XymonHost{Host: host, Port: p}

	//http.HandleFunc("/", handler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/data", dataHandler)

	log.Println("Listening on :" + serverPort + "...")
	server := &http.Server{Addr: ":" + serverPort, Handler: nil}

	//log.Fatal(http.ListenAndServe(":1976", nil))

	// Set up channel to receive OS signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-stopChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)

	fmt.Println("Shutting down gracefully...")

}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body of the POST request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var msg x.StatusMessage
	err = json.Unmarshal(body, &msg)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusInternalServerError)
		return
	}
	msg.Send(h)

	// Now forward this message to the Xymon server

	fmt.Fprintln(w, "Message forwarded successfully to Xymon")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	// Read the body of the POST request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var msg x.DataMessage
	err = json.Unmarshal(body, &msg)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusInternalServerError)
		return
	}
	msg.Send(h)

	// Now forward this message to the Xymon server

	fmt.Fprintln(w, "Message forwarded successfully to Xymon")
}
