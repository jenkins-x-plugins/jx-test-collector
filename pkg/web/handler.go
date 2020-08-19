package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Options the options for a sync operation
type Options struct {
	// Port the port to listen on HTTP
	Port int `env:"PORT"`

	// Sync performs the sync operation
	Sync func() (string, error)
}

const (
	// HealthPath is the URL path for the HTTP endpoint that returns health status.
	HealthPath = "/health"

	// ReadyPath URL path for the HTTP endpoint that returns ready status.
	ReadyPath = "/ready"

	// SyncPath to invoke a sync operation
	SyncPath = "/sync"
)

// Run will implement this command
func (o *Options) Run() error {
	if o.Port == 0 {
		o.Port = 8080
	}
	mux := http.NewServeMux()
	mux.Handle(HealthPath, http.HandlerFunc(o.health))
	mux.Handle(ReadyPath, http.HandlerFunc(o.ready))
	mux.Handle("/", http.HandlerFunc(o.index))
	mux.Handle(SyncPath, http.HandlerFunc(o.sync))

	logrus.Infof("jx-test-collector is now listening port %d", o.Port)
	return http.ListenAndServe(":"+strconv.Itoa(o.Port), mux)
}

// health returns either HTTP 204 if the service is healthy, otherwise nothing ('cos it's dead).
func (o *Options) health(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Health check")
	w.WriteHeader(http.StatusNoContent)
}

// ready returns either HTTP 204 if the service is ready to serve requests, otherwise HTTP 503.
func (o *Options) ready(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Ready check")
	if o.isReady() {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (o *Options) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from jx-test-collector"))
}

func (o *Options) sync(w http.ResponseWriter, r *http.Request) {
	text, err := o.Sync()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to sync: %s", err.Error())))
		return
	}
	w.Write([]byte(text))
}

func (o *Options) isReady() bool {
	return true
}
