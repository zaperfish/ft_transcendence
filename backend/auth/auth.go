package auth

import (
    // Std
	"net/http"
	"strconv"
	"time"

    // External
	"github.com/alexedwards/argon2id"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

var LogoutCookie = http.Cookie {
		Name:		"auth_token",
		Value:		"",
		Path:		"/",
		HttpOnly:	true,
		Secure:		true,
		MaxAge:		-1,
}	

func makeJWT(sub string) (string, error) {
	claims := map[string]any {
		"sub":		sub,
		"exp":		time.Now().Add(jwtExpirationTime).Unix(),
		"iat":		time.Now().Unix(),
	}
    _, ts, err := tokenAuth.Encode(claims)
    if err != nil {
        return "", err
    }
	return ts, nil
}

func MakeJWTCookieFromID(id uint) (http.Cookie, error) {
	return makeJWTCookie(strconv.FormatUint(uint64(id), 10))
}

func makeJWTCookie(sub string) (http.Cookie, error) {
	t, err := makeJWT(sub)
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie {
		Name:		"auth_token",
		Value:		t,
		Path:		"/",
		Expires:	time.Now().Add(jwtExpirationTime),
		HttpOnly:	true,
		Secure:		true,
		SameSite:	http.SameSiteLaxMode,
	}, nil
}

func MatchPassword(pw string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(pw, hash)
}

func CreateHash(pw string) (string, error) {
	return argon2id.CreateHash(pw, argonParams)
}
