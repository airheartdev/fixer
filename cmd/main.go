package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/airheartdev/fixer"
)

func main() {
	fxClient := fixer.NewClient(
		fixer.AccessKey("your-access-key"),
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

	for cur, rate := range resp.Rates {
		fmt.Printf("%s: %f\n", cur, rate)
	}

	fmt.Printf("\nCurrent at: %s\n", time.Time(resp.Date).Format(time.Stamp))
}
