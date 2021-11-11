package Repositories

import (
	uuid "github.com/satori/go.uuid"
	"virologic-serverless/Functions/Shared/Database"
	"virologic-serverless/Functions/Shared/Models"
)

func NewSecureTokenRepository(ds Database.DataStore) *SecureTokenRepository {
	return &SecureTokenRepository{datastore: ds}
}

type SecureTokenRepository struct {
	datastore Database.DataStore
}

/*
func (r *PasswordResetTokenRepository) Get(id string) (*Client, error) {
	var client *Client
	if err := r.datastore.Get(id, &client); err != nil {
		return nil, err
	}
	return client, nil
}

*/

func (r *SecureTokenRepository) Store(secureToken *Models.SecureToken) error {
	id := uuid.NewV4()
	secureToken.ID = id.String()
	return r.datastore.Store(secureToken)
}

/*
func (r *PasswordResetTokenRepository) List() (*[]Client, error) {
	var clients *[]Client
	if err := r.datastore.List(&clients); err != nil {
		return nil, err
	}
	return clients, nil
}

*/

