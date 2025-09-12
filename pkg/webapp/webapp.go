package webapp

import (
	"encoding/json"
	"homeapp/pkg/store"
	"net/http"

	"github.com/gorilla/mux"
)

func StartHTTPServer(port string, storeIns *store.Store) {
	mux := mux.NewRouter()
	endpoints(mux, storeIns)
	mux.Use(headersMiddleware)
	http.ListenAndServe(":"+port, mux)
}

func endpoints(mux *mux.Router, storeIns *store.Store) {
	mux.HandleFunc("/docs", docsHandler(storeIns)).Methods(http.MethodGet)
	mux.HandleFunc("/index", indexHandler(storeIns)).Methods(http.MethodGet)
}

func docsHandler(storeIns *store.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(storeIns.Docs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func indexHandler(storeIns *store.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(storeIns.Index); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
