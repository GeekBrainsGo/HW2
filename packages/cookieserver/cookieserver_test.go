package cookieserver

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func TestApp(t *testing.T) {
	app := NewApp()
	server := httptest.NewServer(app)
	defer server.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Error(err)
	}
	client := &http.Client{Jar: jar}
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	test := client.Jar.Cookies(req.URL)[0].Value
	if test != email {
		t.Fatalf("want: %v, got: %v", email, test)
	}
	t.Logf(string(body))
	t.Logf("want: %v, got: %v", email, test)
}
