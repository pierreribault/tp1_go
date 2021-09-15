package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func notAllowed(w http.ResponseWriter, req *http.Request) {
	response(http.StatusMethodNotAllowed, w, req.Method+" is not allowed.")
}

func badRequest(w http.ResponseWriter, message string) {
	response(http.StatusBadRequest, w, message)
}

func response(status int, w http.ResponseWriter, message string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, message)
}

func addEntry(author, message string) {
	f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.WriteString(author + ":" + message + "\n")

	if err != nil {
		log.Fatal(err)
	}
}

func listEntries() ([]string, error) {
	raw, err := os.ReadFile("data.txt")

	if err != nil {
		return nil, err
	}

	return strings.Split(string(raw), "\n"), nil
}

func index(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		notAllowed(w, req)
	} else {
		response(200, w, time.Now().Format("15:04"))
	}
}

func add(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	if req.Method != http.MethodPost {
		notAllowed(w, req)
	} else {
		author := req.Form.Get("author")
		message := req.Form.Get("message")

		if len(author) > 0 && len(message) > 0 {
			addEntry(author, message)
			response(200, w, author+":"+message)
		} else {
			badRequest(w, "Missing parameters")
		}
	}
}

func entries(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		notAllowed(w, req)
	} else {
		entries, err := listEntries()
		var result string

		if err != nil {
			response(204, w, "Entry is empty for the moment!\n")
		}

		for _, rawEntry := range entries {
			if strings.Contains(rawEntry, ":") {
				entry := strings.Split(rawEntry, ":")
				result += "> " + entry[1] + "\n"
			}
		}

		response(200, w, result)
	}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/add", add)
	http.HandleFunc("/entries", entries)

	log.Println("Server listening on http://localhost:8082")
	server := &http.Server{Addr: ":8082", Handler: nil}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start.")
	}
}