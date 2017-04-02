package discord

import (
	"sync"
	"github.com/bwmarrin/discordgo"
)

// This is a very thin wrapper around the discordgo api for testability purposes
// and to easily keep track of what endpoints are being consumed
type Client interface {
	UpdateMember(guildID, userID string, roles []string) error
	GetAllMembers(guildID, after string, limit int) ([]*discordgo.Member, error)
	GetAllRoles(guildID string) ([]*discordgo.Role, error)
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

func (cl *client) GetAllRoles(guildID string) ([]*discordgo.Role, error) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	return cl.session.GuildRoles(guildID)
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
