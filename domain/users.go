package domain

import (
	"fmt"
	"time"

	"os"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// The first step for auth is to define the User struct.
// next to each element we add the JSON body tag, using ``. Notice that the password is set as "-"

// NOTE: struct tags to control how this information is assigned to the fields of a struct. Struct tags are small pieces of metadata attached to fields of a struct that provide instructions to other Go code that works with the struct.
// Golang will ignore this unles use the encoding/json package.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type JWTToken struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// JWT: JSON Web Token - is used to securely transmit information between parties as a JSON object
// this info can be verified and trusted because is signed.

// How works?:
// 1.The application or client requests authorization to the authorization server. This is performed through one of the different authorization flows.
// 2. When the authorization is granted, the authorization server returns an access token to the application.
// 3. The application uses the access token to access a protected resource (like an API).

func (u *User) GenToken() (*JWTToken, error) {
	jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))

	expiresAt := time.Now().Add(time.Hour * 24 * 7) // valid for 1 week

	// the "body" of our token
	jwtToken.Claims = jwt.MapClaims{
		"id":  u.ID,
		"exp": expiresAt.Unix(),
	}

	//
	accessToken, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return nil, err
	}

	return &JWTToken{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func (u *User) checkPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)

	fmt.Println(password)
	fmt.Println(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func (d *Domain) GetUserByID(id int64) (*User, error) {
	user, err := d.DB.UserRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
