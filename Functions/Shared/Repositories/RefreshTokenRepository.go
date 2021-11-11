package Repositories

import (
	"virologic-serverless/Functions/Shared/Database"
	"virologic-serverless/Functions/Shared/Models"
)

func NewRefreshTokenRepository(ds Database.DataStore) *RefreshTokenRepository {
	return &RefreshTokenRepository{datastore: ds}
}

type RefreshTokenRepository struct {
	datastore Database.DataStore
}

func (r *RefreshTokenRepository) Get(id string, keyName string) (*Models.RefreshToken, error) {
	var refreshToken *Models.RefreshToken
	if err := r.datastore.Get(id, &refreshToken, keyName); err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (r *RefreshTokenRepository) Store(refreshToken *Models.RefreshToken) error {
	return r.datastore.Store(refreshToken)
}

func (r *RefreshTokenRepository) List(projectionList []string) (*[]Models.RefreshToken, error) {
	var clients *[]Models.RefreshToken
	if err := r.datastore.List(&clients, projectionList); err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *RefreshTokenRepository) Delete(id string, keyName string) error {
	return r.datastore.Delete(id, keyName)
}

func (r *RefreshTokenRepository) GetMultiple(key string, keys []string) ([]map[string]interface{},error) {
	return r.datastore.GetMultiple(key, keys)
}

