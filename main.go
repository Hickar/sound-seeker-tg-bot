package main

import (
	"flag"
	"log"

	"gopkg.in/tucnak/telebot.v3"
	"sound-seeker-bot/internal/bot"
	"sound-seeker-bot/internal/config"
)

var (
	configPath = flag.String("config", "", "Path to the configuration file. Defaults to none")
)

func main() {
	flag.Parse()

	if *configPath == "" {
		log.Fatalln("no \"configuration filepath was provided\"")
	}

	conf, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("unable to parse configuration file at %s: %s", *configPath, err.Error())
	}

	bot.NewBot(conf)
}