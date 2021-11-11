package Service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"os"
	"time"
	"virologic-serverless/Functions/Modules/Auth/Model"
	"virologic-serverless/Functions/Shared"
	"virologic-serverless/Functions/Shared/Models"
	"virologic-serverless/Functions/Shared/Repositories"
)

func RegisterUser(userRepository *Repositories.UserRepository,
	jwtRepository *Repositories.JwtRepository,
	refreshRepository *Repositories.RefreshTokenRepository,
	requestBody *string,
	headers Models.HeaderValues) (events.APIGatewayProxyResponse, error) {
	signUpRequest := &Model.SignUpRequest{}
	if err := json.Unmarshal([]byte(*requestBody), signUpRequest); err != nil {
		return Shared.ApiErrorResponseHandler("Invalid request structure", 400)
	}

	user, CreateErr := CreateUser(signUpRequest, userRepository)

	if CreateErr != nil {
		if CreateErr.Error() == "email_exist" {
			return Shared.ApiErrorResponseHandler("Email Exist", 400)
		} else if CreateErr.Error() == "phone_exist" {
			return Shared.ApiErrorResponseHandler("Phone number Exist", 400)
		} else {
			fmt.Println(CreateErr.Error())
			return Shared.ApiErrorResponseHandler("Error Occurred", 400)
		}

	}

	if saveErr := userRepository.Store(&user); saveErr != nil {
		return Shared.ApiErrorResponseHandler("User could not saved", 500)
	}

	return Login(user, jwtRepository, refreshRepository)

}

func CreateUser(request *Model.SignUpRequest, userRepository *Repositories.UserRepository) (Models.User, error) {

	ifNonExist, err := CheckEmailsAndPhones(request.Email, request.Phone, userRepository)

	if !ifNonExist {
		return Models.User{}, err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		return Models.User{}, err
	}
	password := string(passwordHash[:])

	user := Models.User{
		Email:               request.Email,
		Password:            password,
		FirstName:           request.FirstName,
		Surname:             request.Surname,
		Phone:               request.Phone,
		Sex:                 request.Sex,
		CreatedAt:           time.Now().Unix(),
		Id:                  uuid.NewV4().String(),
		Active:              true,
		Role:				 request.Role,
	}

	return user, nil
}

func CheckEmailsAndPhones(email string,
	phone string,
	userRepository *Repositories.UserRepository) (bool, error) {
	users, err := userRepository.ListByEmail(email)
	if err != nil {
		return false, err
	}
	userList := *users
	for i := range userList {
		if userList[i].Email == email {
			return false, fmt.Errorf("email_exist")
		}
		if userList[i].Phone == phone {
			return false, fmt.Errorf("phone_exist")
		}
	}
	return true, nil

}

type TokenDetails struct {
	AccessTokenId string
	AccessToken   string
	RefreshToken  string
	AccessUuid    string
	RefreshUuid   string
	AtExpires     int64
	RtExpires     int64
}

func LogOutUser(
	refreshRepository *Repositories.RefreshTokenRepository,
	userId string) (events.APIGatewayProxyResponse, error) {
	refreshTokens, err := refreshRepository.List([]string{})

	if err != nil {
		return Shared.ApiErrorResponseHandler("Error Occurred", 500)
	}

	for _, r := range *refreshTokens {
		if r.UserID == userId {
			err := refreshRepository.Delete(r.RefreshToken, "refreshToken")
			if err != nil {
				return Shared.ApiErrorResponseHandler("Error occurred logging out", 450)
			}
		}
	}

	return Shared.ApiSuccessResponseHandler("Successfully Logout", 200)
}
func LoginUser(userRepository *Repositories.UserRepository,
	jwtRepository *Repositories.JwtRepository,
	refreshRepository *Repositories.RefreshTokenRepository,
	requestBody *string,
	headers *Models.HeaderValues) (events.APIGatewayProxyResponse, error) {
	signUpRequest := &Model.SignInRequest{}
	if err := json.Unmarshal([]byte(*requestBody), signUpRequest); err != nil {
		return Shared.ApiErrorResponseHandler("Invalid request structure", 400)
	}

	user, GetErr := GetUser(signUpRequest, userRepository)

	if GetErr != nil {
		return Shared.ApiErrorResponseHandler("Email or password is incorrect", 403)
	}

	return Login(user, jwtRepository, refreshRepository)

}
func GetUser(request *Model.SignInRequest, userRepository *Repositories.UserRepository) (Models.User, error) {
	users, err := userRepository.ListByEmail(request.Email)
	userList := *users
	if err != nil {
		return Models.User{}, err
	}
	for i := range userList {
		if userList[i].Email == request.Email {
			match := bcrypt.CompareHashAndPassword([]byte(userList[i].Password), []byte(request.Password))
			if match != nil {
				log.Printf("Password error: %s %s\n", request.Password, userList[i].Password)
				return Models.User{}, fmt.Errorf("email or password is incorrect")
			}
			return userList[i], nil
		}
	}
	return Models.User{}, fmt.Errorf("email does not Exist")
}

func Login(user Models.User,
	jwtRepository *Repositories.JwtRepository,
	refreshRepository *Repositories.RefreshTokenRepository) (events.APIGatewayProxyResponse, error) {

	ts, err := CreateToken(user.Id)
	if err != nil {
		return Shared.ApiErrorResponseHandler("Authentication token could not be created", 500)

	}
	saveErr := CreateAuth(user.Id, ts, jwtRepository, refreshRepository)
	if saveErr != nil {
		fmt.Println(saveErr)
		return Shared.ApiErrorResponseHandler("Authentication token could not be saved", 500)
	}


	response := &Model.AuthenticationResponse{
		AccessToken:  ts.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: ts.RefreshToken,
		Role: 		  user.Role,
	}

	apiResponse, err := Shared.ApiResponseHandler(response, 200)

	apiResponse.Headers = make(map[string]string)
	apiResponse.Headers["Allow"] = "true"
	apiResponse.Headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
	apiResponse.Headers["Access-Control-Allow-Headers"] = "*"
	apiResponse.Headers["Access-Control-Allow-Origin"] = "*"

	return apiResponse, err
}

func RefreshUserToken(jwtRepository *Repositories.JwtRepository,
	refreshRepository *Repositories.RefreshTokenRepository,
	requestBody *string) (events.APIGatewayProxyResponse, error) {
	refreshTokenRequest := &Model.RefreshTokenRequest{}
	if err := json.Unmarshal([]byte(*requestBody), refreshTokenRequest); err != nil {
		return Shared.ApiErrorResponseHandler("Invalid request structure", 401)

	}
	refTok, err := refreshRepository.Get(refreshTokenRequest.RefreshToken, "refreshToken")
	if err != nil {
		return Shared.ApiErrorResponseHandler("Refresh token is not valid", 401)
	}

	token, err := jwt.Parse(refTok.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return Shared.ApiErrorResponseHandler("Refresh token expired", 401)

	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return Shared.ApiErrorResponseHandler("Refresh token is not valid", 401)
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {

		userId, ok := claims["user_id"]
		if !ok {
			return Shared.ApiErrorResponseHandler("Refresh token is not valid", 401)
		}
		if userId != refTok.UserID {
			return Shared.ApiErrorResponseHandler("Refresh token is not valid", 401)

		}
		AccessUuid := uuid.NewV4().String()
		AccessTokenId := uuid.NewV4().String()
		Expires := time.Now().UTC().Add(time.Minute * 1000000).Unix()
		atClaims := jwt.MapClaims{}
		atClaims["authorized"] = true
		atClaims["access_uuid"] = AccessUuid
		atClaims["user_id"] = userId.(string)
		atClaims["id"] = AccessTokenId
		atClaims["exp"] = Expires
		at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
		AccessToken, AtErr := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
		if AtErr != nil {
			return Shared.ApiErrorResponseHandler("Access token could not created", 401)
		}
		accJwt := Models.Jwt{
			TokenUuid: AccessUuid,
			UserId:    userId.(string),
			Expires:   Expires,
			Id: AccessTokenId,
		}
		fmt.Println("accjwt", &accJwt)
		errAccess := jwtRepository.Store(&accJwt)
		fmt.Println("Err access", errAccess)
		if errAccess != nil {
			return Shared.ApiErrorResponseHandler("Db connection returned error", 401)
		}
		response := &Model.AccessRefreshTokenResponse{
			AccessToken: AccessToken,
			TokenType:   "Bearer",
		}
		return Shared.ApiResponseHandler(response, 200)
	} else {
		return Shared.ApiErrorResponseHandler("Refresh token is not valid", 401)
	}
}

func CreateAuth(userid string,
	td *TokenDetails,
	jwtRepository *Repositories.JwtRepository,
	refreshRepository *Repositories.RefreshTokenRepository) error {
	accJwt := Models.Jwt{
		TokenUuid: td.AccessUuid,
		UserId:    userid,
		Expires:   td.AtExpires,
		Id:        td.AccessTokenId,
	}
	errAccess := jwtRepository.Store(&accJwt)
	if errAccess != nil {
		return errAccess
	}
	refToken := Models.RefreshToken{
		UserID:       userid,
		RefreshToken: td.RefreshToken,
		CreatedAt:    time.Now().UTC().Unix(),
		UpdatedAt:    time.Now().UTC().Unix(),
		ID:           uuid.NewV4().String(),
	}
	errRefresh := refreshRepository.Store(&refToken)
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func CreateToken(userid string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().UTC().Add(time.Minute * 75).Unix()
	td.AccessUuid = uuid.NewV4().String()

	//Remove This Later on
	td.RtExpires = time.Now().UTC().Add(time.Hour * 24 * 3500).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + userid
	td.AccessTokenId = uuid.NewV4().String()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["id"] = td.AccessTokenId
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func GetAllUsers(userRepository *Repositories.UserRepository) (events.APIGatewayProxyResponse, error) {

	allUsers, err := userRepository.List([]string{})
	if err != nil {
		return Shared.ApiErrorResponseHandler("Error occurred while getting all users", 400)
	}

	users := make([]Model.AllUsersResponse, 0)
	for _, user := range *allUsers {
		users = append(users, Model.AllUsersResponse{
			Id:		   user.Id,
			FirstName: user.FirstName,
			Surname:   user.Surname,
			Email:     user.Email,
			Phone:     user.Phone,
			Role:      user.Role,
		})
	}

	return Shared.ApiResponseHandler(users, 200)

}

func DeleteUser(userRepository *Repositories.UserRepository, requestBody *string) (events.APIGatewayProxyResponse, error) {

	deleteUserRequest := &Model.DeleteUserRequest{}

	if err := json.Unmarshal([]byte(*requestBody), deleteUserRequest); err != nil {
		return Shared.ApiErrorResponseHandler("Invalid request structure", 401)
	}

	err := userRepository.Delete(deleteUserRequest.Id, "id")
	if err != nil {
		return Shared.ApiErrorResponseHandler("Error occurred while deleting", 400)
	}

	return Shared.ApiSuccessResponseHandler("Successfully deleted", 200)

}

func EditUser(userRepository *Repositories.UserRepository, requestBody *string) (events.APIGatewayProxyResponse, error) {

	editUserRequest := &Model.EditUserRequest{}

	if err := json.Unmarshal([]byte(*requestBody), editUserRequest); err != nil {
		return Shared.ApiErrorResponseHandler("Invalid request structure", 401)
	}

	// TODO: Any update method in dynamodb?
	user, err := userRepository.Get(editUserRequest.Id, "id")
	if err != nil || user.Id == "" {
		return Shared.ApiErrorResponseHandler("User couldn't found by the system", 400)
	}

	editedUser := Models.User{
		Email:     editUserRequest.Email,
		Password:  user.Password,
		FirstName: editUserRequest.FirstName,
		Surname:   editUserRequest.Surname,
		Phone:     editUserRequest.Phone,
		Sex:       user.Sex,
		CreatedAt: user.CreatedAt,
		Id:        editUserRequest.Id,
		Active:    user.Active,
		Role:      editUserRequest.Role,
	}

	if saveErr := userRepository.Store(&editedUser); saveErr != nil {
		return Shared.ApiErrorResponseHandler("User could not saved", 500)
	}

	return Shared.ApiSuccessResponseHandler("Successfully updated", 200)

}

func ResetPassword(userRepository *Repositories.UserRepository, requestBody *string) (events.APIGatewayProxyResponse, error) {

	resetPasswordRequest := &Model.ResetPasswordRequest{}

	if err := json.Unmarshal([]byte(*requestBody), resetPasswordRequest); err != nil {
		return Shared.ApiErrorResponseHandler("Invalid request structure", 401)
	}

	newPassword := GeneratePassword(8)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
	if err != nil {
		return Shared.ApiErrorResponseHandler("Couldn't hash the password", 400)
	}
	password := string(passwordHash[:])

	// TODO: Any update method in dynamodb?
	user, err := userRepository.Get(resetPasswordRequest.Id, "id")
	if err != nil || user.Id == "" {
		return Shared.ApiErrorResponseHandler("User couldn't found by the system", 400)
	}

	editedUser := Models.User{
		Email:     user.Email,
		Password:  password,
		FirstName: user.FirstName,
		Surname:   user.Surname,
		Phone:     user.Phone,
		Sex:       user.Sex,
		CreatedAt: user.CreatedAt,
		Id:        user.Id,
		Active:    user.Active,
		Role:      user.Role,
	}

	if saveErr := userRepository.Store(&editedUser); saveErr != nil {
		return Shared.ApiErrorResponseHandler("User could not saved", 500)
	}

	messageBody := "<!doctype html>\n<html xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:v=\"urn:schemas-microsoft-com:vml\" xmlns:o=\"urn:schemas-microsoft-com:office:office\">\n\n<head>\n  <!-- NAME: 1 COLUMN -->\n  <!--[if gte mso 15]>\n      <xml>\n        <o:OfficeDocumentSettings>\n          <o:AllowPNG/>\n          <o:PixelsPerInch>96</o:PixelsPerInch>\n        </o:OfficeDocumentSettings>\n      </xml>\n    <![endif]-->\n  <meta charset=\"UTF-8\">\n  <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <title>Reset Your Password</title>\n  <!--[if !mso]>\n      <!-- -->\n  <link href='https://fonts.googleapis.com/css?family=Asap:400,400italic,700,700italic' rel='stylesheet' type='text/css'>\n  <!--<![endif]-->\n  <style type=\"text/css\">\n    @media only screen and (min-width:768px){\n          .templateContainer{\n              width:600px !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          body,table,td,p,a,li,blockquote{\n              -webkit-text-size-adjust:none !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          body{\n              width:100% !important;\n              min-width:100% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          #bodyCell{\n              padding-top:10px !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          .mcnImage{\n              width:100% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n         \n   .mcnCaptionTopContent,.mcnCaptionBottomContent,.mcnTextContentContainer,.mcnBoxedTextContentContainer,.mcnImageGroupContentContainer,.mcnCaptionLeftTextContentContainer,.mcnCaptionRightTextContentContainer,.mcnCaptionLeftImageContentContainer,.mcnCaptionRightImageContentContainer,.mcnImageCardLeftTextContentContainer,.mcnImageCardRightTextContentContainer{\n              max-width:100% !important;\n              width:100% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          .mcnBoxedTextContentContainer{\n              min-width:100% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          .mcnImageGroupContent{\n              padding:9px !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          .mcnCaptionLeftContentOuter\n   .mcnTextContent,.mcnCaptionRightContentOuter .mcnTextContent{\n              padding-top:9px !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          .mcnImageCardTopImageContent,.mcnCaptionBlockInner\n   .mcnCaptionTopContent:last-child .mcnTextContent{\n              padding-top:18px !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          .mcnImageCardBottomImageContent{\n              padding-bottom:9px !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          .mcnImageGroupBlockInner{\n              padding-top:0 !important;\n              padding-bottom:0 !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          .mcnImageGroupBlockOuter{\n              padding-top:9px !important;\n              padding-bottom:9px !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          .mcnTextContent,.mcnBoxedTextContentColumn{\n              padding-right:18px !important;\n              padding-left:18px !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          .mcnImageCardLeftImageContent,.mcnImageCardRightImageContent{\n              padding-right:18px !important;\n              padding-bottom:0 !important;\n              padding-left:18px !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n          .mcpreview-image-uploader{\n              display:none !important;\n              width:100% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n      /*\n      @tab Mobile Styles\n      @section Heading 1\n      @tip Make the first-level headings larger in size for better readability\n   on small screens.\n      */\n          h1{\n              /*@editable*/font-size:20px !important;\n              /*@editable*/line-height:150% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n      /*\n      @tab Mobile Styles\n      @section Heading 2\n      @tip Make the second-level headings larger in size for better\n   readability on small screens.\n      */\n          h2{\n              /*@editable*/font-size:20px !important;\n              /*@editable*/line-height:150% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n      /*\n      @tab Mobile Styles\n      @section Heading 3\n      @tip Make the third-level headings larger in size for better readability\n   on small screens.\n      */\n          h3{\n              /*@editable*/font-size:18px !important;\n              /*@editable*/line-height:150% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n      /*\n      @tab Mobile Styles\n      @section Heading 4\n      @tip Make the fourth-level headings larger in size for better\n   readability on small screens.\n      */\n          h4{\n              /*@editable*/font-size:16px !important;\n              /*@editable*/line-height:150% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n      /*\n      @tab Mobile Styles\n      @section Boxed Text\n      @tip Make the boxed text larger in size for better readability on small\n   screens. We recommend a font size of at least 16px.\n      */\n          .mcnBoxedTextContentContainer\n   .mcnTextContent,.mcnBoxedTextContentContainer .mcnTextContent p{\n              /*@editable*/font-size:16px !important;\n              /*@editable*/line-height:150% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n      /*\n      @tab Mobile Styles\n      @section Preheader Visibility\n      @tip Set the visibility of the email's preheader on small screens. You\n   can hide it to save space.\n      */\n          #templatePreheader{\n              /*@editable*/display:block !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n      /*\n      @tab Mobile Styles\n      @section Preheader Text\n      @tip Make the preheader text larger in size for better readability on\n   small screens.\n      */\n          #templatePreheader .mcnTextContent,#templatePreheader\n   .mcnTextContent p{\n              /*@editable*/font-size:12px !important;\n              /*@editable*/line-height:150% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n      /*\n      @tab Mobile Styles\n      @section Header Text\n      @tip Make the header text larger in size for better readability on small\n   screens.\n      */\n          #templateHeader .mcnTextContent,#templateHeader .mcnTextContent p{\n              /*@editable*/font-size:16px !important;\n              /*@editable*/line-height:150% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n      /*\n      @tab Mobile Styles\n      @section Body Text\n      @tip Make the body text larger in size for better readability on small\n   screens. We recommend a font size of at least 16px.\n      */\n          #templateBody .mcnTextContent,#templateBody .mcnTextContent p{\n              /*@editable*/font-size:16px !important;\n              /*@editable*/line-height:150% !important;\n          }\n  \n  }   @media only screen and (max-width: 480px){\n      /*\n      @tab Mobile Styles\n      @section Footer Text\n      @tip Make the footer content text larger in size for better readability\n   on small screens.\n      */\n          #templateFooter .mcnTextContent,#templateFooter .mcnTextContent p{\n              /*@editable*/font-size:12px !important;\n              /*@editable*/line-height:150% !important;\n          }\n  \n  }\n  </style>\n</head>\n\n<body style=\"-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;\n background-color: #4863A0; height: 100%; margin: 0; padding: 0; width: 100%\">\n  <center>\n    <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" height=\"100%\" id=\"bodyTable\" style=\"border-collapse: collapse; mso-table-lspace: 0;\n mso-table-rspace: 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust:\n 100%; background-color: #4863A0; height: 100%; margin: 0; padding: 0; width:\n 100%\" width=\"100%\">\n      <tr>\n        <td align=\"center\" id=\"bodyCell\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; border-top: 0;\n height: 100%; margin: 0; padding: 0; width: 100%\" valign=\"top\">\n          <!-- BEGIN TEMPLATE // -->\n          <!--[if gte mso 9]>\n              <table align=\"center\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\" width=\"600\" style=\"width:600px;\">\n                <tr>\n                  <td align=\"center\" valign=\"top\" width=\"600\" style=\"width:600px;\">\n                  <![endif]-->\n          <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"templateContainer\" style=\"border-collapse: collapse; mso-table-lspace: 0; mso-table-rspace: 0;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; max-width:\n 600px; border: 0\" width=\"100%\">\n            <tr>\n              <td id=\"templatePreheader\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; background-color: #4863A0;\n border-top: 0; border-bottom: 0; padding-top: 16px; padding-bottom: 8px\" valign=\"top\">\n                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"mcnTextBlock\" style=\"border-collapse: collapse; mso-table-lspace: 0;\n mso-table-rspace: 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;\n min-width:100%;\" width=\"100%\">\n                  <tbody class=\"mcnTextBlockOuter\">\n                    <tr>\n                      <td class=\"mcnTextBlockInner\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%\" valign=\"top\">\n                        <table align=\"left\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"mcnTextContentContainer\" style=\"border-collapse: collapse; mso-table-lspace: 0;\n mso-table-rspace: 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust:\n 100%; min-width:100%;\" width=\"100%\">\n                          <tbody>\n                            <tr>\n                              <td class=\"mcnTextContent\" style='mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; word-break: break-word;\n color: #2a2a2a; font-family: \"Asap\", Helvetica, sans-serif; font-size: 12px;\n line-height: 150%; text-align: left; padding-top:9px; padding-right: 18px;\n padding-bottom: 9px; padding-left: 18px;' valign=\"top\">\n                                \n                                  <img align=\"none\" alt=\"Lingo is the best way to\n organize, share and use all your visual assets in one place - all on your desktop.\" height=\"32\" src=\"https://hospitalonmobile.com/logo.png\" style=\"-ms-interpolation-mode: bicubic; border: 0; outline: none;\n text-decoration: none; height: auto; width: 211px; height: 100px; margin: 0px;\" width=\"107\" />\n                                </a>\n                              </td>\n                            </tr>\n                          </tbody>\n                        </table>\n                      </td>\n                    </tr>\n                  </tbody>\n                </table>\n              </td>\n            </tr>\n            <tr>\n              <td id=\"templateHeader\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; background-color: #f7f7ff;\n border-top: 0; border-bottom: 0; padding-top: 16px; padding-bottom: 0\" valign=\"top\">\n                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"mcnImageBlock\" style=\"border-collapse: collapse; mso-table-lspace: 0; mso-table-rspace: 0;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;\n min-width:100%;\" width=\"100%\">\n                  <tbody class=\"mcnImageBlockOuter\">\n                    <tr>\n                      <td class=\"mcnImageBlockInner\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; padding:0px\" valign=\"top\">\n                        <table align=\"left\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"mcnImageContentContainer\" style=\"border-collapse: collapse; mso-table-lspace: 0;\n mso-table-rspace: 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust:\n 100%; min-width:100%;\" width=\"100%\">\n                          <tbody>\n                            <tr>\n                              <td class=\"mcnImageContent\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; padding-right: 0px;\n padding-left: 0px; padding-top: 0; padding-bottom: 0; text-align:center;\" valign=\"top\">\n                                <a class=\"\" style=\"mso-line-height-rule:\n exactly; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; color:\n #f57153; font-weight: normal; text-decoration: none\" target=\"_blank\" title=\"\">\n                                  <a class=\"\" style=\"mso-line-height-rule:\n exactly; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; color:\n #f57153; font-weight: normal; text-decoration: none\" target=\"_blank\" title=\"\">\n                                    <img align=\"center\" alt=\"Forgot your password?\" class=\"mcnImage\" src=\"https://static.lingoapp.com/assets/images/email/il-password-reset@2x.png\" style=\"-ms-interpolation-mode: bicubic; border: 0; height: auto; outline: none;\n text-decoration: none; vertical-align: bottom; max-width:1200px; padding-bottom:\n 0; display: inline !important; vertical-align: bottom;\" width=\"600\"></img>\n                                  </a>\n                                </a>\n                              </td>\n                            </tr>\n                          </tbody>\n                        </table>\n                      </td>\n                    </tr>\n                  </tbody>\n                </table>\n              </td>\n            </tr>\n            <tr>\n              <td id=\"templateBody\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; background-color: #f7f7ff;\n border-top: 0; border-bottom: 0; padding-top: 0; padding-bottom: 0\" valign=\"top\">\n                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"mcnTextBlock\" style=\"border-collapse: collapse; mso-table-lspace: 0; mso-table-rspace: 0;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; min-width:100%;\" width=\"100%\">\n                  <tbody class=\"mcnTextBlockOuter\">\n                    <tr>\n                      <td class=\"mcnTextBlockInner\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%\" valign=\"top\">\n                        <table align=\"left\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"mcnTextContentContainer\" style=\"border-collapse: collapse; mso-table-lspace: 0;\n mso-table-rspace: 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust:\n 100%; min-width:100%;\" width=\"100%\">\n                          <tbody>\n                            <tr>\n                              <td class=\"mcnTextContent\" style='mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; word-break: break-word;\n color: #2a2a2a; font-family: \"Asap\", Helvetica, sans-serif; font-size: 16px;\n line-height: 150%; text-align: center; padding-top:9px; padding-right: 18px;\n padding-bottom: 9px; padding-left: 18px;' valign=\"top\">\n\n                                <h2 class=\"null\" style='color: #2a2a2a; font-family: \"Asap\", Helvetica,\n sans-serif; font-size: 24px; font-style: normal; font-weight: bold; line-height:\n 125%; letter-spacing: 1px; text-align: center; display: block; margin: 0;\n padding: 0'><span style=\"text-transform:uppercase\">your new password is</span></h2>\n\n                              </td>\n                            </tr>\n                          </tbody>\n                        </table>\n                      </td>\n                    </tr>\n                  </tbody>\n                </table>\n                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"mcnTextBlock\" style=\"border-collapse: collapse; mso-table-lspace: 0; mso-table-rspace:\n 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;\n min-width:100%;\" width=\"100%\">\n                  <tbody class=\"mcnTextBlockOuter\">\n                    <tr>\n                      <td class=\"mcnTextBlockInner\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%\" valign=\"top\">\n                        <table align=\"left\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"mcnTextContentContainer\" style=\"border-collapse: collapse; mso-table-lspace: 0;\n mso-table-rspace: 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust:\n 100%; min-width:100%;\" width=\"100%\">\n                          <tbody>\n                            <tr>\n                              <td class=\"mcnTextContent\" style='mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; word-break: break-word;\n color: #2a2a2a; font-family: \"Asap\", Helvetica, sans-serif; font-size: 16px;\n line-height: 150%; text-align: center; padding-top:9px; padding-right: 18px;\n padding-bottom: 25px; padding-left: 18px;' valign=\"top\">" + newPassword + "\n                                <br></br>\n                              </td>\n                            </tr>\n                          </tbody>\n                        </table>\n                      </td>\n                    </tr>\n                  </tbody>\n                </table>\n                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"mcnImageBlock\" style=\"border-collapse: collapse; mso-table-lspace: 0; mso-table-rspace: 0;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; min-width:100%;\" width=\"100%\">\n                  <tbody class=\"mcnImageBlockOuter\">\n                    <tr>\n                      <td class=\"mcnImageBlockInner\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; padding:0px\" valign=\"top\">\n                        <table align=\"left\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"mcnImageContentContainer\" style=\"border-collapse: collapse; mso-table-lspace: 0;\n mso-table-rspace: 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust:\n 100%; min-width:100%;\" width=\"100%\">\n                          <tbody>\n                            <tr>\n                              <td class=\"mcnImageContent\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; padding-right: 0px;\n padding-left: 0px; padding-top: 0; padding-bottom: 0; text-align:center;\" valign=\"top\"></td>\n                            </tr>\n                          </tbody>\n                        </table>\n                      </td>\n                    </tr>\n                  </tbody>\n                </table>\n              </td>\n            </tr>\n            <tr>\n              <td id=\"templateFooter\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; background-color: #4863A0;\n border-top: 0; border-bottom: 0; padding-top: 8px; padding-bottom: 80px\" valign=\"top\">\n                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"mcnTextBlock\" style=\"border-collapse: collapse; mso-table-lspace: 0; mso-table-rspace: 0;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; min-width:100%;\" width=\"100%\">\n                  <tbody class=\"mcnTextBlockOuter\">\n                    <tr>\n                      <td class=\"mcnTextBlockInner\" style=\"mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%\" valign=\"top\">\n                        <table align=\"center\" bgcolor=\"#F7F7FF\" border=\"0\" cellpadding=\"32\" cellspacing=\"0\" class=\"card\" style=\"border-collapse: collapse; mso-table-lspace: 0;\n mso-table-rspace: 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust:\n 100%; background:#F7F7FF; margin:auto; text-align:left; max-width:600px;\n font-family: 'Asap', Helvetica, sans-serif;\" text-align=\"left\" width=\"100%\">\n                          <tr>\n                            <td style=\"mso-line-height-rule: exactly; -ms-text-size-adjust: 100%;\n -webkit-text-size-adjust: 100%\">\n\n                              <h3 style='color: #2a2a2a; font-family: \"Asap\", Helvetica, sans-serif;\n font-size: 20px; font-style: normal; font-weight: normal; line-height: 125%;\n letter-spacing: normal; text-align: center; display: block; margin: 0; padding:\n 0; text-align: left; width: 100%; font-size: 16px; font-weight: bold; '>What\n is Mona?</h3>\n\n                              <p style='margin: 10px 0; padding: 0; mso-line-height-rule: exactly;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; color: #2a2a2a;\n font-family: \"Asap\", Helvetica, sans-serif; font-size: 12px; line-height: 150%;\n text-align: left; text-align: left; font-size: 14px; '>Mona is a monitorizing and analyzing tool for the backend solutions we develop for our customers. It helps us to log and inspect all the traffic in our systems.\n                              </p>\n                              <div style=\"padding-bottom: 18px;\">\n                                <a href=\"https://hospitalonmobile.com\" style=\"mso-line-height-rule: exactly; -ms-text-size-adjust: 100%;\n -webkit-text-size-adjust: 100%; color: #f57153; font-weight: normal; text-decoration: none;\n font-size: 14px; color:#F57153; text-decoration:none;\" target=\"_blank\" title=\"Learn more about Lingo\">Learn More ❯</a>\n                              </div>\n                            </td>\n                          </tr>\n                        </table>\n                        <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" style=\"border-collapse: collapse; mso-table-lspace: 0; mso-table-rspace: 0;\n -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; min-width:100%;\" width=\"100%\">\n                         \n                      </td>\n                    </tr>\n                  </tbody>\n                </table>\n              </td>\n            </tr>\n          </table>\n          <!--[if gte mso 9]>\n                  </td>\n                </tr>\n              </table>\n            <![endif]-->\n          <!-- // END TEMPLATE -->\n        </td>\n      </tr>\n    </table>\n  </center>\n</body>\n\n</html>"

	var emailArr []string
	emailArr = append(emailArr, user.Email)
	_, Eerr := Shared.SendTestEmail("email-service-master-main-module", messageBody, "Your New Password", emailArr)
	if Eerr != nil {
		return Shared.ApiErrorResponseHandler("Couldn't send mail to the user", 400)
	}

	return Shared.ApiSuccessResponseHandler("Successfully updated", 200)

}

func GeneratePassword(digitCount int) string {

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

	b := make([]byte, digitCount)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// TODO: Need to change it. Do not use authenticated user's id,
// might use the mail of the selected user and ListByEmail function?
func SavePermissionList(userRepository *Repositories.UserRepository, userId string, requestBody *string) (events.APIGatewayProxyResponse, error) {
	savePermissionListRequest := &Model.SavePermissionListRequest{}
	if err := json.Unmarshal([]byte(*requestBody), savePermissionListRequest); err != nil {
		return Shared.ApiErrorResponseHandler("Invalid request structure", 401)
	}
	user, err := userRepository.Get(userId, "id")
	if err != nil || user.Id == "" {
		return Shared.ApiErrorResponseHandler("Couldn't find the user", 400)
	}
	user.PermissionList = &savePermissionListRequest.PermissionList

	if saveErr := userRepository.Store(user); saveErr != nil {
		return Shared.ApiErrorResponseHandler("Couldn't save the user", 500)
	}
	return Shared.ApiResponseHandler("Successfully saved", 200)
}

func GetPermissionList(userRepository *Repositories.UserRepository, userId string) (events.APIGatewayProxyResponse, error) {

	user, err := userRepository.Get(userId, "id")
	if err != nil || user.Id == "" {
		return Shared.ApiErrorResponseHandler("Couldn't find the user", 400)
	}

	if user.Role == "admin" {
		res := &Model.GetPermissionListResponse{
			PermissionList: nil,
			IsAdmin:        true,
		}
		return Shared.ApiResponseHandler(res, 200)
	}
	response := &Model.GetPermissionListResponse{
		PermissionList: *user.PermissionList,
		IsAdmin:        false,
	}
	return Shared.ApiResponseHandler(response, 200)
}
