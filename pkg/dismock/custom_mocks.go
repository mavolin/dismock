package dismock

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/dismock/v3/internal/check"
)

// Error simulates an error response for the given path using the given method.
func (m *Mocker) Error(method, path string, e httputil.HTTPError) {
	m.MockAPI("Error", method, path, func(w http.ResponseWriter, r *http.Request, t *testing.T) {
		w.WriteHeader(e.Status)
		err := json.NewEncoder(w).Encode(e)
		require.NoError(t, err)
	})
}

// =============================================================================
// channel.go
// =====================================================================================

// Ack mocks api.Client.Ack.
func (m *Mocker) Ack(channelID discord.ChannelID, messageID discord.MessageID, send, ret api.Ack) {
	m.MockAPI("Ack", http.MethodPost, "channels/"+channelID.String()+"/messages/"+messageID.String()+"/ack",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			check.JSON(t, send, r.Body)
			check.WriteJSON(t, w, ret)
		})
}

// =============================================================================
// guild.go
// =====================================================================================

const maxFetchGuilds = 100

// Guilds mocks api.Client.Guilds.
func (m *Mocker) Guilds(limit uint, g []discord.Guild) {
	if g == nil {
		g = []discord.Guild{}
	}

	if len(g) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent guilds (%d vs. %d)", len(g), limit))
	}

	var after discord.GuildID

	for i := 0; i <= len(g)/maxFetchGuilds; i++ {
		var (
			from = uint(i) * maxFetchGuilds
			to   = uint(math.Min(float64(from+maxFetchGuilds), float64(len(g))))

			fetch = to - from // we use this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or mexFetchGuild, depending on which is smaller, instead.
			if fetch < maxFetchGuilds {
				fetch = uint(math.Min(float64(limit), float64(maxFetchGuilds)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchGuilds
			fetch = maxFetchGuilds
		}

		m.guildsRange(0, after, fmt.Sprintf("Guilds #%d", i+1), fetch, g[from:to])

		if to-from < maxFetchGuilds {
			break
		}

		after = g[to-1].ID
	}
}

// GuildsBefore mocks api.Client.GuildsBefore.
func (m *Mocker) GuildsBefore(before discord.GuildID, limit uint, g []discord.Guild) {
	if g == nil {
		g = []discord.Guild{}
	}

	if len(g) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent guilds (%d vs. %d)", len(g), limit))
	}

	req := len(g)/maxFetchGuilds + 1

	from := uint(math.Min(float64(uint(req)*maxFetchGuilds), float64(len(g))))

	for i := req; i > 0; i-- {
		no := req - i + 1

		to := from
		from = uint(math.Max(float64(0), float64(int(to-maxFetchGuilds))))

		fetch := to - from // we use this as the sent limit

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or mexFetchGuild, depending on which is smaller, instead.
			if fetch < maxFetchGuilds {
				fetch = uint(math.Min(float64(limit), float64(maxFetchGuilds)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchGuilds
			fetch = maxFetchGuilds
		}

		m.guildsRange(before, 0, fmt.Sprintf("GuildsBefore #%d", no), fetch, g[from:to])

		if to-from < maxFetchGuilds {
			break
		}

		before = g[from].ID
	}
}

// GuildsAfter mocks api.Client.GuildsAfter.
func (m *Mocker) GuildsAfter(after discord.GuildID, limit uint, g []discord.Guild) {
	if g == nil {
		g = []discord.Guild{}
	}

	if len(g) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent guilds (%d vs. %d)", len(g), limit))
	}

	for i := 0; i <= len(g)/maxFetchGuilds; i++ {
		var (
			from = uint(i) * maxFetchGuilds
			to   = uint(math.Min(float64(from+maxFetchGuilds), float64(len(g))))

			fetch = to - from // we use this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or maxFetchGuilds, depending on which is smaller, instead.
			if fetch < maxFetchGuilds {
				fetch = uint(math.Min(float64(limit), float64(maxFetchGuilds)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchGuilds
			fetch = maxFetchGuilds
		}

		m.guildsRange(0, after, fmt.Sprintf("GuildsAfter #%d", i+1), fetch, g[from:to])

		if to-from < maxFetchGuilds {
			break
		}

		after = g[to-1].ID
	}
}

// guildsRange mocks a single request to the GET /guilds endpoint.
func (m *Mocker) guildsRange(before, after discord.GuildID, name string, limit uint, g []discord.Guild) {
	m.MockAPI(name, http.MethodGet, "users/@me/guilds",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"limit": {strconv.FormatUint(uint64(limit), 10)},
			}

			if after != 0 {
				expect["after"] = []string{after.String()}
			}

			if before != 0 {
				expect["before"] = []string{before.String()}
			}

			check.Query(t, expect, r.URL.Query())
			check.WriteJSON(t, w, g)
		})
}

// =============================================================================
// interaction.go
// =====================================================================================

// RespondInteraction mocks api.Client.RespondInteraction.
func (m *Mocker) RespondInteraction(id discord.InteractionID, token string, resp api.InteractionResponse) {
	m.MockAPI("RespondInteraction", http.MethodPost, "interactions/"+id.String()+"/"+token+"/callback",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			if resp.NeedsMultipart() {
				files := resp.Data.Files
				resp.Data.Files = nil

				check.Multipart(t, r.Body, r.Header, resp, files)
			} else {
				check.JSON(t, resp, r.Body)
			}
		})
}

// =============================================================================
// members.go
// =====================================================================================

const maxFetchMembers = 1000

// Members mocks aoi.Client.Members.
func (m *Mocker) Members(guildID discord.GuildID, limit uint, members []discord.Member) {
	if members == nil {
		members = []discord.Member{}
	}

	if len(members) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent Members (%d vs. %d)", len(members), limit))
	}

	var after discord.UserID

	for i := 0; i <= len(members)/maxFetchMembers; i++ {
		var (
			from = uint(i) * maxFetchMembers
			to   = uint(math.Min(float64(from+maxFetchMembers), float64(len(members))))

			fetch = to - from // we use this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or maxFetchMembers, depending on which is smaller, instead.
			if fetch < maxFetchMembers {
				fetch = uint(math.Min(float64(limit), float64(maxFetchMembers)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchMembers
			fetch = maxFetchMembers
		}

		m.membersAfter(guildID, after, fmt.Sprintf("Members #%d", i+1), fetch, members[from:to])

		if to-from < maxFetchMembers {
			break
		}

		after = members[to-1].User.ID
	}
}

// MembersAfter mocks api.Client.MembersAfter.
func (m *Mocker) MembersAfter(guildID discord.GuildID, after discord.UserID, limit uint, members []discord.Member) {
	if members == nil {
		members = []discord.Member{}
	}

	if len(members) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent Members (%d vs. %d)", len(members), limit))
	}

	for i := 0; i <= len(members)/maxFetchMembers; i++ {
		var (
			from = uint(i) * maxFetchMembers
			to   = uint(math.Min(float64(from+maxFetchMembers), float64(len(members))))

			fetch = to - from // we use this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or maxFetchMembers, depending on which is smaller, instead.
			if fetch < maxFetchMembers {
				fetch = uint(math.Min(float64(limit), float64(maxFetchMembers)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchMembers
			fetch = maxFetchMembers
		}

		m.membersAfter(guildID, after, fmt.Sprintf("MembersAfter #%d", i+1), fetch, members[from:to])

		if to-from < maxFetchMembers {
			break
		}

		after = members[to-1].User.ID
	}
}

// membersAfter mocks a single request to the GET /Members endpoint.
func (m *Mocker) membersAfter(
	guildID discord.GuildID, after discord.UserID, name string, limit uint, g []discord.Member,
) {
	m.MockAPI(name, http.MethodGet, "guilds/"+guildID.String()+"/members",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"limit": {strconv.FormatUint(uint64(limit), 10)},
			}

			if after != 0 {
				expect["after"] = []string{after.String()}
			}

			check.Query(t, expect, r.URL.Query())
			check.WriteJSON(t, w, g)
		})
}

// =============================================================================
// message.go
// =====================================================================================

const maxFetchMessages = 100

// Messages mocks a Messages request.
func (m *Mocker) Messages(channelID discord.ChannelID, limit uint, messages []discord.Message) {
	if messages == nil {
		messages = []discord.Message{}
	}

	if len(messages) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent messages (%d vs. %d)", len(messages), limit))
	}

	var before discord.MessageID = 0

	for i := 0; i <= len(messages)/maxFetchMessages; i++ {
		var (
			from = uint(i) * maxFetchMessages
			to   = uint(math.Min(float64(from+maxFetchMessages), float64(len(messages))))
		)

		fetch := to - from // we use this as the sent limit

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or maxFetchMessages, depending on which is smaller, instead.
			if fetch < maxFetchMessages {
				fetch = uint(math.Min(float64(limit), float64(maxFetchMessages)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchMessages
			fetch = maxFetchMessages
		}

		m.messagesRange(channelID, before, 0, 0, fmt.Sprintf("MessagesBefore #%d", i+1), fetch, messages[from:to])

		if to-from < maxFetchMessages {
			break
		}

		before = messages[to-1].ID
	}
}

// MessagesAround mocks a MessagesAround request.
func (m *Mocker) MessagesAround(
	channelID discord.ChannelID, around discord.MessageID, limit uint, messages []discord.Message,
) {
	switch {
	case limit == 0:
		limit = 50
	case limit > 100:
		limit = 100
	}

	if messages == nil {
		messages = []discord.Message{}
	}

	if len(messages) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent messages (%d vs. %d)", len(messages), limit))
	}

	m.messagesRange(channelID, 0, 0, around, "MessagesAround", limit, messages)
}

// MessagesBefore mocks a MessagesBefore request.
//
// This method will sanitize Message.ID, Message.ChannelID and
// Message.Author.ID.
func (m *Mocker) MessagesBefore(
	channelID discord.ChannelID, before discord.MessageID, limit uint, messages []discord.Message,
) {
	if messages == nil {
		messages = []discord.Message{}
	}

	if len(messages) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent messages (%d vs. %d)", len(messages), limit))
	}

	for i := 0; i <= len(messages)/maxFetchMessages; i++ {
		var (
			from = uint(i) * maxFetchMessages
			to   = uint(math.Min(float64(from+maxFetchMessages), float64(len(messages))))
		)

		fetch := to - from // we use this as the sent limit

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or maxFetchMessages, depending on which is smaller, instead.
			if fetch < maxFetchMessages {
				fetch = uint(math.Min(float64(limit), float64(maxFetchMessages)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchMessages
			fetch = maxFetchMessages
		}

		m.messagesRange(channelID, before, 0, 0, fmt.Sprintf("MessagesBefore #%d", i+1), fetch, messages[from:to])

		if to-from < maxFetchMessages {
			break
		}

		before = messages[to-1].ID
	}
}

// MessagesAfter mocks a MessagesAfter request.
func (m *Mocker) MessagesAfter(
	channelID discord.ChannelID, after discord.MessageID, limit uint, messages []discord.Message,
) {
	if after == 0 {
		after = 1
	}

	if messages == nil {
		messages = []discord.Message{}
	}

	if len(messages) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent messages (%d vs. %d)", len(messages), limit))
	}

	for i := 0; i <= len(messages)/maxFetchMessages; i++ {
		var (
			to   = len(messages) - i*maxFetchMessages
			from = int(math.Max(float64(to-maxFetchMessages), float64(0)))

			fetch = from - to // we use this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or maxFetchMessages, depending on which is smaller, instead.
			if fetch < maxFetchMessages {
				fetch = int(math.Min(float64(limit), float64(maxFetchMessages)))
			}

			limit -= uint(fetch)
		} else {
			// this means there is no limit, hence we should use
			// maxFetchMessages
			fetch = maxFetchMessages
		}

		m.messagesRange(channelID, 0, after, 0, fmt.Sprintf("MessagesAfter #%d", i+1), uint(fetch), messages[from:to])

		if to-from < maxFetchMessages {
			break
		}

		after = messages[from].ID
	}
}

// messagesRange mocks a single request to the GET /messages endpoint.
func (m *Mocker) messagesRange(
	channelID discord.ChannelID, before, after, around discord.MessageID, name string, limit uint,
	messages []discord.Message,
) {
	m.MockAPI(name, http.MethodGet, "channels/"+channelID.String()+"/messages",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"limit": {strconv.FormatUint(uint64(limit), 10)},
			}

			if after != 0 {
				expect["after"] = []string{after.String()}
			}

			if before != 0 {
				expect["before"] = []string{before.String()}
			}

			if around != 0 {
				expect["around"] = []string{around.String()}
			}

			check.Query(t, expect, r.URL.Query())
			check.WriteJSON(t, w, messages)
		})
}

// SendTextReply mocks api.Client.SendTextReply.
func (m *Mocker) SendTextReply(msg discord.Message) {
	m.sendMessageComplex("SendTextReply", api.SendMessageData{
		Content:   msg.Content,
		Reference: &discord.MessageReference{MessageID: msg.Reference.MessageID},
	}, msg)
}

// SendEmbeds mocks api.Client.SendEmbeds.
func (m *Mocker) SendEmbeds(msg discord.Message) {
	m.sendMessageComplex("SendEmbeds", api.SendMessageData{
		Embeds: msg.Embeds,
	}, msg)
}

// SendEmbedReply mocks api.Client.SendEmbedReply.
func (m *Mocker) SendEmbedReply(msg discord.Message) {
	m.sendMessageComplex("SendEmbedReply", api.SendMessageData{
		Embeds:    msg.Embeds,
		Reference: &discord.MessageReference{MessageID: msg.Reference.MessageID},
	}, msg)
}

// SendMessage mocks api.Client.SendMessage.
func (m *Mocker) SendMessage(msg discord.Message) {
	d := api.SendMessageData{
		Content: msg.Content,
		Embeds:  msg.Embeds,
	}

	m.sendMessageComplex("SendMessage", d, msg)
}

// SendMessageReply mocks api.Client.SendMessageReply.
func (m *Mocker) SendMessageReply(msg discord.Message) {
	d := api.SendMessageData{
		Content:   msg.Content,
		Reference: &discord.MessageReference{MessageID: msg.Reference.MessageID},
		Embeds:    msg.Embeds,
	}

	m.sendMessageComplex("SendMessageReply", d, msg)
}

// EditText mocks api.Client.EditText.
func (m *Mocker) EditText(msg discord.Message) {
	m.editMessageComplex("EditText", api.EditMessageData{
		Content: option.NewNullableString(msg.Content),
	}, msg)
}

// EditEmbeds mocks api.Client.EditEmbeds.
func (m *Mocker) EditEmbeds(msg discord.Message) {
	m.editMessageComplex("EditEmbeds", api.EditMessageData{
		Embeds: &msg.Embeds,
	}, msg)
}

// EditMessage mocks api.Client.EditMessage.
func (m *Mocker) EditMessage(content string, embeds []discord.Embed, msg discord.Message) {
	var data api.EditMessageData

	if len(content) > 0 {
		data.Content = option.NewNullableString(content)
	}

	if len(embeds) > 0 {
		data.Embeds = &embeds
	}

	m.editMessageComplex("EditMessage", data, msg)
}

// EditMessageComplex mocks api.Client.EditMessageComplex.
func (m *Mocker) EditMessageComplex(d api.EditMessageData, msg discord.Message) {
	m.editMessageComplex("EditMessageComplex", d, msg)
}

// editMessageComplex mocks api.Client.EditMessageComplex.
func (m *Mocker) editMessageComplex(name string, d api.EditMessageData, msg discord.Message) {
	if d.Embeds != nil {
		for i, embed := range *d.Embeds {
			if embed.Type == "" {
				(*d.Embeds)[i].Type = discord.NormalEmbed
			}

			if embed.Color == 0 {
				(*d.Embeds)[i].Color = discord.DefaultEmbedColor
			}
		}
	}

	m.MockAPI(name, http.MethodPatch, "channels/"+msg.ChannelID.String()+"/messages/"+msg.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			check.JSON(t, d, r.Body)
			check.WriteJSON(t, w, msg)
		})
}

// DeleteMessages mocks api.Client.DeleteMessages.
func (m *Mocker) DeleteMessages(channelID discord.ChannelID, messageIDs []discord.MessageID) {
	m.MockAPI("DeleteMessages", http.MethodPost, "channels/"+channelID.String()+"/messages/bulk-delete",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := struct {
				Messages []discord.MessageID `json:"messages"`
			}{Messages: messageIDs}

			check.JSON(t, expect, r.Body)
			w.WriteHeader(http.StatusNoContent)
		})
}

// =============================================================================
// message_reaction.go
// =====================================================================================

// Unreact mocks api.Client.Unreact.
func (m *Mocker) Unreact(channelID discord.ChannelID, messageID discord.MessageID, e discord.APIEmoji) {
	m.deleteUserReaction("Unreact", channelID, messageID, 0, e)
}

// Reactions mocks api.Client.Reactions.
func (m *Mocker) Reactions(
	channelID discord.ChannelID, messageID discord.MessageID, limit uint, e discord.APIEmoji, u []discord.User,
) {
	if u == nil {
		u = []discord.User{}
	}

	if len(u) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent users (%d vs. %d)", len(u), limit))
	}

	var after discord.UserID

	for i := 0; i <= len(u)/api.MaxMessageReactionFetchLimit; i++ {
		var (
			from = uint(i) * api.MaxMessageReactionFetchLimit
			to   = uint(math.Min(float64(from+api.MaxMessageReactionFetchLimit), float64(len(u))))

			fetch = to - from // we use this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or api.MaxMessageReactionFetchLimit, depending on which is smaller, instead.
			if fetch < api.MaxMessageReactionFetchLimit {
				fetch = uint(math.Min(float64(limit), float64(api.MaxMessageReactionFetchLimit)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// api.MaxMessageReactionFetchLimit
			fetch = api.MaxMessageReactionFetchLimit
		}

		m.reactionsRange(channelID, messageID, 0, after, fmt.Sprintf("Reactions #%d", i+1), fetch, e, u[from:to])

		if to-from < api.MaxMessageReactionFetchLimit {
			break
		}

		after = u[to-1].ID
	}
}

// ReactionsBefore mocks api.Client.ReactionsBefore.
func (m *Mocker) ReactionsBefore(
	channelID discord.ChannelID, messageID discord.MessageID, before discord.UserID, limit uint, e discord.APIEmoji,
	u []discord.User,
) {
	if u == nil {
		u = []discord.User{}
	}

	if len(u) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent users (%d vs. %d)", len(u), limit))
	}

	req := len(u)/api.MaxMessageReactionFetchLimit + 1

	from := uint(math.Min(float64(uint(req)*api.MaxMessageReactionFetchLimit), float64(len(u))))

	for i := req; i > 0; i-- {
		no := req - i + 1

		to := from
		from = uint(math.Max(float64(0), float64(int(to-api.MaxMessageReactionFetchLimit))))

		fetch := to - from // we use this as the sent limit

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or api.MaxMessageReactionFetchLimit, depending on which is smaller, instead.
			if fetch < api.MaxMessageReactionFetchLimit {
				fetch = uint(math.Min(float64(limit), float64(api.MaxMessageReactionFetchLimit)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// api.MaxMessageReactionFetchLimit
			fetch = api.MaxMessageReactionFetchLimit
		}

		m.reactionsRange(channelID, messageID, before, 0, fmt.Sprintf("ReactionsBefore #%d", no), fetch, e, u[from:to])

		if to-from < api.MaxMessageReactionFetchLimit {
			break
		}

		before = u[from].ID
	}
}

// ReactionsAfter mocks api.Client.ReactionsAfter.
func (m *Mocker) ReactionsAfter(
	channelID discord.ChannelID, messageID discord.MessageID, after discord.UserID, limit uint, e discord.APIEmoji,
	u []discord.User,
) {
	if u == nil {
		u = []discord.User{}
	}

	if len(u) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent users (%d vs. %d)", len(u), limit))
	}

	for i := 0; i <= len(u)/api.MaxMessageReactionFetchLimit; i++ {
		var (
			from = uint(i) * api.MaxMessageReactionFetchLimit
			to   = uint(math.Min(float64(from+api.MaxMessageReactionFetchLimit), float64(len(u))))

			fetch = to - from // we use this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or api.MaxMessageReactionFetchLimit, depending on which is smaller, instead.
			if fetch < api.MaxMessageReactionFetchLimit {
				fetch = uint(math.Min(float64(limit), float64(api.MaxMessageReactionFetchLimit)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// api.MaxMessageReactionFetchLimit
			fetch = api.MaxMessageReactionFetchLimit
		}

		m.reactionsRange(channelID, messageID, 0, after, fmt.Sprintf("ReactionsAfter #%d", i+1), fetch, e, u[from:to])

		if to-from < api.MaxMessageReactionFetchLimit {
			break
		}

		after = u[to-1].ID
	}
}

// reactionsRange mocks a single request to the GET /reactions endpoint.
func (m *Mocker) reactionsRange(
	channelID discord.ChannelID, messageID discord.MessageID, before, after discord.UserID, name string, limit uint,
	e discord.APIEmoji, u []discord.User,
) {
	m.MockAPI(name, http.MethodGet,
		"channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+e.PathString(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"limit": {strconv.FormatUint(uint64(limit), 10)},
			}

			if after != 0 {
				expect["after"] = []string{after.String()}
			}

			if before != 0 {
				expect["before"] = []string{before.String()}
			}

			check.Query(t, expect, r.URL.Query())
			check.WriteJSON(t, w, u)
		})
}

// DeleteUserReaction mocks api.Client.DeleteUserReaction.
func (m *Mocker) DeleteUserReaction(
	channelID discord.ChannelID, messageID discord.MessageID, userID discord.UserID, e discord.APIEmoji,
) {
	m.deleteUserReaction("DeleteUserReaction", channelID, messageID, userID, e)
}

func (m *Mocker) deleteUserReaction(
	name string, channelID discord.ChannelID, messageID discord.MessageID, userID discord.UserID, e discord.APIEmoji,
) {
	user := "@me"
	if userID > 0 {
		user = userID.String()
	}

	m.MockAPI(name, http.MethodDelete,
		"channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+e.PathString()+"/"+user, nil)
}

// =============================================================================
// send.go
// =====================================================================================

// SendMessageComplex mocks a SendMessageComplex request.
//
// The ChannelID field of the passed discord.Message must be set.
func (m *Mocker) SendMessageComplex(d api.SendMessageData, msg discord.Message) {
	m.sendMessageComplex("SendMessageComplex", d, msg)
}

// sendMessageComplex mocks a SendMessageComplex request.
//
// The ChannelID field of the passed discord.Message must be set.
func (m *Mocker) sendMessageComplex(name string, d api.SendMessageData, msg discord.Message) {
	for i, embed := range d.Embeds {
		if embed.Type == "" {
			d.Embeds[i].Type = discord.NormalEmbed
		}

		if embed.Color == 0 {
			d.Embeds[i].Color = discord.DefaultEmbedColor
		}
	}

	for i, e := range msg.Embeds {
		if e.Type == "" {
			msg.Embeds[i].Type = discord.NormalEmbed
		}

		if e.Color == 0 {
			msg.Embeds[i].Color = discord.DefaultEmbedColor
		}
	}

	m.MockAPI(name, http.MethodPost, "channels/"+msg.ChannelID.String()+"/messages",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			if d.NeedsMultipart() {
				files := d.Files
				d.Files = nil
				check.Multipart(t, r.Body, r.Header, &d, files)
			} else {
				check.JSON(t, &d, r.Body)
			}

			check.WriteJSON(t, w, msg)
		})
}

// ExecuteWebhook mocks a ExecuteWebhook request and doesn't "wait" for the
// message to be delivered.
func (m *Mocker) ExecuteWebhook(webhookID discord.WebhookID, token string, d webhook.ExecuteData) {
	m.executeWebhook(webhookID, token, false, d, discord.Message{})
}

// ExecuteWebhookAndWait mocks a ExecuteWebhook request and "waits" for the
// message to be delivered.
func (m *Mocker) ExecuteWebhookAndWait(
	webhookID discord.WebhookID, token string, d webhook.ExecuteData, msg discord.Message,
) {
	m.executeWebhook(webhookID, token, true, d, msg)
}

// executeWebhook mocks a ExecuteWebhook request.
func (m *Mocker) executeWebhook(
	webhookID discord.WebhookID, token string, wait bool, d webhook.ExecuteData, msg discord.Message,
) {
	m.MockAPI("ExecuteWebhook", http.MethodPost, "webhooks/"+webhookID.String()+"/"+token,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			if wait {
				check.Query(t, url.Values{
					"wait": {"true"},
				}, r.URL.Query())
			}

			files := make([]sendpart.File, len(d.Files))
			copy(files, d.Files)

			d.Files = nil

			if len(files) == 0 {
				check.JSON(t, &d, r.Body)
			} else {
				check.Multipart(t, r.Body, r.Header, &d, files)
			}

			if wait {
				check.WriteJSON(t, w, msg)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		})
}

// =============================================================================
// webhook/webhook.go
// =====================================================================================

// WebhookWithToken mocks api.Client.WebhookWithToken.
//
// The ID field and the Token field of the passed discord.Webhook must be set.
//
// This method will sanitize Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) WebhookWithToken(wh discord.Webhook) {
	m.MockAPI("WebhookWithToken", http.MethodGet, "webhooks/"+wh.ID.String()+"/"+wh.Token,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			check.WriteJSON(t, w, wh)
		})
}

// ModifyWebhookWithToken mocks api.Client.ModifyWebhookWithToken.
func (m *Mocker) ModifyWebhookWithToken(d api.ModifyWebhookData, wh discord.Webhook) {
	m.MockAPI("ModifyWebhookWithToken", http.MethodPatch, "webhooks/"+wh.ID.String()+"/"+wh.Token,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			check.JSON(t, &d, r.Body)
			check.WriteJSON(t, w, wh)
		})
}

// DeleteWebhookWithToken mocks api.Client.DeleteWebhookWithToken.
func (m *Mocker) DeleteWebhookWithToken(id discord.WebhookID, token string) {
	m.MockAPI("DeleteWebhookWithToken", http.MethodDelete, "webhooks/"+id.String()+"/"+token, nil)
}
