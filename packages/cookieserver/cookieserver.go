// Package cookieserver provide cookie handling multiplexer.
package cookieserver

/*
	Basics Go.
	Rishat Ishbulatov, dated Sep 17, 2019.

	Write two routes: one will record information in a cookie
	(for example, a name), and the second will receive it
	and display it in response to a request.
*/

import (
	"net/http"
	"time"
)

/*
	Mixed Caps
	See https://golang.org/doc/effective_go.html#mixed-caps.
	This applies even when it breaks conventions in other languages.
	For example an unexported constant is maxLength not MaxLength or MAX_LENGTH.
*/

const (
	name  = "test"
	email = "test@test.com"
)

// App stands for cookie handling multiplexer.
type App struct {
	*http.ServeMux
}

// NewApp return cookie handling multiplexer.
func NewApp() *App {
	router := http.NewServeMux()
	app := &App{router}
	router.HandleFunc("/", app.Root)
	router.HandleFunc("/login", app.SetName)
	return app
}

// Root handle main page.
func (a *App) Root(w http.ResponseWriter, r *http.Request) {
	v, err := r.Cookie(name)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	w.Write([]byte("Hello " + v.Name + ": " + v.Value))
}

// SetName handle cookie creation.
func (a *App) SetName(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   email,
		Expires: time.Now().AddDate(0, 0, 1),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
