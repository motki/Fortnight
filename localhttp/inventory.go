package localhttp

import (
	"net/http"
	"time"

	"github.com/motki/core/evedb"
	"github.com/motki/core/model"

	"strconv"

	"github.com/motki/fortnight/localstore"
)

type inventoryItem struct {
	*model.InventoryItem

	Name string `json:"name"`
}

const inventoryItemTTL = 6 * time.Hour

func (srv *Server) inventoryPurgeHandler(w http.ResponseWriter, req *http.Request) {
	err := srv.store.With(func(s *localstore.Tx) error {
		return s.RemoveBucket(localstore.KindInventoryItem)
	})
	if err != nil {
		srv.logger.Debugf("error purging inventory items from localstorage:\n%s", err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondOK(w, nil)
}

func (srv *Server) inventorySaveHandler(w http.ResponseWriter, req *http.Request) {
	locID, _ := strconv.Atoi(req.PostFormValue("location_id"))
	typeID, _ := strconv.Atoi(req.PostFormValue("type_id"))
	minLevel, _ := strconv.Atoi(req.PostFormValue("minimum_level"))
	if locID <= 0 || typeID <= 0 {
		respond(w, http.StatusInternalServerError, "location and type ID are required")
		return
	}
	it, err := srv.cl.NewInventoryItem(typeID, locID)
	if err != nil {
		srv.logger.Debugf("error creating inventory item:\n%s", err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	it.MinimumLevel = minLevel
	err = srv.cl.SaveInventoryItem(it)
	if err != nil {
		srv.logger.Debugf("error creating inventory item:\n%s", err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = srv.store.With(func(s *localstore.Tx) error {
		b, err := s.Acquire(localstore.KindInventoryItem)
		if err != nil {
			return err
		}
		return b.Put(itID(it), localstore.WithTTL(it, inventoryItemTTL))
	})
	if err != nil {
		// Consider this error non-fatal
		srv.logger.Debugf("error saving new inventory item to localstore:\n%s", err.Error())
	}
	respondOK(w, nil)
}

func (srv *Server) inventoryHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		srv.inventorySaveHandler(w, req)
		return
	}
	var einv []*inventoryItem
	err := srv.store.With(func(s *localstore.Tx) error {
		b, err := s.Acquire(localstore.KindInventoryItem, localstore.WithPrototype(func() localstore.Value {
			return &inventoryItem{}
		}))
		if err != nil {
			return err
		}
		itb, err := s.Acquire(localstore.KindItemType)
		if err != nil {
			return err
		}
		einv, err = srv.inventoryItemsFromLocalstore(b, itb)
		if err != nil || len(einv) == 0 {
			if err != nil {
				srv.logger.Debugf("localhttp: error fetching inventory from localstore: %s", err.Error())
			}
			einv, err = srv.inventoryItemsFromAPI(b, itb)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		srv.logger.Debugf("error: %s", err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondOK(w, einv)
}

func (srv *Server) inventoryItemsFromLocalstore(invItems *localstore.Bucket, itemTypes *localstore.Bucket) ([]*inventoryItem, error) {
	var einv []*inventoryItem
	res, err := invItems.All()
	if err != nil {
		return nil, err
	}
	srv.logger.Debugf("found %d items in localstore", len(res))
	for _, r := range res {
		if rv, ok := r.(*inventoryItem); ok {
			srv.logger.Debugf("item from localstore %v", rv)
			einv = append(einv, rv)
		} else {
			srv.logger.Debugf("item was unexpected type %T", r)
		}
	}
	return einv, nil
}

func itID(it *model.InventoryItem) localstore.Key {
	return localstore.IntPairKey(it.TypeID, it.LocationID)
}

func (srv *Server) enrichItem(itemTypes *localstore.Bucket, it *model.InventoryItem) (*inventoryItem, error) {
	if v, err := itemTypes.Get(itID(it)); err == nil {
		if rv, ok := v.(*evedb.ItemType); ok {
			return &inventoryItem{it, rv.Name}, nil
		}
	}
	itt, err := srv.cl.GetItemType(it.TypeID)
	if err != nil {
		return nil, err
	}
	if err := itemTypes.Put(itID(it), itt); err != nil {
		return nil, err
	}
	return &inventoryItem{it, itt.Name}, nil
}

func (srv *Server) inventoryItemsFromAPI(invItems *localstore.Bucket, itemTypes *localstore.Bucket) ([]*inventoryItem, error) {
	srv.logger.Debugln("fetching items from motki API")
	inv, err := srv.cl.GetInventory()
	if err != nil {
		return nil, err
	}
	var einv []*inventoryItem
	for _, i := range inv {
		it, err := srv.enrichItem(itemTypes, i)
		if err != nil {
			return nil, err
		}
		srv.logger.Debugf("item from api %v", it)
		if err := invItems.Put(itID(it.InventoryItem), localstore.WithTTL(it, inventoryItemTTL)); err != nil {
			return nil, err
		}
		einv = append(einv, it)
	}
	return einv, nil
}
