package gomail

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	cfg "github.com/digisan/go-config"
	lk "github.com/digisan/logkit"
	"github.com/mailgun/mailgun-go/v4"
)

func initMG() string {

	lk.Log("starting... email MG")

	if err := cfg.Init("email", false, cfgMG); err == nil {
		domain = cfg.Val[string]("domain")
		sender = cfg.Val[string]("sender")
		senderEmail = cfg.Val[string]("senderEmail")
		key = translateKey(cfg.Val[string]("apiKey"), []byte(domain))
	}

	lk.FailOnErrWhen(len(senderEmail) == 0, "%v", fmt.Errorf("[senderEmail] is empty, '%s' must be loaded", cfgMG))
	lk.FailOnErrWhen(len(domain) == 0, "%v", fmt.Errorf("[domain] is empty, '%s' must be loaded", cfgMG))
	lk.FailOnErrWhen(len(key) == 0, "%v", fmt.Errorf("[apiKey] is empty, '%s' must be loaded", cfgMG))

	mg = mailgun.NewMailgun(domain, key)

	lk.Log("started... email MG")

	return "mailgun"
}

type resultMG struct {
	recipient string
	msg       string
	id        string
	err       error
}

func (r *resultMG) Recipient() string {
	return r.recipient
}

func (r *resultMG) Err() error {
	return r.err
}

func sendMG(subject, body string, recipients ...string) chan result {
	var (
		chRst = make(chan result)
		nOK   = int32(0)
	)
	for _, recipient := range recipients {
		go func(recipient string) {

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

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
					chRst <- &resultMG{
						recipient: recipient,
						err:       err,
					}
					return
				}
			}

			// The message object allows you to add attachments and Bcc recipients
			message := mg.NewMessage(senderEmail, subject, body, recEmail.(string))

			// Send the message with a 10 second timeout
			if msg, id, err := mg.Send(ctx, message); err != nil {

				lk.Warn("id: %s msg: %s err: %v\n", id, msg, err)
				chRst <- &resultMG{
					recipient: recipient,
					msg:       "",
					id:        "",
					err:       err,
				}

			} else {

				lk.Log("id: %s resp: %s err: %v\n", id, msg, err)
				chRst <- &resultMG{
					recipient: recipient,
					msg:       msg,
					id:        id,
					err:       nil,
				}
				atomic.AddInt32(&nOK, 1)
			}

		}(recipient)
	}
	return chRst
}

