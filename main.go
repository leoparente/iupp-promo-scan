package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {
	opts := append(
		// select all the elements after the third element
		chromedp.DefaultExecAllocatorOptions[3:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.DisableGPU,
		chromedp.Flag("remote-debugging-port", "9222"),
	)
	// create chromedp's context
	parentCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(parentCtx)
	defer cancel()

	var str [3]string
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://atendimentoiupp.zendesk.com//hc/pt-br/sections/4402572414100"),
		chromedp.WaitVisible("#main-content"),
		chromedp.Text("/html/body/main/div[2]/div/section/ul/li[1]", &str[0]),
		chromedp.Text("/html/body/main/div[2]/div/section/ul/li[2]", &str[1]),
		chromedp.Text("/html/body/main/div[2]/div/section/ul/li[3]", &str[2]),
	); err != nil {
		panic(err)
	}

	for _, value := range str {
		b := regexp.MustCompile("azul|smile|latam").MatchString(strings.ToLower(value))
		var stat int
		if b {
			stat = 1
		} else {
			stat = 0
		}
		fmt.Printf("%d, %v\n", stat, value)
	}

	from := mail.NewEmail("iupp promo notify", os.Getenv("SENDGRID_SENDER_EMAIL"))
	subject := "iupp exchange promotion!"
	to := mail.NewEmail(os.Getenv("SENDGRID_TO_NAME"), os.Getenv("SENDGRID_TO_EMAIL"))
	plainTextContent := str[0]
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
	}
}
