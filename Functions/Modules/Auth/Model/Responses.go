package Model

type AuthenticationResponse struct {
	AccessToken  string `json:"accessToken"`
	TokenType 	 string `json:"tokenType"`
	RefreshToken string `json:"refreshToken"`
	Role  		 string	`json:"role"`
}
type AccessRefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	TokenType 	 string `json:"tokenType"`
}
type AllUsersResponse struct {
	Id        string  `json:"id"`
	FirstName string  `json:"firstName"`
	Surname   string  `json:"surname"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Role      string  `json:"role"`
}
type GetPermissionListResponse struct {
	PermissionList []string `json:"permissionList"`
	IsAdmin        bool     `json:"isAdmin"`
}