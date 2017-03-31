package main

import (
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/micro/micro/bot/proto"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/abaeve/auth-bot/command"
	"github.com/abaeve/auth-bot/discord"
	"github.com/abaeve/auth-bot/background"
)

type Configuration struct {
	Application struct {
		BotToken        string `yaml:"botToken"`
		GuildId         string
		Namespace       string
		Name            string
		DiscordServerId string `yaml:"discordServerId"`
	}
}

var configuration Configuration
var checker background.Checker

func main() {
	data, err := ioutil.ReadFile("application.yaml")

	//<editor-fold desc="Configuration Launch Sanity check">
	//TODO: Candidate for shared function for all my services.
	if err != nil {
		panic("Could not read application.yaml for configuration data.")
	}

	err = yaml.Unmarshal([]byte(data), &configuration)

	if err != nil {
		message, _ := fmt.Printf("Parsing application.yaml failed: %s", err)
		panic(message)
	}
	//</editor-fold>

	chatClient, err := discord.NewClient(configuration.Application.BotToken)

	if err != nil {
		message, _ := fmt.Printf("Had an issue starting a discord session: %s", err)
		panic(message)
	}

	service := micro.NewService(
		micro.Name(
			configuration.Application.Namespace +
				"." + configuration.Application.Name,
		),
	)

	service.Init()

	authSvcName := configuration.Application.Namespace + ".auth-srv"
	checker = background.NewChecker(chatClient, authSvcName, service.Client())
	checker.Start()
	proto.RegisterCommandHandler(service.Server(),
		command.NewCommand(
			configuration.Application.GuildId,
			configuration.Application.Name,
			authSvcName,
			service.Client(),
			chatClient,
		),
	)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
