package command

import (
	"golang.org/x/net/context"
	proto "github.com/micro/go-bot/proto"
	uauthsvc "github.com/abaeve/auth-srv/proto"
	"github.com/abaeve/auth-bot/discord"
	"strings"
)

type botError struct {
	message string
}

func (be botError) Error() string {
	return be.message
}

type ClientFactory interface {
	NewClient() uauthsvc.UserAuthenticationClient
}

type Command struct {
	guildID string
	name    string
	factory ClientFactory
	client  discord.Client
}

// Help returns the command usage
func (c *Command) Help(ctx context.Context, req *proto.HelpRequest, rsp *proto.HelpResponse) error {
	// Usage should include the name of the command
	rsp.Usage = c.name
	rsp.Description = "Authenticate your chat user id and link it to the character used to create the given token."
	return nil
}

// Exec executes the command
func (c *Command) Exec(ctx context.Context, req *proto.ExecRequest, rsp *proto.ExecResponse) error {
	sender := strings.Split(req.Sender, ":")

	if len(req.Args) != 2 {
		rsp.Result = []byte("<@" + sender[1] + "> I did not understand your command.")
		return botError{"Could not understand command"}
	}

	client := c.factory.NewClient()

	response, err := client.Confirm(ctx, &uauthsvc.AuthConfirmRequest{UserId: sender[1], AuthenticationCode: req.Args[1]})

	if err != nil {
		rsp.Result = []byte("<@" + sender[1] + "> I had an issue authing your request, please reauth or contact your administrator.")
		return nil
	}

	if response.Roles == nil || len(response.Roles) == 0 {
		rsp.Result = []byte("<@" + sender[1] + "> *Unsure Response*: You have 0 roles assigned")
		return nil
	}

	if len(response.CharacterName) == 0 {
		rsp.Result = []byte("<@" + sender[1] + "> *Unsure Response*: You have no character")
		return nil
	}

	err = c.client.UpdateMember(c.guildID, sender[1], response.Roles)

	if err != nil {
		rsp.Result = []byte("<@" + sender[1] + "> I had an issue talking to the chat service, please try again later.")
		return nil
	}

	rsp.Result = []byte("<@" + sender[1] + "> *Success*: " + response.CharacterName + " has been successfully authed")

	return nil
}

func NewCommand(guildID, myName string, factory ClientFactory, client discord.Client) *Command {
	newCommand := Command{guildID: guildID, name: myName, factory: factory, client: client}
	return &newCommand
}
