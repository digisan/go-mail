package main

import (
	"fmt"

	gm "github.com/digisan/go-mail"
)

func main() {

	ok, sent, failed, errs := gm.SendMail("Fancy subject!", "Hello from Mailgun Go!!!", "cdutwhu@outlook.com", "4987346@qq.com")

	fmt.Println("sent status:", ok)
	if ok {
		fmt.Println("---")
		fmt.Println("sent to:", sent)
		fmt.Println("---")
		fmt.Println("failed on", failed)
		fmt.Println("---")
		fmt.Println("error list:", errs)
	}
}
