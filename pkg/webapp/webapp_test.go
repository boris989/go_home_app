package webapp

import (
	"encoding/json"
	"homeapp/pkg/crawler"
	"homeapp/pkg/index"
	"homeapp/pkg/store"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var testMux *mux.Router

var storeIns = store.New()

func TestMain(m *testing.M) {
	storeIns.Docs = []crawler.Document{
		{
			ID: 1,
		},
	}

	storeIns.Index = index.InvertedIndex{
		"go": []int{1},
	}

	testMux = mux.NewRouter()
	endpoints(testMux, storeIns)
	m.Run()
}

func Test_docsHandler(t *testing.T) {
	// Создаём HTTP=запрос.
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	req.Header.Add("Сontent-type", "application/json")

	// Объект для записи ответа HTTP-сервера.
	rr := httptest.NewRecorder()

	// Вызов маршрутизатора и обслуживание запроса.
	testMux.ServeHTTP(rr, req)

	// Анализ ответа сервера (неверный метод HTTP).
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	t.Log("Response: ", rr.Body)

	gotDocs := []crawler.Document{}

	err := json.Unmarshal(rr.Body.Bytes(), &gotDocs)

	if err != nil {
		log.Fatal(err)
	}

	if len(gotDocs) < 1 && gotDocs[0].ID != storeIns.Docs[0].ID {
		t.Error("ответ неверен")
	}
}

func Test_indexHandler(t *testing.T) {
	// Создаём HTTP=запрос.
	req := httptest.NewRequest(http.MethodGet, "/index", nil)
	req.Header.Add("Сontent-type", "application/json")

	// Объект для записи ответа HTTP-сервера.
	rr := httptest.NewRecorder()

	// Вызов маршрутизатора и обслуживание запроса.
	testMux.ServeHTTP(rr, req)

	// Анализ ответа сервера (неверный метод HTTP).
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	t.Log("Response: ", rr.Body)

	gotIndex := make(index.InvertedIndex)

	err := json.Unmarshal(rr.Body.Bytes(), &gotIndex)

	if err != nil {
		log.Fatal(err)
	}

	if gotIndex["go"][0] != storeIns.Index["go"][0] {
		t.Error("ответ неверен")
	}
}
