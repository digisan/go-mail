package main

import (
	"fmt"

	gm "github.com/digisan/go-mail"
)

func main() {
	ok, sent, failed, errs := gm.SendMG("Fancy subject!", "Hello from Mailgun Go!!!", "cdutwhu@outlook.com", "4987346@qq.com")
	fmt.Println(ok)
	fmt.Println("---")
	fmt.Println(sent)
	fmt.Println("---")
	fmt.Println(failed)
	fmt.Println("---")
	fmt.Println(errs)
}
