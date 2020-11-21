package sanitize

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.User
		expect discord.User
	}{
		{
			name: "none",
			in: discord.User{
				ID: 321,
			},
			expect: discord.User{
				ID: 321,
			},
		},
		{
			name: "id",
			in:   discord.User{},
			expect: discord.User{
				ID: 123,
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := User(c.in, 123)

			assert.Equal(t, c.expect, actual)
		})
	}
}
