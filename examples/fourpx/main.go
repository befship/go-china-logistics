// Example: 4PX estimated cost. Requires FPX_APP_KEY and FPX_APP_SECRET.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/befship/go-china-logistics/fourpx"
)

func main() {
	key := os.Getenv("FPX_APP_KEY")
	secret := os.Getenv("FPX_APP_SECRET")
	if key == "" || secret == "" {
		fmt.Fprintln(os.Stderr, "set FPX_APP_KEY and FPX_APP_SECRET")
		os.Exit(1)
	}

	client := fourpx.New(fourpx.Config{
		AppKey:      key,
		AppSecret:   secret,
		Environment: fourpx.EnvTest,
	})

	items, err := client.EstimatedCost(context.Background(), &fourpx.EstimatedCostRequest{
		CountryCode: "US",
		Weight:      "500",
		Length:      "10",
		Width:       "10",
		Height:      "10",
		CargoCode:   fourpx.CargoCodeParcel,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, it := range items {
		fmt.Printf("%s fee=%.2f CNY eta=%s\n", it.LogisticsProductCode, it.LumpSumFee, it.EstimatedTime)
	}
}
