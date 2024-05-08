package main

import (
	"context"
	"fmt"

	"github.com/rexcfnghk/pricing-store/application"
)

func main() {
	app := application.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}
