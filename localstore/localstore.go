package localstore // import "github.com/motki/fortnight/localstore"

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/coreos/bbolt"
	"github.com/pkg/errors"

	"fmt"

	"github.com/motki/core/evedb"
	"github.com/motki/core/model"
)

const DefaultTTL = 24 * time.Hour

type Kind int

const (
	KindUnknown Kind = iota
	KindLocation
	KindItemType
	KindInventoryItem
)

func (s Kind) String() string {
	switch s {
	case KindLocation:
		return "location"
	case KindItemType:
		return "item_type"
	case KindInventoryItem:
		return "inventory_item"
	default:
		return "unknown"
	}
}

func (s Kind) bucket() []byte {
	return []byte(s.String())
}

func (s Kind) identity() Value {
	switch s {
	case KindLocation:
		return &model.Location{}
	case KindItemType:
		return &evedb.ItemType{}
	case KindInventoryItem:
		return &model.InventoryItem{}
	default:
		return nil
	}
}

type Key []byte

func IntKey(id int) Key {
	return Key(strconv.Itoa(id))
}

func IntPairKey(id1, id2 int) Key {
	return Key(fmt.Sprintf("%d,%d", id1, id2))
}

type Value interface{}

type storedValue struct {
	Value
	CreatedAt time.Time     `json:"__created_at"`
	TTL       time.Duration `json:"__ttl"`
}

func (v storedValue) fresh() bool {
	return v.TTL > 0 && v.CreatedAt.After(time.Now().Add(-1*v.TTL))
}

type Store struct {
	db *bolt.DB
}

func New(dataDir string) (*Store, error) {
	db, err := bolt.Open(dataDir, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return nil, err
	}
	return &Store{db}, nil
}

type BucketOption func(b *Bucket) error

type prototype func() Value

func WithPrototype(p func() Value) BucketOption {
	return func(b *Bucket) error {
		b.prototype = p
		return nil
	}
}

type Tx struct {
	*bolt.Tx
	//s *Store
}

func (tx *Tx) RemoveBucket(kind Kind) error {
	_, err := tx.Acquire(kind)
	if err != nil {
		return err
	}
	// Use the tx that we know was created during Acquire.
	return tx.DeleteBucket(kind.bucket())
}

func (tx *Tx) Acquire(kind Kind, opts ...BucketOption) (*Bucket, error) {
	b, err := tx.CreateBucketIfNotExists(kind.bucket())
	if err != nil {
		return nil, err
	}
	bucket := &Bucket{Bucket: b, kind: kind, prototype: kind.identity}
	for _, o := range opts {
		if err := o(bucket); err != nil {
			return nil, err
		}
	}
	return bucket, nil
}

type withCallback func(*Tx) error

func (s *Store) Begin() (*Tx, error) {
	t, err := s.db.Begin(true)
	if err != nil {
		return nil, err
	}
	return &Tx{t}, nil
}

func (s *Store) With(fn withCallback) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return errors.Wrap(tx.Rollback(), err.Error())
	}
	return tx.Commit()
}

func WithTTL(data Value, ttl time.Duration) Value {
	return storedValue{
		Value: data,
		TTL:   ttl,
	}
}

func withCreatedAt(data Value) Value {
	if v, ok := data.(storedValue); ok {
		v.CreatedAt = time.Now()
		return v
	} else {
		return storedValue{
			Value:     data,
			CreatedAt: time.Now(),
			TTL:       DefaultTTL,
		}
	}
}

func unwrap(data Value) (Value, bool) {
	if v, ok := data.(storedValue); ok {
		return v.Value, v.fresh()
	}
	return data, true
}

type Bucket struct {
	*bolt.Bucket

	kind Kind

	prototype prototype
}

func (bkt *Bucket) Kind() Kind {
	return bkt.kind
}

func (bkt *Bucket) All() ([]Value, error) {
	cur := bkt.Cursor()
	var res []Value
	for k, v := cur.First(); k != nil; k, v = cur.Next() {
		val := storedValue{Value: bkt.prototype()}
		if err := json.Unmarshal(v, &val); err != nil {
			return nil, err
		}
		if v, ok := unwrap(val); ok {
			res = append(res, v)
		}
	}
	return res, nil
}

func (bkt *Bucket) Put(k Key, data Value) error {
	b, err := json.Marshal(withCreatedAt(data))
	if err != nil {
		return err
	}
	return bkt.Bucket.Put([]byte(k), b)
}

var ErrNotFound = errors.New("not found")

func (bkt *Bucket) Get(k Key) (Value, error) {
	v := bkt.Bucket.Get([]byte(k))
	if v == nil {
		return nil, ErrNotFound
	}
	sv := storedValue{Value: bkt.prototype()}
	if err := json.Unmarshal(v, &sv); err != nil {
		return nil, err
	}
	if v, ok := unwrap(sv); ok {
		// Fresh value, return it
		return v, nil
	}
	// Stale, remove it and return not found
	if err := bkt.Delete(k); err != nil {
		return nil, err
	}
	return nil, ErrNotFound
}

func (bkt *Bucket) Delete(k Key) error {
	return bkt.Bucket.Delete([]byte(k))
}
