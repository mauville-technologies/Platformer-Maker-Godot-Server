package controllers

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.HandleFunc("/", s.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", s.SetMiddlewareJSON(s.Login)).Methods("POST")

	// Users Routes
	s.Router.HandleFunc("/users", s.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.GetUsers))).Methods("GET")
	s.Router.HandleFunc("/users/{id}", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.GetUser))).Methods("GET")
	s.Router.HandleFunc("/users/{id}", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}/levels", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.GetLevelsForUser))).Methods("GET")
	s.Router.HandleFunc("/users/{id}", s.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	s.Router.HandleFunc("/levels", s.SetMiddlewareAuthentication(s.PostLevel)).Methods("POST")
	s.Router.HandleFunc("/levels/{id}", s.SetMiddlewareAuthentication(s.GetLevel)).Methods("GET")

	s.Router.HandleFunc("/play/random", s.SetMiddlewareAuthentication(s.GetRandomLevels)).Methods("GET")

}
