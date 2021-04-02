package xoauth2

import (
	"github.com/alexandrevilain/xoauth2/store"
	"golang.org/x/oauth2"
)

// StoreNotifyingSource is an oauth2.TokenSource that store new token when obtained.
type StoreNotifyingSource struct {
	src oauth2.TokenSource
	store store.TokenStore
}

// NewStoreNotifyingSource creates a StoreNotifyingSource from an underlying token source
// and saves new token where they are obtained.
func NewStoreNotifyingSource(store store.TokenStore, src oauth2.TokenSource) *StoreNotifyingSource {
	return &StoreNotifyingSource{
		store: store,
		src: src,
	}
}

// Token fetches a new token from the underlying source and saves it.
func (s *StoreNotifyingSource) Token() (*oauth2.Token, error) {
	t, err := s.src.Token()
	if err != nil {
		return nil, err
	}
	err = s.store.Save(t)
	return t, err
}