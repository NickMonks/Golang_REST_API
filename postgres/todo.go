package postgres

import (
	"todo/domain"

	"github.com/go-pg/pg/v10"
)

type TodoRepo struct {
	DB *pg.DB
}

func (t *TodoRepo) GetByID(id int64) (*domain.Todo, error) {
	todo := new(domain.Todo) // new pointer to a Todo object
	// Here we touched postgres
	err := t.DB.Model(todo).Where("id = ?", id).First()
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (t *TodoRepo) Update(todo *domain.Todo) (*domain.Todo, error) {
	// We return everything from the todo.ID and update it
	_, err := t.DB.Model(todo).Where("id = ?", todo.ID).Returning("*").Update()
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (t *TodoRepo) Delete(todo *domain.Todo) error {
	_, err := t.DB.Model(todo).WherePK().Delete()
	if err != nil {
		return err
	}

	return nil
}

func NewTodoRepo(DB *pg.DB) *TodoRepo {
	return &TodoRepo{DB: DB}
}

func (t *TodoRepo) Create(todo *domain.Todo) (*domain.Todo, error) {
	_, err := t.DB.Model(todo).Returning("*").Insert()
	if err != nil {
		return nil, err
	}

	return todo, nil
}
