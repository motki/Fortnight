package localhttp // import "github.com/motki/fortnight/localhttp"

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/motki/core/log"
	"github.com/motki/core/proto/client"
	"github.com/motki/fortnight/localstore"
)

type Server struct {
	mux *mux.Router
	cl  client.Client

	store  *localstore.Store
	logger log.Logger
}

func NewServer(cl client.Client, l log.Logger, s *localstore.Store, assetsDir string) *Server {
	r := mux.NewRouter()
	srv := &Server{r, cl, s, l}
	r.HandleFunc("/inventory", srv.inventoryHandler)
	r.HandleFunc("/inventory/purge", srv.inventoryPurgeHandler)
	r.HandleFunc("/location/{locationID}", srv.locationHandler)

	l.Debugf("localhttp: serving static assets from %s", assetsDir)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(assetsDir)))

	return srv
}

func (srv *Server) Serve() error {
	logger, closer, err := log.StdLogger(srv.logger, "warn")
	if err != nil {
		srv.logger.Warnf("localhttp: unable to create stdlib Logger: %s", err.Error())
	}
	if closer != nil {
		// TODO: Close returns an error
		defer closer.Close()
	}
	addr := "localhost:10808"
	srv.logger.Debugf("localhttp: listening on %s", addr)
	s := &http.Server{
		Addr:     addr,
		Handler:  srv.mux,
		ErrorLog: logger,
	}
	return s.ListenAndServe()
}
