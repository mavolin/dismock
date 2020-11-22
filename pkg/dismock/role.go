package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/dismock/v2/internal/mockutil"
)

// AddRole mocks a AddRole request.
func (m *Mocker) AddRole(guildID discord.GuildID, userID discord.UserID, roleID discord.RoleID) {
	m.MockAPI("AddRole", http.MethodPut,
		"/guilds/"+guildID.String()+"/members/"+userID.String()+"/roles/"+roleID.String(), nil)
}

// RemoveRole mocks a RemoveRole request.
func (m *Mocker) RemoveRole(guildID discord.GuildID, userID discord.UserID, roleID discord.RoleID) {
	m.MockAPI("RemoveRole", http.MethodDelete,
		"/guilds/"+guildID.String()+"/members/"+userID.String()+"/roles/"+roleID.String(), nil)
}

// Roles mocks a Roles request.
func (m *Mocker) Roles(guildID discord.GuildID, roles []discord.Role) {
	m.MockAPI("Roles", http.MethodGet, "/guilds/"+guildID.String()+"/roles",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, roles)
		})
}

// CreateRole mocks a CreateRole request.
func (m *Mocker) CreateRole(guildID discord.GuildID, d api.CreateRoleData, role discord.Role) {
	m.MockAPI("CreateRole", http.MethodPost, "/guilds/"+guildID.String()+"/roles",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.CreateRoleData), &d)
			mockutil.WriteJSON(t, w, role)
		})
}

// MoveRole mocks a MoveRole request.
func (m *Mocker) MoveRole(guildID discord.GuildID, d []api.MoveRoleData, roles []discord.Role) {
	m.MockAPI("MoveRole", http.MethodPatch, "/guilds/"+guildID.String()+"/roles",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &[]api.MoveRoleData{}, &d)
			mockutil.WriteJSON(t, w, roles)
		})
}

// ModifyRole mocks a ModifyRole request.
func (m *Mocker) ModifyRole(guildID discord.GuildID, d api.ModifyRoleData, role discord.Role) {
	m.MockAPI("ModifyRole", http.MethodPatch, "/guilds/"+guildID.String()+"/roles/"+role.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyRoleData), &d)
			mockutil.WriteJSON(t, w, role)
		})
}

// DeleteRole mocks a DeleteRole request.
func (m *Mocker) DeleteRole(guildID discord.GuildID, roleID discord.RoleID) {
	m.MockAPI("DeleteRole", http.MethodDelete, "/guilds/"+guildID.String()+"/roles/"+roleID.String(), nil)
}
