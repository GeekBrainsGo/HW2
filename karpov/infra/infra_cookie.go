package infra

import (
	"net/http"
)

const (
	COOKIEKEY       = "username"
	SETCOOKIEURL    = "/v1/cookies?username=<your text>"
	DEFAULTUSERNAME = "noname"
)

// SetCookies - cookie setter
func SetCookies(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue(COOKIEKEY)
	if name != "" {
		cookie := &http.Cookie{
			Name:  COOKIEKEY,
			Value: name,
		}
		http.SetCookie(w, cookie)
		return
	}
	cookie := &http.Cookie{
		Name:  COOKIEKEY,
		Value: DEFAULTUSERNAME,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("use correct uri and parameters: " + SETCOOKIEURL + " set default value"))
}

// GetCookies - cokkie getter
func GetCookies(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(COOKIEKEY)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error: " + err.Error()))
		return
	}
	value := cookie.Value
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("cookie value: " + value))
}
