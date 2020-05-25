package sanitize

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
)

func TestWebhook(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.Webhook
		expect discord.Webhook
	}{
		{
			name: "none",
			in: discord.Webhook{
				ID: 321,
				User: discord.User{
					ID: 654,
				},
				ChannelID: 987,
			},
			expect: discord.Webhook{
				ID: 321,
				User: discord.User{
					ID: 654,
				},
				ChannelID: 987,
			},
		},
		{
			name: "id",
			in: discord.Webhook{
				User: discord.User{
					ID: 654,
				},
				ChannelID: 987,
			},
			expect: discord.Webhook{
				ID: 123,
				User: discord.User{
					ID: 654,
				},
				ChannelID: 987,
			},
		},
		{
			name: "user",
			in: discord.Webhook{
				ID:        321,
				ChannelID: 987,
			},
			expect: discord.Webhook{
				ID: 321,
				User: discord.User{
					ID: 456,
				},
				ChannelID: 987,
			},
		},
		{
			name: "channelID",
			in: discord.Webhook{
				ID: 321,
				User: discord.User{
					ID: 654,
				},
			},
			expect: discord.Webhook{
				ID: 321,
				User: discord.User{
					ID: 654,
				},
				ChannelID: 789,
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Webhook(c.in, 123, 456, 789)

			assert.Equal(t, c.expect, actual)
		})
	}
}
