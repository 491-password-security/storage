package Shared

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"os"
	"strings"
	"time"
	"virologic-serverless/Functions/Shared/Models"
	"virologic-serverless/Functions/Shared/Repositories"
	Resource "virologic-serverless/Functions/Shared/Resources"
)

func ExtractLanguageFromHeaders(headers map[string]string) string {
	value, exist := headers["accept-language"]
	language := "en"
	if exist {
		if len(value) > 2 {
			language = value[0:2]
		}
	}
	if language != "en" && language != "es" {
		language = "en"
	}
	return language
}

func ExtractTimeZoneFromHeaders(headers map[string]string) string {
	value, exist := headers["x-timezone"]
	if exist {
		return value
	} else {
		return "UTC"
	}
}

func ExtractOSFromHeaders(headers map[string]string) string {
	value, exist := headers["x-os"]
	if exist {
		if value != "Android" && value != "IOS" {
			return "UNKNOWN"
		}
		return value
	} else {
		return "UNKNOWN"
	}
}

func ExtractVersionFromHeaders(headers map[string]string) string {
	value, exist := headers["x-version"]
	if exist {
		return value
	} else {
		return "0.0.1"
	}
}

func ExtractBuildFromHeaders(headers map[string]string) string {
	value, exist := headers["x-build"]
	if exist {
		return value
	} else {
		return "0"
	}
}

func ExtractValuesFromHeaders(headers map[string]string) Models.HeaderValues {
	for k, s := range headers {
		headers[strings.ToLower(k)] = s
	}

	headerValues := Models.HeaderValues{
		Timestamp:       time.Now().UTC().Unix(),
		Model:           ExtractModelsFromHeaders(headers),
		DeviceDetails:   ExtractDeviceDetailsFromHeaders(headers),
		Language:        ExtractLanguageFromHeaders(headers),
		Timezone:        ExtractTimeZoneFromHeaders(headers),
		OperatingSystem: ExtractOSFromHeaders(headers),
		Version:         ExtractVersionFromHeaders(headers),
		Build:           ExtractBuildFromHeaders(headers),
	}

	return headerValues
}

func ExtractModelsFromHeaders(headers map[string]string) string {
	value, exist := headers["x-model"]
	if exist {
		return value
	} else {
		return ""
	}
}

func ExtractDeviceDetailsFromHeaders(headers map[string]string) string {
	value, exist := headers["x-device-details"]
	if exist {
		return value
	} else {
		return ""
	}
}

func ApiResponseHandler(object interface{}, statusCode int) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(object)
	if err != nil {
		errObj := Resource.ApiResponse{
			Success: false,
			Message: "Internal Server Error",
		}
		errBody, _ := json.Marshal(errObj)
		return events.APIGatewayProxyResponse{
			Body:       string(errBody),
			StatusCode: 500,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: statusCode,
	}, nil
}

func ApiErrorResponseHandler(message string,
	statusCode int) (events.APIGatewayProxyResponse, error) {

	errObj := Resource.ApiResponse{
		Success: false,
		Message: message,
	}
	errBody, _ := json.Marshal(errObj)
	return events.APIGatewayProxyResponse{
		Body:       string(errBody),
		StatusCode: statusCode,
	}, nil
}


func SendSms(smsLambdaName string, message string, phoneNumbers []string) (events.APIGatewayProxyResponse, error) {
	region := os.Getenv("REGION")
	session, err := session.NewSession(&aws.Config{ // Use aws sdk to connect to dynamoDB
		Region: &region,
	})
	if err != nil{
		fmt.Printf("Err new session: %s\n", err.Error())
	}

	svc := lambda.New(session, &aws.Config{Region: aws.String(os.Getenv("REGION"))})

	type SMSRequest struct {
		Message      string    `json:"message"`
		Target       string    `json:"target"`
		PhoneNumbers *[]string `json:"phoneNumber"`
	}

	request := SMSRequest{
		Message:      message,
		Target:       "Phone",
		PhoneNumbers: &phoneNumbers,
	}
	payload, Perr := json.Marshal(request)

	if Perr != nil {
		fmt.Println("Json Marshalling error at going request")
	}
	fmt.Println(string(payload))

	input := &lambda.InvokeInput{
		FunctionName:   aws.String(smsLambdaName),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        payload,
	}
	response, Terr := svc.Invoke(input)
	if Terr != nil {
		fmt.Printf("Error while invoking functions: %s\n", Terr.Error())
	}
	fmt.Printf("Response payload: %s\n", string(response.Payload))

	responseBody := &SmsResponse{}
	if err := json.Unmarshal(response.Payload, responseBody); err != nil {
		fmt.Println("Could not map the data")
	}
	fmt.Printf("Code: %d, Message: %s\n", responseBody.Code, responseBody.Message)

	return events.APIGatewayProxyResponse{
		Body:       "Sms request received",
		StatusCode: 200,
	}, nil

}

func SendTestEmail(emailLambdaName string, message string, subject string, emailAddresses []string) (events.APIGatewayProxyResponse, error) {
	region := os.Getenv("REGION")
	session, err := session.NewSession(&aws.Config{ // Use aws sdk to connect to dynamoDB
		Region: &region,
	})
	if err != nil{
		fmt.Printf("Err new session: %s\n", err.Error())
	}

	svc := lambda.New(session, &aws.Config{Region: aws.String(os.Getenv("REGION"))})

	type EmailRequest struct {
		Message        string   `json:"message"`
		Purpose        string   `json:"purpose"` // must be in ["UserAgreement", "Message", "HtmlBody"]
		Subject        string   `json:"subject"`
		EmailAddresses []string `json:"emailAddresses"`
		Region         *string  `json:"region"`
	}

	request := EmailRequest{
		Message:        message,
		Purpose:        "HtmlBody",
		Subject:        subject,
		EmailAddresses: emailAddresses,
	}
	payload, Perr := json.Marshal(request)

	if Perr != nil {
		fmt.Println("Json Marshalling error at going request")
	}
	input := &lambda.InvokeInput{
		FunctionName:   aws.String(emailLambdaName),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        payload,
	}
	response, Terr := svc.Invoke(input)
	if Terr != nil {
		fmt.Printf("Error while invoking functions: %s\n", Terr.Error())
	}
	fmt.Printf("Response payload: %s\n", string(response.Payload))

	responseBody := &SmsResponse{}
	if err := json.Unmarshal(response.Payload, responseBody); err != nil {
		fmt.Println("Could not map the data")
	}
	fmt.Printf("Code: %d, Message: %s\n", responseBody.Code, responseBody.Message)

	return events.APIGatewayProxyResponse{
		Body:       "Sms request received",
		StatusCode: 200,
	}, nil

}

type SmsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ApiSuccessResponseHandler(message string,
	statusCode int) (events.APIGatewayProxyResponse, error) {
	errObj := Resource.ApiResponse{
		Success: true,
		Message: message,
	}
	errBody, _ := json.Marshal(errObj)
	return events.APIGatewayProxyResponse{
		Body:       string(errBody),
		StatusCode: statusCode,
	}, nil
}

func ExtractUserIdFromRequest(request *events.APIGatewayProxyRequest) string {
	return fmt.Sprintf("%v", request.RequestContext.Authorizer["userId"])
}

func ExtractUserRoleFromRequest(request *events.APIGatewayProxyRequest, userRepository *Repositories.UserRepository) (string, error) {
	userId := fmt.Sprintf("%v", request.RequestContext.Authorizer["userId"])

	user, err := userRepository.Get(userId, "id")
	if err != nil{
		return "", err
	}
	return user.Role, nil

}

func Flatten(m map[string]interface{}) map[string]interface{} {
	o := make(map[string]interface{})
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
			nm := Flatten(child)
			for nk, nv := range nm {
				o[k+"-"+nk] = nv
			}
		default:
			o[k] = v
		}
	}
	return o
}


