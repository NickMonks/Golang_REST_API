package domain

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"golang.org/x/crypto/bcrypt"
)

// This struct will capture the client data request
type RegisterPayload struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Username        string `json:"username"`
}

// Function to validate the request from client
func (r *RegisterPayload) IsValid() (bool, map[string]string) {
	v := NewValidator()

	// Email verification
	v.MustBeNotEmpty("email", r.Email)
	v.MustBeValidEmail("email", r.Email)

	// Password validation
	v.MustBeLongerThan("password", r.Password, 6)
	v.MustBeNotEmpty("password", r.Password)

	//ConfirmPassword validation
	v.MustBeNotEmpty("confirmPassword", r.ConfirmPassword)
	v.MustMatch(ElementMatcher{
		field: "confirmPassword",
		value: r.ConfirmPassword,
	},
		ElementMatcher{
			field: "password",
			value: r.Password,
		})

	// Username validation
	v.MustBeLongerThan("username", r.Username, 3)
	v.MustBeNotEmpty("username", r.Username)

	return v.IsValid(), v.errors
}

func (d *Domain) Register(payload RegisterPayload) (*User, error) {

	// First, we check that user exists, so for that we called the functions of the interface
	userExist, _ := d.DB.UserRepo.GetByEmail(payload.Email)
	if userExist != nil {
		return nil, ErrUserWithEmailAlreadyExist
	}

	userExist, _ = d.DB.UserRepo.GetByUsername(payload.Username)
	if userExist != nil {
		return nil, ErrUserWithUsernameAlreadyExist
	}

	//if defined, we create our password string and the data
	password, err := d.setPassword(payload.Password)
	fmt.Println(*password)
	if err != nil {
		return nil, err
	}

	data := &User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: *password,
	}

	user, err := d.DB.UserRepo.Create(data)
	if err != nil {
		return nil, err
	}

	return user, nil

	// This function returns the user as normal if no errors found AND creates Repo
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginPayload) IsValid() (bool, map[string]string) {
	v := NewValidator()

	v.MustBeNotEmpty("email", l.Email)
	v.MustBeValidEmail("email", l.Email)

	v.MustBeNotEmpty("password", l.Password)

	return v.IsValid(), v.errors

}

func (d *Domain) Login(payload LoginPayload) (*User, error) {

	user, err := d.DB.UserRepo.GetByEmail(payload.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredential
	}

	err = user.checkPassword(payload.Password)

	if err != nil {
		fmt.Printf("no password")
		return nil, ErrInvalidCredential
	}

	return user, nil
}

func (d *Domain) setPassword(password string) (*string, error) {
	bytePassword := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	password = string(passwordHash)
	fmt.Println(password)
	return &password, nil
}

// We need a way to extracft the token from the requested object
// we create a function

func stripBearerPrefixFromToken(token string) (string, error) {
	bearer := "BEARER"
	// The value of the auth is "Bearer" and then the token. we want to make sure that is bigger than token.
	//
	if len(token) > len(bearer) && strings.ToUpper(token[0:len(bearer)]) == bearer {
		// give me anything after that plus one (the splace)
		return token[len(bearer)+1:], nil
	}

	return token, nil
}

// Literally, this extracts the authentication in the header as a Extractor
var authHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"}, // Extractor is what do we want to struct. JWT lives in the Auth header
	Filter:    stripBearerPrefixFromToken,               // a function to filter this token
}

// This can receive multiple extractors - In our case we just want one
var authExtractor = &request.MultiExtractor{
	authHeaderExtractor,
}

// Parse the token from the http request -
// Note: First we auth the user, if so the response will be the auth header. After that, to access to the API
// this token must be used to certify the user is who it is.
func ParseToken(r *http.Request) (*jwt.Token, error) {
	token, err := request.ParseFromRequest(r, authExtractor, func(t *jwt.Token) (interface{}, error) {
		b := []byte(os.Getenv("JWT_SECRET"))
		return b, nil
	})

	return token, err
}
