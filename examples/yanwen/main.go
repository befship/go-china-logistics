// Example: Yanwen calc.list. Requires YANWEN_USER_ID and YANWEN_API_TOKEN.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/befship/go-china-logistics/yanwen"
)

func main() {
	uid := os.Getenv("YANWEN_USER_ID")
	token := os.Getenv("YANWEN_API_TOKEN")
	if uid == "" || token == "" {
		fmt.Fprintln(os.Stderr, "set YANWEN_USER_ID and YANWEN_API_TOKEN")
		os.Exit(1)
	}

	client := yanwen.New(yanwen.Config{
		UserID:      uid,
		APIToken:    token,
		Environment: yanwen.EnvSandbox,
	})

	items, err := client.CalcList(context.Background(), &yanwen.CalcListRequest{
		CityID:    "01",
		CountryID: "US",
		Weight:    500,
		Length:    10,
		Width:     10,
		High:      10,
		PostCode:  "90001",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, it := range items {
		fmt.Printf("%s %s total=%.2f %s\n", it.ProductNumber, it.ProductName, it.TotalMoney, it.Currency)
	}
}
