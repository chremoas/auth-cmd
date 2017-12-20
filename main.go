package main

import (
	"fmt"
	"github.com/chremoas/auth-cmd/background"
	"github.com/chremoas/auth-cmd/command"
	"github.com/chremoas/auth-cmd/discord"
	uauthsvc "github.com/chremoas/auth-srv/proto"
	proto "github.com/chremoas/chremoas/proto"
	"github.com/chremoas/services-common/config"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"time"
)

var Version = "1.0.0"
var service micro.Service
var name = "auth"

func main() {
	service = config.NewService(Version, "cmd", name, initialize)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func initialize(config *config.Configuration) error {
	chatClient, err := discord.NewClient(config.Bot.BotToken)
	if err != nil {
		message, _ := fmt.Printf("Had an issue starting a discord session: %s", err)
		panic(message)
	}

	roleMap := discord.NewRoleMap(config.Bot.DiscordServerId, chatClient)

	err = roleMap.UpdateRoles()
	if err != nil {
		message, _ := fmt.Printf("Had an issue retrieving the discord roles: %s", err)
		panic(message)
	}

	clientFactory := clientFactory{name: config.LookupService("srv", "auth"), client: service.Client()}
	checker := background.NewChecker(config.Bot.DiscordServerId, chatClient, &clientFactory, roleMap, time.Minute*5)

	checker.Start()

	proto.RegisterCommandHandler(service.Server(),
		command.NewCommand(
			config.Bot.DiscordServerId,
			name,
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
