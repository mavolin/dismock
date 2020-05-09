package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_Emojis(t *testing.T) {
	m, s := New(t)

	var guildID discord.Snowflake = 123

	expect := []discord.Emoji{
		{
			ID: 456,
		},
		{
			ID: 789,
		},
	}

	m.Emojis(guildID, expect)

	actual, err := s.Emojis(guildID)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_Emoji(t *testing.T) {
	m, s := New(t)

	var guildID discord.Snowflake = 123

	expect := discord.Emoji{
		ID:   456,
		Name: "abc",
	}

	m.Emoji(guildID, expect)

	actual, err := s.Emoji(guildID, expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_CreateEmoji(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := New(t)

		var (
			guildID discord.Snowflake = 123
			image                     = api.Image{
				ContentType: "image/png",
				Content:     []byte{1, 255, 3},
			}
		)

		expect := discord.Emoji{
			ID:      456,
			Name:    "dismock",
			RoleIDs: []discord.Snowflake{789},
		}

		m.CreateEmoji(guildID, image, expect)

		actual, err := s.CreateEmoji(guildID, expect.Name, image, expect.RoleIDs)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := New(tMock)

		var guildID discord.Snowflake = 123

		expect := discord.Emoji{
			ID:      456,
			Name:    "abc",
			RoleIDs: []discord.Snowflake{789},
		}

		m.CreateEmoji(guildID, api.Image{
			ContentType: "image/png",
			Content:     []byte{1, 255, 3},
		}, expect)

		actual, err := s.CreateEmoji(guildID, expect.Name, api.Image{
			ContentType: "image/png",
			Content:     []byte{255, 0, 8},
		}, expect.RoleIDs)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_ModifyEmoji(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := New(t)

		var (
			guildID discord.Snowflake   = 123
			emojiID discord.Snowflake   = 456
			name                        = "abc"
			roles   []discord.Snowflake = nil
		)

		m.ModifyEmoji(guildID, emojiID, name, roles)

		err := s.ModifyEmoji(guildID, emojiID, name, roles)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := New(tMock)

		var (
			guildID discord.Snowflake = 123
			emojiID discord.Snowflake = 456
		)

		m.ModifyEmoji(guildID, emojiID, "", []discord.Snowflake{789, 012})

		err := s.ModifyEmoji(guildID, emojiID, "", []discord.Snowflake{345, 678})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteEmoji(t *testing.T) {
	m, s := New(t)

	var (
		guildID discord.Snowflake = 123
		emojiID discord.Snowflake = 456
	)

	m.DeleteEmoji(guildID, emojiID)

	err := s.DeleteEmoji(guildID, emojiID)
	require.NoError(t, err)

	m.Eval()
}
