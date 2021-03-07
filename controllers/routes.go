package controllers

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.HandleFunc("/", s.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", s.SetMiddlewareJSON(s.Login)).Methods("POST", "OPTIONS")

	// Users Routes
	s.Router.HandleFunc("/users", s.SetMiddlewareJSON(s.CreateUser)).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/users", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.GetUsers))).Methods("GET")
	s.Router.HandleFunc("/users/{id}", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.GetUser))).Methods("GET")
	s.Router.HandleFunc("/users/{id}", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT", "OPTIONS")
	s.Router.HandleFunc("/users/{id}/levels", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.GetLevelsForUser))).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/users/{id}", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.DeleteUser))).Methods("DELETE")

	s.Router.HandleFunc("/levels", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.PostLevel))).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/levels/{id}", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.GetLevel))).Methods("GET")

	s.Router.HandleFunc("/play/random", s.SetMiddlewareJSON(s.SetMiddlewareAuthentication(s.GetRandomLevels))).Methods("GET", "OPTIONS")

}
