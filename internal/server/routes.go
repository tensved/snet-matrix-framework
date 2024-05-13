package server

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/services", s.GetServices)
	s.App.Get("/orgs", s.GetOrgs)
	s.App.Get("/health", s.healthHandler)
}
