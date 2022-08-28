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

func write(str []string) {
	f, err := os.OpenFile("last_check.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _, value := range str {
		fmt.Fprintln(f, value)
	}
}

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

	str := make([]string, 3)
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://atendimentoiupp.zendesk.com//hc/pt-br/sections/4402572414100"),
		chromedp.WaitVisible("#main-content"),
		chromedp.Text("/html/body/main/div[2]/div/section/ul/li[1]", &str[0]),
		chromedp.Text("/html/body/main/div[2]/div/section/ul/li[2]", &str[1]),
		chromedp.Text("/html/body/main/div[2]/div/section/ul/li[3]", &str[2]),
	); err != nil {
		panic(err)
	}

	var check []string
	cmatch := 0

	for i, v := range str {
		b := regexp.MustCompile("azul|smile|latam").MatchString(strings.ToLower(v))
		var stat int
		if b {
			stat = 1
			cmatch++
			str[i] = fmt.Sprint(stat) + ", " + v
			check = append(check, str[i])
		} else {
			stat = 0
			str[i] = fmt.Sprint(stat) + ", " + v
		}
	}

	if cmatch == 0 {
		write(str)
		fmt.Println("No promotion updates, exit")
		os.Exit(0)
	}

	dat, err := os.ReadFile("last_check.txt")
	if err != nil {
		panic(err)
	}
	sData := strings.Split(string(dat), "\n")

	var content []string
	for _, v := range check {
		sAppend := true
		for _, d := range sData {
			if v == d {
				sAppend = false
				break
			}
		}
		if sAppend {
			content = append(content, v)
		}
	}

	write(str)
	if len(content) == 0 {
		fmt.Println("No promotion updates, exit")
		os.Exit(0)
	}

	from := mail.NewEmail("iupp promo notify", os.Getenv("SENDGRID_SENDER_EMAIL"))
	subject := "iupp exchange promotion!"
	to := mail.NewEmail(os.Getenv("SENDGRID_TO_NAME"), os.Getenv("SENDGRID_TO_EMAIL"))
	plainTextContent := ""
	htmlContent := fmt.Sprintln(content)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("New promotions: ", content)
}
