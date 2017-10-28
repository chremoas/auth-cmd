package background

import (
	"context"
	"fmt"
	"git.maurer-it.net/abaeve/auth-bot/command"
	"git.maurer-it.net/abaeve/auth-bot/discord"
	uauthsvc "git.maurer-it.net/abaeve/auth-srv/proto"
	"github.com/bwmarrin/discordgo"
	"time"
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
				//map[roleId]roleName
				grantedRoleMap := make(map[string]string)
				for _, role := range response.Roles {
					roleId := c.roleMap.GetRoleId(role)
					roles = append(roles, roleId)
					grantedRoleMap[roleId] = role
				}

				//map[roleId]roleName
				currentRoleMap := make(map[string]string)
				for _, roleId := range member.Roles {
					roleName := c.roleMap.GetRoleName(roleId)
					currentRoleMap[roleId] = roleName
				}

				//Check both ways, whats currently assigned and what SHOULD be assigned.  If ANYTHING is missing, just mass update.
				allRolesPresent := len(member.Roles) > 0
				for roleId := range grantedRoleMap {
					if len(currentRoleMap[roleId]) == 0 {
						allRolesPresent = false
					}
				}

				for _, roleId := range member.Roles {
					if len(grantedRoleMap[roleId]) == 0 {
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
