package store

import "golang.org/x/oauth2"

var (
	ErrorTokenNotFound error
)

type TokenStore interface {
	Get() (*oauth2.Token, error)
	Save(*oauth2.Token) error
}