package background

import (
	"github.com/abaeve/auth-bot/discord"
	"github.com/bwmarrin/discordgo"
	"time"
	"github.com/abaeve/auth-bot/command"
	"fmt"
	"context"
	uauthsvc "github.com/abaeve/auth-srv/proto"
)

type Checker interface {
	Start()
	Poll() error
	Update(members []*discordgo.Member) error
	Stop()
}

type checker struct {
	guildID  string
	client   discord.Client
	factory  command.ClientFactory
	roleMap  discord.RoleMap
	tickTime time.Duration
	ticker   *time.Ticker
}

func (c *checker) Start() {
	c.ticker = time.NewTicker(c.tickTime)
	go func() {
		for t := range c.ticker.C {
			fmt.Println("Tick at", t)
			err := c.Poll()
			if err != nil {
				//TODO: Replace with logger object
				fmt.Printf("Received an error while polling: %s\n", err)
			}
		}
	}()
}

func (c *checker) Poll() error {
	err := c.roleMap.UpdateRoles()

	if err != nil {
		return err
	}

	members, err := c.client.GetAllMembers(c.guildID, "", 1000)

	if err != nil {
		return err
	}

	err = c.Update(members)

	if err != nil {
		return err
	}

	return nil
}

func (c *checker) Update(members []*discordgo.Member) error {
	for _, member := range members {
		authClient := c.factory.NewClient()

		response, err := authClient.GetRoles(context.Background(), &uauthsvc.GetRolesRequest{UserId: member.User.ID})

		if err != nil {
			return err
		}

		if len(response.Roles) > 0 {
			var roles []string

			for _, role := range response.Roles {
				roleId := c.roleMap.GetRoleId(role)
				roles = append(roles, roleId)
			}

			err = c.client.UpdateMember(c.guildID, member.User.ID, roles)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *checker) Stop() {
	c.ticker.Stop()
}

func NewChecker(guildID string, client discord.Client, factory command.ClientFactory, roleMap discord.RoleMap, tickTime time.Duration) Checker {
	newChecker := checker{guildID: guildID, client: client, factory: factory, roleMap: roleMap, tickTime: tickTime}
	return &newChecker
}
