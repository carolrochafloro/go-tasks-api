package api

import (
	"go-tasks-api/app/internal/logging"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var HTTPService *HTTPServiceT

// routeMux: ponteiro para instância do ServeMux. criado struct para receber as funções correspondentes.
type HTTPServiceT struct {
	router *mux.Router
}

func NewHTTPService() {

	// se HTTPService não tiver sido inicializado, inicia um
	if HTTPService == nil {
		HTTPService = &HTTPServiceT{
			router: mux.NewRouter(),
		}
	}
}

// adicionar novos endpoints ao serviço recebendo a rota e o handler correspondente
func (h *HTTPServiceT) AddEndpoint(endpoint string, f func(http.ResponseWriter, *http.Request)) {
	h.router.HandleFunc(endpoint, f)
}

func (h *HTTPServiceT) StartServer() {
	
	for route, handler := range Routes {
		h.AddEndpoint(route, handler)
	}

	server := &http.Server{
		Addr: os.Getenv("BASE_URL"),
		ReadTimeout: 60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler: h.router,
	}

	logging.Info("Starting HTTP server:", os.Getenv("BASE_URL"))

	err := server.ListenAndServe()

	if err != nil {
		logging.Error("Unable to start server")
		panic(err)
	}
}



