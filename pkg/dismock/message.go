package dismock

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"

	"github.com/mavolin/dismock/v2/internal/check"
)

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
	m.MockAPI(name, http.MethodGet, "/channels/"+channelID.String()+"/messages",
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

			check.Query(t, r.URL.Query(), expect)
			check.WriteJSON(t, w, messages)
		})
}

// Message mocks a Message request.
//
// The ID field and the ChannelID field of the passed discord.Message must be
// set.
func (m *Mocker) Message(msg discord.Message) {
	m.MockAPI("Message", http.MethodGet, "/channels/"+msg.ChannelID.String()+"/messages/"+msg.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			check.WriteJSON(t, w, msg)
		})
}

// SendText mocks a SendText request.
//
// The ChannelID field and the Content field of the passed discord.Message must
// be set.
func (m *Mocker) SendText(msg discord.Message) {
	m.sendMessageComplex("SendText", api.SendMessageData{
		Content: msg.Content,
	}, msg)
}

// SendEmbed mocks a SendEmbed request.
//
// The ChannelID field and the Embed field of the passed discord.Message must
// be set.
func (m *Mocker) SendEmbed(msg discord.Message) {
	m.sendMessageComplex("SendEmbed", api.SendMessageData{
		Embed: &msg.Embeds[0],
	}, msg)
}

// SendMessage mocks a SendMessage request.
//
// The ChannelID field and the Content field of the passed discord.Message must
// be set.
//
// This method will sanitize Message.ID, Message.Author.ID, Message.Embeds.Type
// and Message.Embeds.Color.
func (m *Mocker) SendMessage(embed *discord.Embed, msg discord.Message) {
	d := api.SendMessageData{
		Content: msg.Content,
	}

	if embed != nil {
		d.Embed = embed

		if len(msg.Embeds) == 0 {
			msg.Embeds = append(msg.Embeds, *d.Embed)
		}
	}

	m.sendMessageComplex("SendMessage", d, msg)
}

// EditText mocks a EditText request.
//
// The ID field, the ChannelID field and the Content field of the passed
// Message must be set.
func (m *Mocker) EditText(msg discord.Message) {
	m.editMessageComplex("EditText", api.EditMessageData{
		Content: option.NewNullableString(msg.Content),
	}, msg)
}

// EditEmbed mocks a EditEmbed request.
//
// The ID field, the ChannelID field and the Embed[0] field of the passed
// discord.Message must be set.
func (m *Mocker) EditEmbed(msg discord.Message) {
	m.editMessageComplex("EditEmbed", api.EditMessageData{
		Embed: &msg.Embeds[0],
	}, msg)
}

// EditMessage mocks a EditMessage request.
//
// The ID field, the ChannelID field, the Content field of the passed
// discord.Message must be set.
func (m *Mocker) EditMessage(embed *discord.Embed, msg discord.Message, suppressEmbeds bool) {
	d := api.EditMessageData{
		Content: option.NewNullableString(msg.Content),
		Embed:   embed,
	}

	if suppressEmbeds {
		d.Flags = &discord.SuppressEmbeds
	}

	m.editMessageComplex("EditMessage", d, msg)
}

// EditMessageComplex mocks a EditMessageComplex request.
//
// The ID field and the ChannelID field of the passed discord.Message must be
// set.
func (m *Mocker) EditMessageComplex(d api.EditMessageData, msg discord.Message) {
	m.editMessageComplex("EditMessageComplex", d, msg)
}

// editMessageComplex mocks a EditMessageComplex request.
//
// The ID field and the ChannelID field of the passed discord.Message must be
// set.
func (m *Mocker) editMessageComplex(name string, d api.EditMessageData, msg discord.Message) {
	if d.Embed != nil {
		if d.Embed.Type == "" {
			d.Embed.Type = discord.NormalEmbed
		}

		if d.Embed.Color == 0 {
			d.Embed.Color = discord.DefaultEmbedColor
		}
	}

	m.MockAPI(name, http.MethodPatch, "/channels/"+msg.ChannelID.String()+"/messages/"+msg.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			check.JSON(t, r.Body, &d)
			check.WriteJSON(t, w, msg)
		})
}

// DeleteMessage mocks a DeleteMessage request.
func (m *Mocker) DeleteMessage(channelID discord.ChannelID, messageID discord.MessageID) {
	m.MockAPI("DeleteMessage", http.MethodDelete, "/channels/"+channelID.String()+"/messages/"+messageID.String(), nil)
}

type deleteMessagesPayload struct {
	Messages []discord.MessageID `json:"messages"`
}

// DeleteMessages mocks a DeleteMessages request.
func (m *Mocker) DeleteMessages(channelID discord.ChannelID, messageIDs []discord.MessageID) {
	m.MockAPI("DeleteMessages", http.MethodPost, "/channels/"+channelID.String()+"/messages/bulk-delete",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := deleteMessagesPayload{
				Messages: messageIDs,
			}

			check.JSON(t, r.Body, &expect)
			w.WriteHeader(http.StatusNoContent)
		})
}
