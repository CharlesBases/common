package auth

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// defaultFormat default time format
const defaultFormat = "2006-01-02 15:04:05"

// errors for auth
var (
	// errInvalidKey is returned when the service provided an invalid private key
	errInvalidKey = errors.New("an invalid private key was provided")
	// errInvalidToken is returned when the token provided is not valid
	errInvalidToken = errors.New("an invalid token was provided")
	// errUnauthorized token not found
	errUnauthorized = errors.New("account unauthorized")
)

var auth Auth

// Auth .
type Auth interface {
	// Init init options
	Init(...Option) error
	// Options allows you to view the current options.
	Options() Options
	// GenToken generate token for account
	GenToken(id string, opts ...GenOption) (string, error)
	// ParToken parse the token
	ParToken(token string, opts ...ParOption) (*Account, error)
	// ParTokenFromRequest parse token from request
	ParTokenFromRequest(r *http.Request, opts ...ParOption) (*Account, error)
}

// Account .
type Account struct {
	UserID   string                 `json:"user_id,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	Created string `json:"created"`
	Expires string `json:"expires"`

	jwt.StandardClaims `json:"-"`
}

// opt .
type opt struct {
	options Options
}

// InitAuth .
func InitAuth(opts ...Option) error {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	var o = new(opt)
	o.options = options

	// verify options
	if o.options.PrivateKey == "" {
		return errInvalidKey
	}

	auth = o
	return nil
}

// Init .
func (opt *opt) Init(opts ...Option) error {
	for _, o := range opts {
		o(&opt.options)
	}
	return nil
}

// Options .
func (opt *opt) Options() Options {
	return opt.options
}

// GenToken .
func (opt *opt) GenToken(id string, opts ...GenOption) (string, error) {
	var options GenOptions
	for _, o := range opts {
		o(&options)
	}

	var acc = new(Account)
	acc.UserID = id
	acc.Created = time.Now().Format(defaultFormat)

	// set expiry
	var now = time.Now()
	if !options.Expiry.IsZero() {
		acc.Expires = options.Expiry.Format(defaultFormat)
		acc.StandardClaims.ExpiresAt = options.Expiry.Unix()
	} else if options.TTL != 0 {
		acc.Expires = now.Add(options.TTL).Format(defaultFormat)
		acc.StandardClaims.ExpiresAt = now.Add(options.TTL).Unix()
	} else if opt.options.TTL != 0 {
		acc.Expires = now.Add(opt.options.TTL).Format(defaultFormat)
		acc.StandardClaims.ExpiresAt = now.Add(opt.options.TTL).Unix()
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, acc).SignedString(decode(opt.options.PrivateKey))
}

// ParToken .
func (opt *opt) ParToken(tokenString string, opts ...ParOption) (*Account, error) {
	if tokenString != "" {
		var options ParOptions
		for _, o := range opts {
			o(&options)
		}

		token, err := jwt.ParseWithClaims(tokenString, new(Account), func(token *jwt.Token) (interface{}, error) {
			return decode(opt.options.PrivateKey), nil
		})
		if err == nil {
			switch token.Valid {
			case true:
				if account, ok := token.Claims.(*Account); ok {
					return account, err
				}
				fallthrough
			default:
				return nil, errInvalidToken
			}
		}
		return nil, errInvalidToken
	}
	return nil, errUnauthorized
}

// ParTokenFromRequest .
func (opt *opt) ParTokenFromRequest(r *http.Request, opts ...ParOption) (*Account, error) {
	var options ParOptions
	for _, o := range opts {
		o(&options)
	}

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		return decode(opt.options.PrivateKey), nil
	})
	if err == nil {
		switch token.Valid {
		case true:
			if account, ok := token.Claims.(*Account); ok {
				return account, err
			}
			fallthrough
		default:
			return nil, errInvalidToken
		}
	}
	return nil, errUnauthorized
}

// decode string to base64
func decode(source string) []byte {
	data, _ := base64.StdEncoding.DecodeString(source)
	return data
}
