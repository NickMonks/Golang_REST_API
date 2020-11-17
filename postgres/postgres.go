package postgres

import "github.com/go-pg/pg/v10"

// here we need to call go-pg for the database. We could directly call go get... but it's preferable to do it like this:
// In order to install all dependencies go to github of each one and go get
func New(opts *pg.Options) *pg.DB {
	// a quick interesting note: the golang compiler will allocate the memory on the heap, so it will return a pointer
	// but the content won't be destroyed after leaving the scope, in contrast with C/C++ where we explicitly allocate memory on the heap (or making the variable static)
	db := pg.Connect(opts)

	return db
}

// NOTE: migrations directory contains the users table. To do so, we "migrate" (i.e update the database) using the https://github.com/golang-migrate/migrate/tree/master/cmd/migrate tool
// then we call migrate create -ext sql -dir postgres/migrations -seq create_users_table, which creates up and down
// What is up and down?
// the up method is a set of directions for running a migration, while the down method is a set of instructions for reverting a migration.
