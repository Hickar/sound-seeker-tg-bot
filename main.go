package main

import (
	"flag"
	"log"

	"sound-seeker-bot/internal/bot"
	"sound-seeker-bot/internal/config"
)

var (
	configPath = flag.String("config", "", "Path to the configuration file. Defaults to none")
)

func main() {
	flag.Parse()

	conf, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("unable to parse configuration file at %s: %s", *configPath, err.Error())
	}

	if err := bot.Start(*conf); err != nil {
		log.Fatalln(err)
	}
}