package discord

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

type RoleMap interface {
	UpdateRoles() error
	GetRoles() map[string]*discordgo.Role
	GetRoleId(roleName string) string
}

type roleMapImpl struct {
	guildID string
	client  Client
	roles   map[string]*discordgo.Role
	mutex   sync.Mutex
}

func (rm *roleMapImpl) UpdateRoles() error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	roles, err := rm.client.GetAllRoles(rm.guildID)

	if err != nil {
		return err
	}

	rm.roles = make(map[string]*discordgo.Role)

	for _, role := range roles {
		rm.roles[role.Name] = role
	}

	return nil
}

func (rm *roleMapImpl) GetRoles() map[string]*discordgo.Role {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	return rm.roles
}

func (rm *roleMapImpl) GetRoleId(roleName string) string {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	role := rm.roles[roleName]

	if role == nil {
		return ""
	}

	return role.ID
}

func NewRoleMap(guildID string, client Client) RoleMap {
	return &roleMapImpl{guildID: guildID, client: client}
}
