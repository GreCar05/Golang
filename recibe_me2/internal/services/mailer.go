package services


import (
"github.com/mailjet/mailjet-apiv3-go"
)

type Recipient struct {
	Name  string
	Email string
}

type Info struct {
	FromRecipient   Recipient
	ToRecipient     Recipient
	ApiKeyPublic    string
	ApiKeyPrivate   string
	Code            string
}

func (i *Info)SendVerificationCode()(*mailjet.ResultsV31, error){
	mailjetClient := mailjet.NewMailjetClient(
		i.ApiKeyPublic,
		i.ApiKeyPrivate,
	)

	var from = mailjet.RecipientV31{
		Email: i.FromRecipient.Email,
		Name:  i.FromRecipient.Name,
	}

	var to = mailjet.RecipientsV31{
		mailjet.RecipientV31{
			Name:  i.ToRecipient.Name,
			Email: i.ToRecipient.Email,
		},
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &from,
			To:   &to,
			Variables: map[string]interface{}{"code": i.Code, "name":i.ToRecipient.Name},
			TemplateLanguage: true,
			Subject: "Verification code: {{var:code}}",
			TextPart: "Hola {{var:name}}, bienvenido a recibeme.cl",
			HTMLPart: "<h3>Hola {{var:name}}, bienvenido a  <a href=\"https://www.mailjet.com/\">Recibeme</a>!</h3><br /> Su codigo de verificacion es: {{var:code}}",
		},
	}

	messages  := mailjet.MessagesV31{Info: messagesInfo}
	resp, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
