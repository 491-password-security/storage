package Model

type SendEmailRequest struct {
	Emails  []string `json:"emails"`
	Subject string   `json:"subject"`
	Message string   `json:"message"`
}

type SendEmailToServiceRequest struct {
	Service string `json:"service"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type SendSmsRequest struct {
	PhoneNumbers []string `json:"phoneNumbers"`
	Message      string   `json:"message"`
}

type SendSmsToServiceRequest struct {
	Service string `json:"service"`
	Message string `json:"message"`
}

type SendNotificationToServiceRequest struct {
	Service  string  `json:"service"`
	Title    string  `json:"title"`
	Body     string  `json:"body"`
	Url      *string `json:"url"`
	Redirect *string `json:"redirect"`
}