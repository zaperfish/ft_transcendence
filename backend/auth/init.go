package auth

import (
	// Std
	"errors"
	"os"
	"runtime"
	"strconv"
	"time"

	// External
	"github.com/alexedwards/argon2id"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth
var jwtExpirationTime time.Duration

var argonParams = &argon2id.Params{
	Memory:      128 * 1024,
	Iterations:  4,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}

func Init() error {
	err := initGlobals()
	if err != nil {
		return err
	}

	err = initJWTAuth()
	if err != nil {
		return err
	}

	return nil
}

func initGlobals() error {
	minutesString, ok := os.LookupEnv("JWT_EXPIRATION_MINUTES")
	if !ok {
		return errors.New("JWT_EXPIRATION_MINUTES not set")
	}

	minutes, err := strconv.ParseUint(minutesString, 10, 64)
	if err != nil {
		return err
	}

	jwtExpirationTime = time.Duration(minutes) * time.Minute
	return nil
}

func initJWTAuth() error {
	algorithm, ok := os.LookupEnv("JWT_ALGORITHM")
	if !ok {
		return errors.New("JWT_ALGORITHM not set")
	}

	key, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		return errors.New("JWT_KEY not set")
	}

	tokenAuth = jwtauth.New(algorithm, []byte(key), nil)

	_, ts, err := tokenAuth.Encode(map[string]any{"test": true})
	if err != nil {
		return errors.New("failed to initialize jwt authenticator")
	}

	_, err = jwtauth.VerifyToken(tokenAuth, ts)
	if err != nil {
		return errors.New("failed to initialize jwt authenticator")
	}

	return nil
}
