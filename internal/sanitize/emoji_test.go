package sanitize

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"
)

func TestEmoji(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.Emoji
		expect discord.Emoji
	}{
		{
			name: "none",
			in: discord.Emoji{
				ID: 321,
				User: discord.User{
					ID: 654,
				},
			},
			expect: discord.Emoji{
				ID: 321,
				User: discord.User{
					ID: 654,
				},
			},
		},
		{
			name: "id",
			in: discord.Emoji{
				User: discord.User{
					ID: 654,
				},
			},
			expect: discord.Emoji{
				ID: 123,
				User: discord.User{
					ID: 654,
				},
			},
		},
		{
			name: "userID",
			in: discord.Emoji{
				ID: 321,
			},
			expect: discord.Emoji{
				ID: 321,
				User: discord.User{
					ID: 456,
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Emoji(c.in, 123, 456)

			assert.Equal(t, c.expect, actual)
		})
	}
}
