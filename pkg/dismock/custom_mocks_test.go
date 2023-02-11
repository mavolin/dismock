package dismock

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_Error(t *testing.T) {
	m, s := NewSession(t)

	sendErr := httputil.HTTPError{
		Status:  http.StatusBadRequest,
		Code:    10011,
		Message: "Unknown guild",
	}

	m.Error(http.MethodGet, "guilds/123", sendErr)

	_, err := s.Guild(123)
	require.IsType(t, new(httputil.HTTPError), err)

	httpErr := err.(*httputil.HTTPError)

	assert.Equal(t, sendErr.Status, httpErr.Status)
	assert.Equal(t, sendErr.Code, httpErr.Code)
	assert.Equal(t, sendErr.Message, httpErr.Message)
}

// =============================================================================
// channel.go
// =====================================================================================

func TestMocker_Ack(t *testing.T) {
	abc := "abc"
	def := "def"
	ghi := "ghi"

	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
			ack                         = api.Ack{Token: &abc}
		)

		expect := api.Ack{Token: &def}
		actual := &ack

		m.Ack(channelID, messageID, ack, expect)

		err := s.Ack(channelID, messageID, actual)
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

		expect := api.Ack{Token: &def}

		m.Ack(channelID, messageID, api.Ack{Token: &abc}, expect)

		actual := &api.Ack{Token: &ghi}

		err := s.Ack(channelID, messageID, actual)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

// =============================================================================
// guild.go
// =====================================================================================

func TestMocker_Guilds(t *testing.T) {
	successCases := []struct {
		name   string
		guilds int
		limit  uint
	}{
		{
			name:   "limited",
			guilds: 130,
			limit:  199,
		},
		{
			name:   "unlimited",
			guilds: 200,
			limit:  0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				expect := make([]discord.Guild, c.guilds)

				for i := 0; i < c.guilds; i++ {
					expect[i] = discord.Guild{
						ID:                     discord.GuildID(i + 1),
						OwnerID:                1,
						RulesChannelID:         2,
						PublicUpdatesChannelID: 3,
						AFKTimeout:             1,
					}
				}

				m.Guilds(c.limit, expect)

				actual, err := s.Guilds(c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("nil guilds", func(t *testing.T) {
		m, s := NewSession(t)

		m.Guilds(100, nil)

		actual, err := s.Guilds(100)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m, _ := NewSession(t)

		assert.Panics(t, func() {
			m.Guilds(1, []discord.Guild{{}, {}})
		})
	})
}

func TestMocker_GuildsBefore(t *testing.T) {
	successCases := []struct {
		name   string
		guilds int
		limit  uint
	}{
		{
			name:   "limited",
			guilds: 130,
			limit:  199,
		},
		{
			name:   "unlimited",
			guilds: 200,
			limit:  0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var before discord.GuildID = 999999999999

				expect := make([]discord.Guild, c.guilds)

				for i := 0; i < c.guilds; i++ {
					expect[i] = discord.Guild{
						ID:                     discord.GuildID(i + 1),
						OwnerID:                1,
						RulesChannelID:         2,
						PublicUpdatesChannelID: 3,
						AFKTimeout:             1,
					}
				}

				m.GuildsBefore(before, c.limit, expect)

				actual, err := s.GuildsBefore(before, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("nil guilds", func(t *testing.T) {
		m, s := NewSession(t)

		m.GuildsBefore(0, 100, nil)

		actual, err := s.GuildsBefore(0, 100)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)
		m, s := NewSession(tMock)

		expect := []discord.Guild{
			{
				ID:                     123,
				Name:                   "abc",
				OwnerID:                1,
				RulesChannelID:         2,
				PublicUpdatesChannelID: 3,
				AFKTimeout:             1,
			},
			{
				ID:                     456,
				Name:                   "def",
				OwnerID:                1,
				RulesChannelID:         2,
				PublicUpdatesChannelID: 3,
				AFKTimeout:             1,
			},
		}

		m.GuildsBefore(890, 100, expect)

		actual, err := s.GuildsBefore(789, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m := New(t)

		assert.Panics(t, func() {
			m.GuildsBefore(0, 1, []discord.Guild{{}, {}})
		})
	})
}

func TestMocker_GuildsAfter(t *testing.T) {
	successCases := []struct {
		name   string
		guilds int
		limit  uint
	}{
		{
			name:   "limited",
			guilds: 130,
			limit:  199,
		},
		{
			name:   "unlimited",
			guilds: 200,
			limit:  0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var after discord.GuildID = 123

				expect := make([]discord.Guild, c.guilds)

				for i := 0; i < c.guilds; i++ {
					expect[i] = discord.Guild{
						ID:                     after + 1,
						OwnerID:                1,
						RulesChannelID:         2,
						PublicUpdatesChannelID: 3,
						AFKTimeout:             1,
					}
				}

				m.GuildsAfter(after, c.limit, expect)

				actual, err := s.GuildsAfter(after, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("nil guilds", func(t *testing.T) {
		m, s := NewSession(t)

		m.GuildsAfter(0, 100, nil)

		actual, err := s.GuildsAfter(0, 100)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)
		m, s := NewSession(tMock)

		expect := []discord.Guild{
			{
				ID:                     456,
				Name:                   "abc",
				OwnerID:                1,
				RulesChannelID:         2,
				PublicUpdatesChannelID: 3,
				AFKTimeout:             1,
			},
			{
				ID:                     789,
				Name:                   "def",
				OwnerID:                1,
				RulesChannelID:         2,
				PublicUpdatesChannelID: 3,
				AFKTimeout:             1,
			},
		}

		m.GuildsAfter(123, 100, expect)

		actual, err := s.GuildsAfter(321, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m := New(t)

		assert.Panics(t, func() {
			m.GuildsAfter(0, 1, []discord.Guild{{}, {}})
		})
	})
}

// =============================================================================
// interaction.go
// =====================================================================================

func TestMocker_RespondInteraction(t *testing.T) {
	successCases := []struct {
		name string
		resp api.InteractionResponse
	}{
		// {
		// 	name: "no files",
		// 	resp: api.InteractionResponse{
		// 		Type: api.MessageInteractionWithSource,
		// 		Data: &api.InteractionResponseData{Content: option.NewNullableString("abc")},
		// 	},
		// },
		{
			name: "with file",
			resp: api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Files: []sendpart.File{
						{Name: "abc", Reader: bytes.NewBufferString("def")},
					},
				},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				cp := c.resp
				dataCp := *c.resp.Data
				cp.Data = &dataCp

				if len(c.resp.Data.Files) > 0 {
					cp.Data.Files = make([]sendpart.File, len(c.resp.Data.Files))
					copy(cp.Data.Files, c.resp.Data.Files) // the readers of the file will be consumed twice

					// The files are copied now, but the reader for them may be a pointer and wasn't
					// deep copied. Therefore, we create two readers using the data from the original
					// reader.
					for i, f := range c.resp.Data.Files {
						b, err := io.ReadAll(f.Reader)
						require.NoError(t, err)

						cp.Data.Files[i].Reader = bytes.NewBuffer(b)
						c.resp.Data.Files[i].Reader = bytes.NewBuffer(b)
					}
				}

				var interactionID discord.InteractionID = 123
				token := "abc"

				m.RespondInteraction(interactionID, token, c.resp)

				err := s.RespondInteraction(interactionID, token, cp)
				fmt.Println(cp.Data.Files)
				require.NoError(t, err)
			})
		}
	})

	failureCases := []struct {
		name  string
		resp1 api.InteractionResponse
		resp2 api.InteractionResponse
	}{
		{
			name: "different content",
			resp1: api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{Content: option.NewNullableString("abc")},
			},
			resp2: api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{Content: option.NewNullableString("cba")},
			},
		},
		{
			name: "different file",
			resp1: api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Files: []sendpart.File{
						{Name: "abc", Reader: bytes.NewBufferString("def")},
					},
				},
			},
			resp2: api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Files: []sendpart.File{
						{Name: "abc", Reader: bytes.NewBufferString("fed")},
					},
				},
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				tMock := new(testing.T)

				m, s := NewSession(tMock)

				var interactionID discord.InteractionID = 123
				token := "abc"

				m.RespondInteraction(interactionID, token, c.resp1)

				err := s.RespondInteraction(interactionID, token, c.resp2)
				require.NoError(t, err)

				assert.True(t, tMock.Failed())
			})
		}
	})
}

// =============================================================================
// members.go
// =====================================================================================

func TestMocker_Members(t *testing.T) {
	successCases := []struct {
		name    string
		members int
		limit   uint
	}{
		{
			name:    "limited",
			members: 1003,
			limit:   2000,
		},
		{
			name:    "unlimited",
			members: 4004,
			limit:   0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var guildID discord.GuildID = 123

				expect := make([]discord.Member, c.members)

				for i := 0; i < c.members; i++ {
					expect[i] = discord.Member{
						User: discord.User{ID: discord.UserID(i + 1)},
					}
				}

				m.Members(guildID, c.limit, expect)

				actual, err := s.Members(guildID, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("nil members", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.GuildID = 123
		m.Members(guildID, 100, nil)

		actual, err := s.Members(guildID, 100)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})

	t.Run("limit smaller than members", func(t *testing.T) {
		m := New(t)

		assert.Panics(t, func() {
			m.Members(123, 1, []discord.Member{{}, {}})
		})
	})
}

func TestMocker_MembersAfter(t *testing.T) {
	successCases := []struct {
		name    string
		members int
		limit   uint
	}{
		{
			name:    "limited",
			members: 1003,
			limit:   2000,
		},
		{
			name:    "unlimited",
			members: 4004,
			limit:   0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var (
					guildID discord.GuildID = 123
					after   discord.UserID  = 456
				)

				expect := make([]discord.Member, c.members)

				for i := 0; i < c.members; i++ {
					expect[i] = discord.Member{
						User: discord.User{
							ID: discord.UserID(int(after) + i + 1),
						},
					}
				}

				m.MembersAfter(guildID, after, c.limit, expect)

				actual, err := s.MembersAfter(guildID, after, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("nil members", func(t *testing.T) {
		m, s := NewSession(t)
		var guildID discord.GuildID = 123

		m.MembersAfter(guildID, 0, 100, nil)

		actual, err := s.MembersAfter(guildID, 0, 100)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)
		m, s := NewSession(tMock)

		expect := []discord.Member{
			{User: discord.User{ID: 456}},
			{User: discord.User{ID: 789}},
		}

		m.MembersAfter(123, 456, 100, expect)

		actual, err := s.MembersAfter(123, 654, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m := New(t)

		assert.Panics(t, func() {
			m.MembersAfter(123, 0, 1, []discord.Member{{}, {}})
		})
	})
}

// =============================================================================
// messages.go
// =====================================================================================

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

		var channelID discord.ChannelID = 123
		m.Messages(channelID, 100, nil)

		actual, err := s.Messages(channelID, 100)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})

	t.Run("limit smaller than messages", func(t *testing.T) {
		m := New(t)

		assert.Panics(t, func() {
			m.Messages(123, 1, []discord.Message{{}, {}})
		})
	})
}

func TestMocker_MessagesAround(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

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

		var channelID discord.ChannelID = 123
		m.MessagesAround(channelID, 0, 100, nil)

		actual, err := s.MessagesAround(channelID, 0, 100)
		require.NoError(t, err)

		assert.Empty(t, actual)
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
		m := New(t)

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

		var channelID discord.ChannelID = 123

		m.MessagesBefore(channelID, 0, 100, nil)

		actual, err := s.MessagesBefore(channelID, 0, 100)
		require.NoError(t, err)

		assert.Empty(t, actual)
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
		m := New(t)

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
						Author:    discord.User{ID: 12},
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

		var (
			channelID discord.ChannelID = 123
			expect    []discord.Message
		)

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
		m := New(t)

		assert.Panics(t, func() {
			m.MessagesAfter(123, 0, 1, []discord.Message{{}, {}})
		})
	})
}

func TestMocker_SendTextReply(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		expect := discord.Message{
			ID:        123,
			ChannelID: 456,
			Author:    discord.User{ID: 789},
			Content:   "abc",
			Reference: &discord.MessageReference{MessageID: 12},
		}

		m.SendTextReply(expect)

		actual, err := s.SendTextReply(expect.ChannelID, expect.Content, expect.Reference.MessageID)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)
		m, s := NewSession(tMock)

		var channelID discord.ChannelID = 123

		m.SendTextReply(discord.Message{
			ChannelID: channelID,
			Content:   "abc",
			Reference: &discord.MessageReference{MessageID: 456},
		})

		_, err := s.SendTextReply(channelID, "cba", 456)
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_SendEmbeds(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		expect := discord.Message{
			ID:        123,
			ChannelID: 456,
			Author:    discord.User{ID: 789},
			Embeds:    []discord.Embed{{Title: "def", Description: "ghi"}},
		}

		m.SendEmbeds(expect)

		actual, err := s.SendEmbeds(expect.ChannelID, expect.Embeds[0])
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

		m.SendEmbeds(discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds:    []discord.Embed{{Title: "def", Description: "ghi"}},
		})

		_, err := s.SendEmbeds(channelID, discord.Embed{
			Title:       "fed",
			Description: "ihg",
		})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_SendEmbedReply(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		expect := discord.Message{
			ID:        123,
			ChannelID: 456,
			Author:    discord.User{ID: 789},
			Embeds:    []discord.Embed{{Title: "def", Description: "ghi"}},
			Reference: &discord.MessageReference{MessageID: 12},
		}

		m.SendEmbedReply(expect)

		actual, err := s.SendEmbedReply(expect.ChannelID, expect.Reference.MessageID, expect.Embeds...)
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

		m.SendEmbeds(discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds:    []discord.Embed{{Title: "def", Description: "ghi"}},
		})

		_, err := s.SendEmbedReply(channelID, 789, discord.Embed{
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

				m.SendMessage(c.msg)

				actual, err := s.SendMessage(c.msg.ChannelID, c.msg.Content, c.msg.Embeds...)
				require.NoError(t, err)

				assert.Equal(t, c.msg, *actual)
			})
		}

		t.Run("param embed", func(t *testing.T) {
			m, s := NewSession(t)

			expect := discord.Message{
				ID:        123,
				ChannelID: 456,
				Author:    discord.User{ID: 789},
				Content:   "abc",
				Embeds:    []discord.Embed{{Title: "def", Description: "ghi"}},
			}

			m.SendMessage(expect)

			actual, err := s.SendMessage(expect.ChannelID, expect.Content, expect.Embeds...)
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

		m.SendMessage(discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds:    []discord.Embed{embed},
		})

		_, err := s.SendMessage(channelID, "cba")
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_EditText(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

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

func TestMocker_EditEmbeds(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

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

		m.EditEmbeds(expect)

		actual, err := s.EditEmbeds(expect.ChannelID, expect.ID, expect.Embeds[0])
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

		m.EditEmbeds(discord.Message{
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

		_, err := s.EditEmbeds(channelID, messageID, discord.Embed{
			Title:       "fed",
			Description: "ihg",
		})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_EditMessage(t *testing.T) {
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

				m.EditMessage(c.msg.Content, c.msg.Embeds, c.msg)

				actual, err := s.EditMessage(c.msg.ChannelID, c.msg.ID, c.msg.Content, c.msg.Embeds...)
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
			content                     = "abc"
			embed                       = discord.Embed{
				Title:       "def",
				Description: "ghi",
			}
		)

		m.EditMessage(content, []discord.Embed{embed}, discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   content,
			Embeds:    []discord.Embed{embed},
		})

		_, err := s.EditMessage(channelID, messageID, "cba", embed)
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_EditMessageComplex(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

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

// =============================================================================
// message_reaction.go
// =====================================================================================

func TestMocker_Unreact(t *testing.T) {
	m, s := NewSession(t)

	var (
		channelID discord.ChannelID = 123
		messageID discord.MessageID = 456
		emoji     discord.APIEmoji  = "🍆"
	)

	m.Unreact(channelID, messageID, emoji)

	err := s.Unreact(channelID, messageID, emoji)
	require.NoError(t, err)
}

func TestMocker_Reactions(t *testing.T) {
	successCases := []struct {
		name      string
		reactions int
		limit     uint
	}{
		{
			name:      "limited",
			reactions: 130,
			limit:     199,
		},
		{
			name:      "unlimited",
			reactions: 200,
			limit:     0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var (
					channelID discord.ChannelID = 123
					messageID discord.MessageID = 456
					emoji     discord.APIEmoji  = "🍆"
				)

				expect := make([]discord.User, c.reactions)

				for i := 1; i < c.reactions+1; i++ {
					expect[i-1] = discord.User{ID: discord.UserID(i)}
				}

				m.Reactions(channelID, messageID, c.limit, emoji, expect)

				actual, err := s.Reactions(channelID, messageID, emoji, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("nil users", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
			emoji     discord.APIEmoji  = "🍆"
		)

		m.Reactions(channelID, messageID, 100, emoji, nil)

		actual, err := s.Reactions(channelID, messageID, emoji, 100)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m := New(t)

		assert.Panics(t, func() {
			m.Reactions(123, 456, 1, "abc", []discord.User{{}, {}})
		})
	})
}

func TestMocker_ReactionsBefore(t *testing.T) {
	successCases := []struct {
		name      string
		reactions int
		limit     uint
	}{
		{
			name:      "limited",
			reactions: 130,
			limit:     199,
		},
		{
			name:      "unlimited",
			reactions: 200,
			limit:     0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var (
					channelID discord.ChannelID = 123
					messageID discord.MessageID = 456
					emoji     discord.APIEmoji  = "🍆"

					before discord.UserID = 999999999999
				)

				expect := make([]discord.User, c.reactions)

				for i := 1; i < c.reactions+1; i++ {
					expect[i-1] = discord.User{ID: discord.UserID(i)}
				}

				m.ReactionsBefore(channelID, messageID, before, c.limit, emoji, expect)

				actual, err := s.ReactionsBefore(channelID, messageID, before, emoji, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("nil users", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
			emoji     discord.APIEmoji  = "🍆"
		)

		m.ReactionsBefore(channelID, messageID, 0, 100, emoji, nil)

		actual, err := s.ReactionsBefore(channelID, messageID, 0, emoji, 100)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)
		m, s := NewSession(tMock)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
			emoji     discord.APIEmoji  = "abc"
		)

		expect := []discord.User{{ID: 123}, {ID: 456}}

		m.ReactionsBefore(channelID, messageID, 890, 100, emoji, expect)

		actual, err := s.ReactionsBefore(channelID, messageID, 789, emoji, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m := New(t)

		assert.Panics(t, func() {
			m.ReactionsBefore(123, 456, 0, 1, "abc", []discord.User{{}, {}})
		})
	})
}

func TestMocker_ReactionsAfter(t *testing.T) {
	successCases := []struct {
		name      string
		reactions int
		limit     uint
	}{
		{
			name:      "limited",
			reactions: 130,
			limit:     199,
		},
		{
			name:      "unlimited",
			reactions: 200,
			limit:     0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var (
					channelID discord.ChannelID = 123
					messageID discord.MessageID = 456
					emoji     discord.APIEmoji  = "🍆"

					after discord.UserID = 123
				)

				expect := make([]discord.User, c.reactions)

				for i := 0; i < c.reactions; i++ {
					expect[i] = discord.User{ID: discord.UserID(int(after) + i + 1)}
				}

				m.ReactionsAfter(channelID, messageID, after, c.limit, emoji, expect)

				actual, err := s.ReactionsAfter(channelID, messageID, after, emoji, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("nil users", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
			emoji     discord.APIEmoji  = "🍆"
		)

		m.ReactionsAfter(channelID, messageID, 0, 100, emoji, nil)

		actual, err := s.ReactionsAfter(channelID, messageID, 0, emoji, 100)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)
		m, s := NewSession(tMock)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
			emoji     discord.APIEmoji  = "🍆"
		)

		expect := []discord.User{{ID: 456}, {ID: 789}}

		m.ReactionsAfter(channelID, messageID, 123, 100, emoji, expect)

		actual, err := s.ReactionsAfter(channelID, messageID, 321, emoji, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m := New(t)

		assert.Panics(t, func() {
			m.ReactionsAfter(123, 456, 0, 1, "abc", []discord.User{{}, {}})
		})
	})
}

func TestMocker_DeleteUserReaction(t *testing.T) {
	m, s := NewSession(t)

	var (
		channelID discord.ChannelID = 123
		messageID discord.MessageID = 456
		userID    discord.UserID    = 789
		emoji     discord.APIEmoji  = "🍆"
	)

	m.DeleteUserReaction(channelID, messageID, userID, emoji)

	err := s.DeleteUserReaction(channelID, messageID, userID, emoji)
	require.NoError(t, err)
}

// =============================================================================
// send.go
// =====================================================================================

func TestMocker_SendMessageComplex(t *testing.T) {
	successCases := []struct {
		name string
		data api.SendMessageData
	}{
		{name: "no files", data: api.SendMessageData{Content: "abc"}},
		{
			name: "with file",
			data: api.SendMessageData{
				Files: []sendpart.File{
					{Name: "abc", Reader: bytes.NewBufferString("def")},
				},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				expect := discord.Message{
					ID:        123,
					ChannelID: 456,
					Author:    discord.User{ID: 789},
				}

				cp := c.data

				cp.Files = make([]sendpart.File, len(c.data.Files))
				copy(cp.Files, c.data.Files) // the readers of the file will be consumed twice

				// the files are copied now, but the reader for them may be a pointer and wasn't
				// deep copied. therefore we create two readers using the data from the original
				// reader
				for i, f := range c.data.Files {
					b, err := io.ReadAll(f.Reader)
					require.NoError(t, err)

					cp.Files[i].Reader = bytes.NewBuffer(b)
					c.data.Files[i].Reader = bytes.NewBuffer(b)
				}

				m.SendMessageComplex(c.data, expect)

				actual, err := s.SendMessageComplex(expect.ChannelID, cp)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)
			})
		}
	})

	failureCases := []struct {
		name  string
		data1 api.SendMessageData
		data2 api.SendMessageData
	}{
		{
			name:  "different content",
			data1: api.SendMessageData{Content: "abc"},
			data2: api.SendMessageData{Content: "cba"},
		},
		{
			name: "different file",
			data1: api.SendMessageData{
				Files: []sendpart.File{
					{Name: "abc", Reader: bytes.NewBufferString("def")},
				},
			},
			data2: api.SendMessageData{
				Files: []sendpart.File{
					{Name: "abc", Reader: bytes.NewBufferString("fed")},
				},
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				tMock := new(testing.T)

				m, s := NewSession(tMock)

				expect := discord.Message{
					ID:        123,
					ChannelID: 456,
					Author:    discord.User{ID: 789},
				}

				m.SendMessageComplex(c.data1, expect)

				actual, err := s.SendMessageComplex(expect.ChannelID, c.data2)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)
				assert.True(t, tMock.Failed())
			})
		}
	})
}

func TestMocker_ExecuteWebhook(t *testing.T) {
	successCases := []struct {
		name string
		data webhook.ExecuteData
	}{
		{name: "no files", data: webhook.ExecuteData{Content: "abc"}},
		{
			name: "with file",
			data: webhook.ExecuteData{
				Files: []sendpart.File{
					{Name: "abc", Reader: bytes.NewBufferString("def")},
				},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m := New(t)

				var (
					webhookID discord.WebhookID = 123
					token                       = "abc"
				)

				cp := c.data

				cp.Files = make([]sendpart.File, len(c.data.Files))
				copy(cp.Files, c.data.Files) // the readers of the file will be consumed twice

				// the files are copied now, but the reader for them may be a pointer and wasn't
				// deep copied. therefore we create two readers using the data from the original
				// reader
				for i, f := range c.data.Files {
					b, err := io.ReadAll(f.Reader)
					require.NoError(t, err)

					cp.Files[i].Reader = bytes.NewBuffer(b)
					c.data.Files[i].Reader = bytes.NewBuffer(b)
				}

				m.ExecuteWebhook(webhookID, token, c.data)

				err := webhook.NewCustom(webhookID, token, m.HTTPClient()).Execute(cp)
				require.NoError(t, err)
			})
		}
	})

	failureCases := []struct {
		name  string
		data1 webhook.ExecuteData
		data2 webhook.ExecuteData
	}{
		{
			name:  "different content",
			data1: webhook.ExecuteData{Content: "abc"},
			data2: webhook.ExecuteData{Content: "cba"},
		},
		{
			name: "different file",
			data1: webhook.ExecuteData{
				Files: []sendpart.File{
					{Name: "abc", Reader: bytes.NewBufferString("def")},
				},
			},
			data2: webhook.ExecuteData{
				Files: []sendpart.File{
					{Name: "abc", Reader: bytes.NewBufferString("fed")},
				},
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				tMock := new(testing.T)
				m := New(tMock)

				var (
					webhookID discord.WebhookID = 123
					token                       = "abc"
				)

				m.ExecuteWebhook(webhookID, token, c.data1)

				err := webhook.NewCustom(webhookID, token, m.HTTPClient()).Execute(c.data2)
				require.NoError(t, err)

				assert.True(t, tMock.Failed())
			})
		}
	})
}

func TestMocker_ExecuteWebhookAndWait(t *testing.T) {
	successCases := []struct {
		name string
		data webhook.ExecuteData
	}{
		{
			name: "no files",
			data: webhook.ExecuteData{Content: "abc"},
		},
		{
			name: "with file",
			data: webhook.ExecuteData{
				Files: []sendpart.File{
					{Name: "abc", Reader: bytes.NewBufferString("def")},
				},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m := New(t)

				var (
					webhookID discord.WebhookID = 123
					token                       = "abc"
				)

				expect := discord.Message{
					ID:        123,
					ChannelID: 456,
					Author:    discord.User{ID: 789},
				}

				cp := c.data

				cp.Files = make([]sendpart.File, len(c.data.Files))
				copy(cp.Files, c.data.Files) // the readers of the file will be consumed twice

				// the files are copied now, but the reader for them may be a pointer and wasn't
				// deep copied. therefore we create two readers using the data from the original
				// reader
				for i, f := range c.data.Files {
					b, err := io.ReadAll(f.Reader)
					require.NoError(t, err)

					cp.Files[i].Reader = bytes.NewBuffer(b)
					c.data.Files[i].Reader = bytes.NewBuffer(b)
				}

				m.ExecuteWebhookAndWait(webhookID, token, c.data, expect)

				actual, err := webhook.NewCustom(webhookID, token, m.HTTPClient()).ExecuteAndWait(cp)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)
			})
		}
	})

	failureCases := []struct {
		name  string
		data1 webhook.ExecuteData
		data2 webhook.ExecuteData
	}{
		{
			name:  "different content",
			data1: webhook.ExecuteData{Content: "abc"},
			data2: webhook.ExecuteData{Content: "cba"},
		},
		{
			name: "different file",
			data1: webhook.ExecuteData{
				Files: []sendpart.File{
					{Name: "abc", Reader: bytes.NewBufferString("def")},
				},
			},
			data2: webhook.ExecuteData{
				Files: []sendpart.File{
					{Name: "abc", Reader: bytes.NewBufferString("fed")},
				},
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				tMock := new(testing.T)

				m := New(tMock)

				var (
					webhookID discord.WebhookID = 123
					token                       = "abc"
				)

				expect := discord.Message{
					ID:        123,
					ChannelID: 456,
					Author:    discord.User{ID: 789},
				}

				m.ExecuteWebhookAndWait(webhookID, token, c.data1, expect)

				actual, err := webhook.NewCustom(webhookID, token, m.HTTPClient()).ExecuteAndWait(c.data2)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)
				assert.True(t, tMock.Failed())
			})
		}
	})
}

// =============================================================================
// webhook/webhook.go
// =====================================================================================

func TestMocker_WebhookWithToken(t *testing.T) {
	m := New(t)

	expect := discord.Webhook{
		ID:            123,
		ChannelID:     456,
		User:          &discord.User{ID: 789},
		Token:         "abc",
		ApplicationID: 1,
	}

	m.WebhookWithToken(expect)

	actual, err := webhook.NewCustom(expect.ID, expect.Token, m.HTTPClient()).Get()
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_ModifyWebhookWithToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := New(t)

		data := api.ModifyWebhookData{Name: option.NewString("abc")}

		expect := discord.Webhook{
			ID:            123,
			Name:          "abc",
			Token:         "def",
			ChannelID:     456,
			User:          &discord.User{ID: 789},
			ApplicationID: 1,
		}

		m.ModifyWebhookWithToken(data, expect)

		actual, err := webhook.NewCustom(expect.ID, expect.Token, m.HTTPClient()).Modify(data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m := New(tMock)

		expect := discord.Webhook{
			ID:            123,
			Name:          "abc",
			Token:         "def",
			ChannelID:     456,
			User:          &discord.User{ID: 789},
			ApplicationID: 1,
		}

		m.ModifyWebhookWithToken(api.ModifyWebhookData{
			Name: option.NewString("abc"),
		}, expect)

		actual, err := webhook.NewCustom(expect.ID, expect.Token, m.HTTPClient()).Modify(api.ModifyWebhookData{
			Name: option.NewString("cba"),
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteWebhookWithToken(t *testing.T) {
	m := New(t)

	var (
		id    discord.WebhookID = 123
		token                   = "abc"
	)

	m.DeleteWebhookWithToken(id, token)

	err := webhook.NewCustom(id, token, m.HTTPClient()).Delete()
	require.NoError(t, err)
}
