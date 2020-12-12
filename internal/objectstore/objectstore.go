package objectstore

import (
	"database/sql"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/yozel/otrera/internal/log"

	_ "github.com/mattn/go-sqlite3"
)

type Object struct {
	Key               string
	Labels            map[string]string
	CreationTimestamp time.Time
	Data              interface{}
}

var logger zerolog.Logger

func init() {
	logger = log.Log().With().Logger()
}

type ObjectStore struct {
	store     map[string]Object
	db        *sql.DB
	storeLock sync.RWMutex
}

func NewObjectStore() (*ObjectStore, error) {
	err := os.MkdirAll("/tmp/.otrera", 0755)
	if err != nil {
		return nil, err // TODO: wrap error
	}

	db, err := sql.Open("sqlite3", "file:/tmp/.otrera/objectstore.db?cache=shared&mode=rwc")
	if err != nil {
		return nil, err // TODO: wrap error
	}
	db.SetMaxOpenConns(1)

	sqlStmt := `CREATE TABLE IF NOT EXISTS kv (key text NOT NULL PRIMARY KEY, value text NOT NULL);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		logger.Err(err).Str("statement", sqlStmt).Msg(`Can't create "kv" table`)
		return nil, nil
	}

	return &ObjectStore{
		store:     make(map[string]Object),
		db:        db,
		storeLock: sync.RWMutex{},
	}, nil
}

func (s *ObjectStore) Keys() []string {
	s.storeLock.RLock()
	defer s.storeLock.RUnlock()
	keys := make([]string, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k)
	}
	return keys
}

func (s *ObjectStore) Clear() error {
	stmt, err := s.db.Prepare("DELETE FROM kv")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (s *ObjectStore) Set(key string, l map[string]string, c time.Time, d interface{}) error {
	labels := make(map[string]string, len(l))
	for k, v := range l {
		labels[k] = v
	}

	object := Object{Key: key, Labels: l, CreationTimestamp: c, Data: &d}

	stmt, err := s.db.Prepare("INSERT OR REPLACE INTO kv(key, value) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	value, err := json.Marshal(object)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(key, value)
	if err != nil {
		return err
	}

	return nil
}

func (s *ObjectStore) Get(key string) (*Object, error) {

	rows, err := s.db.Query("SELECT value FROM kv WHERE key = ?", key)
	if err != nil {

		return nil, err
	}
	defer rows.Close()
	var value []byte
	if !rows.Next() {
		return nil, nil
	}
	err = rows.Scan(&value)
	if err != nil {
		return nil, err
	}

	var o Object
	err = json.Unmarshal(value, &o)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (s *ObjectStore) GetAll(key string, l map[string]string) ([]Object, error) {
	if l == nil {
		l = map[string]string{}
	}

	r := []Object{}

	rows, err := s.db.Query("SELECT key, value FROM kv WHERE key LIKE ? || '%'", key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
eachobject:
	for rows.Next() {
		var k string
		var vraw []byte

		err = rows.Scan(&k, &vraw)
		if err != nil {
			return nil, err
		}

		var o Object
		err = json.Unmarshal(vraw, &o)
		if err != nil {
			return nil, err
		}

		for lk, lv := range l {
			if o.Labels[lk] != lv {
				continue eachobject
			}
		}
		r = append(r, o)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return r, nil
}
