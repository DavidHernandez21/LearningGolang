package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"src/github.com/DavidHernandez21/dependencyInjection/cache"
	"src/github.com/DavidHernandez21/dependencyInjection/database"
	"src/github.com/DavidHernandez21/dependencyInjection/provider"

	"github.com/gorilla/mux"
)

type defaultRepositoryProvider struct {
	db           database.Executor
	cacheStorage cache.Executor
}

func (d *defaultRepositoryProvider) Database() database.Executor {
	return d.db
}

func (d *defaultRepositoryProvider) Cache() cache.Executor {
	return d.cacheStorage
}

type apiServer struct {
	provider provider.RepositoryProvider
}

func (a *apiServer) CreateProviderMiddleware(next http.Handler) http.Handler {
	providers := map[interface{}]interface{}{
		provider.ContextKey:         a.provider,
		database.ExecutorContextKey: a.provider.Database(),
		cache.ExecutorContextKey:    a.provider.Cache(),
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for key, value := range providers {
			req = req.WithContext(context.WithValue(req.Context(), key, value))
		}
		next.ServeHTTP(w, req)
	})

}

func InsertSample(rw http.ResponseWriter, req *http.Request) {
	repositories := req.Context().Value(provider.ContextKey).(provider.RepositoryProvider)

	value := fmt.Sprintf("Entry Created At: %v", time.Now())
	tableName := mux.Vars(req)["table"]
	err := repositories.Database().Insert(req.Context(), tableName, value)
	if err != nil {
		// Do something about the error (log, alert, etc)
		log.Printf("Failed to add item to db: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := repositories.Cache().Delete("all-records"); err != nil {
		// Do something about the error (log, alert, etc)
		// Maybe even revert the insert/transaction
		log.Printf("Failed to clear cache: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "text/html")
	response := fmt.Sprintf("Inserted a value into table %v", tableName)
	_, err = rw.Write([]byte(response))

	if err != nil {
		log.Printf("error writing back to the client: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CreateTable(rw http.ResponseWriter, req *http.Request) {
	repositories := req.Context().Value(provider.ContextKey).(provider.RepositoryProvider)

	tableName := mux.Vars(req)["table"]

	err := repositories.Database().CreateTable(req.Context(), tableName)
	if err != nil {
		// Do something about the error (log, alert, etc)
		log.Printf("Failed to create Table: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "text/html")
	response := fmt.Sprintf("Table '%v' created", tableName)
	_, err = rw.Write([]byte(response))

	if err != nil {
		log.Printf("error writing back to the client: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Get all records from a specific table
func GetAllRecords(rw http.ResponseWriter, req *http.Request) {

	repositories := req.Context().Value(provider.ContextKey).(provider.RepositoryProvider)

	cached, err := repositories.Cache().Get("all-records")
	if err != nil {
		// Do something about the error (log, alert, etc)
		fmt.Println("Failed to get info from cache")
	}

	if cached != nil {
		log.Println("Record fetched from the cache")
		_, err = rw.Write(cached)
		if err != nil {
			log.Printf("error writing back to the client: %v\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	records, err := repositories.Database().LookupAll(req.Context(), "records")
	if err != nil {
		if errors.Is(err, database.ErrNotFoundTable) || errors.Is(err, database.ErrNotFoundRecord) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var buffer bytes.Buffer

	if err = json.NewEncoder(&buffer).Encode(records); err != nil {
		fmt.Println("Failed to encode json", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = repositories.Cache().Set("all-records", buffer.Bytes())

	_, err = rw.Write(buffer.Bytes())
	if err != nil {
		log.Printf("error writing back to the client: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func NewDefaultRepositoryProvider(db database.Executor, cache cache.Executor) *defaultRepositoryProvider {

	return &defaultRepositoryProvider{
		db:           db,
		cacheStorage: cache,
	}
}

func NewApiServer(provider provider.RepositoryProvider) *apiServer {
	return &apiServer{
		provider: provider,
	}
}
