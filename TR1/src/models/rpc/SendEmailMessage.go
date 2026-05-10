package rm

type SendEmailMessageParams struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
}

type SendEmailMessageResult struct {
	CommonResult
}
