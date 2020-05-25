package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/dismock/internal/sanitize"
)

func TestMocker_CreateWebhook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		data := api.CreateWebhookData{
			Name: "abc",
		}

		expect := sanitize.Webhook(discord.Webhook{
			ID:        123,
			Name:      "abc",
			ChannelID: 456,
		}, 1, 1, 1)

		m.CreateWebhook(data, expect)

		actual, err := s.CreateWebhook(expect.ChannelID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := sanitize.Webhook(discord.Webhook{
			ID:        123,
			Name:      "abc",
			ChannelID: 456,
		}, 1, 1, 1)

		m.CreateWebhook(api.CreateWebhookData{
			Name: "abc",
		}, expect)

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

	var channelID discord.Snowflake = 123

	expect := []discord.Webhook{
		{
			ID: 456,
		},
		{
			ID: 789,
		},
	}

	for i, w := range expect {
		expect[i] = sanitize.Webhook(w, 1, 1, channelID)
	}

	m.ChannelWebhooks(channelID, expect)

	actual, err := s.ChannelWebhooks(channelID)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_GuildWebhooks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.Snowflake = 123

		expect := []discord.Webhook{
			{
				ID: 456,
			},
			{
				ID: 789,
			},
		}

		for i, w := range expect {
			expect[i] = sanitize.Webhook(w, 1, 1, 1)

			if w.GuildID <= 0 {
				expect[i].GuildID = guildID
			}
		}

		m.GuildWebhooks(guildID, expect)

		actual, err := s.GuildWebhooks(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("auto guild id", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.Snowflake = 123

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

		m.Eval()
	})
}

func TestMocker_Webhook(t *testing.T) {
	m, s := NewSession(t)

	expect := sanitize.Webhook(discord.Webhook{
		ID: 456,
	}, 1, 1, 1)

	m.Webhook(expect)

	actual, err := s.Webhook(expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_WebhookWithToken(t *testing.T) {
	m, s := NewSession(t)

	expect := sanitize.Webhook(discord.Webhook{
		ID:    456,
		Token: "abc",
	}, 1, 1, 1)

	m.WebhookWithToken(expect)

	actual, err := s.WebhookWithToken(expect.ID, expect.Token)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_ModifyWebhook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		data := api.ModifyWebhookData{
			Name: option.NewString("abc"),
		}

		expect := sanitize.Webhook(discord.Webhook{
			ID:   456,
			Name: "abc",
		}, 1, 1, 1)

		m.ModifyWebhook(data, expect)

		actual, err := s.ModifyWebhook(expect.ID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := sanitize.Webhook(discord.Webhook{
			ID:   456,
			Name: "abc",
		}, 1, 1, 1)

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
		m, s := NewSession(t)

		data := api.ModifyWebhookData{
			Name: option.NewString("abc"),
		}

		expect := sanitize.Webhook(discord.Webhook{
			ID:    456,
			Name:  "abc",
			Token: "def",
		}, 1, 1, 1)

		m.ModifyWebhookWithToken(data, expect)

		actual, err := s.ModifyWebhookWithToken(expect.ID, expect.Token, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := sanitize.Webhook(discord.Webhook{
			ID:    456,
			Name:  "abc",
			Token: "def",
		}, 1, 1, 1)

		m.ModifyWebhookWithToken(api.ModifyWebhookData{
			Name: option.NewString("abc"),
		}, expect)

		actual, err := s.ModifyWebhookWithToken(expect.ID, expect.Token, api.ModifyWebhookData{
			Name: option.NewString("cba"),
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteWebhook(t *testing.T) {
	m, s := NewSession(t)

	var id discord.Snowflake = 123

	m.DeleteWebhook(id)

	err := s.DeleteWebhook(id)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_DeleteWebhookWithToken(t *testing.T) {
	m, s := NewSession(t)

	var (
		id    discord.Snowflake = 123
		token                   = "abc"
	)

	m.DeleteWebhookWithToken(id, token)

	err := s.DeleteWebhookWithToken(id, token)
	require.NoError(t, err)

	m.Eval()
}
