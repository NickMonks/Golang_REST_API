package handlers

import (
	"net/http"
	"todo/domain"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
)

type authResponse struct {
	User  *domain.User     `json:"user"`
	Token *domain.JWTToken `json:token`
}

// users handlers. Think of it as a controller

func (s *Server) registerUser() http.HandlerFunc {
	// Here we decode the JSON. Good news is that we dont need to marshall Go to json
	var payload domain.RegisterPayload

	// we return a http handler that has a context of the payload validated

	// We test in postman the received request decoded
	// in http://localhost:8081/api/v1/users, where we actually post
	// a raw JSON
	// We want to do a validation of the payload (i.e, anything can be pass and that's not good!)
	// We will create our own validator library, returning an error for each key (i.e, if Email was bad written - EmailError)

	return validatePayload(func(w http.ResponseWriter, r *http.Request) {
		// once validated, we start to register our model
		user, err := s.domain.Register(payload)
		if err != nil {
			badRequestResponse(w, err)
			return
		}
		// generate jwt token:
		token, err := user.GenToken()
		if err != nil {
			badRequestResponse(w, err)
			return
		}

		jsonResponse(w, &authResponse{
			// now if the user exists, will output both User and Token
			User:  user,
			Token: token,
		}, http.StatusCreated)

	}, &payload)

}

// Login User function:

func (s *Server) loginUser() http.HandlerFunc {
	// Here we decode the JSON. Good news is that we dont need to marshall Go to json
	var payload domain.LoginPayload

	// we return a http handler that has a context of the payload validated

	// We test in postman the received request decoded
	// in http://localhost:8081/api/v1/users, where we actually post
	// a raw JSON
	// We want to do a validation of the payload (i.e, anything can be pass and that's not good!)
	// We will create our own validator library, returning an error for each key (i.e, if Email was bad written - EmailError)

	return validatePayload(func(w http.ResponseWriter, r *http.Request) {
		// once validated, we start to register our model
		user, err := s.domain.Login(payload)
		if err != nil {
			badRequestResponse(w, err)
			return
		}
		// generate jwt token:
		token, err := user.GenToken()
		if err != nil {
			badRequestResponse(w, err)
			return
		}

		jsonResponse(w, &authResponse{
			// now if the user exists, will output both User and Token
			User:  user,
			Token: token,
		}, http.StatusOK)

	}, &payload)

}

func (s *Server) currentUserFromCTX(r *http.Request) *domain.User {
	currentUser := r.Context().Value("currentUser").(*domain.User) // we cast the value returned, since if we just return it directly it will complain since it is a interface{}
	return currentUser
}

// Middleware to add the token into the context. to do so we use http.handler and it's a receiver function of server.
func (s *Server) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := domain.ParseToken(r)

		if err != nil {
			unauthorizedResponse(w)
			return
		}

		// we check the claims of the JWT - VWhat we pass to the users

		// 									\/ Map the JWT
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			userID := int64(claims["id"].(float64)) // the claim is in float, but our userID is int64

			user, err := s.domain.GetUserByID(userID)

			if err != nil {
				unauthorizedResponse(w)
				return
			}

			// If everything correct, return a context to the handler wrapper (middleware)
			ctx := context.WithValue(r.Context(), "currentUser", user)

			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			unauthorizedResponse(w)
			return
		}
	})

}
