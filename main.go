package main

import (
	"fmt"
	"strings"

	"github.com/micro/go-micro"
	"golang.org/x/net/context"

	proto "github.com/abaeve/micro/bot/proto"
	"github.com/bwmarrin/discordgo"
	"sync"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Command struct{}

type Discord struct {
	session *discordgo.Session
	mutex   sync.Mutex
}

type Configuration struct {
	Application struct {
		BotToken        string `yaml:"botToken"`
		Namespace       string
		Name            string
		DiscordServerId string `yaml:"discordServerId"`
	}
}

var discord Discord
var configuration Configuration

func init() {
	data, err := ioutil.ReadFile("application.yaml")

	//<editor-fold desc="Configuration Launch Sanity check">
	//TODO: Candidate for shared function for all my services.
	if err != nil {
		panic("Could not read application.yaml for configuration data.")
	}

	err = yaml.Unmarshal([]byte(data), &configuration)

	if err != nil {
		message, _ := fmt.Printf("Parsing application.yaml failed: %s", err)
		panic(message)
	}
	//</editor-fold>

	discord.session, err = discordgo.New("Bot " + configuration.Application.BotToken)

	if err != nil {
		message, _ := fmt.Printf("Had an issue starting a discord session: %s", err)
		panic(message)
	}

	members, err := discord.session.GuildMembers(configuration.Application.DiscordServerId, "", 1000)
	if err != nil {
		message, _ := fmt.Printf("Houston we had a problem: %s", err)
		fmt.Println(message)
	}

	for _, member := range members {
		fmt.Print(member.User.ID + " ")
		fmt.Println(member.User.Username)
	}
}

func main() {
	service := micro.NewService(
		micro.Name("com.aba-eve.bot.aba"),
	)

	service.Init()

	proto.RegisterCommandHandler(service.Server(), new(Command))

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// Help returns the command usage
func (c *Command) Help(ctx context.Context, req *proto.HelpRequest, rsp *proto.HelpResponse) error {
	// Usage should include the name of the command
	rsp.Usage = "aba"
	rsp.Description = "This is an example bot command as a micro service which echos the message"
	return nil
}

// Exec executes the command
func (c *Command) Exec(ctx context.Context, req *proto.ExecRequest, rsp *proto.ExecResponse) error {
	rsp.Result = []byte(strings.Join(req.Args, " "))
	fmt.Println(req)
	// rsp.Error could be set to return an error instead
	// the function error would only be used for service level issues
	return nil
}
