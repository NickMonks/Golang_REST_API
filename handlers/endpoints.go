package handlers

import (
	"github.com/go-chi/chi"
)

func (s *Server) setupEndpoints(r *chi.Mux) {
	// This function allocates all endpoints from our rest api
	// the client will have to write the pattern to get something
	r.Route("/api/v1", func(r chi.Router) {
		// this function adds another child route to the main api
		// we set it as users
		r.Route("/users", func(r chi.Router) {
			// this function adds another child route to the main api
			// we set it as users

			// here we add a post method, calling the registerUser function
			r.Post("/register", s.registerUser())

			r.Post("/login", s.loginUser())

		})

		r.Route("/todos", func(r chi.Router) {
			// Use the middleware we created
			r.Use(s.withUser)
			r.Post("/", s.createTodo())

			// extract the id from the context
			r.Route("/{id}", func(r chi.Router) {
				// and now use the todo context in the middleware
				r.Use(s.todoCtx)

				// verify that the user of the context is the owner of the todo id:
				// we passs the subject type. In our case, is the "todo"
				r.Use(s.withOwner("todo"))

				r.Patch("/", s.updateTodo())
				r.Delete("/", s.deleteTodo())
			})
		})

	})

}
