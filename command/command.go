package command

import (
	"golang.org/x/net/context"
	proto "github.com/micro/micro/bot/proto"
	uauthsvc "github.com/abaeve/auth-srv/proto"
	"github.com/abaeve/auth-bot/discord"
)

type Command struct {
	Name    string
	Service uauthsvc.UserAuthenticationClient
	Client  discord.Client
}

// Help returns the command usage
func (c *Command) Help(ctx context.Context, req *proto.HelpRequest, rsp *proto.HelpResponse) error {
	// Usage should include the name of the command
	rsp.Usage = c.Name
	rsp.Description = "Authenticate your chat user id and link it to the character used to create the given token."
	return nil
}

// Exec executes the command
func (c *Command) Exec(ctx context.Context, req *proto.ExecRequest, rsp *proto.ExecResponse) error {
	return nil
}
