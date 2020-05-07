package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/dismock/internal/mockutil"
)

// Emojis mocks a Emojis request.
func (m *Mocker) Emojis(guildID discord.Snowflake, e []discord.Emoji) {
	m.Mock("Emojis", http.MethodGet, "/guilds/"+guildID.String()+"/emojis",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, e)
		})
}

// Emoji mocks a Emoji request.
// The ID field of the passed emoji is required.
func (m *Mocker) Emoji(guildID discord.Snowflake, e discord.Emoji) {
	m.Mock("Emoji", http.MethodGet, "/guilds/"+guildID.String()+"/emojis/"+e.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, e)
		})
}

type createEmojiPayload struct {
	Name  string              `json:"name"`
	Image api.Image           `json:"image"`
	Roles []discord.Snowflake `json:"roles"`
}

// CreateEmoji mocks a CreateEmoji request.
// The fields Name and RoleIDs of the passed Emoji must be set.
func (m *Mocker) CreateEmoji(guildID discord.Snowflake, i api.Image, e discord.Emoji) {
	m.Mock("CreateEmoji", http.MethodPost, "/guilds/"+guildID.String()+"/emojis",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := &createEmojiPayload{
				Name:  e.Name,
				Image: i,
				Roles: e.RoleIDs,
			}

			mockutil.CheckJSONBody(t, r.Body, new(createEmojiPayload), expect)
			mockutil.WriteJSON(t, w, e)
		})
}

type modifyEmojiPayload struct {
	Name  string              `json:"name,omitempty"`
	Roles []discord.Snowflake `json:"roles,omitempty"`
}

// ModifyEmoji mocks a ModifyEmoji request.
func (m *Mocker) ModifyEmoji(guildID, emojiID discord.Snowflake, name string, roles []discord.Snowflake) {
	m.Mock("ModifyEmoji", http.MethodPatch, "/guilds/"+guildID.String()+"/emojis/"+emojiID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := &modifyEmojiPayload{
				Name:  name,
				Roles: roles,
			}

			mockutil.CheckJSONBody(t, r.Body, new(modifyEmojiPayload), expect)
			w.WriteHeader(http.StatusNoContent)
		})
}

// DeleteEmoji mocks a DeleteEmoji request.
func (m *Mocker) DeleteEmoji(guildID, emojiID discord.Snowflake) {
	m.Mock("DeleteEmoji", http.MethodDelete, "/guilds/"+guildID.String()+"/emojis/"+emojiID.String(), nil)
}
