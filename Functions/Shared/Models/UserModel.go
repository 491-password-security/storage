package Models

type User struct {
	Email          string      `json:"email"`
	Password       string      `json:"password"`
	FirstName      string      `json:"firstName"`
	Surname        string      `json:"surname"`
	Phone          string      `json:"phone"`
	Sex            *string     `json:"sex,omitempty"`
	CreatedAt      int64       `json:"createdAt"`
	Id             string      `json:"id"`
	Active         bool        `json:"active,omitempty"`
	Role           string      `json:"role"`
	ViewList       *[][]string `json:"viewList"`
	PermissionList *[]string   `json:"permissionList"`
}

type UserEmails struct {
	Email string `json:"email"`
}

type UserPhoneNumbers struct {
	Phone string `json:"phone"`
}

type UserNotificationTokens struct {
	Token string `json:"token"`
}

type ResponseStatusCodes struct {
	StatusCode *int `json:"http-response-status"`
}