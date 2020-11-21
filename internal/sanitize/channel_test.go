package sanitize

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.Channel
		expect discord.Channel
	}{
		{
			name: "none",
			in: discord.Channel{
				ID: 321,
			},
			expect: discord.Channel{
				ID: 321,
			},
		},
		{
			name: "id",
			in:   discord.Channel{},
			expect: discord.Channel{
				ID: 123,
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Channel(c.in, 123)

			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestMessage(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.Message
		expect discord.Message
	}{
		{
			name: "none",
			in: discord.Message{
				ID:        321,
				ChannelID: 654,
				Author: discord.User{
					ID: 987,
				},
			},
			expect: discord.Message{
				ID:        321,
				ChannelID: 654,
				Author: discord.User{
					ID: 987,
				},
			},
		},
		{
			name: "id",
			in: discord.Message{
				ChannelID: 654,
				Author: discord.User{
					ID: 987,
				},
			},
			expect: discord.Message{
				ID:        123,
				ChannelID: 654,
				Author: discord.User{
					ID: 987,
				},
			},
		},
		{
			name: "channelID",
			in: discord.Message{
				ID: 321,
				Author: discord.User{
					ID: 987,
				},
			},
			expect: discord.Message{
				ID:        321,
				ChannelID: 456,
				Author: discord.User{
					ID: 987,
				},
			},
		},
		{
			name: "authorID",
			in: discord.Message{
				ID:        321,
				ChannelID: 654,
			},
			expect: discord.Message{
				ID:        321,
				ChannelID: 654,
				Author: discord.User{
					ID: 789,
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Message(c.in, 123, 456, 789)

			assert.Equal(t, c.expect, actual)
		})
	}
}
