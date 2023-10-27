package gomail

import (
	"fmt"
	"sync/atomic"

	cfg "github.com/digisan/go-config"
	lk "github.com/digisan/logkit"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func initSG() string {

	lk.Log("starting... email SG")

	if err := cfg.Init("email", false, cfgSG...); err == nil {
		domain = cfg.Val[string]("domain")
		sender = cfg.Val[string]("sender")
		senderEmail = cfg.Val[string]("senderEmail")
		key = translateKey(cfg.Val[string]("apiKey"), []byte(senderEmail))
	}

	if len(senderEmail) == 0 || len(key) == 0 {
		lk.Warn("[senderEmail] or [apiKey] is empty, check [%v]", cfg.CurrentCfgFile())
		return ""
	}

	sg = sendgrid.NewSendClient(key)

	lk.Log("started... email SG @ %s", cfg.CurrentCfgFile())

	return "sendgrid"
}

type resultSG struct {
	recipient string
	resp      string
	err       error
}

func (r *resultSG) Recipient() string {
	return r.recipient
}

func (r *resultSG) Err() error {
	return r.err
}

func sendSG(subject, body string, recipients ...string) chan result {
	var (
		chRst = make(chan result)
		nOK   = int32(0)
	)
	for _, recipient := range recipients {
		go func(recipient string) {

			from := mail.NewEmail(sender, senderEmail)

			// if recipient is email, directly use it, otherwise, fetch registered email
			var (
				recEmail any
				ok       bool
			)
			if validEmail(recipient) {
				recEmail = recipient
			} else {
				if recEmail, ok = mRecipient.Load(recipient); !ok {
					err := fmt.Errorf("recipient %v has no email", recipient)
					lk.WarnOnErr("%v", err)
					chRst <- &resultSG{
						recipient: recipient,
						err:       err,
					}
					return
				}
			}
			to := mail.NewEmail(recipient, recEmail.(string))

			message := mail.NewSingleEmail(from, subject, to, body, body)

			if resp, err := sg.Send(message); err != nil {

				lk.Warn("resp: %v err: %v\n", resp, err)
				chRst <- &resultSG{
					recipient: recipient,
					resp:      fmt.Sprintf("%v", resp),
					err:       err,
				}

			} else {

				lk.Log("resp: %v err: %v\n", resp, err)
				chRst <- &resultSG{
					recipient: recipient,
					resp:      fmt.Sprintf("%v", resp),
					err:       nil,
				}
				atomic.AddInt32(&nOK, 1)
			}

		}(recipient)
	}
	return chRst
}
