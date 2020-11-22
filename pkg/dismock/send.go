package dismock

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/httputil/httpdriver"
	"github.com/diamondburned/arikawa/v2/webhook"

	"github.com/mavolin/dismock/v2/internal/mockutil"
)

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
	if d.Embed != nil {
		if d.Embed.Type == "" {
			d.Embed.Type = discord.NormalEmbed
		}

		if d.Embed.Color == 0 {
			d.Embed.Color = discord.DefaultEmbedColor
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

	m.MockAPI(name, http.MethodPost, "/channels/"+msg.ChannelID.String()+"/messages",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			files := make([]api.SendMessageFile, len(d.Files))
			copy(files, d.Files)

			d.Files = nil

			if len(files) == 0 {
				mockutil.CheckJSON(t, r.Body, new(api.SendMessageData), &d)
			} else {
				mockutil.CheckMultipart(t, r.Body, r.Header, new(api.SendMessageData), &d, files)
			}

			mockutil.WriteJSON(t, w, msg)
		})
}

// ExecuteWebhook mocks a ExecuteWebhook request and doesn't "wait" for the
// message to be delivered.
func (m *Mocker) ExecuteWebhook(webhookID discord.WebhookID, token string, d api.ExecuteWebhookData) {
	m.executeWebhook(webhookID, token, false, d, discord.Message{})
}

// ExecuteWebhookAndWait mocks a ExecuteWebhook request and "waits" for the
// message to be delivered.
func (m *Mocker) ExecuteWebhookAndWait(
	webhookID discord.WebhookID, token string, d api.ExecuteWebhookData, msg discord.Message,
) {
	m.executeWebhook(webhookID, token, true, d, msg)
}

// executeWebhook mocks a ExecuteWebhook request.
func (m *Mocker) executeWebhook(
	webhookID discord.WebhookID, token string, wait bool, d api.ExecuteWebhookData, msg discord.Message,
) {
	webhook.DefaultHTTPClient.Client = httpdriver.WrapClient(*m.Client)

	m.MockAPI("ExecuteWebhook", http.MethodPost, "/webhooks/"+webhookID.String()+"/"+token,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			if wait {
				mockutil.CheckQuery(t, r.URL.Query(), url.Values{
					"wait": {"true"},
				})
			}

			files := make([]api.SendMessageFile, len(d.Files))
			copy(files, d.Files)

			d.Files = nil

			if len(files) == 0 {
				mockutil.CheckJSON(t, r.Body, new(api.ExecuteWebhookData), &d)
			} else {
				mockutil.CheckMultipart(t, r.Body, r.Header, new(api.ExecuteWebhookData), &d, files)
			}

			if wait {
				mockutil.WriteJSON(t, w, msg)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		})
}
