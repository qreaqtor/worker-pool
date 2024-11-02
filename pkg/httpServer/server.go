package httpserver

import (
	"log/slog"
	"net"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

// return http server with added recovery middleware
func NewHTTPServer(handler http.Handler) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Handler: panic(handler),
		},
	}
}

func (h *HTTPServer) Serve(l net.Listener) error {
	slog.Info("Start http server at " + l.Addr().String())
	return h.server.Serve(l)
}

func (h *HTTPServer) Close() error {
	slog.Info("Stop http server")
	return h.server.Close()
}
