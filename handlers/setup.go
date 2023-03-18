package handlers

import (
	"net/http"
	"os"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	fPath    string
}

// SetupHandlers constructs a http.Handler with all application routes defined.
func SetupHandlers(cfg APIMuxConfig) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/", secretHandler)

	return mux
}
