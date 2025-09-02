package transport

import (
	"fmt"
	"net/http"

	"effective_mobile/docs"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

func NewServer(log *zap.Logger, handlers map[HandlerMetadata]http.HandlerFunc) *HTTPServer {
	return &HTTPServer{
		log:      log,
		server:   &http.Server{},
		handlers: handlers,
	}
}

type HTTPServer struct {
	log      *zap.Logger
	server   *http.Server
	handlers map[HandlerMetadata]http.HandlerFunc
}

func (inst *HTTPServer) Start(address, port string) error {
	docs.SwaggerInfo.Title = "API Swagger"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("127.0.0.1:%s", port)
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r := mux.NewRouter()

	apiRoutes := r.PathPrefix("/api/v1").Subrouter()
	for meta, handler := range inst.handlers {
		apiRoutes.HandleFunc(meta.Path, handler).Methods(meta.Method)
	}

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	inst.server.Handler = r
	inst.server.Addr = fmt.Sprintf("%s:%s", address, port)
	inst.log.Debug("start listen on", zap.String("address", address), zap.String("port", port))

	return inst.server.ListenAndServe()
}

func (inst *HTTPServer) Stop() error {
	return inst.server.Close()
}
