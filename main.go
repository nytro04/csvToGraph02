package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Statement model.
type Statement struct {
	Credit int `json:"credit"`
	Debit  int `json:"debit"`
}

func parseStatementsCSV(r io.Reader) ([]Statement, error) {
	
	csvFile := csv.NewReader(r)

	records, err := csvFile.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("could not parse csv: missing header")
	}

	headers := records[0]

	if len(headers) != 4 {
		return nil, fmt.Errorf("could not parse csv: invalid headers")
	}

	if headers[0] != "DATE" {
		return nil, fmt.Errorf("could not parse csv: missing DATE")
	}

	if headers[1] != "DESCRIPTION" {
		return nil, fmt.Errorf("could not parse csv: missing DESCRIPTION")
	}

	if headers[2] != "TYPE" {
		return nil, fmt.Errorf("could not parse csv: missing TYPE")
	}

	if headers[3] != "AMOUNT" {
		return nil, fmt.Errorf("could not parse csv: missing AMOUNT")
	}

	records = records[1:] // Skip's the firt row which are the label.
	credit := 0
	debit := 0
	ss := make([]Statement, 0)

	for _, record := range records {
		amount, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, fmt.Errorf("could not parse amount as integer %v", err)
		}

		if record[2] == "C" {
			credit += amount
		} else if record[2] == "D" {
			debit += amount
		}

	
	}

	ss = append(ss, Statement{credit, debit})

	return ss, nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/form-data") {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	// 32 MB
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ss, err := parseStatementsCSV(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(ss)

}

func main() {
	var port int
	flag.IntVar(&port, "p", intEnv("PORT", 3000), "Port")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/upload", uploadHandler)
	mux.Handle("/", http.FileServer(http.Dir("static")))

	log.Printf("starting server on http://localhost:%d/\n", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func intEnv(key string, fallbackValue int) int {
	s, ok := os.LookupEnv(key)
	if !ok {
		return fallbackValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return fallbackValue
	}

	return i
}
