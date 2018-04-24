package command

import (
	"bytes"
	"fmt"
	uauthsvc "github.com/chremoas/auth-srv/proto"
	proto "github.com/chremoas/chremoas/proto"
	"golang.org/x/net/context"
	"strings"
)

//TODO: Refactor this elsewhere... too tired right now and I want to start the checker tests.
type ClientFactory interface {
	NewClient() uauthsvc.UserAuthenticationClient
}

type Command struct {
	name    string
	factory ClientFactory
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
		rsp.Result = []byte("<@" + sender[1] + ">, :octagonal_sign: I did not understand your command.")
		return nil
	}

	client := c.factory.NewClient()

	if req.Args[1] == "sync" {
		synced, err := client.SyncToRoleService(ctx, &uauthsvc.NilRequest{})
		if err != nil {
			return err
		}

		if len(synced.Roles) == 0 {
			rsp.Result = []byte("```Nothing to sync```\n")
			return nil
		}

		var buffer bytes.Buffer

		buffer.WriteString("Synced:\n")
		for s := range synced.Roles {
			buffer.WriteString(fmt.Sprintf("\t%s: %s\n",
				synced.Roles[s].Name,
				synced.Roles[s].Description,
			))
		}

		rsp.Result = []byte(fmt.Sprintf("```%s```\n", buffer.String()))
		return nil
	}

	response, err := client.Confirm(ctx, &uauthsvc.AuthConfirmRequest{UserId: sender[1], AuthenticationCode: req.Args[1]})

	if err != nil {
		rsp.Result = []byte("<@" + sender[1] + ">, :octagonal_sign: I had an issue authing your request, please reauth or contact your administrator.")
		return nil
	}

	if len(response.CharacterName) == 0 {
		rsp.Result = []byte("<@" + sender[1] + ">, :no_entry_sign: **Unsure Response**: You have no character.")
		return nil
	}

	rsp.Result = []byte("<@" + sender[1] + ">, :white_check_mark: **Success**: " + response.CharacterName + " has been successfully authed.")
	client.SyncToRoleService(ctx, &uauthsvc.NilRequest{})

	return nil
}

func NewCommand(myName string, factory ClientFactory) *Command {
	newCommand := Command{name: myName, factory: factory}
	return &newCommand
}
