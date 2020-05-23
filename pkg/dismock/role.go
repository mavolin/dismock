package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	. "github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// AddRole mocks a AddRole request.
func (m *Mocker) AddRole(guildID, userID, roleID discord.Snowflake) {
	m.Mock("AddRole", http.MethodPut,
		"/guilds/"+guildID.String()+"/members/"+userID.String()+"/roles/"+roleID.String(), nil)
}

// RemoveRole mocks a RemoveRole request.
func (m *Mocker) RemoveRole(guildID, userID, roleID discord.Snowflake) {
	m.Mock("RemoveRole", http.MethodDelete,
		"/guilds/"+guildID.String()+"/members/"+userID.String()+"/roles/"+roleID.String(), nil)
}

// Roles mocks a Roles request.
//
// This method will sanitize Roles.ID.
func (m *Mocker) Roles(guildID discord.Snowflake, roles []discord.Role) {
	for i, r := range roles {
		roles[i] = sanitize.Role(r, 1)
	}

	m.Mock("Roles", http.MethodGet, "/guilds/"+guildID.String()+"/roles",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			WriteJSON(t, w, roles)
		})
}

// CreateRole mocks a CreateRole request.
//
// This method will sanitize Role.ID.
func (m *Mocker) CreateRole(guildID discord.Snowflake, d api.CreateRoleData, role discord.Role) {
	role = sanitize.Role(role, 1)

	m.Mock("CreateRole", http.MethodPost, "/guilds/"+guildID.String()+"/roles",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			CheckJSON(t, r.Body, new(api.CreateRoleData), &d)
			WriteJSON(t, w, role)
		})
}

// MoveRole mocks a MoveRole request.
//
// This method will sanitize Roles.ID.
func (m *Mocker) MoveRole(guildID discord.Snowflake, d []api.MoveRoleData, roles []discord.Role) {
	for i, r := range roles {
		roles[i] = sanitize.Role(r, 1)
	}

	m.Mock("MoveRole", http.MethodPatch, "/guilds/"+guildID.String()+"/roles",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			CheckJSON(t, r.Body, &[]api.MoveRoleData{}, &d)
			WriteJSON(t, w, roles)
		})
}

// ModifyRole mocks a ModifyRole request.
//
// This method will sanitize Role.ID.
func (m *Mocker) ModifyRole(guildID discord.Snowflake, d api.ModifyRoleData, role discord.Role) {
	role = sanitize.Role(role, 1)

	m.Mock("ModifyRole", http.MethodPatch, "/guilds/"+guildID.String()+"/roles/"+role.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			CheckJSON(t, r.Body, new(api.ModifyRoleData), &d)
			WriteJSON(t, w, role)
		})
}

// DeleteRole mocks a DeleteRole request.
func (m *Mocker) DeleteRole(guildID, roleID discord.Snowflake) {
	m.Mock("DeleteRole", http.MethodDelete, "/guilds/"+guildID.String()+"/roles/"+roleID.String(), nil)
}
