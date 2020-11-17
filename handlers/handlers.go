package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"todo/domain"
)

type Server struct {
	// inside our server struct we call a domain type (which contains: Repo)
	// domain will be a pointer to Domain type - Why? we don't want to create another variable
	// type domain.Domain, we want to pass it by reference.
	domain *domain.Domain
}

// we add a constructor type function
func setupMiddleware(r *chi.Mux) {

	// our middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Compress(6, "application/json"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))

}

func NewServer(domain *domain.Domain) *Server {
	//As explained above, we dont want to pass by value (i.e, copy the variable)
	// we want to pass by address as much as possible to reduce space complexity - the same for the other functions
	return &Server{domain}
}

func SetupRouter(domain *domain.Domain) *chi.Mux {

	// for info about chi: https://github.com/go-chi/chi
	// Mux is a simple HTTP route multiplexer that parses a request path,
	// records any URL params, and executes an end handler. It implements
	// the http.Handler interface and is friendly with the standard library.
	// again, import with go get.

	server := NewServer(domain)

	r := chi.NewRouter()

	server.setupEndpoints(r)

	return r
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {

	// JSON response to error, we send it as a json in this way:
	//
	w.Header().Set("Content-Type", "application/json")
	// After setting the header, we write a BadRequest status exit
	w.WriteHeader(statusCode)

	if data == nil {
		data = map[string]string{}
	}

	//In the Request handler:
	//									\/ Encode using this data type
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func badRequestResponse(w http.ResponseWriter, err error) {
	// ----- After receiving and decoding the body, we send a response to w ----------
	// The next step is to encode this response to a json file and catch the error if happens.
	response := map[string]string{"error": err.Error()}

	jsonResponse(w, response, http.StatusBadRequest)
}

func forbiddenResponse(w http.ResponseWriter) {
	// ----- After receiving and decoding the body, we send a response to w ----------
	// The next step is to encode this response to a json file and catch the error if happens.
	response := map[string]string{"error": "Forbidden"}

	jsonResponse(w, response, http.StatusForbidden)
}

func unauthorizedResponse(w http.ResponseWriter) {
	// ----- After receiving and decoding the body, we send a response to w ----------
	// The next step is to encode this response to a json file and catch the error if happens.
	response := map[string]string{"error": "Unauthorized"}

	jsonResponse(w, response, http.StatusUnauthorized)
}

// Validation of the payload in the middleware
// We define a interface PayloadValidation which follows the contract IsValid()
type PayloadValidation interface {
	IsValid() (bool, map[string]string)
}

// validatePyaload will decode the body to json and also validate each field
// It takes the original http.HandlerFunc and handles the requests using validatePayload content
func validatePayload(next http.HandlerFunc, payload PayloadValidation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// To understand better http.Handler wrapper: https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702

		// the Decode method will look at the JSON tag directly, in that way we avoid marshalling
		// Basically this can be read as: a request handler is received. So, we create a new enconder (which is the body sent)
		// and decode it using the payload struct, which as JSON tags
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			badRequestResponse(w, err)
			return
		}
		// only if we dont receive an error (i.e, the whole function must be executed), we want to close the body
		defer r.Body.Close()

		// In case the pauload is NOT valid send a response error
		if isValid, errs := payload.IsValid(); !isValid {
			jsonResponse(w, errs, http.StatusBadRequest)
			return
		}

		// context

		//A way to think about context package in go is that it allows you to pass values without being a global variable
		// is very common to see code where middleware are added to HTTP pipeline and the results are added to the http.Request,
		// like we do here. https://blog.golang.org/context

		ctx := context.WithValue(r.Context(), "payload", payload)

		// we serve the payload to the HTTP pipeline only if isValid is true
		// serveHTTP calls next(w,r)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
}

//												\/ returns a function with a middleware and the handler
// 												This is done so we can use the output easily (see "endpoints")
func (s *Server) withOwner(subjectType string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentUser := s.currentUserFromCTX(r)
			subject := r.Context().Value(subjectType).(domain.HaveOwner)

			if !subject.IsOwner(currentUser) {
				forbiddenResponse(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
