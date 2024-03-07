package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	x "github.com/bizmonx/bizmonx-gateway/xymon"
)

var h *x.XymonHost

func main() {
	host := os.Getenv("XYMON_HOST")
	port := os.Getenv("XYMON_PORT")
	if host == "" || port == "" {
		log.Fatal("XYMON_HOST and XYMON_PORT must be set")
		os.Exit(1)
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	h = &x.XymonHost{Host: host, Port: p}

	http.HandleFunc("/", handler)

	log.Println("Listening on :1976...")
	log.Fatal(http.ListenAndServe(":1976", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {
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

	var msg x.XymonMessage
	err = json.Unmarshal(body, &msg)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusInternalServerError)
		return
	}
	msg.Send(h)

	// Now forward this message to the Xymon server

	fmt.Fprintln(w, "Message forwarded successfully to Xymon")
}
