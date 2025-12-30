package main

import (
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
)

func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	app := api.New(cfg)
	app.Start()
}
