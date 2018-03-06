package captcha

import (
	"fay-blog/config"
	"time"
)

type Store struct {
	Expiration time.Duration
}

func (s *Store) Get(id string, clear bool) (digits []byte) {
	if config.Cache.Has(id) {
		digits = []byte(config.Cache.Get(id).(string))
		if clear {
			config.Cache.Forget(id)
		}
	}
	return
}

func (s *Store) Set(id string, digits []byte) {
	config.Cache.Put(id, digits, s.Expiration)
}
