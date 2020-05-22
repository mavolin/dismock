package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	. "github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// Channels mocks a channels request.
//
// This method will sanitize Channel.ID.
func (m *Mocker) Channels(guildID discord.Snowflake, c []discord.Channel) {
	for i, channel := range c {
		c[i] = sanitize.Channel(channel, 1)
	}

	m.Mock("Channels", http.MethodGet, "/guilds/"+guildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			WriteJSON(t, w, c)
		})
}

// CreateChannel mocks a CreateChannel request.
//
// The GuildID field of the passed Channel must be set.
//
// This method will sanitize Channel.ID.
func (m *Mocker) CreateChannel(d api.CreateChannelData, c discord.Channel) {
	c = sanitize.Channel(c, 1)

	m.Mock("CreateChannel", http.MethodPost, "/guilds/"+c.GuildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			CheckJSON(t, r.Body, new(api.CreateChannelData), &d)
			WriteJSON(t, w, c)
		})
}

// MoveChannel mocks a MoveChannel request.
func (m *Mocker) MoveChannel(guildID discord.Snowflake, d []api.MoveChannelData) {
	m.Mock("CreateChannel", http.MethodPatch, "/guilds/"+guildID.String()+"/channels",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			CheckJSON(t, r.Body, &[]api.MoveChannelData{}, &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// Channel mocks a Channel request.
//
// The ID field of the passed Channel must be set.
//
// This method will sanitize Channel.ID.
func (m *Mocker) Channel(c discord.Channel) {
	c = sanitize.Channel(c, 1)

	m.Mock("CreateChannel", http.MethodGet, "/channels/"+c.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			WriteJSON(t, w, c)
		})
}

// ModifyChannel mocks a ModifyChannel request.
func (m *Mocker) ModifyChannel(id discord.Snowflake, d api.ModifyChannelData) {
	m.Mock("ModifyChannel", http.MethodPatch, "/channels/"+id.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			CheckJSON(t, r.Body, new(api.ModifyChannelData), &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// DeleteChannel mocks a DeleteChannel request.
func (m *Mocker) DeleteChannel(id discord.Snowflake) {
	m.Mock("DeleteChannel", http.MethodDelete, "/channels/"+id.String(), nil)
}

// EditChannelPermission mocks a EditChannelPermission request.
//
// The ID field of the Overwrite must be set.
func (m *Mocker) EditChannelPermission(channelID discord.Snowflake, o discord.Overwrite) {
	m.Mock("EditChannelPermission", http.MethodPut, "/channels/"+channelID.String()+"/permissions/"+o.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			o.ID = 0

			CheckJSON(t, r.Body, new(discord.Overwrite), &o)
			w.WriteHeader(http.StatusNoContent)
		})
}

// DeleteChannelPermission mocks a DeleteChannelPermission request.
func (m *Mocker) DeleteChannelPermission(channelID, overwriteID discord.Snowflake) {
	m.Mock("DeleteChannelPermission", http.MethodDelete,
		"/channels/"+channelID.String()+"/permissions/"+overwriteID.String(), nil)
}

// Typing mocks a Typing request.
func (m *Mocker) Typing(channelID discord.Snowflake) {
	m.Mock("Typing", http.MethodPost, "/channels/"+channelID.String()+"/typing", nil)
}

// PinnedMessages mocks a PinnedMessages request.
//
// This method will sanitize Message.ID, Message.ChannelID and
// Message.Author.ID.
func (m *Mocker) PinnedMessages(channelID discord.Snowflake, messages []discord.Message) {
	for i, message := range messages {
		messages[i] = sanitize.Message(message, 1, channelID, 1)
	}

	m.Mock("PinnedMessages", http.MethodGet, "/channels/"+channelID.String()+"/pins",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			WriteJSON(t, w, messages)
		})
}

// PinMessage mocks a PinMessage request.
func (m *Mocker) PinMessage(channelID, messageID discord.Snowflake) {
	m.Mock("PinMessage", http.MethodPut, "/channels/"+channelID.String()+"/pins/"+messageID.String(), nil)
}

type addRecipientPayload struct {
	AccessToken string `json:"access_token"`
	Nickname    string `json:"nickname"`
}

// AddRecipient mocks a AddRecipient request.
func (m *Mocker) AddRecipient(channelID, userID discord.Snowflake, accessToken, nickname string) {
	m.Mock("PinMessage", http.MethodPut, "/channels/"+channelID.String()+"/recipients/"+userID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := addRecipientPayload{
				AccessToken: accessToken,
				Nickname:    nickname,
			}

			CheckJSON(t, r.Body, new(addRecipientPayload), &expect)

			w.WriteHeader(http.StatusNoContent)
		})
}

// RemoveRecipient mocks a RemoveRecipient request.
func (m *Mocker) RemoveRecipient(channelID, userID discord.Snowflake) {
	m.Mock("RemoveRecipient", http.MethodDelete, "/channels/"+channelID.String()+"/recipients/"+userID.String(), nil)
}

// Ack mocks a Ack request.
func (m *Mocker) Ack(channelID, messageID discord.Snowflake, send, ret api.Ack) {
	m.Mock("Ack", http.MethodPost, "/channels/"+channelID.String()+"/messages/"+messageID.String()+"/ack",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			CheckJSON(t, r.Body, new(api.Ack), &send)
			WriteJSON(t, w, ret)
		})
}
