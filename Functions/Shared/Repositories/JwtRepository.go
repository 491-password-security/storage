package Repositories

import (
	"virologic-serverless/Functions/Shared/Database"
	"virologic-serverless/Functions/Shared/Models"
)

func NewJwtRepository(ds Database.DataStore) *JwtRepository {
	return &JwtRepository{datastore: ds}
}

type JwtRepository struct {
	datastore Database.DataStore
}

func (r *JwtRepository) Get(id string, keyName string) (*Models.Jwt, error) {
	var jwt *Models.Jwt
	if err := r.datastore.Get(id, &jwt, keyName); err != nil {
		return nil, err
	}
	return jwt, nil
}

func (r *JwtRepository) Store(user *Models.Jwt) error {
	return r.datastore.Store(user)
}

/*
func (r *UserRepository) List() (*[]Models.User, error) {
	var users *[]Models.User
	if err := r.datastore.List(&users); err != nil {
		return nil, err
	}
	return users, nil
}

*/

func (r *JwtRepository) Delete(id string, keyName string) error {
	return r.datastore.Delete(id, keyName)
}

func (r *JwtRepository) GetMultiple(key string, keys []string) ([]map[string]interface{},error) {
	return r.datastore.GetMultiple(key, keys)
}

