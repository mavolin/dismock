package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/diamondburned/arikawa/v2/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_CreateWebhook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		data := api.CreateWebhookData{
			Name: "abc",
		}

		expect := discord.Webhook{
			ID:        123,
			Name:      "abc",
			ChannelID: 456,
			User:      discord.User{ID: 789},
		}

		m.CreateWebhook(data, expect)

		actual, err := s.CreateWebhook(expect.ChannelID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := discord.Webhook{
			ID:        123,
			Name:      "abc",
			ChannelID: 456,
			User:      discord.User{ID: 789},
		}

		m.CreateWebhook(api.CreateWebhookData{Name: "abc"}, expect)

		actual, err := s.CreateWebhook(expect.ChannelID, api.CreateWebhookData{
			Name: "cba",
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_ChannelWebhooks(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	var channelID discord.ChannelID = 123

	expect := []discord.Webhook{
		{
			ID:        456,
			ChannelID: channelID,
			User:      discord.User{ID: 789},
		},
		{
			ID:        789,
			ChannelID: channelID,
			User:      discord.User{ID: 012},
		},
	}

	m.ChannelWebhooks(channelID, expect)

	actual, err := s.ChannelWebhooks(channelID)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)
}

func TestMocker_GuildWebhooks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		expect := []discord.Webhook{
			{
				ID:        456,
				ChannelID: 789,
				GuildID:   guildID,
				User:      discord.User{ID: 012},
			},
			{
				ID:        789,
				ChannelID: 012,
				GuildID:   guildID,
				User:      discord.User{ID: 345},
			},
		}

		m.GuildWebhooks(guildID, expect)

		actual, err := s.GuildWebhooks(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})

	t.Run("auto guild id", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		in := []discord.Webhook{
			{
				ID: 456,
			},
			{
				ID: 789,
			},
		}

		m.GuildWebhooks(guildID, in)

		actual, err := s.GuildWebhooks(guildID)
		require.NoError(t, err)

		for _, w := range actual {
			assert.Equal(t, guildID, w.GuildID)
		}
	})
}

func TestMocker_Webhook(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	expect := discord.Webhook{
		ID:        123,
		ChannelID: 456,
		User:      discord.User{ID: 789},
	}

	m.Webhook(expect)

	actual, err := s.Webhook(expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_WebhookWithToken(t *testing.T) {
	m := New(t)

	expect := discord.Webhook{
		ID:        123,
		Token:     "abc",
		ChannelID: 456,
		User:      discord.User{ID: 789},
	}

	m.WebhookWithToken(expect)

	actual, err := webhook.Get(expect.ID, expect.Token)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_ModifyWebhook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		data := api.ModifyWebhookData{
			Name: option.NewString("abc"),
		}

		expect := discord.Webhook{
			ID:        456,
			Name:      "abc",
			ChannelID: 456,
			User:      discord.User{ID: 789},
		}

		m.ModifyWebhook(data, expect)

		actual, err := s.ModifyWebhook(expect.ID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := discord.Webhook{
			ID:        123,
			Name:      "abc",
			ChannelID: 456,
			User:      discord.User{ID: 789},
		}

		m.ModifyWebhook(api.ModifyWebhookData{
			Name: option.NewString("abc"),
		}, expect)

		actual, err := s.ModifyWebhook(expect.ID, api.ModifyWebhookData{
			Name: option.NewString("cba"),
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_ModifyWebhookWithTokenWithToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := New(t)

		data := api.ModifyWebhookData{
			Name: option.NewString("abc"),
		}

		expect := discord.Webhook{
			ID:        123,
			Name:      "abc",
			Token:     "def",
			ChannelID: 456,
			User:      discord.User{ID: 789},
		}

		m.ModifyWebhookWithToken(data, expect)

		actual, err := webhook.Modify(expect.ID, expect.Token, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m := New(tMock)

		expect := discord.Webhook{
			ID:        123,
			Name:      "abc",
			Token:     "def",
			ChannelID: 456,
			User:      discord.User{ID: 789},
		}

		m.ModifyWebhookWithToken(api.ModifyWebhookData{
			Name: option.NewString("abc"),
		}, expect)

		actual, err := webhook.Modify(expect.ID, expect.Token, api.ModifyWebhookData{
			Name: option.NewString("cba"),
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteWebhook(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	var id discord.WebhookID = 123

	m.DeleteWebhook(id)

	err := s.DeleteWebhook(id)
	require.NoError(t, err)
}

func TestMocker_DeleteWebhookWithToken(t *testing.T) {
	m := New(t)

	var (
		id    discord.WebhookID = 123
		token                   = "abc"
	)

	m.DeleteWebhookWithToken(id, token)

	err := webhook.Delete(id, token)
	require.NoError(t, err)
}
