package Repositories

import (
	"virologic-serverless/Functions/Shared/Database"
)

func NewRequestRecordRepository(ds Database.DataStore) *RequestRecordRepository {
	return &RequestRecordRepository{datastore: ds}
}

type RequestRecordRepository struct {
	datastore Database.DataStore
}

func (r *RequestRecordRepository) Get(id string, keyName string) (map[string]interface{}, error) {
	var req map[string]interface{}
	if err := r.datastore.Get(id, &req, keyName); err != nil {
		return nil, err
	}
	return req, nil
}

func (r *RequestRecordRepository) Query(indexName string, partitionKey string, sortingKey string, partitionKeyValue string, timeFrameStart int64, timeFrameEnd int64, projectionList []string) ([]map[string]interface{}, error) {
	var req []map[string]interface{}
	if err := r.datastore.Query(indexName, partitionKey, sortingKey, partitionKeyValue, timeFrameStart, timeFrameEnd, &req, projectionList); err != nil {
		return nil, err
	}
	return req, nil
}

func (r *RequestRecordRepository) Store(request map[string]interface{}) error {
	return r.datastore.Store(request)
}

func (r *RequestRecordRepository) StoreMultiple(request []map[string]interface{}) error {
	return r.datastore.StoreMultiple(request)
}

func (r *RequestRecordRepository) Delete(id string, keyName string) error {
	return r.datastore.Delete(id, keyName)
}

func (r *RequestRecordRepository) GetMultiple(key string, keys []string) ([]map[string]interface{},error) {
	return r.datastore.GetMultiple(key, keys)
}

func (r *RequestRecordRepository) QueryByPartitionKey(partitionKey string, partitionKeyValue string, projectionList []string) ([]map[string]interface{}, error) {
	var req []map[string]interface{}
	if err := r.datastore.QueryByPartitionKey(partitionKey, partitionKeyValue, &req, projectionList); err != nil {
		return nil, err
	}
	return req, nil
}

func (r *RequestRecordRepository) QueryByPrimaryKey(partitionKey string, partitionKeyValue string, projectionList []string)  ([]map[string]interface{}, error)  {
	var req []map[string]interface{}
	if err := r.datastore.QueryByPrimaryKey(partitionKey, partitionKeyValue, &req, projectionList); err != nil {
		return nil, err
	}
	return req, nil
}