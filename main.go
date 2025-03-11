package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// AnalyticsData represents the structure of the incoming analytics data
type AnalyticsData struct {
	Timestamp     int64  `json:"timestamp"`
	OperationName string `json:"operation_name"`
	OperationBody string `json:"operation_body"`
	ClientName    string `json:"client_name"`
	ClientVersion string `json:"client_version"`
}

// GraphQLField represents a field in a GraphQL query
type GraphQLField struct {
	Name  string
	Count int
}

// GraphQLQuery represents a parsed GraphQL query
type GraphQLQuery struct {
	Fields map[string]int
}

// ParseGraphQLQuery parses a GraphQL query string into a GraphQLQuery data structure
func ParseGraphQLQuery(query string) GraphQLQuery {
	query = strings.TrimSpace(query)
	fields := make(map[string]int)
	parseFields(query, fields)
	return GraphQLQuery{Fields: fields}
}

// parseFields is a helper function to parse fields from a GraphQL query string
func parseFields(query string, fields map[string]int) {
	// Improved parsing logic to count field usage
	fieldRegex := regexp.MustCompile(`(?m)^\s*(\w+)\s*\(`)
	matches := fieldRegex.FindAllStringSubmatch(query, -1)
	for _, match := range matches {
		if len(match) > 1 {
			fields[match[1]]++
		}
	}
}

// Example GraphQL query with variables
const exampleQuery = `query GetUser($id: ID!) {
  user(id: $id) {
    id
    name
    email
  }
}`

var (
	eventQueue = make(chan AnalyticsData, 100) // Buffered channel for events
	wg         sync.WaitGroup
)

// worker function to process events
func worker(id int) {
	defer wg.Done()
	for event := range eventQueue {
		log.Printf("Worker %d processing event at %d", id, event.Timestamp)
		parsedQuery := ParseGraphQLQuery(event.OperationBody)
		log.Printf("Parsed query: %+v", parsedQuery)
		// Simulate processing time
		// time.Sleep(time.Second)
	}
}

// handler function to process incoming analytics data
func handler(w http.ResponseWriter, r *http.Request) {
	var data AnalyticsData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Send event to the queue
	eventQueue <- data
	fmt.Fprintf(w, "Data received")
}

func main() {
	// Start worker pool
	numWorkers := 5
	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker(i)
	}

	http.HandleFunc("/analytics", handler)

	// Call ParseMain
	LexerMain()
	ParseMain()
	// Parse and log the example query
	parsedQuery := ParseGraphQLQuery(exampleQuery)
	log.Printf("Parsed example query: %+v", parsedQuery)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}

	// Close the event queue and wait for workers to finish
	close(eventQueue)
	wg.Wait()

}
