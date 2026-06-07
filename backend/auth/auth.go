package auth

import (
	// Std
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	// External
	"github.com/alexedwards/argon2id"
	"github.com/go-chi/jwtauth/v5"
)

// makeLogoutCookie()
//
// MaxAge: -1
// instructs browser to delete matching cookie
func MakeLogoutCookie() http.Cookie {
	return http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/api",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	}
}

func makeJWT(sub string) (string, error) {
	claims := map[string]any{
		"sub": sub,
		"exp": time.Now().Add(jwtExpirationTime).Unix(),
		"iat": time.Now().Unix(),
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

// makeJWTCookie()
//
// Path: "/api"
// browser sends cookie when accessing this path
//
// HttpOnly: true
// prevents JavaScript from accessing the Set-Cookie header
//
// Secure: true
// browser will only send this cookie with HTTPS not HTTP
//
// SameSite: http.SameSiteStrictMode
// browser only sends cookie when accessing from the same site
func makeJWTCookie(sub string) (http.Cookie, error) {
	t, err := makeJWT(sub)
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     "jwt",
		Value:    t,
		Path:     "/api",
		Expires:  time.Now().Add(jwtExpirationTime),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}, nil
}

func MatchPassword(pw string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(pw, hash)
}

func CreateHash(pw string) (string, error) {
	return argon2id.CreateHash(pw, argonParams)
}

func ClaimFromCtx(ctx context.Context) (string, error) {
	claims, ok := ctx.Value("claims").(map[string]any)
	if !ok {
		return "", errors.New("no claims in context")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("sub not in claims")
	}

	return sub, nil
}

func UidFromCtx(ctx context.Context) (uint, error) {
	sub, err := ClaimFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	u64, err := strconv.ParseUint(sub, 10, strconv.IntSize)
	if err != nil {
		return 0, err
	}
	// next line is safe because we used strconv.IntSize above
	return uint(u64), nil
}

func UidFromRequest(r *http.Request) (uint, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return 0, fmt.Errorf("read jwt cookie: %w", err)
	}

	token, err := jwtauth.VerifyToken(tokenAuth, cookie.Value)
	if err != nil {
		return 0, fmt.Errorf("verify jwt token: %w", err)
	}

	sub, ok := token.Subject()
	if !ok {
		return 0, errors.New("read jwt subject: sub not in claims")
	}

	u64, err := strconv.ParseUint(sub, 10, strconv.IntSize)
	if err != nil {
		return 0, fmt.Errorf("parse jwt subject: %w", err)
	}

	return uint(u64), nil
}
