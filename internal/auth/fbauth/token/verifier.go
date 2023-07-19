package token

import (
	"d.kin-app/internal/auth/fbauth/keystore"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"sync"
	"time"
)

func NewVerifier(cache keystore.KeyStore) *Verifier {
	if cache == nil {
		cache = keystore.DefaultKeyCache
	}
	return &Verifier{
		Cache: cache,
	}
}

type Verifier struct {
	rl    sync.Mutex
	Cache keystore.KeyStore
	keys  *keystore.PublicKeys
}

func (v *Verifier) Verify(idToken string) (res Token, err error) {
	res.Token, err = jwt.ParseWithClaims(idToken, new(Claims), v.KeyFunc)
	return
}

func (v *Verifier) KeyFunc(token *jwt.Token) (interface{}, error) {
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("must need kid")
	}

	if err := v.loadOrRefreshKeys(); err != nil {
		return nil, err
	}

	return v.keys.GetKey(kid)
}

func (v *Verifier) loadOrRefreshKeys() (err error) {
	if v.keys == nil {
		err = v.loadKeys()
	}

	if err != nil {
		if errors.Is(err, keystore.ErrNotFound) {
			return v.refreshKeys()
		}

		return
	}

	if v.keys.IsExpire() {
		return v.refreshKeys()
	}

	return
}

func (v *Verifier) loadKeys() (err error) {
	v.keys, err = v.Cache.Get()
	return
}

func (v *Verifier) refreshKeys() (err error) {
	v.rl.Lock()
	defer v.rl.Unlock()
	if v.keys != nil && !v.keys.IsExpire() {
		return
	}

	resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	exp, err := time.Parse(time.RFC1123, resp.Header.Get("Expires"))
	if err != nil {
		return
	}

	var m map[string]string
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return
	}

	v.keys = &keystore.PublicKeys{
		KidMap:    m,
		ExpiresAt: exp.UnixMilli(),
	}

	err = v.Cache.Set(v.keys)
	if err != nil {
		v.keys = nil
	}
	return
}

//type PublicKeys struct {
//	publicKeyMap sync.Map
//
//	KidMap    map[string]string
//	ExpiresAt int64 // unix millis
//}
//
//func (p *PublicKeys) GetKey(kid string) (*rsa.PublicKey, error) {
//	res, ok := p.publicKeyMap.Load(kid)
//	if ok {
//		return res.(*rsa.PublicKey), nil
//	}
//
//	return p.parseKey(kid)
//}
//
//func (p *PublicKeys) parseKey(kid string) (res *rsa.PublicKey, err error) {
//	res, err = jwt.ParseRSAPublicKeyFromPEM([]byte(p.KidMap[kid]))
//	if err != nil {
//		return
//	}
//
//	p.publicKeyMap.Store(kid, res)
//	return
//}
//
//func (p *PublicKeys) IsExpire() bool {
//	if p.ExpiresAt <= time.Now().UnixMilli() {
//		return true
//	}
//
//	return false
//}
//
//type KeyStore interface {
//	Get() (*PublicKeys, error)
//	Set(newKeys *PublicKeys) error
//}
