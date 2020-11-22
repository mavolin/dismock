package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_Messages(t *testing.T) {
	successCases := []struct {
		name     string
		messages int
		limit    uint
	}{
		{
			name:     "limited",
			messages: 130,
			limit:    199,
		},
		{
			name:     "unlimited",
			messages: 200,
			limit:    0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)
				defer m.Eval()

				var channelID discord.ChannelID = 123

				expect := make([]discord.Message, c.messages)

				for i := 0; i < c.messages; i++ {
					expect[i] = discord.Message{
						ID:        discord.MessageID(c.messages - i + 1),
						ChannelID: channelID,
						GuildID:   456,
						Author:    discord.User{ID: 789},
					}
				}

				m.Messages(channelID, c.limit, expect)

				actual, err := s.Messages(channelID, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("nil messages", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var channelID discord.ChannelID = 123

		m.Messages(channelID, 100, nil)

		actual, err := s.Messages(channelID, 100)
		require.NoError(t, err)

		assert.Len(t, actual, 0)
	})

	t.Run("limit smaller than messages", func(t *testing.T) {
		m, _ := NewSession(t)

		assert.Panics(t, func() {
			m.Messages(123, 1, []discord.Message{{}, {}})
		})
	})
}

func TestMocker_MessagesAround(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var (
			channelID discord.ChannelID = 123
			around    discord.MessageID = 456
		)

		expect := make([]discord.Message, 100)

		for i := 0; i < len(expect); i++ {
			expect[i] = discord.Message{
				ID:        discord.MessageID(int(around) - i + 1),
				ChannelID: channelID,
				GuildID:   456,
				Author:    discord.User{ID: 789},
			}
		}

		m.MessagesAround(channelID, around, 100, expect)

		actual, err := s.MessagesAround(channelID, around, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})

	limitCases := []struct {
		name     string
		messages int
		limit    uint
	}{
		{
			name:     "limit within range",
			messages: 70,
			limit:    70,
		},
		{
			name:     "limit > 100",
			messages: 100,
			limit:    199,
		},
		{
			name:     "limit 0",
			messages: 50,
			limit:    0,
		},
	}

	for _, c := range limitCases {
		t.Run(c.name, func(t *testing.T) {
			m, s := NewSession(t)
			defer m.Eval()

			var (
				channelID discord.ChannelID = 123
				around    discord.MessageID = 456
			)

			expect := make([]discord.Message, c.messages)

			for i := 0; i < c.messages; i++ {
				expect[i] = discord.Message{
					ID:        discord.MessageID(c.messages - i + 1),
					ChannelID: channelID,
					GuildID:   456,
					Author:    discord.User{ID: 789},
				}
			}

			m.MessagesAround(channelID, around, c.limit, expect)

			actual, err := s.MessagesAround(channelID, around, c.limit)
			require.NoError(t, err)

			assert.Equal(t, expect, actual)
		})
	}

	t.Run("nil messages", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var channelID discord.ChannelID = 123

		//noinspection GoPreferNilSlice
		expect := []discord.Message{}

		m.MessagesAround(channelID, 0, 100, nil)

		actual, err := s.MessagesAround(channelID, 0, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var channelID discord.ChannelID = 123

		expect := []discord.Message{
			{
				ID:        123,
				ChannelID: channelID,
				GuildID:   456,
				Author:    discord.User{ID: 789},
			},
			{
				ID:        456,
				ChannelID: channelID,
				GuildID:   456,
				Author:    discord.User{ID: 789},
			},
		}

		m.MessagesAround(channelID, 890, 100, expect)

		actual, err := s.MessagesAround(channelID, 789, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than messages", func(t *testing.T) {
		m, _ := NewSession(t)

		assert.Panics(t, func() {
			m.MessagesAround(123, 0, 1, []discord.Message{{}, {}})
		})
	})
}

func TestMocker_MessagesBefore(t *testing.T) {
	successCases := []struct {
		name     string
		messages int
		limit    uint
	}{
		{
			name:     "limited",
			messages: 130,
			limit:    199,
		},
		{
			name:     "unlimited",
			messages: 200,
			limit:    0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)
				defer m.Eval()

				var (
					channelID discord.ChannelID = 123
					before    discord.MessageID = 9999
				)

				expect := make([]discord.Message, c.messages)

				for i := 0; i < c.messages; i++ {
					expect[i] = discord.Message{
						ID:        discord.MessageID(c.messages - i + 1),
						ChannelID: channelID,
						GuildID:   456,
						Author:    discord.User{ID: 789},
					}
				}

				m.MessagesBefore(channelID, before, c.limit, expect)

				actual, err := s.MessagesBefore(channelID, before, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("nil messages", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var channelID discord.ChannelID = 123

		m.MessagesBefore(channelID, 0, 100, nil)

		actual, err := s.MessagesBefore(channelID, 0, 100)
		require.NoError(t, err)

		assert.Len(t, actual, 0)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var channelID discord.ChannelID = 123

		expect := []discord.Message{
			{
				ID:        123,
				ChannelID: channelID,
				GuildID:   456,
				Author:    discord.User{ID: 789},
			},
			{
				ID:        456,
				ChannelID: channelID,
				GuildID:   456,
				Author:    discord.User{ID: 789},
			},
		}

		m.MessagesBefore(channelID, 890, 100, expect)

		actual, err := s.MessagesBefore(channelID, 789, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than messages", func(t *testing.T) {
		m, _ := NewSession(t)

		assert.Panics(t, func() {
			m.MessagesBefore(123, 0, 1, []discord.Message{{}, {}})
		})
	})
}

func TestMocker_MessagesAfter(t *testing.T) {
	successCases := []struct {
		name     string
		messages int
		limit    uint
	}{
		{
			name:     "limited",
			messages: 130,
			limit:    199,
		},
		{
			name:     "unlimited",
			messages: 200,
			limit:    0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)
				defer m.Eval()

				var (
					channelID discord.ChannelID = 123
					after     discord.MessageID = 456
				)

				expect := make([]discord.Message, c.messages)

				for i := 0; i < c.messages; i++ {
					expect[i] = discord.Message{
						ID:        discord.MessageID(int(after) + c.messages + 1),
						ChannelID: channelID,
						GuildID:   789,
						Author:    discord.User{ID: 012},
					}
				}

				m.MessagesAfter(channelID, after, c.limit, expect)

				actual, err := s.MessagesAfter(channelID, after, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("nil guilds", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var channelID discord.ChannelID = 123

		var expect []discord.Message

		m.MessagesAfter(channelID, 0, 100, expect)

		actual, err := s.MessagesAfter(channelID, 0, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var channelID discord.ChannelID = 123

		expect := []discord.Message{
			{
				ID:        456,
				ChannelID: channelID,
				GuildID:   456,
				Author:    discord.User{ID: 789},
			},
			{
				ID:        789,
				ChannelID: channelID,
				GuildID:   456,
				Author:    discord.User{ID: 789},
			},
		}

		m.MessagesAfter(channelID, 123, 100, expect)

		actual, err := s.MessagesAfter(channelID, 321, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than messages", func(t *testing.T) {
		m, _ := NewSession(t)

		assert.Panics(t, func() {
			m.MessagesAfter(123, 0, 1, []discord.Message{{}, {}})
		})
	})
}

func TestMocker_Message(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	expect := discord.Message{
		ID:        123,
		ChannelID: 465,
		GuildID:   456,
		Author:    discord.User{ID: 789},
	}

	m.Message(expect)

	actual, err := s.Message(expect.ChannelID, expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_SendText(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		expect := discord.Message{
			ID:        123,
			ChannelID: 456,
			Author:    discord.User{ID: 789},
			Content:   "abc",
		}

		m.SendText(expect)

		actual, err := s.SendText(expect.ChannelID, expect.Content)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var channelID discord.ChannelID = 123

		m.SendText(discord.Message{
			ChannelID: channelID,
			Content:   "abc",
		})

		_, err := s.SendText(channelID, "cba")
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_SendEmbed(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		expect := discord.Message{
			ID:        123,
			ChannelID: 456,
			Author:    discord.User{ID: 789},
			Embeds: []discord.Embed{
				{
					Title:       "def",
					Description: "ghi",
				},
			},
		}

		m.SendEmbed(expect)

		actual, err := s.SendEmbed(expect.ChannelID, expect.Embeds[0])
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
		)

		m.SendEmbed(discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds: []discord.Embed{
				{
					Title:       "def",
					Description: "ghi",
				},
			},
		})

		_, err := s.SendEmbed(channelID, discord.Embed{
			Title:       "fed",
			Description: "ihg",
		})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_SendMessage(t *testing.T) {
	successCases := []struct {
		name string
		msg  discord.Message
	}{
		{
			name: "with embed",
			msg: discord.Message{
				ID:        123,
				ChannelID: 456,
				Author:    discord.User{ID: 789},
				Content:   "abc",
				Embeds: []discord.Embed{
					{
						Title:       "def",
						Description: "ghi",
					},
				},
			},
		},
		{
			name: "without embed",
			msg: discord.Message{
				ID:        123,
				ChannelID: 456,
				Author:    discord.User{ID: 789},
				Content:   "abc",
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)
				defer m.Eval()

				var embed *discord.Embed
				if len(c.msg.Embeds) > 0 {
					embed = &c.msg.Embeds[0]
				}

				m.SendMessage(embed, c.msg)

				actual, err := s.SendMessage(c.msg.ChannelID, c.msg.Content, embed)
				require.NoError(t, err)

				assert.Equal(t, c.msg, *actual)
			})
		}

		t.Run("param embed", func(t *testing.T) {
			m, s := NewSession(t)
			defer m.Eval()

			var (
				msg = discord.Message{
					ID:        123,
					ChannelID: 456,
					Author:    discord.User{ID: 789},
					Content:   "abc",
				}

				embed = discord.Embed{
					Title:       "def",
					Description: "ghi",
				}
			)

			require.NoError(t, embed.Validate())

			expect := msg
			expect.Embeds = append(expect.Embeds, embed)

			m.SendMessage(&embed, msg)

			actual, err := s.SendMessage(expect.ChannelID, expect.Content, &embed)
			require.NoError(t, err)

			assert.Equal(t, expect, *actual)
		})
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
			embed                       = discord.Embed{
				Title:       "def",
				Description: "ghi",
			}
		)

		m.SendMessage(&embed, discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds:    []discord.Embed{embed},
		})

		_, err := s.SendMessage(channelID, "cba", &embed)
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_EditText(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		expect := discord.Message{
			ID:        123,
			ChannelID: 456,
			Author:    discord.User{ID: 789},
			Content:   "abc",
		}

		m.EditText(expect)

		actual, err := s.EditText(expect.ChannelID, expect.ID, expect.Content)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
			embed                       = discord.Embed{
				Title:       "def",
				Description: "ghi",
			}
		)

		m.EditText(discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds:    []discord.Embed{embed},
		})

		_, err := s.EditText(channelID, messageID, "cba")
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_EditEmbed(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		expect := discord.Message{
			ID:        123,
			ChannelID: 456,
			Author:    discord.User{ID: 789},
			Embeds: []discord.Embed{
				{
					Title:       "def",
					Description: "ghi",
				},
			},
		}

		m.EditEmbed(expect)

		actual, err := s.EditEmbed(expect.ChannelID, expect.ID, expect.Embeds[0])
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
		)

		m.EditEmbed(discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds: []discord.Embed{
				{
					Title:       "def",
					Description: "ghi",
				},
			},
		})

		_, err := s.EditEmbed(channelID, messageID, discord.Embed{
			Title:       "fed",
			Description: "ihg",
		})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_EditMessage(t *testing.T) {
	successCases := []struct {
		name           string
		suppressEmbeds bool
		msg            discord.Message
	}{
		{
			name:           "suppressEmbeds",
			suppressEmbeds: true,
			msg: discord.Message{
				ID:        123,
				ChannelID: 456,
				Author:    discord.User{ID: 789},
				Content:   "abc",
			},
		},
		{
			name:           "don't suppressEmbeds",
			suppressEmbeds: false,
			msg: discord.Message{
				ID:        123,
				ChannelID: 456,
				Author:    discord.User{ID: 789},
				Content:   "abc",
			},
		},
		{
			name: "with embed",
			msg: discord.Message{
				ID:        123,
				ChannelID: 456,
				Author:    discord.User{ID: 789},
				Content:   "abc",
				Embeds: []discord.Embed{
					{
						Title:       "def",
						Description: "ghi",
					},
				},
			},
		},
		{
			name: "without embed",
			msg: discord.Message{
				ID:        123,
				ChannelID: 456,
				Author:    discord.User{ID: 789},
				Content:   "abc",
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)
				defer m.Eval()

				var embed *discord.Embed
				if len(c.msg.Embeds) > 0 {
					embed = &c.msg.Embeds[0]
				}

				m.EditMessage(embed, c.msg, c.suppressEmbeds)

				actual, err := s.EditMessage(c.msg.ChannelID, c.msg.ID, c.msg.Content, embed, c.suppressEmbeds)
				require.NoError(t, err)

				assert.Equal(t, c.msg, *actual)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
			embed                       = discord.Embed{
				Title:       "def",
				Description: "ghi",
			}
		)

		m.EditMessage(&embed, discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds:    []discord.Embed{embed},
		}, false)

		_, err := s.EditMessage(channelID, messageID, "cba", &embed, false)
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_EditMessageComplex(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		data := api.EditMessageData{
			Content: option.NewNullableString("abc"),
		}

		expect := discord.Message{
			ID:        123,
			ChannelID: 456,
			Author:    discord.User{ID: 789},
			Content:   "abc",
		}

		m.EditMessageComplex(data, expect)

		actual, err := s.EditMessageComplex(expect.ChannelID, expect.ID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := discord.Message{
			ID:        123,
			ChannelID: 456,
			Author:    discord.User{ID: 789},
			Content:   "abc",
		}

		m.EditMessageComplex(api.EditMessageData{
			Content: option.NewNullableString("abc"),
		}, expect)

		actual, err := s.EditMessageComplex(expect.ChannelID, expect.ID, api.EditMessageData{
			Content: option.NewNullableString("cba"),
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteMessage(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	var (
		channelID discord.ChannelID = 123
		messageID discord.MessageID = 456
	)

	m.DeleteMessage(channelID, messageID)

	err := s.DeleteMessage(channelID, messageID)
	require.NoError(t, err)
}

func TestMocker_DeleteMessages(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var (
			channelID  discord.ChannelID = 123
			messageIDs                   = []discord.MessageID{456, 789}
		)

		m.DeleteMessages(channelID, messageIDs)

		err := s.DeleteMessages(channelID, messageIDs)
		require.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var channelID discord.ChannelID = 123

		m.DeleteMessages(channelID, []discord.MessageID{456, 789})

		err := s.DeleteMessages(channelID, []discord.MessageID{654, 987})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}
