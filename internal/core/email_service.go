package core

import (
	"fmt"
	"net/smtp"
	"strings"
)

// InvokeSendMail invokes Gmail SMTP configuration to be email sending process.
func (co *Core) InvokeSendMail(from string, to []string, subject, body string) error {
	// Use of HTML Meta tags for MIME context setups.
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	// Prepare the mail body to be sent.
	mailBody := fmt.Sprintf("To: %s\nSubject: %s\n%s\n\n%s\n\n", strings.Join(to, ","), subject, headers, body)

	// Sending email.
	err := smtp.SendMail(
		fmt.Sprintf("%s:587", co.smtpAddr),
		co.smtpAuth,
		from,
		to,
		[]byte(mailBody),
	)

	return err
}
