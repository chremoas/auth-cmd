package main

import (
	"fmt"
	proto "github.com/micro/go-bot/proto"
	"github.com/abaeve/auth-common/config"
	"github.com/abaeve/auth-bot/command"
	"github.com/abaeve/auth-bot/discord"
	"github.com/abaeve/auth-bot/background"
	"github.com/micro/go-micro/client"
	uauthsvc "github.com/abaeve/auth-srv/proto"
	"time"
)

var version string = "1.0.0"
var checker background.Checker

func main() {
	configuration := config.Configuration{}
	// These needs to be a commandline argument eventually
	configuration.Load("application.yaml")

	chatClient, err := discord.NewClient(configuration.Application.BotToken)

	if err != nil {
		message, _ := fmt.Printf("Had an issue starting a discord session: %s", err)
		panic(message)
	}

	service, err := configuration.NewService(version)
	service.Init()

	authSvcName := configuration.Application.AuthSrvNamespace + ".auth-srv"
	roleMap := discord.NewRoleMap(configuration.Application.DiscordServerId, chatClient)

	err = roleMap.UpdateRoles()
	if err != nil {
		message, _ := fmt.Printf("Had an issue retrieving the discord roles: %s", err)
		panic(message)
	}

	clientFactory := clientFactory{name: authSvcName, client: service.Client()}
	checker = background.NewChecker(configuration.Application.DiscordServerId, chatClient, &clientFactory, roleMap, time.Minute*5)
	checker.Start()
	proto.RegisterCommandHandler(service.Server(),
		command.NewCommand(
			configuration.Application.DiscordServerId,
			configuration.Name,
			&clientFactory,
			chatClient,
			roleMap,
		),
	)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

type clientFactory struct {
	name   string
	client client.Client
}

func (c clientFactory) NewClient() uauthsvc.UserAuthenticationClient {
	return uauthsvc.NewUserAuthenticationClient(c.name, c.client)
}
