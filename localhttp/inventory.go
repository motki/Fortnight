package localhttp

import (
	"net/http"

	"github.com/motki/core/model"
	"github.com/motki/core/proto/client"
)

type inventoryItem struct {
	*model.InventoryItem

	Name string `json:"name"`
}

func (srv *Server) inventoryHandler(w http.ResponseWriter, req *http.Request) {
	inv, err := srv.cl.GetInventory()
	if err != nil {

		if err == client.ErrNotAuthenticated {
			respond(w, http.StatusForbidden, "Forbidden")
			return
		}
		srv.logger.Warnf("locathttp: error getting inventory: %s", err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	var einv []*inventoryItem
	for _, i := range inv {
		var name string
		if it, berr := srv.cl.GetItemType(i.TypeID); err != nil {
			srv.logger.Warnf("localhttp: error getting type ID name: %s", berr.Error())
			name = "???"
		} else {
			name = it.Name
		}
		einv = append(einv, &inventoryItem{i, name})
	}
	respondOK(w, einv)
}
