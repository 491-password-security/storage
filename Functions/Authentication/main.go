package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"strings"
	"time"
	"virologic-serverless/Functions/Shared/Database"
	"virologic-serverless/Functions/Shared/Repositories"
)

type AccessDetails struct {
	AccessUuid string
	UserId     string
	Id         string
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	fmt.Println("Verify Token: ", err)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenMetadata(tokenString string) (*AccessDetails, error) {
	token, err := VerifyToken(tokenString)
	fmt.Println("Token Err", err)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}
		Id, ok := claims["id"].(string)
		if !ok {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
			Id:         Id,
		}, nil
	}
	return nil, err
}

func FetchAuth(authD *AccessDetails, jwtRepo Repositories.JwtRepository) (string, error) {
	fmt.Printf("AuthD: %s , %s\n", authD.UserId, authD.AccessUuid)

	jwtToken, err := jwtRepo.Get(authD.Id, "id")
	if err != nil {
		fmt.Println("JwtRepoErr: ", err)
		return "", err
	}

	fmt.Printf("Jwt: %s , %s\n", jwtToken.UserId, jwtToken.TokenUuid)
	if authD.AccessUuid != jwtToken.TokenUuid {
		return "", errors.New("unauthorized")
	}
	if jwtToken.Expires < time.Now().UTC().Unix() {
		fmt.Printf("Expires")
		return "", errors.New("unauthorized")
	}
	return authD.UserId, nil
}

type RequestHandleFunc func(ctx context.Context,
	event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)

func Router(jwtRepo *Repositories.JwtRepository) RequestHandleFunc {
	return func(ctx context.Context,
		event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
		fmt.Println(event)

		token := event.AuthorizationToken
		tokenSlice := strings.Split(token, " ")
		var bearerToken string
		if len(tokenSlice) > 1 {
			bearerToken = tokenSlice[len(tokenSlice)-1]
		}
		ad, err := ExtractTokenMetadata(bearerToken)
		fmt.Println("ET: ", err)
		if err != nil {
			return events.APIGatewayCustomAuthorizerResponse{}, errors.New("unauthorized") // Return a 401 Unauthorized response
		}
		userId, uIDErr := FetchAuth(ad, *jwtRepo)
		if uIDErr != nil {
			return events.APIGatewayCustomAuthorizerResponse{}, errors.New("unauthorized") // Return a 401 Unauthorized response

		}

		return generatePolicy("user", "Allow", event.MethodArn, userId), nil
	}
}

// Help function to generate an IAM policy
func generatePolicy(principalId, effect, resource string, userId string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	// Optional output with custom properties of the String, Number or Boolean type.
	authResponse.Context = map[string]interface{}{
		"userId": userId,
		//"language": Shared.ExtractLanguageFromHeaders(headers),
	}
	return authResponse
}

func main() {
	conn, err := Database.CreateConnection(os.Getenv("REGION"))
	if err != nil {
		log.Panic(err)
	}
	jwtDB := Database.NewDynamoDB(conn, os.Getenv("JWT_DB"))
	jwtRepo := Repositories.NewJwtRepository(jwtDB)

	router := Router(jwtRepo)

	lambda.Start(router)
}
