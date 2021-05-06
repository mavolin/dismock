package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/dismock/v2/internal/mockutil"
)

// Emojis mocks a Emojis request.
func (m *Mocker) Emojis(guildID discord.GuildID, e []discord.Emoji) {
	if e == nil {
		e = []discord.Emoji{}
	}

	m.MockAPI("Emojis", http.MethodGet, "/guilds/"+guildID.String()+"/emojis",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, e)
		})
}

// Emoji mocks a Emoji request.
//
// The ID field of the passed discord.Emoji is required.
func (m *Mocker) Emoji(guildID discord.GuildID, e discord.Emoji) {
	m.MockAPI("Emoji", http.MethodGet, "/guilds/"+guildID.String()+"/emojis/"+e.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, e)
		})
}

// CreateEmoji mocks a CreateEmoji request.
//
// The fields Name and RoleIDs of the passed discord.Emoji must be set.
//
// This method will sanitize Emoji.ID and Emoji.User.ID.
func (m *Mocker) CreateEmoji(guildID discord.GuildID, d api.CreateEmojiData, e discord.Emoji) {
	m.MockAPI("CreateEmoji", http.MethodPost, "/guilds/"+guildID.String()+"/emojis",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &d)
			mockutil.WriteJSON(t, w, e)
		})
}

// ModifyEmoji mocks a ModifyEmoji request.
func (m *Mocker) ModifyEmoji(guildID discord.GuildID, emojiID discord.EmojiID, d api.ModifyEmojiData) {
	m.MockAPI("ModifyEmoji", http.MethodPatch, "/guilds/"+guildID.String()+"/emojis/"+emojiID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// DeleteEmoji mocks a DeleteEmoji request.
func (m *Mocker) DeleteEmoji(guildID discord.GuildID, emojiID discord.EmojiID) {
	m.MockAPI("DeleteEmoji", http.MethodDelete, "/guilds/"+guildID.String()+"/emojis/"+emojiID.String(), nil)
}
