package main

import (
	_ "github.com/mebr0/tiny-url/docs"
	"github.com/mebr0/tiny-url/internal/app"
)

const configPath = "configs/main.yml"

func main() {
	app.Run(configPath)
}
