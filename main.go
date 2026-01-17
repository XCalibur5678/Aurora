package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	cmd := &cli.Command{
		Name:  "Aurora",
		Usage: "A simple yet effective AUR helper",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("Hello friend!")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
