package keystore

import "sync"

var DefaultKeyCache KeyStore = new(defaultKeyCache)

type defaultKeyCache struct {
	locker sync.RWMutex
	stored *PublicKeys
}

func (c *defaultKeyCache) Get() (*PublicKeys, error) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	if c.stored == nil {
		return nil, ErrNotFound
	}

	return c.stored, nil
}

func (c *defaultKeyCache) Set(newKeys *PublicKeys) error {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.stored = newKeys
	return nil
}
