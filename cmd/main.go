package main

import (
	"context"
	"log"

	"github.com/airheartdev/fixer"
)

func main() {
	fxClient := fixer.NewClient(
		fixer.AccessKey("BQYyyb0xO1b72OX8NivRjTpyq41XhK5z"),
	)

	ctx := context.Background()

	resp, err := fxClient.Latest(ctx,
		fixer.Base(fixer.USD),
		fixer.Symbols(
			fixer.EUR,
			fixer.AUD,
			fixer.NZD,
			fixer.GBP,
			fixer.CAD,
			fixer.MXN,
			fixer.JPY,
			fixer.SGD,
			fixer.IDR,
			fixer.INR,
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", resp)
}
