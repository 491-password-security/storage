package Models

type Jwt struct {
	TokenUuid string `json:"tokenUuid"`
	UserId    string `json:"userId"`
	Expires   int64  `json:"expires"`
	Id        string `json:"id"`
}

type RefreshToken struct {
	UserID       string `json:"userId"`
	RefreshToken string `json:"refreshToken"`
	CreatedAt    int64  `json:"createdAt"`
	UpdatedAt    int64  `json:"updatedAt"`
	ID           string `json:"id"`
}

type SecureToken struct {
	UserID     string          `json:"userId"`
	Token      SecureTokenType `json:"token"`
	Type       string          `json:"type"`
	ExpiryDate int64           `json:"expiryDate"`
	ID         string          `json:"id"`
}

type UserHeaderValues struct {
	UserId       string         `json:"userId"`
	HeaderValues []HeaderValues `json:"headerValues"`
}

type HeaderValues struct {
	Timestamp       int64  `json:"timestamp"`
	Language        string `json:"language"`
	Timezone        string `json:"timezone"`
	OperatingSystem string `json:"operatingSystem"`
	Version         string `json:"version"`
	Build           string `json:"build"`
	Model           string `json:"model"`
	DeviceDetails   string `json:"deviceDetails"`
}

type SecureTokenType string
