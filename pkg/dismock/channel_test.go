package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/dismock/internal/sanitize"
)

func TestMocker_Channels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.GuildID = 123

		expect := []discord.Channel{
			{
				ID: 456,
			},
			{
				ID: 789,
			},
		}

		for i, c := range expect {
			expect[i] = sanitize.Channel(c, 1)
		}

		m.Channels(guildID, expect)

		actual, err := s.Channels(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("nil channels", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.GuildID = 123

		m.Channels(guildID, nil)

		actual, err := s.Channels(guildID)
		require.NoError(t, err)
		assert.Nil(t, actual)

		m.Eval()
	})
}

func TestMocker_CreateChannel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		data := api.CreateChannelData{
			Name: "abc",
			Permissions: []discord.Overwrite{
				{
					ID:    789,
					Type:  discord.OverwriteRole,
					Allow: 012,
					Deny:  345,
				},
			},
		}

		expect := sanitize.Channel(discord.Channel{
			ID:      456,
			GuildID: 123,
			Name:    "abc",
			Permissions: []discord.Overwrite{
				{
					ID:    789,
					Type:  discord.OverwriteRole,
					Allow: 012,
					Deny:  345,
				},
			},
		}, 1)

		m.CreateChannel(data, expect)

		actual, err := s.CreateChannel(expect.GuildID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := sanitize.Channel(discord.Channel{
			ID:      456,
			GuildID: 123,
			Name:    "abc",
		}, 1)

		m.CreateChannel(api.CreateChannelData{
			Name: "abc",
		}, expect)

		actual, err := s.CreateChannel(expect.GuildID, api.CreateChannelData{
			Name: "def",
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_MoveChannel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			guildID discord.GuildID = 123
			data                    = []api.MoveChannelData{
				{
					ID:       123,
					Position: option.NewInt(0),
				},
				{
					ID:       456,
					Position: option.NewInt(1),
				},
			}
		)

		m.MoveChannel(guildID, data)

		err := s.MoveChannel(guildID, data)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var guildID discord.GuildID = 123

		m.MoveChannel(guildID, []api.MoveChannelData{
			{
				ID:       123,
				Position: option.NewInt(0),
			},
			{
				ID:       456,
				Position: option.NewInt(1),
			},
		})

		err := s.MoveChannel(guildID, []api.MoveChannelData{
			{
				ID:       789,
				Position: option.NewInt(0),
			},
			{
				ID:       012,
				Position: option.NewInt(1),
			},
		})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_Channel(t *testing.T) {
	m, s := NewSession(t)

	expect := sanitize.Channel(discord.Channel{
		ID: 123,
		Permissions: []discord.Overwrite{
			{
				ID:    456,
				Type:  discord.OverwriteRole,
				Allow: 789,
				Deny:  012,
			},
		},
	}, 1)

	m.Channel(expect)

	actual, err := s.Channel(expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_ModifyChannel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			channelID discord.ChannelID = 123
			data                        = api.ModifyChannelData{
				Name: "abc",
				Permissions: &[]discord.Overwrite{
					{
						ID:    456,
						Type:  discord.OverwriteRole,
						Allow: 798,
						Deny:  012,
					},
				},
			}
		)

		m.ModifyChannel(channelID, data)

		err := s.ModifyChannel(channelID, data)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var channelID discord.ChannelID = 123

		m.ModifyChannel(channelID, api.ModifyChannelData{
			Name: "abc",
		})

		err := s.ModifyChannel(channelID, api.ModifyChannelData{
			Name: "def",
		})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteChannel(t *testing.T) {
	m, s := NewSession(t)

	var channelID discord.ChannelID = 123

	m.DeleteChannel(channelID)

	err := s.DeleteChannel(channelID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_EditChannelPermission(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			channelID   discord.ChannelID = 123
			overwriteID discord.Snowflake = 456
			data                          = api.EditChannelPermissionData{
				Type:  discord.OverwriteMember,
				Allow: 1,
				Deny:  0,
			}
		)

		m.EditChannelPermission(channelID, overwriteID, data)

		err := s.EditChannelPermission(channelID, overwriteID, data)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			channelID   discord.ChannelID = 123
			overwriteID discord.Snowflake = 456
		)

		m.EditChannelPermission(channelID, overwriteID, api.EditChannelPermissionData{
			Type:  discord.OverwriteMember,
			Allow: 1,
			Deny:  0,
		})

		err := s.EditChannelPermission(channelID, overwriteID, api.EditChannelPermissionData{
			Type:  discord.OverwriteMember,
			Allow: 0,
			Deny:  0,
		})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteChannelPermission(t *testing.T) {
	m, s := NewSession(t)

	var (
		channelID   discord.ChannelID = 123
		overwriteID discord.Snowflake = 456
	)

	m.DeleteChannelPermission(channelID, overwriteID)

	err := s.DeleteChannelPermission(channelID, overwriteID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_Typing(t *testing.T) {
	m, s := NewSession(t)

	var channelID discord.ChannelID = 123

	m.Typing(channelID)

	err := s.Typing(channelID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_PinnedMessages(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			channelID discord.ChannelID = 123
			userID    discord.UserID    = 456
		)

		expect := []discord.Message{
			{
				ID: 789,
			},
			{
				ID: 012,
			},
		}

		for i, m := range expect {
			expect[i] = sanitize.Message(m, 1, channelID, userID)
		}

		m.PinnedMessages(channelID, expect)

		actual, err := s.PinnedMessages(channelID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("nil messages", func(t *testing.T) {
		m, s := NewSession(t)

		var channelID discord.ChannelID = 123

		//noinspection GoPreferNilSlice
		expect := []discord.Message{}

		m.PinnedMessages(channelID, nil)

		actual, err := s.PinnedMessages(channelID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})
}

func TestMocker_PinMessage(t *testing.T) {
	m, s := NewSession(t)

	var (
		channelID discord.ChannelID = 123
		messageID discord.MessageID = 456
	)

	m.PinMessage(channelID, messageID)

	err := s.PinMessage(channelID, messageID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_UnpinMessage(t *testing.T) {
	m, s := NewSession(t)

	var (
		channelID discord.ChannelID = 123
		messageID discord.MessageID = 456
	)

	m.UnpinMessage(channelID, messageID)

	err := s.UnpinMessage(channelID, messageID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_AddRecipient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			channelID   discord.ChannelID = 123
			userID      discord.UserID    = 456
			accessToken                   = "abc"
			nickname                      = "Ragnar89"
		)

		m.AddRecipient(channelID, userID, accessToken, nickname)

		err := s.AddRecipient(channelID, userID, accessToken, nickname)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			channelID discord.ChannelID = 123
			userID    discord.UserID    = 456
		)

		m.AddRecipient(channelID, userID, "abc", "Ragnar89")

		err := s.AddRecipient(channelID, userID, "def", "Ragnar89")
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_RemoveRecipient(t *testing.T) {
	m, s := NewSession(t)

	var (
		channelID discord.ChannelID = 123
		userID    discord.UserID    = 456
	)

	m.RemoveRecipient(channelID, userID)

	err := s.RemoveRecipient(channelID, userID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_Ack(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
			ack                         = api.Ack{
				Token: "abc",
			}
		)

		expect := api.Ack{
			Token: "def",
		}

		actual := &ack

		m.Ack(channelID, messageID, ack, expect)

		err := s.Ack(channelID, messageID, actual)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
		)

		expect := api.Ack{
			Token: "def",
		}

		m.Ack(channelID, messageID, api.Ack{
			Token: "abc",
		}, expect)

		actual := &api.Ack{
			Token: "ghi",
		}

		err := s.Ack(channelID, messageID, actual)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}
