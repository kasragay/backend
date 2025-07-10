package services

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

const telecomCaller = packageCaller + ".Telecom"

type Telecom struct {
	logger  *utils.Logger
	noReply string
	client  *twilio.RestClient
}

func NewTelecomService(logger *utils.Logger) ports.TelecomService {
	noReplyPhone := os.Getenv("NOREPLY_PHONE")
	if noReplyPhone == "" {
		logger.Fatal(context.Background(), "NOREPLY_PHONE is not set")
	}
	if !ports.PhoneValidator(noReplyPhone) {
		logger.Fatal(context.Background(), "NOREPLY_PHONE is not valid")
	}
	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	if twilioAccountSID == "" {
		logger.Fatal(context.Background(), "TWILIO_ACCOUNT_SID is not set")
	}
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if twilioAuthToken == "" {
		logger.Fatal(context.Background(), "TWILIO_AUTH_TOKEN is not set")
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilioAccountSID,
		Password: twilioAuthToken,
	})
	return &Telecom{
		logger:  logger,
		noReply: noReplyPhone,
		client:  client,
	}
}

func (s *Telecom) NoReplySend(ctx context.Context, dst []string, message string) (success []string, err error) {
	defer func() { err = utils.FuncPipe(telecomCaller+".NoReplySend", err) }()
	type resp struct {
		addr string
		err  error
	}
	respChan := make(chan *resp)
	var wg sync.WaitGroup

	for _, d := range dst {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			respChan <- &resp{
				addr: d,
				err:  s.send(ctx, s.noReply, d, message),
			}
		}(d)
	}

	go func() {
		wg.Wait()
		close(respChan)
	}()

	for resp_ := range respChan {
		if resp_.err != nil {
			err = resp_.err
		} else {
			success = append(success, resp_.addr)
		}
	}
	return
}

func (s *Telecom) send(ctx context.Context, src, dst, message string) (err error) {
	defer func() { err = utils.FuncPipe(telecomCaller+".send", err) }()
	// TODO_DEL
	if dst[:11] == "+98920240012" {
		return
	}
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(dst)
	params.SetFrom(src)
	params.SetBody(message)

	result := make(chan error, 1)
	go func() {
		defer close(result)
		_, err := s.client.Api.CreateMessage(params)
		if err != nil {
			result <- fmt.Errorf("failed to send message: %v", err)
			return
		}
		result <- nil
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-result:
		return err
	}
}

const otpMessageHeader string = `Kasragay

This message contains sensitive information.

`

const (
	otpAdminSignupKeyMessage     string = otpMessageHeader + "Admin Signup key: %s"
	otpSignupMessage             string = otpMessageHeader + "Signup code:\n%s"
	otpSigninMessage             string = otpMessageHeader + "Signin code:\n%s"
	otpChangePasswordMessage     string = otpMessageHeader + "Change password code:\n%s"
	otpChangePhoneMessage        string = otpMessageHeader + "Change phone code:\n%s"
	otpDeleteAccountPhoneMessage string = otpMessageHeader + "Delete account code:\n%s"
)

var otpMessages = map[ports.OtpType]string{
	ports.AdminSignupKeyOtpType: otpAdminSignupKeyMessage,
	ports.SignupOtpType:         otpSignupMessage,
	ports.SigninOtpType:         otpSigninMessage,
	ports.ChangePasswordOtpType: otpChangePasswordMessage,
	ports.ChangePhoneOtpType:    otpChangePhoneMessage,
	ports.DeleteAccountOtpType:  otpDeleteAccountPhoneMessage,
}
