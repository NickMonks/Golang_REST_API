package postgres

import (
	"errors"
	"todo/domain"

	"github.com/go-pg/pg/v10"
)

// Once we created the interface that we want the user to follow, we create our struct type
// UserRepo which is a DB type.
type UserRepo struct {
	DB *pg.DB
}

func (u *UserRepo) GetByEmail(email string) (*domain.User, error) {

	// the "new" directly creates a new pointer to the variable specified
	// inside domain/users.go type struct
	user := new(domain.User)

	// Now the interesting part, we will use SQL queries to extract the model
	// using golang by using the DB userrepo field.

	err := u.DB.Model(user).Where("email = ?", email).First()
	if err != nil {
		// errors.Is will check if err is pg.ErrNoRows

		if errors.Is(err, pg.ErrNoRows) {
			// since we want to return a error, we add on errors.go that one.
			return nil, domain.ErrNoResult
		}
		return nil, err
	}
	return user, nil
}
func (u *UserRepo) GetByUsername(username string) (*domain.User, error) {
	user := new(domain.User)
	err := u.DB.Model(user).Where("username = ?", username).First()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, domain.ErrNoResult
		}
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) GetByID(id int64) (*domain.User, error) {
	user := new(domain.User)
	err := u.DB.Model(user).Where("id = ?", id).First()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, domain.ErrNoResult
		}
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) Create(user *domain.User) (*domain.User, error) {
	//							return everything as sql
	_, err := u.DB.Model(user).Returning("*").Insert()
	// because we have a pointer, we need to be careful that it's not injecting the pointer
	if err != nil {

		return nil, err
	}
	return user, nil
}

func NewUserRepo(DB *pg.DB) *UserRepo {
	return &UserRepo{DB: DB}
}
