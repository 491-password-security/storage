package Database

import "github.com/aws/aws-sdk-go/service/dynamodb"

type DataStore interface {
	List(castTo interface{}, projectionList []string) error
	ListByInput(castTo interface{}, input *dynamodb.ScanInput) error
	Get(key string, castTo interface{}, keyName string) error
	Query(indexName string, partitionKey string, sortingKey string, partitionKeyValue string, timeFrameStart int64, timeFrameEnd int64, castTo interface{}, projectionList []string) error
	QueryByPartitionKey(partitionKey string, partitionKeyValue string, castTo interface{}, projectionList []string) error
	QueryByPrimaryKey(partitionKey string, partitionKeyValue string, castTo interface{}, projectionList []string) error
	Store(item interface{}) error
	StoreUnstructured(item map[string]interface{}) error
	StoreMultiple(items []map[string]interface{}) error
	Delete(id string, keyName string) error
	GetMultiple(key string, keys []string) ([]map[string]interface{}, error)
	ListTableNames() ([]*string, error)
}
