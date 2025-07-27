package core

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

// InvokeSendMail invokes Gmail SMTP configuration to send an email with an optional attachment.
func (co *Core) InvokeSendMailWithAttachment(from string, to []string, subject, body, tailoredResumeID, attachmentPath string) error {
	// Create a new multipart writer
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the email headers
	headers := map[string]string{
		"From":         from,
		"To":           strings.Join(to, ","),
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary()),
	}

	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString("\r\n")

	// Add the email body
	bodyPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	if err != nil {
		return err
	}
	bodyPart.Write([]byte(body))

	// Add the attachment if provided
	if attachmentPath != "" {
		attachmentPart, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(attachmentPath))},
			"Content-Type":              {"application/octet-stream"},
			"Content-Transfer-Encoding": {"base64"},
		})
		if err != nil {
			return err
		}

		attachmentData, err := os.ReadFile(attachmentPath)
		if err != nil {
			return err
		}

		encodedAttachment := base64.StdEncoding.EncodeToString(attachmentData)
		attachmentPart.Write([]byte(encodedAttachment))
	}

	// Close the multipart writer
	writer.Close()

	// Send the email
	err = smtp.SendMail(
		fmt.Sprintf("%s:587", co.opts.SmtpAddr),
		co.smtpAuth,
		from,
		to,
		buf.Bytes(),
	)
	if err != nil {
		return err
	}

	// Store the email into the database
	err = co.DB.CreateEmailInMailbox(from, to, subject, body, tailoredResumeID)
	if err != nil {
		co.Lo.Error("error saving referral email into mailbox", "error", err)
		return err
	}

	return nil
}

// InvokeSendMail invokes Gmail SMTP configuration to be email sending process.
func (co *Core) InvokeSendMail(from string, to []string, subject, body, tailoredResumeID string) error {
	// Use of HTML Meta tags for MIME context setups.
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	// Prepare the mail body to be sent.
	mailBody := fmt.Sprintf("To: %s\nSubject: %s\n%s\n\n%s\n\n", strings.Join(to, ","), subject, headers, body)

	// // Sending email.
	err := smtp.SendMail(
		fmt.Sprintf("%s:587", co.opts.SmtpAddr),
		co.smtpAuth,
		from,
		to,
		[]byte(mailBody),
	)

	if err != nil {
		return err
	}

	// Store the email into database.
	err = co.DB.CreateEmailInMailbox(from, to, subject, body, tailoredResumeID)
	if err != nil {
		co.Lo.Error("error saving referral email into mailbox", "error", err)
		return err
	}

	return err
}
