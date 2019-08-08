package main

import (
	"net/http"

	"bitbucket.org/SlothNinja/log"
	"bitbucket.org/SlothNinja/sn"
	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	logKind       = "Log"
	batch   int64 = 5
)

type logEntry map[string]interface{}
type Log []logEntry

// type logEntry struct {
// 	Template string  `json:"template"`
// 	Data     logData `json:"data"`
// }

// func newLogEntry(t string, d logData) *logEntry {
// 	return &logEntry{
// 		Template: t,
// 		Data:     d,
// 	}
// }

// type glog struct {
// 	// Key       *datastore.Key `datastore:"__key__" json:"-"`
// 	Entries []*logEntry `datastore:"-" json:"entries"`
// 	Encoded string      `datastore:",noindex" json:"-"`
// 	// CreatedAt time.Time      `json:"createdAt"`
// }

// func newGLog(id, count int64) *glog {
// 	l := new(glog)
// 	l.Key = newLogKey(id, count)
// 	return l
// }

// func newLogKey(id, count int64) *datastore.Key {
// 	return datastore.IDKey(logKind, count, newKey(id))
// }
//
// func (l *glog) Load(ps []datastore.Property) error {
// 	err := datastore.LoadStruct(l, ps)
// 	if err != nil {
// 		return err
// 	}
// 	var entries []*logEntry
// 	err = json.Unmarshal([]byte(l.Encoded), &entries)
// 	if err != nil {
// 		return err
// 	}
// 	l.Entries = entries
// 	return nil
// }
//
// func (l *glog) Save() ([]datastore.Property, error) {
// 	encoded, err := json.Marshal(l.Entries)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	l.Encoded = string(encoded)
// 	return datastore.SaveStruct(l)
// }
//
// func (l *glog) LoadKey(k *datastore.Key) error {
// 	l.Key = k
// 	return nil
// }

// func (l *glog) addEntry(t string, d logData) (*glog, *logEntry) {
// 	e := newLogEntry(t, d)
// 	l = append(l, e)
// 	return l, e
// }

// func (l *glog) addData(d logData) {
// 	entry := l.last()
// 	if entry == nil {
// 		return
// 	}
// 	for k, v := range d {
// 		l.last().Data[k] = v
// 	}
// }
//
func (l Log) last() (logEntry, bool) {
	if len(l) == 0 {
		return logEntry{}, false
	}
	return l[len(l)-1], true
}

func (s *server) getCount(c *gin.Context) (int64, error) {
	count, err := sn.Int64Param(c, countParam)
	if err != nil {
		return -1, errors.WithMessage(err, "unable to get count")
	}
	return count, nil
}

func (s server) getOffset(c *gin.Context) (int64, error) {
	offset, err := sn.Int64Param(c, offsetParam)
	if err != nil {
		return -1, errors.WithMessage(err, "unable to get offset")
	}
	return offset, nil
}

func (s server) getGLog(hidParam, countParam, offsetParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		s, err := s.init(c)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		param, err := s.getParam(c, param)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		count, err := s.getCount(c)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		offset, err := s.getOffset(c)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		if count < 1 {
			c.JSON(http.StatusOK, struct{}{})
			return
		}

		id, l := count-offset, batch
		if id < batch {
			l = id
		}

		histories := make([]history, l)
		ks := make([]*datastore.Key, l)
		for i := range histories {
			ks[i] = newHistoryKey(param, id)
			id--
		}

		err = ignore(s.GetMulti(c, ks, histories), datastore.ErrNoSuchEntity)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": errUnexpected.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"offset": offset + l, "logs": histories})
	}
}

func ignore(err error, ignore error) error {
	if err == nil {
		return nil
	}

	merr, ok := err.(datastore.MultiError)
	if !ok {
		return err
	}

	for _, e := range merr {
		if e != nil && e != ignore {
			return err
		}
	}
	return nil
}
