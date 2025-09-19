package webapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"homeapp/pkg/crawler"
	"homeapp/pkg/index"
	"homeapp/pkg/store"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

var testMux *mux.Router

var storeIns = store.New()

func TestMain(m *testing.M) {
	setup()

	testMux = mux.NewRouter()
	endpoints(testMux, storeIns)
	m.Run()
}

func setup() {
	storeIns.Docs = []crawler.Document{
		{
			ID: 1,
		},
	}

	storeIns.Index = index.InvertedIndex{
		"go": []int{1},
	}
}

func Test_docsHandler(t *testing.T) {
	setup()

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
	setup()

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

func Test_createDocHandler(t *testing.T) {
	setup()

	// Создаём HTTP=запрос.
	url := "http://google.com"
	title := "test title"
	id := storeIns.Docs[len(storeIns.Docs)-1].ID + 1
	body := bytes.NewReader([]byte(fmt.Sprintf(`{"URL": "%s","Title": "%s"}`, url, title)))

	req := httptest.NewRequest(http.MethodPost, "/docs", body)
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

	gotDoc := crawler.Document{}

	err := json.Unmarshal(rr.Body.Bytes(), &gotDoc)

	if err != nil {
		log.Fatal(err)
	}

	if gotDoc.URL != url || gotDoc.Title != title || gotDoc.ID != id {
		t.Error("ответ неверен")
	}

	if len(storeIns.Docs) < 2 {
		t.Error("документ не добавлен")
	}
}

func Test_deleteDocHandler(t *testing.T) {
	setup()

	// Создаём HTTP=запрос.
	id := strconv.Itoa(storeIns.Docs[len(storeIns.Docs)-1].ID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/docs/%s", id), nil)
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

	if len(storeIns.Docs) != 0 {
		t.Error("документ не удален")
	}
}

func Test_updateDocHandler(t *testing.T) {
	setup()

	// Создаём HTTP=запрос.
	id := strconv.Itoa(storeIns.Docs[len(storeIns.Docs)-1].ID)
	url := "new url"
	title := "new title"
	body := bytes.NewReader([]byte(fmt.Sprintf(`{"URL": "%s","Title": "%s"}`, url, title)))

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/docs/%s", id), body)
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

	gotDoc := crawler.Document{}

	err := json.Unmarshal(rr.Body.Bytes(), &gotDoc)

	if err != nil {
		log.Fatal(err)
	}

	if gotDoc.URL != url || gotDoc.Title != title {
		t.Error("ответ неверен")
	}

	doc := &storeIns.Docs[0]

	if doc.URL != url || doc.Title != title {
		t.Error("документ не обновлен")
	}
}
