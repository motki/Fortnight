package localhttp

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/motki/core/model"
	"github.com/motki/core/proto/client"

	"github.com/motki/fortnight/localstore"
)

func (srv *Server) locationHandler(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	locID, _ := strconv.Atoi(vars["locationID"])
	if locID == 0 {
		srv.logger.Debugf("localhttp: couldnt parse location ID from URL: %d", locID)
		respond(w, http.StatusBadRequest, "Invalid location ID")
		return
	}
	var loc *model.Location
	err := srv.store.With(func(s *localstore.Tx) error {
		k := localstore.IntKey(locID)
		b, err := s.Acquire(localstore.KindLocation)
		if err != nil {
			return err
		}
		v, err := b.Get(k)
		if err != nil {
			if err != localstore.ErrNotFound {
				return err
			}
		}
		if v == nil {
			loc, err = srv.cl.GetLocation(locID)
			if err != nil {
				return err
			}
			return b.Put(k, loc)
		}
		var ok bool
		if loc, ok = v.(*model.Location); !ok {
			return errors.New("boom")
		}
		return nil
	})
	if err != nil {
		if err == client.ErrNotAuthenticated {
			respond(w, http.StatusForbidden, "Forbidden")
			return
		}
		srv.logger.Debugf("localhttp: couldnt get location from localstore: %s", err.Error())
		respond(w, http.StatusInternalServerError, "Error in localstore")
		return
	}
	respondOK(w, loc)
}
