package htmlhouse

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/juju/errgo"
	"io/ioutil"
	"net/http"
)

const (
	tokenHeader = "Authorization"
)

type sessionManager interface {
	readToken(*http.Request) (string, error)
	createToken(string) (string, error)
	writeToken(http.ResponseWriter, string) error
}

// Basic user session info
type sessionInfo struct {
	ID string `json:"id"`
}

func newSessionInfo(houseID string) *sessionInfo {
	return &sessionInfo{houseID}
}

func newSessionManager(cfg *config) (sessionManager, error) {
	mgr := &defaultSessionManager{}

	var err error

	mgr.signKey, err = ioutil.ReadFile(cfg.PrivateKey)
	if err != nil {
		return mgr, errgo.Mask(err)
	}

	mgr.verifyKey, err = ioutil.ReadFile(cfg.PublicKey)
	if err != nil {
		return mgr, errgo.Mask(err)
	}
	return mgr, nil
}

type defaultSessionManager struct {
	verifyKey, signKey []byte
}

func (m *defaultSessionManager) readToken(r *http.Request) (string, error) {
	tokenString := r.Header.Get(tokenHeader)
	if tokenString == "" {
		return "", nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return m.verifyKey, nil
	})
	switch err.(type) {
	case nil:
		if !token.Valid {
			return "", nil
		}

		houseID := token.Claims["houseID"].(string)

		return houseID, nil
	case *jwt.ValidationError:
		return "", nil
	default:
		return "", errgo.Mask(err)
	}
}

func (m *defaultSessionManager) createToken(houseID string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims["houseID"] = houseID

	tokenString, err := token.SignedString(m.signKey)
	if err != nil {
		return tokenString, errgo.Mask(err)
	}

	return tokenString, nil
}

func (m *defaultSessionManager) writeToken(w http.ResponseWriter, houseID string) error {
	tokenString, err := m.createToken(houseID)
	if err != nil {
		return err
	}

	w.Header().Set(tokenHeader, tokenString)
	return nil
}
