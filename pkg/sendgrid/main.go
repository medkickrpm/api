package sendgrid

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var client *sendgrid.Client
var from *mail.Email

func Setup() {
	client = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	from = mail.NewEmail("MedKick Robot", "robot@air.business")
}

func SendEmail(toName string, toEmail string, subject string, body string) error {
	to := mail.NewEmail(toName, toEmail)
	htmlContent := body

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
