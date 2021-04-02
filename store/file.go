package store

import (
	"encoding/json"
	"io/ioutil"

	"golang.org/x/oauth2"
)


type FileStore struct {
	filename string
}

// NewFileStore creates a new instance of a FileStore.
func NewFileStore(filename string) TokenStore{
	return &FileStore{
		filename: filename,
	}
}

// Get returns the token stored in the file
func (s *FileStore) Get() (*oauth2.Token, error) {
	content, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}
	token := &oauth2.Token{}
	err = json.Unmarshal(content, token)
	return token, err
}

// Save stores the token in the file
func (s *FileStore) Save(token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.filename, data, 0777)
}

