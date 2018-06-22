package localhttp

import (
	"net/http"

	"github.com/motki/core/evedb"
	"github.com/motki/core/model"

	"github.com/motki/fortnight/localstore"
)

type inventoryItem struct {
	*model.InventoryItem

	Name string `json:"name"`
}

func (srv *Server) inventoryHandler(w http.ResponseWriter, req *http.Request) {
	var einv []*inventoryItem
	err := srv.store.With(func(s *localstore.Store) error {
		b, err := s.Acquire(localstore.KindInventoryItem)
		if err != nil {
			return err
		}
		b.SetPrototype(func() localstore.Value {
			return &inventoryItem{}
		})
		itb, err := s.Acquire(localstore.KindItemType)
		if err != nil {
			return err
		}
		einv, err = srv.inventoryItemsFromLocalstore(b, itb)
		if err != nil || len(einv) == 0 {
			srv.logger.Debugf("localhttp: error fetching inventory from localstore: %s", err.Error())
			einv, err = srv.inventoryItemsFromAPI(b, itb)
			if err != nil {
				srv.logger.Debugf("localhttp: error loading inventory: %s", err.Error())
				respond(w, http.StatusInternalServerError, err.Error())
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

func itID(it *model.InventoryItem) int {
	rid := it.TypeID << 32
	rid |= it.LocationID << 0
	return rid
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
		einv = append(einv, it)
		if err := invItems.Put(itID(it.InventoryItem), it); err != nil {
			return nil, err
		}
	}
	return einv, nil
}
