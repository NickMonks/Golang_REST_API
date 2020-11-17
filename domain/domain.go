package domain

type UserRepo interface {
	// As a refresher: Golang can return pointers to a var because is allocated in the heap
	// https://www.geeksforgeeks.org/returning-pointer-from-a-function-in-go/
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(user *User) (*User, error)
	GetByID(id int64) (*User, error)
}

// We create a TODO repo for us:
type TodoRepo interface {
	// CRUD operations for the TODO database
	Create(todo *Todo) (*Todo, error)
	GetByID(id int64) (*Todo, error)
	Update(todo *Todo) (*Todo, error)
	Delete(todo *Todo) error
}

// In order to avoid that other users can delete TODO id's from other users, we created the following interface
type HaveOwner interface {
	IsOwner(user *User) bool
}

//I also create a Domain struct who will keep the DB instance.
//So we can make sure we have only one instance of this last one (the struct Domain will only have one value)
// This will also make life easier and no cycle dependencies issue.

type DB struct {
	UserRepo UserRepo // DB has a UserRepo, which can be any time as long as the methods provided above are implemented
	TodoRepo TodoRepo
}
type Domain struct {
	DB DB // Same for this
	// IMPORTANT: We do DB.UserRepo to create dependency injection.
}
