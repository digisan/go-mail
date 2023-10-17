package main

import (
	"fmt"
	"os"

	. "github.com/digisan/go-generics/v2"
	gm "github.com/digisan/go-mail"
)

func main() {

	recipients := []string{"cdutwhu@outlook.com", "cdutwhu@qq.com"}
	for _, r := range os.Args[1:] {
		if IsEmail(r) {
			recipients = append(recipients, r)
		}
	}

	ok, sent, failed, errs := gm.SendMail("Fancy subject!", "Hello from digisan/go-mail Go!!!", recipients...)

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

////////////////////////////////////////////////////////////////////////////

// using SendGrid's Go Library
// https://github.com/sendgrid/sendgrid-go
// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/sendgrid/sendgrid-go"
// 	"github.com/sendgrid/sendgrid-go/helpers/mail"
// )

// func main() {
// 	from := mail.NewEmail("Vhub.Wismed", "wismed.cn@gmail.com")
// 	subject := "Sending with SendGrid is Fun"
// 	to := mail.NewEmail("Qing Miao", "4987346@qq.com")
// 	plainTextContent := "and easy to do anywhere, even with Go"
// 	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
// 	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
// 	client := sendgrid.NewSendClient("SG.*****************************************")
// 	response, err := client.Send(message)
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		fmt.Println(response.StatusCode)
// 		fmt.Println(response.Body)
// 		fmt.Println(response.Headers)
// 	}
// }
