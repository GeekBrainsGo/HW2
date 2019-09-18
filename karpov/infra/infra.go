package infra

import (
	"fmt"
	"net/http"

	finder "github.com/art-frela/lightfinder"
	"github.com/go-chi/render"
)

// [CUSTOM MIDDLEWARE]

// FilterContentType - middleware to check content type as JSON
func FilterContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Filtering requests by MIME type
		if r.Method == "POST" { // filter for POST request
			if r.Header.Get("Content-type") != "application/json" {
				render.Render(w, r, ErrUnsupportedFormat)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// [HANDLER FUNCS]

// SearchText - handler func for search query text at the Sites
func SearchText(w http.ResponseWriter, r *http.Request) {
	params := &searchRequest{}
	if err := render.Bind(r, params); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	sq := finder.NewSingleQuery(params.Search, params.Sites)
	containSites, err := sq.QuerySearch()
	if err != nil {
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	if len(containSites) == 0 {
		err = fmt.Errorf("No one resource doesn't contains search text (%s)", params.Search)
		render.Render(w, r, ErrNotFound(err))
		return
	}
	response := &searchResponse{
		Result: containSites,
	}
	render.Render(w, r, response)
}

type searchRequest struct {
	Search string   `json:"search"`
	Sites  []string `json:"sites"`
}

func (sq *searchRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if sq.Search == "" || len(sq.Sites) == 0 {
		return fmt.Errorf("missing query text <%s> or list of resource %v", sq.Search, sq.Sites)
	}
	return nil
}

type searchResponse struct {
	Result []string `json:"result"`
}

func (sr *searchResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render - implement method Render for render.Renderer
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest - wrapper for make err structure
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrServerInternal - wrapper for make err structure
func ErrServerInternal(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal server error.",
		ErrorText:      err.Error(),
	}
}

// ErrNotFound - wrapper for make err structure for empty result
func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     http.StatusText(http.StatusNotFound),
		ErrorText:      err.Error(),
	}
}

// ErrUnsupportedFormat - 415 error implementation
var ErrUnsupportedFormat = &ErrResponse{HTTPStatusCode: http.StatusUnsupportedMediaType, StatusText: "415 - Unsupported Media Type. Please send JSON"}
