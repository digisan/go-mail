package gomail

import (
	"fmt"
	"testing"

	cfg "github.com/digisan/go-config"
)

// comment out init() if run this test
func TestGenKeyCodeMG(t *testing.T) {
	if err := cfg.Init("email", false, cfgMG); err == nil {
		domain := cfg.Val[string]("domain")
		sender := cfg.Val[string]("sender")
		key := cfg.Val[string]("apiKey")
		fmt.Println(sender)
		// fmt.Println(genCode(key, []byte(domain)))
		fmt.Println(translateKey(key, []byte(domain)))
	}
}

// comment out init() if run this test
func TestGenKeyCodeSG(t *testing.T) {
	if err := cfg.Init("email", false, cfgSG); err == nil {
		name := cfg.Val[string]("sender")
		email := cfg.Val[string]("senderEmail")
		key := cfg.Val[string]("apiKey")
		fmt.Println(name)
		// fmt.Println(genCode(key, []byte(email)))   // original api code => encoded api code
		fmt.Println(translateKey(key, []byte(email))) // encoded api code => original api code
	}
}

func TestSendMail(t *testing.T) {
	ok, sent, failed, errs := SendMail("Fancy subject!", "Hello from Go!", "cdutwhu@outlook.com", "4987346@qq.com")
	fmt.Println("sent status:", ok)
	if ok {
		fmt.Println("sent to:", sent)
		fmt.Println("---")
		fmt.Println("failed on", failed)
		fmt.Println("---")
		fmt.Println("error list:", errs)
	}
}
