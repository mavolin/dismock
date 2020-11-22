package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_Invite(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	expect := discord.Invite{
		Code:    "abc",
		Channel: discord.Channel{ID: 123},
	}

	m.Invite(expect)

	actual, err := s.Invite(expect.Code)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_InviteWithCounts(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	expect := discord.Invite{
		Code:                 "abc",
		Channel:              discord.Channel{ID: 123},
		ApproximatePresences: 456,
		ApproximateMembers:   789,
	}

	m.InviteWithCounts(expect)

	actual, err := s.InviteWithCounts(expect.Code)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_ChannelInvites(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var channelID discord.ChannelID = 123

		expect := []discord.Invite{
			{
				Code:    "abc",
				Channel: discord.Channel{ID: channelID},
			},
			{
				Code:    "def",
				Channel: discord.Channel{ID: channelID},
			},
		}

		m.ChannelInvites(channelID, expect)

		actual, err := s.ChannelInvites(channelID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})

	t.Run("nil invites", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var channelID discord.ChannelID = 123

		//noinspection GoPreferNilSlice
		expect := []discord.Invite{}

		m.ChannelInvites(channelID, nil)

		actual, err := s.ChannelInvites(channelID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})
}

func TestMocker_GuildInvites(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		expect := []discord.Invite{
			{
				Code:    "abc",
				Channel: discord.Channel{ID: 456, GuildID: guildID},
			},
			{
				Code:    "def",
				Channel: discord.Channel{ID: 456, GuildID: guildID},
			},
		}

		m.GuildInvites(guildID, expect)

		actual, err := s.GuildInvites(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})

	t.Run("nil invites", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		//noinspection GoPreferNilSlice
		expect := []discord.Invite{}

		m.GuildInvites(guildID, nil)

		actual, err := s.GuildInvites(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})
}

func TestMocker_CreateInvite(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		data := api.CreateInviteData{MaxAge: option.NewUint(12)}

		expect := discord.Invite{
			Code:    "abc",
			Channel: discord.Channel{ID: 123},
		}

		m.CreateInvite(data, expect)

		actual, err := s.CreateInvite(expect.Channel.ID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := discord.Invite{
			Code:    "abc",
			Channel: discord.Channel{ID: 123},
		}

		m.CreateInvite(api.CreateInviteData{
			MaxAge: option.NewUint(12),
		}, expect)

		actual, err := s.CreateInvite(expect.Channel.ID, api.CreateInviteData{
			MaxAge: option.NewUint(21),
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteInvite(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	expect := discord.Invite{
		Code:                 "abc",
		Channel:              discord.Channel{ID: 123},
		ApproximatePresences: 456,
		ApproximateMembers:   789,
	}

	m.DeleteInvite(expect)

	actual, err := s.DeleteInvite(expect.Code)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}
