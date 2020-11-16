package webserver

func (s *WebServer) RegisterRoutes() {
	docRouter := s.router.Group("/api/document/")
	docRouter.POST("set", s.Recovery(s.setDocHandler()))
	docRouter.GET("get", s.Recovery(s.getDocHandler()))
}
