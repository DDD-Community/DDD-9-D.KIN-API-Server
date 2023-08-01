package keystore

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v4"
	"sync"
	"time"
)

type PublicKeys struct {
	RawPublicKeys    map[string]string
	ExpiresAt        int64 // unix millis
	cachedPublicKeys sync.Map
}

func (p *PublicKeys) GetKey(kid string) (*rsa.PublicKey, error) {
	res, ok := p.cachedPublicKeys.Load(kid)
	if ok {
		return res.(*rsa.PublicKey), nil
	}

	return p.parseKey(kid)
}

func (p *PublicKeys) parseKey(kid string) (res *rsa.PublicKey, err error) {
	res, err = jwt.ParseRSAPublicKeyFromPEM([]byte(p.RawPublicKeys[kid]))
	if err != nil {
		return
	}

	p.cachedPublicKeys.Store(kid, res)
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
