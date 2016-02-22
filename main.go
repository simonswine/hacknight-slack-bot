package main

import (
	"encoding/csv"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/blevesearch/bleve"
	"io"
	"net/http"
	"os"
)

type Quote struct {
	Text   string
	Author string
	Genre  string
}

func getQuotes() ([]Quote, error) {

	var quotes []Quote

	file, err := os.Open("quotes_all.csv")
	if err != nil {
		return []Quote{}, err
	}
	// automatically call Close() at the end of current method
	defer file.Close()
	//
	reader := csv.NewReader(file)
	reader.Comma = ';'
	lineCount := 0
	for {
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return []Quote{}, err
		}

		lineCount += 1

		// skip first 2 lines
		if lineCount < 3 {
			continue
		}

		// record is an array of string so is directly printable
		quotes = append(quotes, Quote{Text: record[0], Genre: record[2], Author: record[1]})

		if lineCount >= 1000 {
			break
		}
	}

	return quotes, nil

}

func indexQuotes(quotes []Quote) {
	log.Info("Indexing all quotes")
	mapping := bleve.NewIndexMapping()

	path := "example.bleve"

	os.RemoveAll(path)

	index, err := bleve.New(path, mapping)
	if err != nil {
		panic(err)
	}

	for i, quote := range quotes {
		log.Infof("Indexing %d", i)
		index.Index(fmt.Sprintf("%d", i), quote)
	}
}

func main() {
	log.Info("Parse quotes from csv")
	quotes, err := getQuotes()
	if err != nil {
		log.Fatalf("Error while parsing quotes: %s", err)
	}

	//indexQuotes(quotes)

	log.Info("First quotes %s", quotes[0])

	log.Info("Serving on http")
	http.HandleFunc("/", queryIndex)
	http.ListenAndServe(":8008", nil)
}

func queryIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "A Go WebServer")
	w.Header().Set("Content-Type", "text/plain")

	text := r.URL.Query().Get("text")

	index, _ := bleve.Open("example.bleve")
	query := bleve.NewQueryStringQuery(text)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)

	var output string

	output = fmt.Sprintf("request=%s response=%+v", text, searchResult)
	w.Write([]byte(output))
	index.Close()

}
