package worker

import (
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"strings"
	"tilank/utils/logger"
)

const (
	envGmailAccount   = "GMAIL_TILANK"
	envGmailCCAccount = "GMAIL_CC_HSSE"
	envGmailPassword  = "GMAIL_PASSWORD_TILANK"
	configSMTPHost    = "smtp.gmail.com"
	configSMTPPort    = 587
)

type MailInfo struct {
	ViolID        string
	TruckIdentity string
	ToEmail       string
}

var mailInfoCh = make(chan *MailInfo, 30)

func init() {
	go func() {
		for info := range mailInfoCh {
			sendEmailGmail(info.ViolID, info.TruckIdentity, info.ToEmail)
		}
	}()
	logger.Info("email worker dijalankan")
}

func RegSendEmail(data *MailInfo) {
	mailInfoCh <- data
}

func sendEmailGmail(violID string, violIdentity string, toEmail string) {
	email := strings.TrimSpace(os.Getenv(envGmailAccount))
	emailCC := strings.TrimSpace(os.Getenv(envGmailCCAccount))
	password := strings.TrimSpace(os.Getenv(envGmailPassword))

	if email == "" || password == "" {
		logger.Error("konfigurasi email salah", errors.New("environment variable not set"))
		return
	}
	senderName := fmt.Sprintf("PT. Pelabuhan Indonesia III TPKB <%s>", email)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", senderName)
	mailer.SetHeader("To", toEmail)
	mailer.SetAddressHeader("Cc", emailCC, "HSSE TPKB")
	mailer.SetHeader("Subject", "Pemberitahuan ETI TPKB")
	mailer.SetBody("text/html", fmt.Sprintf("Pemberitahuan, truck anda dengan Nomor Lambung <b>%s</b> telah melakukan pelanggaran di area TPKB. terlampir surat Elektronik tilang TPKB", violIdentity))
	mailer.Attach(fmt.Sprintf("static/pdf/%s.pdf", violID))

	dialer := gomail.NewDialer(
		configSMTPHost,
		configSMTPPort,
		email,
		password,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		logger.Error("email gagal dikirim", err)
	}
	logger.Info(fmt.Sprintf("email dikirim ke %s", toEmail))
}
