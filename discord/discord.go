package discord

import (
	"sync"
	"github.com/bwmarrin/discordgo"
)

type Client interface {
	UpdateMember(guildID, userID string, roles []string) error
}

type client struct {
	session *discordgo.Session
	mutex   sync.Mutex
}

func (cl *client) UpdateMember(guildID, userID string, roles []string) error {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	return cl.session.GuildMemberEdit(guildID, userID, roles)
}

func NewClient(token string) (Client, error) {
	session, err := discordgo.New("Bot " + token)
	var newClient client
	if err != nil {
		newClient = client{session: session}
	}
	return &newClient, err
}
