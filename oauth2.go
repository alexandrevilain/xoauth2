package xoauth2

import (
	"context"
	"net/http"

	"github.com/alexandrevilain/xoauth2/store"
	"golang.org/x/oauth2"
)

// Config wraps the oauth2.Config but adds a store to save retrieved tokens.
type Config struct {
	Config *oauth2.Config
	Store store.TokenStore
}

// RestoreToken tries to save the token from the store.
func (c *Config) SaveToken(token *oauth2.Token) error {
	return c.Store.Save(token)
}

// RestoreToken tries to retrieve the token from the store.
func (c *Config) RestoreToken() (*oauth2.Token, error) {
	return c.Store.Get()	
}

// AuthCodeURL wraps oauth2.AuthCodeURL.
func (c *Config) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return c.Config.AuthCodeURL(state, opts...)
}

// Client wraps oauth2.Client but with a specific token source which allows the 
// store to be notified when token is updated.
// Note that the first instanciation of the client doesn't save the token, you must make 
// an explicit call to SaveToken.
func (c *Config) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	initialTokenSource := c.Config.TokenSource(ctx, t)
	tokenSourceWithStore := NewStoreNotifyingSource(c.Store, initialTokenSource)
	tokenSource := oauth2.ReuseTokenSource(t, tokenSourceWithStore)
	return oauth2.NewClient(ctx, tokenSource)
}

// Exchange wraps oauth2.Exchange.
func (c *Config) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return c.Config.Exchange(ctx, code, opts...)
}

// PasswordCredentialsToken wraps oauth2.PasswordCredentialsToken.
func (c *Config) PasswordCredentialsToken(ctx context.Context, username, password string) (*oauth2.Token, error) {
	return c.Config.PasswordCredentialsToken(ctx, username, password)
}

// TokenSource wraps oauth2.TokenSource.
func (c *Config) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	return c.Config.TokenSource(ctx, t)
}
