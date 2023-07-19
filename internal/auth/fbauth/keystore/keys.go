package keystore

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v4"
	"sync"
	"time"
)

type PublicKeys struct {
	cachedKeyMap sync.Map

	KidMap    map[string]string
	ExpiresAt int64 // unix millis
}

func (p *PublicKeys) GetKey(kid string) (*rsa.PublicKey, error) {
	res, ok := p.cachedKeyMap.Load(kid)
	if ok {
		return res.(*rsa.PublicKey), nil
	}

	return p.parseKey(kid)
}

func (p *PublicKeys) parseKey(kid string) (res *rsa.PublicKey, err error) {
	res, err = jwt.ParseRSAPublicKeyFromPEM([]byte(p.KidMap[kid]))
	if err != nil {
		return
	}

	p.cachedKeyMap.Store(kid, res)
	return
}

func (p *PublicKeys) IsExpire() bool {
	if p.ExpiresAt <= time.Now().UnixMilli() {
		return true
	}

	return false
}

type KeyStore interface {
	Get() (*PublicKeys, error)
	Set(newKeys *PublicKeys) error
}
