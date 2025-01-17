package core

import (
	"fmt"
	"net/smtp"
	"strings"
)

func (co *Core) InvokeSendMail(from string, to []string, subject, body string) error {

	// auth := smtp.PlainAuth(
	// 	"",
	// 	co.mailAddr,
	// 	co.mailSecret,
	// 	co.smtpAddr,
	// )

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	mailBody := fmt.Sprintf("To: %s\nSubject: %s\n%s\n\n%s\n\n", strings.Join(to, ","), subject, headers, body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:587", co.smtpAddr),
		co.smtpAuth,
		from,
		to,
		[]byte(mailBody),
	)

	return err
}
