package postserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testCases = []struct {
	description string
	method      string
	search      []byte
	expected    string
}{
	{
		"first case",
		http.MethodPost,
		[]byte(`{"search": "search", "sites": ["https://mail.ru","https://yandex.ru","https://brie3.github.io/page"]}`),
		`["https://mail.ru","https://yandex.ru"]`,
	},
	{
		"second case",
		http.MethodPost,
		[]byte(`{"search": "Поиск в интернете", "sites": ["https://yandex.ru","https://mail.ru","https://www.google.com","https://brie3.github.io/page"]}`),
		`["https://mail.ru","https://yandex.ru"]`,
	},
	{
		"third case",
		http.MethodGet,
		[]byte(`{"search": "href", "sites": ["https://yandex.ru","https://mail.ru"]}`),
		"only post method supported.\n",
	},
}

func TestSearch(t *testing.T) {
	msg := `
	Description: %s
	Search Query: %q
	Expected: %q
	Got: %q
	`
	for _, test := range testCases {
		request, _ := http.NewRequest(test.method, "/", bytes.NewReader(test.search))
		response := httptest.NewRecorder()

		Search(response, request)

		got := response.Body.String()
		want := test.expected

		if got != want {
			t.Errorf(msg, test.description, test.search, want, got)
		}
		t.Logf(msg, test.description, test.search, want, got)
	}
}
