package command

import (
	"golang.org/x/net/context"
	proto "github.com/micro/micro/bot/proto"
	uauthsvc "github.com/abaeve/auth-srv/proto"
	"github.com/abaeve/auth-bot/discord"
	"github.com/micro/go-micro/client"
	"strings"
)

type botError struct {
	message string
}

func (be botError) Error() string {
	return be.message
}

type Command struct {
	guildID string
	name    string
	authSvc uauthsvc.UserAuthenticationClient
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
		rsp.Error = "@" + sender[1] + " I did not understand your command."
		return botError{"Could not understand command"}
	}

	response, err := c.authSvc.Confirm(context.Background(), &uauthsvc.AuthConfirmRequest{UserId: sender[1], AuthenticationCode: req.Args[1]}, nil)

	if err != nil {
		rsp.Error = "@" + sender[1] + " I had an issue authing your request, please reauth or contact your administrator."
		return botError{"Received an error from the auth service: " + err.Error()}
	}

	err = c.client.UpdateMember(c.guildID, sender[1], response.Roles)

	if err != nil {
		rsp.Error = "@" + sender[1] + " I had an issue talking to the chat service, please try again later."
		return botError{"Received (" + err.Error() + ") from chat service."}
	}

	rsp.Result = []byte("@" + sender[1] + " *Success*: " + response.CharacterName + " has been successfully authed")

	return nil
}

func NewCommand(guildID, myName, authSvcName string, microClient client.Client, client discord.Client) *Command {
	newCommand := Command{guildID: guildID, name: myName, authSvc: uauthsvc.NewUserAuthenticationClient(authSvcName, microClient), client: client}
	return &newCommand
}
