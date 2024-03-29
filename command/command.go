package command

import (
	uauthsvc "github.com/chremoas/auth-srv/proto"
	proto "github.com/chremoas/chremoas/proto"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"strings"
)

//TODO: Refactor this elsewhere... too tired right now and I want to start the checker tests.
type ClientFactory interface {
	NewClient() uauthsvc.UserAuthenticationService
}

type Command struct {
	name    string
	factory ClientFactory
}

var logger *zap.Logger

// Help returns the command usage
func (c *Command) Help(ctx context.Context, req *proto.HelpRequest, rsp *proto.HelpResponse) error {
	// Usage should include the name of the command
	rsp.Usage = c.name
	rsp.Description = "Authenticate your chat user id and link it to the character used to create the given token."
	return nil
}

// Exec executes the command
func (c *Command) Exec(ctx context.Context, req *proto.ExecRequest, rsp *proto.ExecResponse) error {
	sugar := logger.Sugar()
	sender := strings.Split(req.Sender, ":")

	if len(req.Args) != 2 {
		rsp.Result = []byte("<@" + sender[1] + ">, :octagonal_sign: I did not understand your command. Please look in #welcome for instructions.")
		return nil
	}

	client := c.factory.NewClient()

	if len(req.Args) == 1 || req.Args[1] == "help" {
		rsp.Result = []byte("<@" + sender[1] + ">, :no_entry_sign: Please use `!auth <auth_key_from_web>`")
		return nil
	}

	if req.Args[1] == "sync" {
		sugar.Info("Performing Sync")
		_, err := client.SyncToRoleService(ctx, &uauthsvc.SyncRequest{ChannelId: sender[0], UserId: sender[1]})
		if err != nil {
			return err
		}
		sugar.Info("Call to SyncToRolesService completed")

		return nil
	}

	response, err := client.Confirm(ctx, &uauthsvc.AuthConfirmRequest{UserId: sender[1], AuthenticationCode: req.Args[1]})

	if err != nil {
		rsp.Result = []byte("<@" + sender[1] + ">, :octagonal_sign: I had an issue authing your request, please reauth or contact your administrator. " + err.Error())
		return nil
	}

	if len(response.CharacterName) == 0 {
		rsp.Result = []byte("<@" + sender[1] + ">, :no_entry_sign: **Unsure Response**: You have no character.")
		return nil
	}

	rsp.Result = []byte("<@" + sender[1] + ">, :white_check_mark: **Success**: " + response.CharacterName + " has been successfully authed.")
	client.SyncToRoleService(ctx, &uauthsvc.SyncRequest{ChannelId: sender[0], UserId: sender[1], SendMessage: false})

	return nil
}

func NewCommand(myName string, factory ClientFactory, log *zap.Logger) *Command {
	logger = log
	return &Command{name: myName, factory: factory}
}
