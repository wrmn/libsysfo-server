package email

type data struct {
	Sender      senderData `json:"sender"`
	To          []ToData   `json:"to"`
	Subject     string     `json:"subject"`
	HtmlContent string     `json:"htmlContent"`
}

type senderData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ToData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Content struct {
	Subject     string
	HtmlContent string
}
