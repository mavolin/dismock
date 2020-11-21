package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_User(t *testing.T) {
	m, s := NewSession(t)

	expect := discord.User{
		ID: 123,
	}

	m.User(expect)

	actual, err := s.User(expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_Me(t *testing.T) {
	m, s := NewSession(t)

	expect := discord.User{
		ID: 123,
	}

	m.Me(expect)

	actual, err := s.Me()
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_ModifyMe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		data := api.ModifySelfData{
			Username: option.NewString("abc"),
		}

		expect := discord.User{
			ID:       123,
			Username: "abc",
		}

		m.ModifyMe(data, expect)

		actual, err := s.ModifyMe(data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := discord.User{
			ID:       123,
			Username: "abc",
		}

		m.ModifyMe(api.ModifySelfData{
			Username: option.NewString("abc"),
		}, expect)

		actual, err := s.ModifyMe(api.ModifySelfData{
			Username: option.NewString("cba"),
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_ChangeOwnNickname(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			guildID discord.GuildID = 123
			nick                    = "abc"
		)

		m.ChangeOwnNickname(guildID, nick)

		err := s.ChangeOwnNickname(guildID, nick)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var guildID discord.GuildID = 123

		m.ChangeOwnNickname(guildID, "abc")

		err := s.ChangeOwnNickname(guildID, "cba")
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_PrivateChannels(t *testing.T) {
	m, s := NewSession(t)

	expect := []discord.Channel{
		{
			ID: 123,
		},
		{
			ID: 456,
		},
	}

	m.PrivateChannels(expect)

	actual, err := s.PrivateChannels()
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_CreatePrivateChannel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		expect := discord.Channel{
			ID: 123,
			DMRecipients: []discord.User{
				{
					ID: 456,
				},
			},
		}

		m.CreatePrivateChannel(expect)

		actual, err := s.CreatePrivateChannel(expect.DMRecipients[0].ID)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := discord.Channel{
			ID: 123,
			DMRecipients: []discord.User{
				{
					ID: 456,
				},
			},
		}

		m.CreatePrivateChannel(discord.Channel{
			ID: 123,
			DMRecipients: []discord.User{
				{
					ID: 456,
				},
			},
		})

		actual, err := s.CreatePrivateChannel(654)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_UserConnections(t *testing.T) {
	m, s := NewSession(t)

	expect := []discord.Connection{
		{
			ID: "123",
		},
		{
			ID: "456",
		},
	}

	m.UserConnections(expect)

	actual, err := s.UserConnections()
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}
