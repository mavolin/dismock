package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_Emojis(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		expect := []discord.Emoji{
			{
				ID:   456,
				User: discord.User{ID: 1},
			},
			{
				ID:   789,
				User: discord.User{ID: 1},
			},
		}

		m.Emojis(guildID, expect)

		actual, err := s.Emojis(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})

	t.Run("nil emojis", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		//noinspection GoPreferNilSlice
		expect := []discord.Emoji{}

		m.Emojis(guildID, nil)

		actual, err := s.Emojis(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})
}

func TestMocker_Emoji(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	var guildID discord.GuildID = 123

	expect := discord.Emoji{
		ID:   456,
		Name: "abc",
		User: discord.User{ID: 1},
	}

	m.Emoji(guildID, expect)

	actual, err := s.Emoji(guildID, expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_CreateEmoji(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var (
			guildID discord.GuildID = 123
			data                    = api.CreateEmojiData{
				Name: "dismock",
				Image: api.Image{
					ContentType: "image/png",
					Content:     []byte{1, 255, 3},
				},
			}
		)

		expect := discord.Emoji{
			ID:      456,
			Name:    data.Name,
			RoleIDs: []discord.RoleID{789},
			User:    discord.User{ID: 1},
		}

		m.CreateEmoji(guildID, data, expect)

		actual, err := s.CreateEmoji(guildID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var guildID discord.GuildID = 123

		expect := discord.Emoji{
			ID:   456,
			Name: "abc",
			User: discord.User{ID: 1},
		}

		m.CreateEmoji(guildID, api.CreateEmojiData{
			Name: expect.Name,
			Image: api.Image{
				ContentType: "image/png",
				Content:     []byte{0, 255, 100},
			},
		}, expect)

		actual, err := s.CreateEmoji(guildID, api.CreateEmojiData{
			Name: expect.Name,
			Image: api.Image{
				ContentType: "image/png",
				Content:     []byte{1, 255, 3},
			},
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_ModifyEmoji(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var (
			guildID discord.GuildID = 123
			emojiID discord.EmojiID = 456
			data                    = api.ModifyEmojiData{Name: "abc"}
		)

		m.ModifyEmoji(guildID, emojiID, data)

		err := s.ModifyEmoji(guildID, emojiID, data)
		require.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			guildID discord.GuildID = 123
			emojiID discord.EmojiID = 456
		)

		m.ModifyEmoji(guildID, emojiID, api.ModifyEmojiData{
			Roles: &[]discord.RoleID{789, 012},
		})

		err := s.ModifyEmoji(guildID, emojiID, api.ModifyEmojiData{
			Roles: &[]discord.RoleID{345, 678},
		})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteEmoji(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	var (
		guildID discord.GuildID = 123
		emojiID discord.EmojiID = 456
	)

	m.DeleteEmoji(guildID, emojiID)

	err := s.DeleteEmoji(guildID, emojiID)
	require.NoError(t, err)
}
