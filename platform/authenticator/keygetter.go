package authenticator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-jose/go-jose"
	"github.com/golang-jwt/jwt"
)

func NewKeyGetter(wellKnownUrl string) JwtKeyGetter {
	return JwtKeyGetter{
		cache:        nil,
		wellKnownUrl: wellKnownUrl,
		client:       http.Client{},
	}
}

type JwtKeyGetter struct {
	wellKnownUrl string
	cache        *jose.JSONWebKeySet
	client       http.Client
}

func (kg *JwtKeyGetter) GetKey(token *jwt.Token) (any, error) {
	if kg.cache != nil {
		fmt.Printf("Returning cached key\n")
		return kg.cache.Keys[0].Key, nil
	}
	fmt.Printf("Fetching key from %s\n", kg.wellKnownUrl)
	client := http.Client{}
	resp, err := client.Get(kg.wellKnownUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	keySet := jose.JSONWebKeySet{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &keySet)
	if err != nil {
		return nil, err
	}
	if len(keySet.Keys) == 0 {
		return nil, fmt.Errorf("no keys found in JWKS from %s", kg.wellKnownUrl)
	}
	kg.cache = &keySet
	return keySet.Keys[0].Key, nil
}
