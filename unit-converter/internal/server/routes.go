package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"unit-converter/cmd/web"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()
	r.HandlerFunc(http.MethodGet, "/", s.HomeHandler)

	r.HandlerFunc(http.MethodGet, "/health", s.healthHandler)

	fileServer := http.FileServer(http.FS(web.Files))
	r.Handler(http.MethodGet, "/assets/*filepath", fileServer)

	return r
}

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// home.templ
	err := s.RenderTemplate(w, r, "home.page.tmpl", nil)
	if err != nil {
		log.Fatalf("error rendering template. Err: %v", err)
	}
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["status"] = "healthy"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
