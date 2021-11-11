package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
	AuthService "virologic-serverless/Functions/Modules/Auth/Service"
	"virologic-serverless/Functions/Shared"
	"virologic-serverless/Functions/Shared/Database"
	"virologic-serverless/Functions/Shared/Repositories"
)

//TODO: Handle all returns that may have err defined but have response nil.

var globalRepositoryHandler *Repositories.RepositoryHandler
var tableNames []string

func Router(h *Repositories.RepositoryHandler, conn *dynamodb.DynamoDB) RequestHandleFunc {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		path := request.Resource

		userId := Shared.ExtractUserIdFromRequest(&request)
		fmt.Println("Request Body: " + request.Body + " User Id: " + userId)

		var response events.APIGatewayProxyResponse
		var err error
		switch path {
		case "/auth/sign-up":
			headers := Shared.ExtractValuesFromHeaders(request.Headers)
			response, err = AuthService.RegisterUser(&h.UserRepository,
				&h.JwtRepository, &h.RefreshTokenRepository, &request.Body,
				headers)

			response.Headers = make(map[string]string)
			response.Headers["Allow"] = "true"
			response.Headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
			response.Headers["Access-Control-Allow-Headers"] = "*"
			response.Headers["Access-Control-Allow-Origin"] = "*"

		case "/auth/refresh":
			response, err = AuthService.RefreshUserToken(&h.JwtRepository, &h.RefreshTokenRepository, &request.Body)

			response.Headers = make(map[string]string)
			response.Headers["Allow"] = "true"
			response.Headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
			response.Headers["Access-Control-Allow-Headers"] = "*"
			response.Headers["Access-Control-Allow-Origin"] = "*"

		case "/auth/sign-in":
			headers := Shared.ExtractValuesFromHeaders(request.Headers)
			response, err = AuthService.LoginUser(&h.UserRepository, &h.JwtRepository, &h.RefreshTokenRepository, &request.Body, &headers)

			response.Headers = make(map[string]string)
			response.Headers["Allow"] = "true"
			response.Headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
			response.Headers["Access-Control-Allow-Headers"] = "*"
			response.Headers["Access-Control-Allow-Origin"] = "*"

		case "/auth/revoke":
			response, err = AuthService.LogOutUser(&h.RefreshTokenRepository, userId)

			response.Headers = make(map[string]string)
			response.Headers["Allow"] = "true"
			response.Headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
			response.Headers["Access-Control-Allow-Headers"] = "*"
			response.Headers["Access-Control-Allow-Origin"] = "*"

		case "/auth/all-users":
			response, err = AuthService.GetAllUsers(&h.UserRepository)

			response.Headers = make(map[string]string)
			response.Headers["Allow"] = "true"
			response.Headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
			response.Headers["Access-Control-Allow-Headers"] = "*"
			response.Headers["Access-Control-Allow-Origin"] = "*"

		case "/auth/delete-user":
			response, err = AuthService.DeleteUser(&h.UserRepository, &request.Body)

			response.Headers = make(map[string]string)
			response.Headers["Allow"] = "true"
			response.Headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
			response.Headers["Access-Control-Allow-Headers"] = "*"
			response.Headers["Access-Control-Allow-Origin"] = "*"

		case "/auth/edit-user":
			response, err = AuthService.EditUser(&h.UserRepository, &request.Body)

			response.Headers = make(map[string]string)
			response.Headers["Allow"] = "true"
			response.Headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
			response.Headers["Access-Control-Allow-Headers"] = "*"
			response.Headers["Access-Control-Allow-Origin"] = "*"

		case "/auth/reset-password":
			response, err = AuthService.ResetPassword(&h.UserRepository, &request.Body)

			response.Headers = make(map[string]string)
			response.Headers["Allow"] = "true"
			response.Headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
			response.Headers["Access-Control-Allow-Headers"] = "*"
			response.Headers["Access-Control-Allow-Origin"] = "*"

		case "/save-permissions":
			response, err = AuthService.SavePermissionList(&h.UserRepository, userId, &request.Body)

			response.Headers = make(map[string]string)
			response.Headers["Allow"] = "true"
			response.Headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
			response.Headers["Access-Control-Allow-Headers"] = "*"
			response.Headers["Access-Control-Allow-Origin"] = "*"

		case "/get-permissions":
			response, err = AuthService.GetPermissionList(&h.UserRepository, userId)

			response.Headers = make(map[string]string)
			response.Headers["Allow"] = "true"
			response.Headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
			response.Headers["Access-Control-Allow-Headers"] = "*"
			response.Headers["Access-Control-Allow-Origin"] = "*"

		case "/test":
			response, err = Shared.ApiSuccessResponseHandler("Test succeeded", 200)

		default:
			response, err = events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Invalid payload",
			}, nil
		}
		return response, err
	}

}

func GetDatabaseConnections(conn *dynamodb.DynamoDB) *Repositories.RepositoryHandler {
	//We can make those connections run in parallel

	userDB := Database.NewDynamoDB(conn, os.Getenv("USER_DB"))
	userRepo := Repositories.NewUserRepository(userDB)

	refreshTokenDB := Database.NewDynamoDB(conn, os.Getenv("REFRESH_TOKEN_DB"))
	refreshTokenRepo := Repositories.NewRefreshTokenRepository(refreshTokenDB)

	secureTokenDB := Database.NewDynamoDB(conn, os.Getenv("SECURE_TOKEN_DB"))
	secureTokenRepo := Repositories.NewSecureTokenRepository(secureTokenDB)

	jwtDB := Database.NewDynamoDB(conn, os.Getenv("JWT_DB"))
	jwtRepo := Repositories.NewJwtRepository(jwtDB)

	requestRecordDb := Database.NewDynamoDB(conn, os.Getenv("REQUEST_RECORD_DB"))
	requestRecordRepo := Repositories.NewRequestRecordRepository(requestRecordDb)


	repoHandler := &Repositories.RepositoryHandler{
		RefreshTokenRepository:          *refreshTokenRepo,
		SecureTokenRepository:           *secureTokenRepo,
		JwtRepository:                   *jwtRepo,
		UserRepository:                  *userRepo,
		RequestRecordRepository:         *requestRecordRepo,
	}
	return repoHandler
}

type RequestHandleFunc func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func main() {

	// Create a connection to the datastore, in this case, DynamoDB
	conn, err := Database.CreateConnection(os.Getenv("REGION"))
	if err != nil {
		log.Panic(err)
	}
	globalRepositoryHandler = GetDatabaseConnections(conn)
	router := Router(globalRepositoryHandler, conn)

	if os.Getenv("ENVIRONMENT") == "local" {
		LocalDevelopment()
	} else {
		lambda.Start(router)
	}

}
func localHandler(c echo.Context) error {
	// Create a connection to the datastore, in this case, DynamoDB
	conn, err := Database.CreateConnection(os.Getenv("REGION"))
	if err != nil {
		log.Panic(err)
	}
	router := Router(globalRepositoryHandler, conn)
	request := c.Request()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(request.Body)
	if err != nil {
		return err
	} //TODO: Handle error here.
	requestBody := buf.String()
	fmt.Printf("body: %s\n", requestBody)
	fmt.Printf("method: %s\n", request.Method)
	fmt.Printf("path: %s\n", request.URL.Path)

	response, err := router(context.TODO(), events.APIGatewayProxyRequest{
		Resource:                        request.URL.Path,
		Path:                            request.URL.Path,
		HTTPMethod:                      request.Method,
		Headers:                         nil,
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  nil,
		StageVariables:                  nil,
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{"userId": "27770cd2-2196-444a-b543-6cab4ab92e4f"},
		},
		Body:            requestBody,
		IsBase64Encoded: false,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	return c.JSONBlob(response.StatusCode, []byte(response.Body))
}
func LocalDevelopment() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	// Routes
	e.POST("/auth/sign-up", localHandler)
	e.POST("/auth/refresh", localHandler)
	e.POST("/auth/sign-in", localHandler)
	e.DELETE("/auth/revoke", localHandler)
	e.POST("/save-request", localHandler)
	e.GET("/auth/all-users", localHandler)
	e.POST("/auth/delete-user", localHandler)
	e.POST("/auth/edit-user", localHandler)
	e.POST("/auth/reset-password", localHandler)
	e.GET("/test",localHandler)

	// Start server
	e.Logger.Fatal(e.Start(":6662"))
}

