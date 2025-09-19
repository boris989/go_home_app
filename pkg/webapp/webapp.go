package webapp

import (
	"encoding/json"
	"homeapp/pkg/crawler"
	"homeapp/pkg/store"
	"log"
	"net/http"
	"strconv"

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
	mux.HandleFunc("/docs", createDocHandler(storeIns)).Methods(http.MethodPost)
	mux.HandleFunc("/docs/{id}", deleteDocHandler(storeIns)).Methods(http.MethodDelete)
	mux.HandleFunc("/docs/{id}", updateDocHandler(storeIns)).Methods(http.MethodPatch)
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

func createDocHandler(storeIns *store.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var d crawler.Document

		// Декодируем JSON
		err := json.NewDecoder(r.Body).Decode(&d)

		if err != nil {
			http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
			return
		}

		if len(storeIns.Docs) > 0 {
			d.ID = storeIns.Docs[len(storeIns.Docs)-1].ID + 1
		} else {
			d.ID = 1
		}

		storeIns.Docs = append(storeIns.Docs, d)

		if err := json.NewEncoder(w).Encode(d); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func deleteDocHandler(storeIns *store.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			log.Fatal(err)
		}

		for index, doc := range storeIns.Docs {
			if doc.ID == id {
				storeIns.Docs = append(storeIns.Docs[:index], storeIns.Docs[index+1:]...)
				break
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func updateDocHandler(storeIns *store.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			log.Fatal(err)
		}

		var d *crawler.Document

		for index, doc := range storeIns.Docs {
			if doc.ID == id {
				d = &storeIns.Docs[index]
				break
			}
		}

		if d == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Декодируем JSON
		err = json.NewDecoder(r.Body).Decode(d)

		if err != nil {
			http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(d); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
