package common

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"fmt"
)

// AppClaims provides custom claim for JWT
type AppClaims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// claimKey is the key for common.AppClaims values in Contexts. It is
// unexported; clients use common.NewContext and common.FromContext
// instead of using this key directly.
var claimKey key = 0

// using asymmetric crypto/RSA keys
// location of private/public key files
const (
	// openssl genrsa -out app.rsa 1024
	privKeyPath = "keys/app.rsa"
	// openssl rsa -in app.rsa -pubout > app.rsa.pub
	pubKeyPath = "keys/app.rsa.pub"
)

// Private key for signing and public key for verification
var (
	verifyKey, signKey []byte
)

// Read the key files before starting http handlers
func initKeys() {
	var err error

	signKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatalf("[initKeys]: %s\n", err)
	}

	verifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("[initKeys]: %s\n", err)
		panic(err)
	}
	log.Println("initKeys completed.")
}

// Generate JWT token
func GenerateJWT(userName, role string) (string, error) {
	// create a signer for rsa 256
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, AppClaims{
		Username: userName,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
		},
	})
	tokenString, err := t.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Middleware for validating JWT tokens
func Authorize(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// validate the token
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Verify the token with public key, which is the counter part of private key
		return signKey, nil
	})
	log.Printf("token: %s", token)

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError: // JWT validation error
			vErr := err.(*jwt.ValidationError)

			switch vErr.Errors {
			case jwt.ValidationErrorExpired: // JWT expired
				DisplayAppError(
					w,
					err,
					"Access Token is expired, get a new Token",
					401,
				)
				return

			default:
				DisplayAppError(
					w,
					err,
					"Error while parsing the Access Token!",
					500,
				)
				return
			}
		default:
			DisplayAppError(
				w,
				err,
				"Error while parsing Access Token!",
				500,
			)
			return
		}
	}
	if token.Valid {
		ctx := NewContext(context.Background(), token.Claims.(*AppClaims))
		rCtx := r.WithContext(ctx)
		next(w, rCtx)
	} else {
		DisplayAppError(
			w,
			err,
			"Invalid Access Token",
			401,
		)
	}
}

func NewContext(ctx context.Context, appClaims *AppClaims) context.Context {
	return context.WithValue(ctx, claimKey, appClaims)
}

func FromContext(ctx context.Context) (*AppClaims, bool) {
	aC, ok := ctx.Value(claimKey).(*AppClaims)
	return aC, ok
}
