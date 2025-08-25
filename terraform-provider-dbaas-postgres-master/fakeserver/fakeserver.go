package fakeserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Define a constant for the job type
const JOB_TYPE = "database"

// Structure to hold database properties
type DatabaseProperties struct {
	Name             string `json:"name"`
	Tenant           string `json:"tenant"`
	SizeGB           int    `json:"size_gb"`
	MaxConnections   int    `json:"max_connections"`
	CurrentDataUsage int    `json:"current_data_usage"`
}

// Fakeserver struct to encapsulate the server and its properties
type Fakeserver struct {
	server  *http.Server
	dbs     map[string]DatabaseProperties
	debug   bool
	running bool
}

// Middleware to log HTTP requests for debugging
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// Initialization function that sets up a new FakeServer instance
func NewFakeServer(port int, databases map[string]DatabaseProperties, start, debug bool) *Fakeserver {
	serverMux := http.NewServeMux()

	srv := &Fakeserver{
		debug:   debug,
		dbs:     databases,
		running: false,
	}

	// Handle requests directed at the job type endpoint
	serverMux.HandleFunc(fmt.Sprintf("/api/%s/", JOB_TYPE), srv.handleApiDatabase)

	// Configure the server with the designated port and request handler
	apiServer := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: logRequest(serverMux),
	}

	srv.server = apiServer

	// Start the server in the background if the 'start' parameter is true
	if start {
		srv.StartInBackground()
	}

	if srv.debug {
		log.Printf("fakeserver: Set up fakeserver on port %d\n", port)
	}

	return srv
}

// Method to start the fake server asynchronously
func (srv *Fakeserver) StartInBackground() {
	go func() {
		if err := srv.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("fakeserver: ListenAndServe() failed: %v", err)
		}
	}()
	srv.running = true
}

// Method to gracefully shut down the server
func (srv *Fakeserver) Shutdown() {
	if err := srv.server.Close(); err != nil {
		log.Fatalf("fakeserver: Server shutdown failed: %v", err)
	}
	srv.running = false
}

// Method to access the underlying HTTP server, useful for direct manipulation
func (srv *Fakeserver) GetServer() *http.Server {
	return srv.server
}

// Primary request handler for database operations
func (srv *Fakeserver) handleApiDatabase(w http.ResponseWriter, r *http.Request) {
	var db DatabaseProperties
	var name string
	var ok bool

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if srv.debug {
		log.Printf("fakeserver: Request received: %s %s\n", r.Method, r.URL.EscapedPath())
		log.Printf("fakeserver: BODY: %s\n", string(body))
	}

	path := r.URL.EscapedPath()
	parts := strings.Split(path, "/")

	if len(parts) == 4 {
		name = parts[3]
		if name == "" {
			http.Error(w, "Empty database name in URI", http.StatusBadRequest)
			return
		}
		db, ok = srv.dbs[name]

		if r.Method == "POST" && ok {
			http.Error(w, fmt.Sprintf("Database %s already exists", name), http.StatusConflict)
			return
		}

		if r.Method != "POST" && !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		db.Name = name

		if r.Method == "GET" {
			b, _ := json.Marshal(db)
			w.Write(b)
			return
		}

		if r.Method == "DELETE" {
			delete(srv.dbs, name)
			w.Write([]byte(fmt.Sprintf("Database %s deleted", name)))
			return
		}

		if (r.Method == "POST" || r.Method == "PUT") && string(body) != "" {
			json.Unmarshal(body, &db)
			srv.dbs[name] = db
			b, _ := json.Marshal(db)
			w.Write(b)
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
