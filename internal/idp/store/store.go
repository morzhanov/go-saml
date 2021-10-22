package store

import (
	"encoding/json"

	"github.com/crewjam/saml/samlidp"
	"github.com/jmoiron/sqlx"
)

type idpStore struct {
	db *sqlx.DB
}

// Get fetches the data stored in `key` and unmarshals it into `value`.
func (s *idpStore) Get(key string, value interface{}) error {
	var val string
	if err := s.db.QueryRow(`SELECT value FROM idp WHERE key=$1`, key).Scan(&val); err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), &value)
}

// Put marshals `value` and stores it in `key`.
func (s *idpStore) Put(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = s.db.Query(`UPDATE idp SET value=$1`, string(b))
	return err
}

// Delete removes `key`
func (s *idpStore) Delete(key string) error {
	_, err := s.db.Query(`DELETE FROM idp WHERE key=$1`, key)
	return err
}

// List returns all the keys that start with `prefix`. The prefix is
// stripped from each returned value. So if keys are ["aa", "ab", "cd"]
// then List("a") would produce []string{"a", "b"}
func (s *idpStore) List(prefix string) ([]string, error) {
	rows, err := s.db.Query(`SELECT key FROM idp WHERE key LIKE $1`, prefix)
	if err != nil {
		return nil, err
	}
	var res []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		res = append(res, key)
	}
	return res, nil
}

func NewIDPStore(db *sqlx.DB) samlidp.Store {
	return &idpStore{db}
}
