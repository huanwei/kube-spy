package pkg

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	Router *mux.Router
}

func (s *Server) Initialize() {
	glog.Infof("Intializing server")
	s.Router = mux.NewRouter()
	s.Router.StrictSlash(true)
	glog.Infof("Server initialized")
}

func (s *Server) AddResponseHandler(config RequestConfig,service string) {
	s.Router.Methods(config.Method).
		Path(config.URL).
		HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Server",service)
			w.Write([]byte(fmt.Sprintf("I am %s\n", service)))
			w.Write([]byte(fmt.Sprintf("I received a request at %s\n", time.Now().Format(time.StampMicro))))
		})
}

func (s *Server) AddSendToNextHandler(config RequestConfig, host, service, nextService string, port int) {
	s.Router.Methods(config.Method).
		Path(config.URL).
		HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server",service)
		w.Write([]byte(fmt.Sprintf("I am %s\n", service)))
		w.Write([]byte(fmt.Sprintf("I received a request at %s\n", time.Now().Format(time.StampMicro))))
		response := SendRequest(config, fmt.Sprintf("%s:%d", host, port))
		w.Write([]byte(fmt.Sprintf("I received a response at %s\n", time.Now().Format(time.StampMicro))))
		w.Write([]byte(fmt.Sprintf("Response from %s: {\n", nextService)))
		w.Write(response.Body())
		w.Write([]byte(fmt.Sprintf("}")))
	})
}

func (s *Server) ListenAndServe(port int) {
	http.Handle("/", s.Router)
	glog.Infof("Start listening on %d", port)
	server := &http.Server{Addr:fmt.Sprintf(":%d", port), Handler:s.Router}
	server.ReadTimeout = 2000*time.Millisecond
	server.ListenAndServe()
}
