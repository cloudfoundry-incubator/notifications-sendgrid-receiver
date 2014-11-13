package services

type RequestParams struct {
	To      string
	Subject string
	From    string
	ReplyTo string
	Text    string
	HTML    string
	Kind    string
}

func NewRequestParamsFromMap(paramsMap map[string]string) RequestParams {
	return RequestParams{
		To:      paramsMap["to"],
		Subject: paramsMap["subject"],
		From:    paramsMap["from"],
		ReplyTo: paramsMap["reply_to"],
		Text:    paramsMap["text"],
		HTML:    paramsMap["html"],
		Kind:    paramsMap["kind"],
	}
}

func (params RequestParams) ToMap() map[string]string {
	return map[string]string{
		"to":       params.To,
		"subject":  params.Subject,
		"from":     params.From,
		"reply_to": params.ReplyTo,
		"text":     params.Text,
		"html":     params.HTML,
		"kind":     params.Kind,
	}
}
