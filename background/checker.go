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
	myId     string
}

func (c *checker) Start() {
	c.ticker = time.NewTicker(c.tickTime)
	user, _ := c.client.GetUser("@me")

	c.myId = user.ID

	go func() {
		for range c.ticker.C {
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
		//Don't attempt to remove our own roles
		if c.myId != member.User.ID {
			authClient := c.factory.NewClient()

			response, err := authClient.GetRoles(context.Background(), &uauthsvc.GetRolesRequest{UserId: member.User.ID})

			if err != nil {
				return err
			}

			if len(response.Roles) > 0 {
				var roles []string

				//Map the roles to their id's for quick lookup right below
				roleMap := make(map[string]string)
				for _, role := range response.Roles {
					roleId := c.roleMap.GetRoleId(role)
					roles = append(roles, roleId)
					roleMap[roleId] = role
				}

				//If the user has some roles, it's POSSIBLE that he has tall the one's she should have
				allRolesPresent := len(member.Roles) > 0
				for _, role := range member.Roles {
					if len(roleMap[role]) == 0 {
						//FOUND ONE!
						allRolesPresent = false
					}
				}

				//I want to prematurely cut out any web requests that I can, don't make updates if we're all in sync
				if !allRolesPresent {
					err = c.client.UpdateMember(c.guildID, member.User.ID, roles)

					if err != nil {
						return err
					}
				}
			} else {
				//I want to prematurely cut out any web requests that I can, don't make updates if we're all in sync
				if len(member.Roles) > 0 {
					for _, assignedRole := range member.Roles {
						err = c.client.RemoveMemberRole(c.guildID, member.User.ID, assignedRole)

						if err != nil {
							return err
						}
					}
				}
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
