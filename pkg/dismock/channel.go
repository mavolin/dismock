package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/dismock/v2/internal/mockutil"
)

// Channels mocks a channels request.
func (m *Mocker) Channels(guildID discord.GuildID, c []discord.Channel) {
	m.MockAPI("Channels", http.MethodGet, "/guilds/"+guildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, c)
		})
}

// CreateChannel mocks a CreateChannel request.
//
// The GuildID field of the passed discord.Channel must be set.
func (m *Mocker) CreateChannel(d api.CreateChannelData, c discord.Channel) {
	m.MockAPI("CreateChannel", http.MethodPost, "/guilds/"+c.GuildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.CreateChannelData), &d)
			mockutil.WriteJSON(t, w, c)
		})
}

// MoveChannel mocks a MoveChannel request.
func (m *Mocker) MoveChannel(guildID discord.GuildID, d []api.MoveChannelData) {
	m.MockAPI("CreateChannel", http.MethodPatch, "/guilds/"+guildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &[]api.MoveChannelData{}, &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// Channel mocks a Channel request.
//
// The ID field of the passed discord.Channel must be set.
func (m *Mocker) Channel(c discord.Channel) {
	m.MockAPI("Channel", http.MethodGet, "/channels/"+c.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, c)
		})
}

// ModifyChannel mocks a ModifyChannel request.
func (m *Mocker) ModifyChannel(id discord.ChannelID, d api.ModifyChannelData) {
	m.MockAPI("ModifyChannel", http.MethodPatch, "/channels/"+id.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyChannelData), &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// DeleteChannel mocks a DeleteChannel request.
func (m *Mocker) DeleteChannel(id discord.ChannelID) {
	m.MockAPI("DeleteChannel", http.MethodDelete, "/channels/"+id.String(), nil)
}

// EditChannelPermission mocks a EditChannelPermission request.
func (m *Mocker) EditChannelPermission(
	channelID discord.ChannelID, overwriteID discord.Snowflake, d api.EditChannelPermissionData,
) {
	m.MockAPI("EditChannelPermission", http.MethodPut,
		"/channels/"+channelID.String()+"/permissions/"+overwriteID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.EditChannelPermissionData), &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// DeleteChannelPermission mocks a DeleteChannelPermission request.
func (m *Mocker) DeleteChannelPermission(channelID discord.ChannelID, overwriteID discord.Snowflake) {
	m.MockAPI("DeleteChannelPermission", http.MethodDelete,
		"/channels/"+channelID.String()+"/permissions/"+overwriteID.String(), nil)
}

// Typing mocks a Typing request.
func (m *Mocker) Typing(channelID discord.ChannelID) {
	m.MockAPI("Typing", http.MethodPost, "/channels/"+channelID.String()+"/typing", nil)
}

// PinnedMessages mocks a PinnedMessages request.
func (m *Mocker) PinnedMessages(channelID discord.ChannelID, messages []discord.Message) {
	if messages == nil {
		messages = []discord.Message{}
	}

	m.MockAPI("PinnedMessages", http.MethodGet, "/channels/"+channelID.String()+"/pins",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, messages)
		})
}

// PinMessage mocks a PinMessage request.
func (m *Mocker) PinMessage(channelID discord.ChannelID, messageID discord.MessageID) {
	m.MockAPI("PinMessage", http.MethodPut, "/channels/"+channelID.String()+"/pins/"+messageID.String(), nil)
}

// UnpinMessage mocks a UnpinMessage request.
func (m *Mocker) UnpinMessage(channelID discord.ChannelID, messageID discord.MessageID) {
	m.MockAPI("UnpinMessage", http.MethodDelete, "/channels/"+channelID.String()+"/pins/"+messageID.String(), nil)
}

type addRecipientPayload struct {
	AccessToken string `json:"access_token"`
	Nickname    string `json:"nickname"`
}

// AddRecipient mocks a AddRecipient request.
func (m *Mocker) AddRecipient(channelID discord.ChannelID, userID discord.UserID, accessToken, nickname string) {
	m.MockAPI("PinMessage", http.MethodPut, "/channels/"+channelID.String()+"/recipients/"+userID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := addRecipientPayload{
				AccessToken: accessToken,
				Nickname:    nickname,
			}

			mockutil.CheckJSON(t, r.Body, new(addRecipientPayload), &expect)

			w.WriteHeader(http.StatusNoContent)
		})
}

// RemoveRecipient mocks a RemoveRecipient request.
func (m *Mocker) RemoveRecipient(channelID discord.ChannelID, userID discord.UserID) {
	m.MockAPI("RemoveRecipient", http.MethodDelete, "/channels/"+channelID.String()+"/recipients/"+userID.String(), nil)
}

// Ack mocks a Ack request.
func (m *Mocker) Ack(channelID discord.ChannelID, messageID discord.MessageID, send, ret api.Ack) {
	m.MockAPI("Ack", http.MethodPost, "/channels/"+channelID.String()+"/messages/"+messageID.String()+"/ack",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.Ack), &send)
			mockutil.WriteJSON(t, w, ret)
		})
}
