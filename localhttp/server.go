package localhttp // import "github.com/motki/fortnight/localhttp"

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/motki/core/log"
	"github.com/motki/core/proto/client"
)

type Server struct {
	mux *mux.Router
	cl  client.Client

	logger log.Logger
}

func NewServer(cl client.Client, l log.Logger, assetsDir string) *Server {
	r := mux.NewRouter()
	srv := &Server{r, cl, l}
	r.HandleFunc("/inventory", srv.inventoryHandler)
	r.PathPrefix("/location/").Handler(http.StripPrefix("/location/", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		locStr := strings.TrimLeft(request.URL.Path, "/")
		locID, _ := strconv.Atoi(locStr)
		if locID == 0 {
			l.Debugf("localhttp: couldnt parse location ID from URL: %s", locStr)
			respond(w, http.StatusBadRequest, "Invalid location ID")
			return
		}
		loc, err := cl.GetLocation(locID)
		if err != nil {
			if err == client.ErrNotAuthenticated {
				respond(w, http.StatusForbidden, "Forbidden")
				return
			}
			respond(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondOK(w, loc)
	})))
	l.Debugf("localhttp: serving static assets from %s", assetsDir)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(assetsDir)))

	return srv
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
