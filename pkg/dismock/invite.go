package dismock

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/dismock/v2/internal/mockutil"
)

// Invite mocks a Invite request.
//
// The Code field of the passed discord.Invite must be set.
func (m *Mocker) Invite(i discord.Invite) {
	m.MockAPI("Invite", http.MethodGet, "/invites/"+i.Code,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, i)
		})
}

// InviteWithCounts mocks a InviteWithCounts request.
//
// The Code field of the passed discord.Invite must be set.
func (m *Mocker) InviteWithCounts(i discord.Invite) {
	m.MockAPI("Invite", http.MethodGet, "/invites/"+i.Code,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckQuery(t, r.URL.Query(), url.Values{
				"with_counts": {"true"},
			})
			mockutil.WriteJSON(t, w, i)
		})
}

// ChannelInvites mocks a ChannelInvites request.
func (m *Mocker) ChannelInvites(channelID discord.ChannelID, invites []discord.Invite) {
	if invites == nil {
		invites = []discord.Invite{}
	}

	m.MockAPI("ChannelInvites", http.MethodGet, "/channels/"+channelID.String()+"/invites",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, invites)
		})
}

// GuildInvites mocks a GuildInvites request.
func (m *Mocker) GuildInvites(guildID discord.GuildID, invites []discord.Invite) {
	if invites == nil {
		invites = []discord.Invite{}
	}

	m.MockAPI("GuildInvites", http.MethodGet, "/guilds/"+guildID.String()+"/invites",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, invites)
		})
}

// CreateInvite mocks a CreateInvite request.
//
// The Channel.ID field of the passed discord.Invite must be set.
func (m *Mocker) CreateInvite(d api.CreateInviteData, i discord.Invite) {
	m.MockAPI("CreateInvite", http.MethodPost, "/channels/"+i.Channel.ID.String()+"/invites",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.CreateInviteData), &d)
			mockutil.WriteJSON(t, w, i)
		})
}

// DeleteInvite mocks a DeleteInvite request.
//
// The Code field of the passed discord.Invite must be set.
func (m *Mocker) DeleteInvite(i discord.Invite) {
	m.MockAPI("DeleteInvite", http.MethodDelete, "/invites/"+i.Code,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, i)
		})
}
