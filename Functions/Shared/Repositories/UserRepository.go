package Repositories

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"virologic-serverless/Functions/Shared/Database"
	"virologic-serverless/Functions/Shared/Models"
)

func NewUserRepository(ds Database.DataStore) *UserRepository {
	return &UserRepository{datastore: ds}
}

type UserRepository struct {
	datastore Database.DataStore
}

func (r *UserRepository) Get(id string, keyName string) (*Models.User, error) {
	var user *Models.User
	if err := r.datastore.Get(id, &user, keyName); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Store(user *Models.User) error {
	//id := uuid.NewV4()
	//passwordResetToken.ID = id.String()
	fmt.Println("Came to Repository Store")
	return r.datastore.Store(user)
}

func (r *UserRepository) List(projectionList []string) (*[]Models.User, error) {
	var users *[]Models.User
	if err := r.datastore.List(&users, projectionList); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) ListByEmail(email string) (*[]Models.User, error) {
	var users *[]Models.User
	filter := expression.Name("email").Equal(expression.Value(email))
	projection := expression.NamesList(expression.Name("email"),
		expression.Name("password"),
		expression.Name("firstName"),
		expression.Name("surname"),
		expression.Name("phone"),
		expression.Name("sex"),
		expression.Name("id"),
		expression.Name("active"))

	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		fmt.Println(err)
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Users"),
	}

	if err := r.datastore.ListByInput(&users, input); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) ListByPhone(phone string) (*[]Models.User, error) {
	var users *[]Models.User
	filter := expression.Name("phone").Equal(expression.Value(phone))
	projection := expression.NamesList(expression.Name("email"),
		expression.Name("password"),
		expression.Name("firstName"),
		expression.Name("surname"),
		expression.Name("phone"),
		expression.Name("sex"),
		expression.Name("id"),
		expression.Name("active"))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		fmt.Println(err)
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Music"),
	}

	if err := r.datastore.ListByInput(&users, input); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Delete(id string, keyName string) error {
	return r.datastore.Delete(id, keyName)
}