package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/dismock/v2/internal/mockutil"
)

// CreateWebhook mocks a CreateWebhook request.
func (m *Mocker) CreateWebhook(d api.CreateWebhookData, wh discord.Webhook) {
	m.MockAPI("CreateWebhook", http.MethodPost, "/channels/"+wh.ChannelID.String()+"/webhooks",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.CreateWebhookData), &d)
			mockutil.WriteJSON(t, w, wh)
		})
}

// ChannelWebhooks mocks a ChannelWebhooks request.
func (m *Mocker) ChannelWebhooks(channelID discord.ChannelID, webhooks []discord.Webhook) {
	m.MockAPI("ChannelWebhooks", http.MethodGet, "/channels/"+channelID.String()+"/webhooks",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, webhooks)
		})
}

// GuildWebhooks mocks a GuildWebhooks request.
func (m *Mocker) GuildWebhooks(guildID discord.GuildID, webhooks []discord.Webhook) {
	for i, w := range webhooks {
		if w.GuildID == 0 {
			webhooks[i].GuildID = guildID
		}
	}

	m.MockAPI("GuildWebhooks", http.MethodGet, "/guilds/"+guildID.String()+"/webhooks",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, webhooks)
		})
}

// Webhook mocks a Webhook request.
//
// The ID field of the passed discord.Webhook must be set.
func (m *Mocker) Webhook(webhook discord.Webhook) {
	m.MockAPI("Webhook", http.MethodGet, "/webhooks/"+webhook.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, webhook)
		})
}

// WebhookWithToken mocks a WebhookWithToken request.
//
// The ID field and the Token field of the passed discord.Webhook must be set.
//
// This method will sanitize Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) WebhookWithToken(wh discord.Webhook) {
	m.MockAPI("WebhookWithToken", http.MethodGet, "/webhooks/"+wh.ID.String()+"/"+wh.Token,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, wh)
		})
}

// ModifyWebhook mocks a ModifyWebhook request.
//
// The ID field of the passed discord.Webhook must be set.
//
// This method will sanitize Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) ModifyWebhook(d api.ModifyWebhookData, wh discord.Webhook) {
	m.MockAPI("ModifyWebhook", http.MethodPatch, "/webhooks/"+wh.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyWebhookData), &d)
			mockutil.WriteJSON(t, w, wh)
		})
}

// ModifyWebhookWithToken mocks a ModifyWebhookWithToken request.
//
// The ID field and the Token field of the passed discord.Webhook must be set.
//
// This method will sanitize Webhook.User.ID and Webhook.ChannelID.
func (m *Mocker) ModifyWebhookWithToken(d api.ModifyWebhookData, wh discord.Webhook) {
	m.MockAPI("ModifyWebhookWithToken", http.MethodPatch, "/webhooks/"+wh.ID.String()+"/"+wh.Token,
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyWebhookData), &d)
			mockutil.WriteJSON(t, w, wh)
		})
}

// DeleteWebhook mocks a DeleteWebhook request.
func (m *Mocker) DeleteWebhook(id discord.WebhookID) {
	m.MockAPI("DeleteWebhook", http.MethodDelete, "/webhooks/"+id.String(), nil)
}

// DeleteWebhookWithToken mocks a DeleteWebhookWithToken request.
func (m *Mocker) DeleteWebhookWithToken(id discord.WebhookID, token string) {
	m.MockAPI("DeleteWebhookWithToken", http.MethodDelete, "/webhooks/"+id.String()+"/"+token, nil)
}
