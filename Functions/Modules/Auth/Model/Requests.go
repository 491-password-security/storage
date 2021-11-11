package Model

type SignUpRequest struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName string  `json:"firstName"`
	Surname   string  `json:"surname"`
	Phone     string  `json:"phone"`
	Sex       *string `json:"sex,omitempty"`
	Role      string  `json:"role"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DeleteUserRequest struct {
	Id string `json:"id"`
}

type EditUserRequest struct {
	Id        string  `json:"id"`
	FirstName string  `json:"firstName"`
	Surname   string  `json:"surname"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Role      string  `json:"role"`
}

type ResetPasswordRequest struct {
	Id string `json:"id"`
}

type SavePermissionListRequest struct {
	PermissionList []string `json:"permissionList"`
	Email          string   `json:"email"`
}