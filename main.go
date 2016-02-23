package main

import (
	"encoding/csv"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
	"gopkg.in/olivere/elastic.v3"
	"io"
	"net/http"
	"os"
	"reflect"
)

type Quote struct {
	Text   string
	Author string
	Genre  string
}

type SlackBot struct {
	Quotes      *[]Quote
	esClient    *elastic.Client
	IndexName   string
	slackApiKey string
}

func NewSlackBot() *SlackBot {
	return &SlackBot{
		IndexName:   "quotes",
		slackApiKey: "xoxb-22757066566-LtgYbLxQcDpedIKOe11PPFFH",
	}
}

func (sb *SlackBot) connectEs() (err error) {
	// Create a client
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.7.188:9200"))
	if err != nil {
		log.Warn("Connection error: ", err)
		return err
	}
	log.Info("Established connection to elasticsearch")
	sb.esClient = client

	return nil
}

func (sb *SlackBot) ensureIndexExists() (err error) {
	// check if index exists
	exists, err := sb.esClient.IndexExists(sb.IndexName).Do()
	if err != nil {
		return err
	}

	// return if index exists
	if exists {
		log.Infof("Index '%s' already exists in elasticsearch", sb.IndexName)
		return nil
	}

	// build up index otherwise

	// create index
	_, err = sb.esClient.CreateIndex(sb.IndexName).Do()
	if err != nil {
		return err
	}
	log.Infof("Created index '%s' in elasticsearch", sb.IndexName)

	quotes, err := getQuotes()
	if err != nil {
		return err
	}

	// how many requests per bulk
	bulkEach := 100
	quotesCount := len(quotes)
	bulkCount := quotesCount / bulkEach

	// split into bulks of 100 requests
	for posBulk := 0; posBulk <= bulkCount; posBulk++ {
		minQuote := posBulk * bulkEach
		maxQuote := (1 + posBulk) * bulkEach
		if maxQuote >= quotesCount {
			maxQuote = quotesCount
		}

		log.Infof(
			"Indexing bulk %d/%d quotes %d - %d",
			posBulk,
			bulkCount,
			minQuote,
			maxQuote-1,
		)

		bulkRequest := sb.esClient.Bulk()

		for posQuote := minQuote; posQuote < maxQuote; posQuote++ {
			bulkRequest = bulkRequest.Add(
				elastic.NewBulkIndexRequest().
					Index(sb.IndexName).
					Type("quote").
					Doc(quotes[posQuote]))
		}

		_, err := bulkRequest.Do()
		if err != nil {
			log.Warnf(
				"Indexing of bulk %d quotes %d - %d failed: %s",
				posBulk,
				bulkCount,
				minQuote,
				maxQuote,
				err,
			)
		}
	}

	return nil

}

func (sb *SlackBot) Query(text string) (Quote, error) {
	matchQuery := elastic.NewMatchQuery("Text", text)
	searchResult, _ := sb.esClient.Search().
		Index(sb.IndexName).
		Query(matchQuery).
		From(0).Size(1).
		Pretty(true).
		Do()
	log.Infof("query elasticsearch for '%s' hits=%d", text, searchResult.TotalHits())

	if searchResult.TotalHits() == 0 {
		return Quote{}, fmt.Errorf("Not quote found")
	}

	var quote Quote

	for _, item := range searchResult.Each(reflect.TypeOf(quote)) {
		if q, ok := item.(Quote); ok {
			quote = q
			break
		}
	}

	return quote, nil

}

func (sb *SlackBot) httpQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go WebServer")
	w.Header().Set("Content-Type", "text/plain")

	text := r.URL.Query().Get("text")

	quote, _ := sb.Query(text)

	output := fmt.Sprintf("request=%s response=%+v", text, quote)

	w.Write([]byte(output))

}

func (sb *SlackBot) connectSlack() {
	api := slack.New(sb.slackApiKey)

	rtm := api.NewRTM()
	log.Info("Connecting to slack")
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				log.Infof("Incoming message '%s'", ev.Text)
				quote, err := sb.Query(ev.Text)
				message := "Sorry but I am speechless from your input"
				if err == nil {
					message = fmt.Sprintf("%s (by %s, %s)", quote.Text, quote.Author, quote.Genre)
				}
				rtm.SendMessage(rtm.NewOutgoingMessage(message, ev.Channel))

			case *slack.InvalidAuthEvent:
				log.Errorf("Invalid credentials")
				break Loop

			default:
				// ignore everything else
			}
		}
	}
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

	}

	return quotes, nil

}

func main() {
	log.Info("Initializing slack quote bot...")

	sb := NewSlackBot()

	err := sb.connectEs()
	if err != nil {
		log.Fatalf("Error while initializing connection to elasticsearch: %s", err)
	}

	err = sb.ensureIndexExists()
	if err != nil {
		log.Fatalf("Error while ensuring index exists: %s", err)
	}

	// connect to slack
	sb.connectSlack()
}
