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
	Cache             keystore.KeyStore
	refreshKeysLocker sync.Mutex
	keys              *keystore.PublicKeys
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
	v.refreshKeysLocker.Lock()
	defer v.refreshKeysLocker.Unlock()
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
		RawPublicKeys: m,
		ExpiresAt:     exp.UnixMilli(),
	}

	err = v.Cache.Set(v.keys)
	if err != nil {
		v.keys = nil
	}
	return
}
