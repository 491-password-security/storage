package Service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"os"
	"regexp"
	"virologic-serverless/Functions/Modules/Email/Model"
	"virologic-serverless/Functions/Shared"
	"virologic-serverless/Functions/Shared/Database"
	"virologic-serverless/Functions/Shared/Models"
)

func SendEmailToEmailAddresses(requestBody *string) (events.APIGatewayProxyResponse, error) {

	sendEmailRequest := &Model.SendEmailRequest{}

	if err := json.Unmarshal([]byte(*requestBody), sendEmailRequest); err != nil {
		return Shared.ApiErrorResponseHandler("Invalid request structure", 401)
	}

	_, err := Shared.SendTestEmail("email-service-master-main-module", sendEmailRequest.Message,
		sendEmailRequest.Subject, sendEmailRequest.Emails)
	if err != nil {
		return Shared.ApiErrorResponseHandler("Couldn't send mail to the user", 400)
	}

	return Shared.ApiSuccessResponseHandler("Successfully sent", 200)

}

func SendEmailToService(requestBody *string) (events.APIGatewayProxyResponse, error) {

	sendEmailToServiceRequest := &Model.SendEmailToServiceRequest{}

	if err := json.Unmarshal([]byte(*requestBody), sendEmailToServiceRequest); err != nil {
		return Shared.ApiErrorResponseHandler("Invalid request structure", 401)
	}

	conn, err := Database.CreateConnection(os.Getenv("REGION"))
	if err != nil {
		return Shared.ApiErrorResponseHandler("Couldn't create connection to the database", 400)
	}
	userDB := Database.NewDynamoDB(conn, sendEmailToServiceRequest.Service)


	var userEmails *[]Models.UserEmails
	if err := userDB.List(&userEmails, []string{}); err != nil {
		return Shared.ApiErrorResponseHandler("Error occurred while listing user emails", 400)
	}

	var emailArr []string
	for _, n := range *userEmails {
		// If the email attribute is not empty
		if n.Email != "" {
			emailArr = append(emailArr, n.Email)
		}
	}

	// TODO: herhangi bir mail verified degilse mail gondermiyor, onun disinda calisiyor
	_, SErr := Shared.SendTestEmail("email-service-master-main-module", sendEmailToServiceRequest.Message,
		sendEmailToServiceRequest.Subject, emailArr)
	if SErr != nil {
		return Shared.ApiErrorResponseHandler("Couldn't send mail to the user", 400)
	}

	return Shared.ApiSuccessResponseHandler("Successfully sent", 200)

}

func FilterTable(regex string, tableNames []*string) *[]string {
	var filteredTables []string
	for _, n := range tableNames {
		if matchedTable, _ := regexp.MatchString(regex, *n); matchedTable {
			filteredTables = append(filteredTables, *n)
		}
	}
	return &filteredTables
}

func GetTableNamesWithRegex(regex string) (events.APIGatewayProxyResponse, error) {
	conn, err := Database.CreateConnection(os.Getenv("REGION"))
	if err != nil {
		return Shared.ApiErrorResponseHandler("Couldn't create connection to the database", 400)
	}
	userEmailDB := Database.NewDynamoDB(conn, os.Getenv("USER_DB"))
	allTables, err := userEmailDB.ListTableNames()
	filteredTables := FilterTable(regex, allTables)

	tableNames := Model.TableNamesResponse{
		TableNames: *filteredTables,
	}
	return Shared.ApiResponseHandler(tableNames, 200)
}

func SendSmsToPhoneNumber(requestBody *string) (events.APIGatewayProxyResponse, error) {

	sendSmsRequest := &Model.SendSmsRequest{}

	if err := json.Unmarshal([]byte(*requestBody), sendSmsRequest); err != nil {
		return Shared.ApiErrorResponseHandler("Invalid request structure", 401)
	}

	_, err := Shared.SendSms("sms-service-SMSServiceFunction-1SHOGN9IF0NVB",
		sendSmsRequest.Message, sendSmsRequest.PhoneNumbers)
	if err != nil {
		return Shared.ApiErrorResponseHandler("Couldn't send sms to the user", 400)
	}

	return Shared.ApiSuccessResponseHandler("Successfully sent", 200)

}

func SendSmsToService(requestBody *string) (events.APIGatewayProxyResponse, error) {

	sendSmsToServiceRequest := &Model.SendSmsToServiceRequest{}

	if err := json.Unmarshal([]byte(*requestBody), sendSmsToServiceRequest); err != nil {
		return Shared.ApiErrorResponseHandler("Invalid request structure", 401)
	}

	conn, err := Database.CreateConnection(os.Getenv("REGION"))
	if err != nil {
		return Shared.ApiErrorResponseHandler("Couldn't create connection to the database", 400)
	}
	userDB := Database.NewDynamoDB(conn, sendSmsToServiceRequest.Service)


	var userPhoneNumbers *[]Models.UserPhoneNumbers
	if err := userDB.List(&userPhoneNumbers, []string{}); err != nil {
		return Shared.ApiErrorResponseHandler("Error occurred while listing user phone numbers", 400)
	}

	var phoneArr []string
	for _, n := range *userPhoneNumbers {
		// If the phone attribute is not empty
		if n.Phone != "" {
			phoneArr = append(phoneArr, n.Phone)
		}
	}
	fmt.Println(phoneArr)

	_, SErr := Shared.SendSms("sms-service-SMSServiceFunction-1SHOGN9IF0NVB",
		sendSmsToServiceRequest.Message, phoneArr)
	if SErr != nil {
		return Shared.ApiErrorResponseHandler("Couldn't send sms to the user", 400)
	}

	return Shared.ApiSuccessResponseHandler("Successfully sent", 200)
}

