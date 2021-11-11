package Database

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"strconv"
	"strings"
)

// CreateConnection to dynamodb
func CreateConnection(region string) (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}
	return dynamodb.New(sess), nil
}

// DynamoDB is a concrete implementation
// to interface with common DynamoDB operations
type DynamoDB struct {
	table string
	conn  *dynamodb.DynamoDB
}

func NewDynamoDB(conn *dynamodb.DynamoDB, table string) *DynamoDB {
	return &DynamoDB{
		conn: conn, table: table,
	}
}

func (ddb *DynamoDB) List(castTo interface{}, projectionList []string) error {

	var proj expression.ProjectionBuilder
	if len(projectionList) != 0 {
		nameBuilder := expression.Name(projectionList[0])
		var nbList []expression.NameBuilder
		for _, projAttr := range projectionList[1:] {
			nbList = append(nbList, expression.Name(projAttr))
		}
		proj = expression.NamesList(nameBuilder, nbList...)
	}


	var input *dynamodb.ScanInput
	if len(projectionList) != 0 {
		expr, err := expression.NewBuilder().WithProjection(proj).Build()
		if err != nil {
			return err
		}
		input = &dynamodb.ScanInput{
			TableName:            aws.String(ddb.table),
			ProjectionExpression: expr.Projection(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		}
	} else {
		input = &dynamodb.ScanInput{
			TableName: aws.String(ddb.table),
		}
	}

	mappedValues := make([]map[string]*dynamodb.AttributeValue, 0)
	results, err := ddb.conn.Scan(input)
	if err != nil {
		return err
	}
	mappedValues = append(mappedValues, results.Items...)
	for results.LastEvaluatedKey != nil {
		input.ExclusiveStartKey = results.LastEvaluatedKey
		results, err = ddb.conn.Scan(input)
		if err != nil {
			return err
		}
		mappedValues = append(mappedValues, results.Items...)
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(mappedValues, &castTo); err != nil {
		return err
	}
	return nil
}

func (ddb *DynamoDB) ListByInput(castTo interface{}, input *dynamodb.ScanInput) error {
	mappedValues := make([]map[string]*dynamodb.AttributeValue, 0)
	input.TableName = aws.String(ddb.table)
	results, err := ddb.conn.Scan(input)
	if err != nil {
		return err
	}
	mappedValues = append(mappedValues, results.Items...)
	for results.LastEvaluatedKey != nil {
		input.ExclusiveStartKey = results.LastEvaluatedKey
		results, err = ddb.conn.Scan(input)
		if err != nil {
			return err
		}
		mappedValues = append(mappedValues, results.Items...)
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(mappedValues, &castTo); err != nil {
		return err
	}
	return nil
}

func (ddb *DynamoDB) Store(item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(ddb.table),
	}
	_, err = ddb.conn.PutItem(input)
	if err != nil {
		fmt.Println("Store error:", err)
	}
	if err != nil {
		return err
	}
	return err
}

func (ddb *DynamoDB) StoreUnstructured(data map[string]interface{}) error {

	var vv = make(map[string]*dynamodb.AttributeValue)
	for k, v := range data {
		x := v.(string) //assert string type
		xx := &(x)
		vv[k] = &dynamodb.AttributeValue{S: xx}
	}

	input := &dynamodb.PutItemInput{
		Item:      vv,
		TableName: aws.String(ddb.table),
	}
	_, err := ddb.conn.PutItem(input)
	if err != nil {
		return err
	}
	return err
}

func (ddb *DynamoDB) Get(key string, castTo interface{}, keyName string) error {
	result, err := ddb.conn.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ddb.table),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		return err
	}

	if len(result.Item) == 0 {
		castTo = nil
		return nil
	}

	if err := dynamodbattribute.UnmarshalMap(result.Item, &castTo); err != nil {
		return err
	}
	return nil
}

func (ddb *DynamoDB) QueryByPrimaryKey(partitionKey string, partitionKeyValue string, castTo interface{}, projectionList []string) error {

	keyCond := expression.Key(partitionKey).Equal(expression.Value(partitionKeyValue))

	var proj expression.ProjectionBuilder
	if len(projectionList) != 0 {
		nameBuilder := expression.Name(projectionList[0])
		var nbList []expression.NameBuilder
		for _, projAttr := range projectionList[1:] {
			nbList = append(nbList, expression.Name(projAttr))
		}
		proj = expression.NamesList(nameBuilder, nbList...)
	}

	var input *dynamodb.QueryInput
	if len(projectionList) != 0 {
		expr, err := expression.NewBuilder().WithKeyCondition(keyCond).WithProjection(proj).Build()
		if err != nil {
			return err
		}

		input = &dynamodb.QueryInput{
			TableName:                 aws.String(ddb.table),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			ProjectionExpression:      expr.Projection(),
		}
	} else {
		expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
		if err != nil {
			return err
		}

		input = &dynamodb.QueryInput{
			TableName:                 aws.String(ddb.table),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		}
	}

	results, err := ddb.conn.Query(input)
	if err != nil {
		return err
	}
	mappedValues := make([]map[string]*dynamodb.AttributeValue, 0)

	mappedValues = append(mappedValues, results.Items...)
	for results.LastEvaluatedKey != nil {
		input.ExclusiveStartKey = results.LastEvaluatedKey
		results, err = ddb.conn.Query(input)
		if err != nil {
			return err
		}
		mappedValues = append(mappedValues, results.Items...)
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(mappedValues, &castTo); err != nil {
		return err
	}
	return nil
}

func (ddb *DynamoDB) QueryByPartitionKey(partitionKey string, partitionKeyValue string, castTo interface{}, projectionList []string) error {

	keyCond := expression.Key(partitionKey).Equal(expression.Value(partitionKeyValue))

	var proj expression.ProjectionBuilder
	if len(projectionList) != 0 {
		nameBuilder := expression.Name(projectionList[0])
		var nbList []expression.NameBuilder
		for _, projAttr := range projectionList[1:] {
			nbList = append(nbList, expression.Name(projAttr))
		}
		proj = expression.NamesList(nameBuilder, nbList...)
	}

	var input *dynamodb.QueryInput
	if len(projectionList) != 0 {
		expr, err := expression.NewBuilder().WithKeyCondition(keyCond).WithProjection(proj).Build()
		if err != nil {
			return err
		}

		input = &dynamodb.QueryInput{
			TableName:                 aws.String(ddb.table),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			IndexName:                 aws.String(partitionKey + "-index"),
			ProjectionExpression:      expr.Projection(),
		}
	} else {
		expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
		if err != nil {
			return err
		}

		input = &dynamodb.QueryInput{
			TableName:                 aws.String(ddb.table),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			IndexName:                 aws.String(partitionKey + "-index"),
		}
	}

	results, err := ddb.conn.Query(input)

	mappedValues := make([]map[string]*dynamodb.AttributeValue, 0)
	if err != nil {
		return err
	}
	mappedValues = append(mappedValues, results.Items...)
	for results.LastEvaluatedKey != nil {
		input.ExclusiveStartKey = results.LastEvaluatedKey
		results, err = ddb.conn.Query(input)
		if err != nil {
			return err
		}
		mappedValues = append(mappedValues, results.Items...)
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(mappedValues, &castTo); err != nil {
		return err
	}
	return nil
}

func (ddb *DynamoDB) Query(indexName string, partitionKey string, sortingKey string, partitionKeyValue string, timeFrameStart int64, timeFrameEnd int64, castTo interface{}, projectionList []string) error {

	// The measurement id index is the only one that does not use the time as its sorting key. It uses the id as the sorting key.

	var keyCond expression.KeyConditionBuilder
	if strings.Compare(partitionKey, "http-response-status") == 0 {
		val, err := strconv.Atoi(partitionKeyValue)
		if err != nil {
			return err
		}
		keyCond = expression.Key(partitionKey).Equal(expression.Value(val)).And(expression.Key(sortingKey).Between(expression.Value(timeFrameStart), expression.Value(timeFrameEnd)))
	} else {
		keyCond = expression.Key(partitionKey).Equal(expression.Value(partitionKeyValue)).And(expression.Key(sortingKey).Between(expression.Value(timeFrameStart), expression.Value(timeFrameEnd)))
	}

	var proj expression.ProjectionBuilder
	if len(projectionList) != 0 {
		nameBuilder := expression.Name(projectionList[0])
		var nbList []expression.NameBuilder
		for _, projAttr := range projectionList[1:] {
			nbList = append(nbList, expression.Name(projAttr))
		}
		proj = expression.NamesList(nameBuilder, nbList...)
	}

	var input *dynamodb.QueryInput
	if len(projectionList) != 0 {
		expr, err := expression.NewBuilder().WithKeyCondition(keyCond).WithProjection(proj).Build()
		if err != nil {
			return err
		}
		input = &dynamodb.QueryInput{
			TableName:                 aws.String(ddb.table),
			IndexName:                 aws.String(indexName),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			ProjectionExpression:      expr.Projection(),
		}
	} else {
		expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
		if err != nil {
			return err
		}
		input = &dynamodb.QueryInput{
			TableName:                 aws.String(ddb.table),
			IndexName:                 aws.String(indexName),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		}
	}

	results, err := ddb.conn.Query(input)

	mappedValues := make([]map[string]*dynamodb.AttributeValue, 0)
	if err != nil {
		return err
	}
	mappedValues = append(mappedValues, results.Items...)
	for results.LastEvaluatedKey != nil {
		input.ExclusiveStartKey = results.LastEvaluatedKey
		results, err = ddb.conn.Query(input)
		if err != nil {
			return err
		}
		mappedValues = append(mappedValues, results.Items...)
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(mappedValues, &castTo); err != nil {
		return err
	}
	return nil
}

func (ddb *DynamoDB) Delete(id string, keyName string) error {
	_, err := ddb.conn.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(id),
			},
		},
		TableName: aws.String(ddb.table),
	})
	return err
}

func (ddb *DynamoDB) StoreMultiple(items []map[string]interface{}) error {

	var writeRequests []*dynamodb.WriteRequest
	for _, item := range items {
		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			return err
		}
		writeRequests = append(writeRequests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{Item: av},
		})
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			*aws.String(ddb.table): writeRequests,
		},
	}

	_, err := ddb.conn.BatchWriteItem(input)
	if err != nil {
		return err
	}

	return nil
}

//goland:noinspection ALL
func (ddb *DynamoDB) GetMultiple(key string, keys []string) ([]map[string]interface{}, error) {
	resps := []map[string]interface{}{}
	listToBeRead := []map[string]*dynamodb.AttributeValue{}

	for _, k := range keys {

		listToBeRead = append(listToBeRead, map[string]*dynamodb.AttributeValue{
			key: {
				S: aws.String(k),
			},
		})
	}

	requestItems := map[string]*dynamodb.KeysAndAttributes{
		*aws.String(ddb.table): {
			Keys: listToBeRead,
		},
	}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: requestItems,
	}

	result, err := ddb.conn.BatchGetItem(input)
	if err != nil {
		return resps, err
	}

	reqBody := make(map[string]interface{})
	for _, r := range result.Responses {
		for _, rr := range r {
			if err := dynamodbattribute.UnmarshalMap(rr, &reqBody); err != nil {
				return resps, err
			}
			resps = append(resps, reqBody)
		}
	}

	return resps, nil

}

func (ddb *DynamoDB) ListTableNames() ([]*string, error) {

	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	var tables []*string
	//fmt.Printf("Tables:\n")

	for {
		// Get the list of tables
		result, err := ddb.conn.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return tables, err
		}

		tables = append(tables, result.TableNames...)

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}

	return tables, nil
}

func (ddb *DynamoDB) ListByStatusCode(castTo interface{}) error {

	proj := expression.NamesList(expression.Name("http-response-status"))

	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(ddb.table),
	}

	mappedValues := make([]map[string]*dynamodb.AttributeValue, 0)
	results, err := ddb.conn.Scan(input)
	if err != nil {
		return err
	}
	mappedValues = append(mappedValues, results.Items...)
	for results.LastEvaluatedKey != nil {
		input.ExclusiveStartKey = results.LastEvaluatedKey
		results, err = ddb.conn.Scan(input)
		if err != nil {
			return err
		}
		mappedValues = append(mappedValues, results.Items...)
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(mappedValues, &castTo); err != nil {
		return err
	}
	return nil

}
