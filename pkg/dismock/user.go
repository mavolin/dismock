package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/dismock/v2/internal/check"
)

// User mocks a User request.
//
// The ID field of the passed User must be set.
func (m *Mocker) User(u discord.User) {
	m.MockAPI("User", http.MethodGet, "/users/"+u.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			check.WriteJSON(t, w, u)
		})
}

// Me mocks a Me request.
func (m *Mocker) Me(u discord.User) {
	m.MockAPI("Me", http.MethodGet, "/users/@me",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			check.WriteJSON(t, w, u)
		})
}

// ModifyMe mocks a ModifyMe request.
func (m *Mocker) ModifyMe(d api.ModifySelfData, u discord.User) {
	m.MockAPI("ModifyMe", http.MethodPatch, "/users/@me",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			check.JSON(t, r.Body, &d)
			check.WriteJSON(t, w, u)
		})
}

type changeOwnNicknamePayload struct {
	Nick string `json:"nick"`
}

// ChangeOwnNickname mocks a ChangeOwnNickname request.
func (m *Mocker) ChangeOwnNickname(guildID discord.GuildID, nick string) {
	m.MockAPI("ChangeOwnNickname", http.MethodPatch, "/guilds/"+guildID.String()+"/members/@me/nick",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := changeOwnNicknamePayload{
				Nick: nick,
			}

			check.JSON(t, r.Body, &expect)
			w.WriteHeader(http.StatusNoContent)
		})
}

// PrivateChannels mocks a PrivateChannels request.
func (m *Mocker) PrivateChannels(c []discord.Channel) {
	m.MockAPI("PrivateChannels", http.MethodGet, "/users/@me/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			check.WriteJSON(t, w, c)
		})
}

type createPrivateChannelPayload struct {
	RecipientID discord.UserID `json:"recipient_id"`
}

// CreatePrivateChannel mocks a CreatePrivateChannel request.
//
// The c.DMRecipients[0] field of the passed discord.Channel must be set.
func (m *Mocker) CreatePrivateChannel(c discord.Channel) {
	m.MockAPI("CreatePrivateChannel", http.MethodPost, "/users/@me/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := createPrivateChannelPayload{
				RecipientID: c.DMRecipients[0].ID,
			}

			check.JSON(t, r.Body, &expect)
			check.WriteJSON(t, w, c)
		})
}

// UserConnections mocks a UserConnections request.
func (m *Mocker) UserConnections(c []discord.Connection) {
	m.MockAPI("UserConnections", http.MethodGet, "/users/@me/connections",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			check.WriteJSON(t, w, c)
		})
}
