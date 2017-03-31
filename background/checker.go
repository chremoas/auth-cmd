package background

import (
	"github.com/abaeve/auth-bot/discord"
	"github.com/bwmarrin/discordgo"
	uauthsvc "github.com/abaeve/auth-srv/proto"
	"github.com/micro/go-micro/client"
)

type Checker interface {
	Start() error
	Poll() error
	Update(members []*discordgo.Member) error
	Stop() error
}

type checker struct {
	client  discord.Client
	authSvc uauthsvc.UserAuthenticationClient
}

func (c checker) Start() error {
	return nil
}

func (c checker) Poll() error {
	return nil
}

func (c checker) Update(members []*discordgo.Member) error {
	return nil
}

func (c checker) Stop() error {
	return nil
}

func NewChecker(client discord.Client, serviceName string, microClient client.Client) Checker {
	newChecker := checker{client: client, authSvc: uauthsvc.NewUserAuthenticationClient(serviceName, microClient)}
	return newChecker
}
