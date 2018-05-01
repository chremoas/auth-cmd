package main

import (
	"fmt"
	"github.com/chremoas/auth-cmd/command"
	uauthsvc "github.com/chremoas/auth-srv/proto"
	proto "github.com/chremoas/chremoas/proto"
	"github.com/chremoas/services-common/config"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
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
	clientFactory := clientFactory{name: config.LookupService("srv", "auth"), client: service.Client()}


	proto.RegisterCommandHandler(service.Server(),
		command.NewCommand(
			name,
			&clientFactory,
		),
	)

	return nil
}

type clientFactory struct {
	name   string
	client client.Client
}

func (c clientFactory) NewClient() uauthsvc.UserAuthenticationService {
	return uauthsvc.NewUserAuthenticationService(c.name, c.client)
}
