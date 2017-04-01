package discord

import (
	"sync"
	"github.com/bwmarrin/discordgo"
)

type Client interface {
	UpdateMember(guildID, userID string, roles []string) error
	GetAllMembers(guildID, after string, limit int) ([]*discordgo.Member, error)
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

func (cl *client) GetAllMembers(guildID, after string, limit int) ([]*discordgo.Member, error) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	return cl.session.GuildMembers(guildID, after, limit)
}

func NewClient(token string) (Client, error) {
	session, err := discordgo.New("Bot " + token)
	var newClient client
	if err != nil {
		return nil, err
	}
	newClient = client{session: session}
	return &newClient, nil
}
