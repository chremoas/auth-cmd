package main

import (
	"fmt"
	"github.com/chremoas/auth-cmd/command"
	uauthsvc "github.com/chremoas/auth-srv/proto"
	proto "github.com/chremoas/chremoas/proto"
	"github.com/chremoas/services-common/config"
	chremoasPrometheus "github.com/chremoas/services-common/prometheus"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"go.uber.org/zap"
)

var Version = "1.0.0"
var service micro.Service
var name = "auth"
var logger *zap.Logger

func main() {
	service = config.NewService(Version, "cmd", name, initialize)
	var err error

	// TODO pick stuff up from the config
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("Initialized logger")

	go chremoasPrometheus.PrometheusExporter(logger)

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
			logger,
		),
	)

	return nil
}

type clientFactory struct {
	name   string
	client client.Client
}

func (c clientFactory) NewClient() uauthsvc.UserAuthenticationService {
	return uauthsvc.UserAuthenticationServiceClient(c.name, c.client)
}
