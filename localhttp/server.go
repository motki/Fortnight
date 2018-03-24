package localhttp // import "github.com/motki/fortnight/localhttp"

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/motki/core/log"
	"github.com/motki/core/proto/client"
)

type Server struct {
	mux *http.ServeMux

	logger log.Logger
}

func NewServer(cl client.Client, l log.Logger, assetsDir string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/inventory", func(w http.ResponseWriter, request *http.Request) {
		inv, err := cl.GetInventory()
		if err != nil {
			if err == client.ErrNotAuthenticated {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			l.Warnf("locathttp: error getting inventory: %s", err.Error())
			b, berr := json.Marshal(err)
			if berr != nil {
				l.Warnf("localhttp: unable to marshal error: %s", berr.Error())
			}
			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}
		b, berr := json.Marshal(inv)
		if berr != nil {
			l.Warnf("localhttp: unable to marshal inventory: %s", berr.Error())
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
	mux.Handle("/location/", http.StripPrefix("/location/", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		locStr := strings.TrimLeft(request.URL.Path, "/")
		locID, _ := strconv.Atoi(locStr)
		if locID == 0 {
			l.Debugf("localhttp: couldnt parse location ID from URL: %s", locStr)
			w.WriteHeader(400)
			return
		}
		loc, err := cl.GetLocation(locID)
		if err != nil {
			if err == client.ErrNotAuthenticated {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			b, berr := json.Marshal(err)
			l.Warnf("locathttp: error getting location details: %s", err.Error())
			if berr != nil {
				l.Warnf("localhttp: unable to marshal error: %s", berr.Error())
			}
			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}
		b, berr := json.Marshal(loc)
		if berr != nil {
			l.Warnf("localhttp: unable to marshal location: %s", berr.Error())
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})))
	l.Debugf("localhttp: serving static assets from %s", assetsDir)
	mux.Handle("/", http.FileServer(http.Dir(assetsDir)))

	return &Server{mux, l}
}

func (s *Server) Serve() error {
	logger, closer, err := log.StdLogger(s.logger, "warn")
	if err != nil {
		s.logger.Warnf("localhttp: unable to create stdlib Logger: %s", err.Error())
	}
	if closer != nil {
		// TODO: Close returns an error
		defer closer.Close()
	}
	addr := "localhost:10808"
	s.logger.Debugf("localhttp: listening on %s", addr)
	srv := &http.Server{
		Addr:     addr,
		Handler:  s.mux,
		ErrorLog: logger,
	}
	return srv.ListenAndServe()
}
