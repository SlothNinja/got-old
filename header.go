package main

import (
	"encoding/json"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn"
	"github.com/SlothNinja/user"
	"github.com/gin-gonic/gin"
)

const headerKind = sn.HeaderKind

func newHeaderEntity(g game) headerEntity {
	return headerEntity{Header: g.Header, Key: newHeaderKey(g.ID())}
}

func newHeaderKey(id int64) *datastore.Key {
	return datastore.IDKey(headerKind, id, nil)
}

// func New(ctx context.Context, id int64) (h Header) {
// 	h = &Header{
// 		ID:   id,
// 		Kind: kind,
// 		GLog: (glog.NewGLog(ctx, id)),
// 	}
// 	return
// }

type headerEntity struct {
	Key *datastore.Key `datastore:"__key__" json:"-"`
	Header
}

func (e headerEntity) Accept(c *gin.Context, cu user.User) (headerEntity, bool, error) {
	h, start, err := e.Header.Header.Accept(c, cu)
	e.Header.Header = h
	return e, start, err
}

func (e headerEntity) Drop(cu user.User) (headerEntity, error) {
	h, err := e.Header.Header.Drop(cu)
	e.Header.Header = h
	return e, err
}

func (e headerEntity) AddUser(cu user.User) headerEntity {
	h := e.Header.Header.AddUser(cu)
	e.Header.Header = h
	return e
}

// func (e headerEntity) LastUpdated() string {
// 	return sn.LastUpdated(e.UpdatedAt)
// }

func (e headerEntity) LastUpdate() time.Time {
	return e.UpdatedAt
}

func (e headerEntity) ID() int64 {
	return e.Key.ID
}

func (e headerEntity) MarshalJSON() ([]byte, error) {
	status := "Public"
	if e.Password != "" {
		status = "Private"
	}

	type JEntity headerEntity
	return json.Marshal(struct {
		JEntity
		ID          int64  `json:"id"`
		Public      string `json:"public"`
		LastUpdated string `json:"lastUpdated"`
	}{
		JEntity: JEntity(e),
		ID:      e.Key.ID,
		Public:  status,
		// LastUpdated: e.LastUpdated(),
	})
}

// Header provides game/invitation header data
type Header struct {
	TwoThiefVariant bool `json:"twoThief"`
	sn.Header
}

var getHID = sn.GetHID
