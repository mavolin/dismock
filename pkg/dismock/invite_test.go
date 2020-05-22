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

func TestMocker_Invite(t *testing.T) {
	m, s := NewArikawaSession(t)

	expect := sanitize.Invite(discord.Invite{
		Code: "abc",
	}, 1, 1, 1, 1, 1)

	m.Invite(expect)

	actual, err := s.Invite(expect.Code)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_InviteWithCounts(t *testing.T) {
	m, s := NewArikawaSession(t)

	expect := sanitize.Invite(discord.Invite{
		Code:                 "abc",
		ApproximatePresences: 123,
		ApproximateMembers:   456,
	}, 1, 1, 1, 1, 1)

	m.InviteWithCounts(expect)

	actual, err := s.InviteWithCounts(expect.Code)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_ChannelInvites(t *testing.T) {
	m, s := NewArikawaSession(t)

	var channelID discord.Snowflake = 123

	expect := []discord.Invite{
		{
			Code: "abc",
		},
		{
			Code: "def",
		},
	}

	for i, invite := range expect {
		expect[i] = sanitize.Invite(invite, 1, 1, 1, 1, 1)
	}

	m.ChannelInvites(channelID, expect)

	actual, err := s.ChannelInvites(channelID)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_GuildInvites(t *testing.T) {
	m, s := NewArikawaSession(t)

	var guildID discord.Snowflake = 123

	expect := []discord.Invite{
		{
			Code: "abc",
		},
		{
			Code: "def",
		},
	}

	for i, invite := range expect {
		expect[i] = sanitize.Invite(invite, 1, 1, 1, 1, 1)
	}

	m.GuildInvites(guildID, expect)

	actual, err := s.GuildInvites(guildID)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_CreateInvite(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewArikawaSession(t)

		data := api.CreateInviteData{
			MaxAge: option.NewUint(12),
		}

		expect := sanitize.Invite(discord.Invite{
			Code: "abc",
			Channel: discord.Channel{
				ID: 123,
			},
		}, 1, 1, 1, 1, 1)

		m.CreateInvite(data, expect)

		actual, err := s.CreateInvite(expect.Channel.ID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewArikawaSession(tMock)

		expect := sanitize.Invite(discord.Invite{
			Code: "abc",
			Channel: discord.Channel{
				ID: 123,
			},
		}, 1, 1, 1, 1, 1)

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
	m, s := NewArikawaSession(t)

	expect := sanitize.Invite(discord.Invite{
		Code:                 "abc",
		ApproximatePresences: 123,
		ApproximateMembers:   456,
	}, 1, 1, 1, 1, 1)

	m.DeleteInvite(expect)

	actual, err := s.DeleteInvite(expect.Code)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}
