module github.com/chremoas/auth-cmd

go 1.14

require (
	github.com/bwmarrin/discordgo v0.19.0
	github.com/chremoas/auth-srv v1.3.0
	github.com/chremoas/chremoas v1.3.0
	github.com/chremoas/services-common v1.3.2
	github.com/golang/mock v1.3.1
	github.com/micro/go-micro v1.9.1
	go.uber.org/zap v1.10.0
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80
)

replace github.com/chremoas/auth-cmd => ../auth-cmd
replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1