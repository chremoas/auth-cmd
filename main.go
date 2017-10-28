package main

import (
	"fmt"
	"git.maurer-it.net/abaeve/auth-bot/background"
	"git.maurer-it.net/abaeve/auth-bot/command"
	"git.maurer-it.net/abaeve/auth-bot/discord"
	uauthsvc "git.maurer-it.net/abaeve/auth-srv/proto"
	proto "git.maurer-it.net/abaeve/chremoas/proto"
	"git.maurer-it.net/abaeve/services-common/config"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"time"
)

var version string = "1.0.0"
var service micro.Service

func main() {
	service = config.NewService(version, "auth", initialize)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func initialize(configuration *config.Configuration) error {
	chatClient, err := discord.NewClient(configuration.Bot.BotToken)
	if err != nil {
		message, _ := fmt.Printf("Had an issue starting a discord session: %s", err)
		panic(message)
	}

	authSvcName := configuration.Bot.AuthSrvNamespace + "." + configuration.ServiceNames.AuthSrv
	roleMap := discord.NewRoleMap(configuration.Bot.DiscordServerId, chatClient)

	err = roleMap.UpdateRoles()
	if err != nil {
		message, _ := fmt.Printf("Had an issue retrieving the discord roles: %s", err)
		panic(message)
	}

	clientFactory := clientFactory{name: authSvcName, client: service.Client()}
	checker := background.NewChecker(configuration.Bot.DiscordServerId, chatClient, &clientFactory, roleMap, time.Minute*5)

	checker.Start()

	proto.RegisterCommandHandler(service.Server(),
		command.NewCommand(
			configuration.Bot.DiscordServerId,
			configuration.Name,
			&clientFactory,
			chatClient,
			roleMap,
		),
	)

	return nil
}

type clientFactory struct {
	name   string
	client client.Client
}

func (c clientFactory) NewClient() uauthsvc.UserAuthenticationClient {
	return uauthsvc.NewUserAuthenticationClient(c.name, c.client)
}
