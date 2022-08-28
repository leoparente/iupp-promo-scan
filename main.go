package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/chromedp/chromedp"
)

func main() {
	opts := append(
		// select all the elements after the third element
		chromedp.DefaultExecAllocatorOptions[3:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
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
}
