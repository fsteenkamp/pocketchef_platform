package conf

import (
	"errors"
	"fmt"

	"github.com/ardanlabs/conf/v3"
)

func ParseAndPrint(cfg any) error {
	help, err := conf.Parse("", cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return err
	}

	confStr, err := conf.String(cfg)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("CONFIG")
	fmt.Print(confStr)
	fmt.Println()
	fmt.Println()

	return nil
}

func Parse(cfg any) error {
	help, err := conf.Parse("", cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return err
	}

	return nil
}
