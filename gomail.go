package gomail

import (
	"fmt"
	"sync"
	"time"

	cfg "github.com/digisan/go-config"
	lk "github.com/digisan/logkit"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sendgrid/sendgrid-go"
)

const (
	cfgMG = "mailgun-config.json"
	cfgSG = "sendgrid-config.json"
)

var (
	cfgOK       = true
	sendBy      = ""
	mg          *mailgun.MailgunImpl // mailgun
	sg          *sendgrid.Client     // sendgrid
	sender      = ""                 // sender name, 		 both
	senderEmail = ""                 // sender email, 		 sendgrid
	domain      = ""                 // sender domain, 	     mailgun
	key         = ""                 // api key (encrypted), both
	mRecipient  = sync.Map{}         // both
	timeout     = 12 * time.Second   // both
)

// SendGrid is the first priority
func init() {
	if err := cfg.Init("sendgrid", false, cfgSG); err == nil {
		sendBy = initSG()
		lk.Log("using %v", sendBy)
		return
	}
	if err := cfg.Init("mailgun", false, cfgMG); err == nil {
		sendBy = initMG()
		lk.Log("using %v", sendBy)
		return
	}
	lk.Warn("%v", fmt.Errorf("cannot find any of email agent config [%v, %v], sending email doesn't work", cfgMG, cfgSG))
	cfgOK = false
}

func RegisterRecipient(name, email string) error {
	if cfgOK {
		if validEmail(email) {
			mRecipient.Store(name, email)
			return nil
		}
		return fmt.Errorf("[%v] is invalid email format", email)
	}
	return fmt.Errorf("cannot find any of email agent config [%v, %v], cannot register recipient", cfgMG, cfgSG)
}

type result interface {
	Recipient() string
	Err() error
}

func SendMail(subject, body string, recipients ...string) (OK bool, sent []string, failed []string, errs []error) {

	if cfgOK {

		var (
			chRst chan result
			nOK   = 0
			done  = make(chan bool)
		)

		switch sendBy {
		case "sendgrid":
			chRst = sendSG(subject, body, recipients...)
		case "mailgun":
			chRst = sendMG(subject, body, recipients...)
		default:
			panic("only [mailgun, sendgrid] are supported")
		}

		go func() {
			for rst := range chRst {
				if rst.Err() == nil {
					sent = append(sent, rst.Recipient())
					nOK++
				} else {
					failed = append(failed, rst.Recipient())
					errs = append(errs, rst.Err())
				}
				if nOK == len(recipients) {
					close(chRst)
				}
			}
			done <- true
		}()
		select {
		case <-done:
			return nOK == len(recipients), sent, failed, errs

		case <-time.After(timeout):
			errs = append(errs, fmt.Errorf("timeout @%vs", timeout/time.Second))
			return false, nil, nil, errs
		}

	} else {

		errs = append(errs, fmt.Errorf("cannot find any of email agent config [%v, %v], cannot send email", cfgMG, cfgSG))
		return false, nil, nil, errs
	}
}
