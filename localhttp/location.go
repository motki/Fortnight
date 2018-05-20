package localhttp

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/motki/core/proto/client"
)

func (srv *Server) locationHandler(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	locID, _ := strconv.Atoi(vars["locationID"])
	if locID == 0 {
		srv.logger.Debugf("localhttp: couldnt parse location ID from URL: %d", locID)
		respond(w, http.StatusBadRequest, "Invalid location ID")
		return
	}
	loc, err := srv.cl.GetLocation(locID)
	if err != nil {
		if err == client.ErrNotAuthenticated {
			respond(w, http.StatusForbidden, "Forbidden")
			return
		}
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondOK(w, loc)
}
