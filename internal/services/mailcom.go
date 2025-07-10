package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"sync"

	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"

	"gopkg.in/gomail.v2"
)

const mailcomCaller = packageCaller + ".Mailcom"

type emailType string

const (
	noReplyEmailType emailType = "no-reply"
)

type Mailcom struct {
	logger    *utils.Logger
	host      string
	port      int
	emails    map[emailType][2]string
	revEmails map[string]emailType
}

func NewMailcomService(logger *utils.Logger) ports.MailcomService {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		logger.Fatal(context.Background(), "SMTP_HOST is not set")
	}
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		logger.Fatal(context.Background(), "SMTP_PORT is not set")
	}
	port_, err := strconv.Atoi(port)
	if err != nil {
		logger.Fatal(context.Background(), "SMTP_PORT is not valid")
	}
	noReplyEmail := os.Getenv("NOREPLY_EMAIL")
	if noReplyEmail == "" {
		logger.Fatal(context.Background(), "NOREPLY_EMAIL is not set")
	}
	if !ports.EmailValidator(noReplyEmail) {
		logger.Fatal(context.Background(), "NOREPLY_EMAIL is not valid")
	}
	noReplyPassword := os.Getenv("NOREPLY_EMAIL_PASSWORD")
	if noReplyPassword == "" {
		logger.Fatal(context.Background(), "NOREPLY_EMAIL_PASSWORD is not set")
	}
	return &Mailcom{
		logger: logger,
		host:   host,
		port:   port_,
		emails: map[emailType][2]string{
			noReplyEmailType: {noReplyEmail, noReplyPassword},
		},
		revEmails: map[string]emailType{
			noReplyEmail: noReplyEmailType,
		},
	}
}

func (s *Mailcom) NoReplySend(ctx context.Context, dst []string, subject, message string) (success []string, err error) {
	defer func() { err = utils.FuncPipe(mailcomCaller+".NoReplySend", err) }()
	type resp struct {
		addr string
		err  error
	}
	respChan := make(chan *resp) // Unbuffered channel
	var wg sync.WaitGroup

	for _, d := range dst {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			respChan <- &resp{
				addr: d,
				err:  s.send(ctx, s.emails[noReplyEmailType], d, subject, message),
			}
		}(d)
	}

	go func() {
		wg.Wait()
		close(respChan)
	}()
	for resp_ := range respChan {
		if resp_.err != nil {
			err = resp_.err // Consider collecting all errors
		} else {
			success = append(success, resp_.addr)
		}
	}
	return
}

func (s *Mailcom) send(ctx context.Context, src [2]string, dst, subject, htmlBody string) (err error) {
	defer func() { err = utils.FuncPipe(mailcomCaller+".send", err) }()
	// TODO_DEL
	if true {
		return nil
	}
	m := gomail.NewMessage()
	m.SetHeader("From", src[0])
	m.SetHeader("To", dst)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)
	d := gomail.NewDialer(s.host, s.port, src[0], src[1])
	result := make(chan error, 1)
	go func() {
		defer close(result)
		if err := d.DialAndSend(m); err != nil {
			result <- fmt.Errorf("failed to send email: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-result:
		return err
	}
}

type emailStruct struct {
	Domain       string
	Version      string
	OtpType_     string
	Token        string
	SupportEmail string
}

func emailBody(otpType ports.OtpType) func(input emailStruct) (string, error) {
	// Parse the template once during initialization.
	t, err := template.ParseFiles("./templates/otp_email.html")
	if err != nil {
		// Panic during initialization if template file is missing or invalid.
		// In a production system, you might want to log this and handle it gracefully.
		panic(fmt.Sprintf("failed to parse template file: %v", err))
	}
	return func(input emailStruct) (string, error) {
		// Validate input fields to prevent rendering issues.
		if input.Domain == "" || input.Version == "" || input.Token == "" || input.SupportEmail == "" {
			return "", fmt.Errorf("missing required emailStruct fields")
		}

		switch otpType {
		case ports.AdminSignupKeyOtpType:
			input.OtpType_ = "Admin Signup Key"
		case ports.SignupOtpType:
			input.OtpType_ = "Signup"
		case ports.SigninOtpType:
			input.OtpType_ = "Signin"
		case ports.ChangePasswordOtpType:
			input.OtpType_ = "Change Password"
		case ports.ChangePhoneOtpType:
			input.OtpType_ = "Change Phone"
		case ports.DeleteAccountOtpType:
			input.OtpType_ = "Delete Account"
		}
		var buf bytes.Buffer
		err = t.Execute(&buf, input)
		if err != nil {
			return "", fmt.Errorf("failed to execute template: %w", err)
		}
		return buf.String(), nil
	}
}

// otpSubjects maps OTP types to email subject lines.
var otpSubjects = map[ports.OtpType]string{
	ports.AdminSignupKeyOtpType: "Admin Signup Key - Kasragay",
	ports.SignupOtpType:         "Signup Otp - Kasragay",
	ports.SigninOtpType:         "Signin Otp - Kasragay",
	ports.ChangePhoneOtpType:    "Change Phone Otp - Kasragay",
	ports.ChangePasswordOtpType: "Change Password Otp - Kasragay",
	ports.DeleteAccountOtpType:  "Delete Account Otp - Kasragay",
}

// otpMails maps OTP types to email body generator functions.
var otpMails = map[ports.OtpType]func(input emailStruct) (string, error){
	ports.AdminSignupKeyOtpType: emailBody(ports.AdminSignupKeyOtpType),
	ports.SignupOtpType:         emailBody(ports.SignupOtpType),
	ports.SigninOtpType:         emailBody(ports.SigninOtpType),
	ports.ChangePhoneOtpType:    emailBody(ports.ChangePhoneOtpType),
	ports.ChangePasswordOtpType: emailBody(ports.ChangePasswordOtpType),
	ports.DeleteAccountOtpType:  emailBody(ports.DeleteAccountOtpType),
}
