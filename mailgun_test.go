package gomail

import (
	"fmt"
	"testing"
)

func TestMail(t *testing.T) {
	ok, sent, failed, errs := SendMG("Fancy subject!", "Hello from Mailgun Go!", "cdutwhu@outlook.com", "4987346@qq.com")

	fmt.Println("sent status:", ok)
	if ok {
		fmt.Println("sent to:", sent)
		fmt.Println("---")
		fmt.Println("failed on", failed)
		fmt.Println("---")
		fmt.Println("error list:", errs)
	}
}
