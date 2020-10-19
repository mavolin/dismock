package dismock

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/dismock/internal/sanitize"
)

func TestMocker_CreateGuild(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		data := api.CreateGuildData{
			Name: "abc",
		}

		expect := sanitize.Guild(discord.Guild{
			ID:   123,
			Name: "abc",
		}, 1, 1, 1, 1)

		m.CreateGuild(data, expect)

		actual, err := s.CreateGuild(data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := sanitize.Guild(discord.Guild{
			ID:   123,
			Name: "abc",
		}, 1, 1, 1, 1)

		m.CreateGuild(api.CreateGuildData{
			Name: "abc",
		}, expect)

		actual, err := s.CreateGuild(api.CreateGuildData{
			Name: "def",
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_Guild(t *testing.T) {
	m, s := NewSession(t)

	expect := sanitize.Guild(discord.Guild{
		ID:   123,
		Name: "abc",
	}, 1, 1, 1, 1)

	m.Guild(expect)

	actual, err := s.Guild(expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_GuildWithCount(t *testing.T) {
	m, s := NewSession(t)

	expect := sanitize.Guild(discord.Guild{
		ID:                   123,
		Name:                 "abc",
		ApproximateMembers:   3,
		ApproximatePresences: 2,
	}, 1, 1, 1, 1)

	m.GuildWithCount(expect)

	actual, err := s.GuildWithCount(expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_GuildPreview(t *testing.T) {
	m, s := NewSession(t)

	expect := sanitize.GuildPreview(discord.GuildPreview{
		ID:   123,
		Name: "abc",
	}, 1)

	m.GuildPreview(expect)

	actual, err := s.GuildPreview(expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_Guilds(t *testing.T) {
	successCases := []struct {
		name  string
		limit uint
	}{
		{
			name:  "limited",
			limit: 199,
		},
		{
			name:  "unlimited",
			limit: 0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				expect := []discord.Guild{ // more than 100 entries so multiple requests are mocked
					{ID: 1234567890}, {ID: 2345678901}, {ID: 3456789012},
					{ID: 4567890123}, {ID: 5678901234}, {ID: 6789012345}, {ID: 7890123456}, {ID: 8901234567},
					{ID: 9012345678}, {ID: 123456789}, {ID: 234567890}, {ID: 345678901}, {ID: 456789012},
					{ID: 567890123}, {ID: 678901234}, {ID: 789012345}, {ID: 890123456}, {ID: 901234567},
					{ID: 12345678}, {ID: 23456789}, {ID: 34567890}, {ID: 45678901}, {ID: 56789012},
					{ID: 67890123}, {ID: 78901234}, {ID: 89012345}, {ID: 90123456}, {ID: 1234567},
					{ID: 2345678}, {ID: 3456789}, {ID: 4567890}, {ID: 5678901}, {ID: 6789012},
					{ID: 7890123}, {ID: 8901234}, {ID: 9012345}, {ID: 123456}, {ID: 234567},
					{ID: 345678}, {ID: 456789}, {ID: 567890}, {ID: 678901}, {ID: 789012},
					{ID: 890123}, {ID: 901234}, {ID: 12345}, {ID: 23456}, {ID: 34567},
					{ID: 45678}, {ID: 56789}, {ID: 67890}, {ID: 78901}, {ID: 89012},
					{ID: 90123}, {ID: 1234}, {ID: 2345}, {ID: 3456}, {ID: 4567},
					{ID: 5678}, {ID: 6789}, {ID: 7890}, {ID: 8901}, {ID: 9012},
					{ID: 123}, {ID: 234}, {ID: 345}, {ID: 456}, {ID: 567},
					{ID: 678}, {ID: 789}, {ID: 890}, {ID: 901}, {ID: 12},
					{ID: 23}, {ID: 45}, {ID: 56}, {ID: 67}, {ID: 78},
					{ID: 89}, {ID: 90}, {ID: 98}, {ID: 87}, {ID: 76},
					{ID: 65}, {ID: 54}, {ID: 43}, {ID: 32}, {ID: 21},
					{ID: 10}, {ID: 987}, {ID: 876}, {ID: 765}, {ID: 654},
					{ID: 543}, {ID: 432}, {ID: 321}, {ID: 210}, {ID: 109},
					{ID: 9876}, {ID: 8765}, {ID: 7654}, {ID: 6543}, {ID: 5432},
					{ID: 4321}, {ID: 3210}, {ID: 2109}, {ID: 1098}, {ID: 98765},
					{ID: 87654}, {ID: 76543}, {ID: 65432}, {ID: 54321}, {ID: 43210},
					{ID: 32109}, {ID: 21098}, {ID: 10987}, {ID: 987654}, {ID: 876543},
					{ID: 765432}, {ID: 654321}, {ID: 543210}, {ID: 432109}, {ID: 321098},
					{ID: 210987}, {ID: 109876}, {ID: 9876543}, {ID: 8765432}, {ID: 7654321},
					{ID: 6543210}, {ID: 5432109}, {ID: 4321098}, {ID: 3210987}, {ID: 2109876},
					{ID: 1098765}, {ID: 98765432}, {ID: 87654321}, {ID: 76543210}, {ID: 65432109},
					{ID: 54321098}, {ID: 43210987}, {ID: 32109876}, {ID: 21098765}, {ID: 10987654},
					{ID: 987654321}, {ID: 876543210}, {ID: 765432109}, {ID: 654321098}, {ID: 543210987},
					{ID: 432109876}, {ID: 321098765}, {ID: 210987654}, {ID: 109876543}, {ID: 9876543210},
					{ID: 8765432109}, {ID: 7654321098}, {ID: 6543210987}, {ID: 5432109876}, {ID: 4321098765},
					{ID: 3210987654}, {ID: 2109876543}, {ID: 1098765432},
				}

				for i, g := range expect {
					expect[i] = sanitize.Guild(g, 1, 1, 1, 1)
				}

				m.Guilds(c.limit, expect)

				actual, err := s.Guilds(c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("nil guilds", func(t *testing.T) {
		m, s := NewSession(t)

		m.Guilds(100, nil)

		actual, err := s.Guilds(100)
		require.NoError(t, err)

		assert.Len(t, actual, 0)

		m.Eval()
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
		name  string
		limit uint
	}{
		{
			name:  "limited",
			limit: 199,
		},
		{
			name:  "unlimited",
			limit: 0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var before discord.GuildID = 999999999999

				expect := []discord.Guild{ // more than 100 entries so multiple requests are mocked
					{ID: 1234567890}, {ID: 2345678901}, {ID: 3456789012},
					{ID: 4567890123}, {ID: 5678901234}, {ID: 6789012345}, {ID: 7890123456}, {ID: 8901234567},
					{ID: 9012345678}, {ID: 123456789}, {ID: 234567890}, {ID: 345678901}, {ID: 456789012},
					{ID: 567890123}, {ID: 678901234}, {ID: 789012345}, {ID: 890123456}, {ID: 901234567},
					{ID: 12345678}, {ID: 23456789}, {ID: 34567890}, {ID: 45678901}, {ID: 56789012},
					{ID: 67890123}, {ID: 78901234}, {ID: 89012345}, {ID: 90123456}, {ID: 1234567},
					{ID: 2345678}, {ID: 3456789}, {ID: 4567890}, {ID: 5678901}, {ID: 6789012},
					{ID: 7890123}, {ID: 8901234}, {ID: 9012345}, {ID: 123456}, {ID: 234567},
					{ID: 345678}, {ID: 456789}, {ID: 567890}, {ID: 678901}, {ID: 789012},
					{ID: 890123}, {ID: 901234}, {ID: 12345}, {ID: 23456}, {ID: 34567},
					{ID: 45678}, {ID: 56789}, {ID: 67890}, {ID: 78901}, {ID: 89012},
					{ID: 90123}, {ID: 1234}, {ID: 2345}, {ID: 3456}, {ID: 4567},
					{ID: 5678}, {ID: 6789}, {ID: 7890}, {ID: 8901}, {ID: 9012},
					{ID: 123}, {ID: 234}, {ID: 345}, {ID: 456}, {ID: 567},
					{ID: 678}, {ID: 789}, {ID: 890}, {ID: 901}, {ID: 12},
					{ID: 23}, {ID: 45}, {ID: 56}, {ID: 67}, {ID: 78},
					{ID: 89}, {ID: 90}, {ID: 98}, {ID: 87}, {ID: 76},
					{ID: 65}, {ID: 54}, {ID: 43}, {ID: 32}, {ID: 21},
					{ID: 10}, {ID: 987}, {ID: 876}, {ID: 765}, {ID: 654},
					{ID: 543}, {ID: 432}, {ID: 321}, {ID: 210}, {ID: 109},
					{ID: 9876}, {ID: 8765}, {ID: 7654}, {ID: 6543}, {ID: 5432},
					{ID: 4321}, {ID: 3210}, {ID: 2109}, {ID: 1098}, {ID: 98765},
					{ID: 87654}, {ID: 76543}, {ID: 65432}, {ID: 54321}, {ID: 43210},
					{ID: 32109}, {ID: 21098}, {ID: 10987}, {ID: 987654}, {ID: 876543},
					{ID: 765432}, {ID: 654321}, {ID: 543210}, {ID: 432109}, {ID: 321098},
					{ID: 210987}, {ID: 109876}, {ID: 9876543}, {ID: 8765432}, {ID: 7654321},
					{ID: 6543210}, {ID: 5432109}, {ID: 4321098}, {ID: 3210987}, {ID: 2109876},
					{ID: 1098765}, {ID: 98765432}, {ID: 87654321}, {ID: 76543210}, {ID: 65432109},
					{ID: 54321098}, {ID: 43210987}, {ID: 32109876}, {ID: 21098765}, {ID: 10987654},
					{ID: 987654321}, {ID: 876543210}, {ID: 765432109}, {ID: 654321098}, {ID: 543210987},
					{ID: 432109876}, {ID: 321098765}, {ID: 210987654}, {ID: 109876543}, {ID: 9876543210},
					{ID: 8765432109}, {ID: 7654321098}, {ID: 6543210987}, {ID: 5432109876}, {ID: 4321098765},
					{ID: 3210987654}, {ID: 2109876543}, {ID: 1098765432},
				}

				for i, g := range expect {
					expect[i] = sanitize.Guild(g, 1, 1, 1, 1)
				}

				m.GuildsBefore(before, c.limit, expect)

				actual, err := s.GuildsBefore(before, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("nil guilds", func(t *testing.T) {
		m, s := NewSession(t)

		m.GuildsBefore(0, 100, nil)

		actual, err := s.GuildsBefore(0, 100)
		require.NoError(t, err)

		assert.Len(t, actual, 0)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := []discord.Guild{
			{
				ID:   123,
				Name: "abc",
			},
			{
				ID:   456,
				Name: "def",
			},
		}

		m.GuildsBefore(890, 100, expect)

		actual, err := s.GuildsBefore(789, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m, _ := NewSession(t)

		assert.Panics(t, func() {
			m.GuildsBefore(0, 1, []discord.Guild{{}, {}})
		})
	})
}

func TestMocker_GuildsAfter(t *testing.T) {
	successCases := []struct {
		name  string
		limit uint
	}{
		{
			name:  "limited",
			limit: 199,
		},
		{
			name:  "unlimited",
			limit: 0,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var after discord.GuildID = 123

				expect := []discord.Guild{ // more than 100 entries so multiple requests are mocked
					{ID: 1234567890}, {ID: 2345678901}, {ID: 3456789012},
					{ID: 4567890123}, {ID: 5678901234}, {ID: 6789012345}, {ID: 7890123456}, {ID: 8901234567},
					{ID: 9012345678}, {ID: 123456789}, {ID: 234567890}, {ID: 345678901}, {ID: 456789012},
					{ID: 567890123}, {ID: 678901234}, {ID: 789012345}, {ID: 890123456}, {ID: 901234567},
					{ID: 12345678}, {ID: 23456789}, {ID: 34567890}, {ID: 45678901}, {ID: 56789012},
					{ID: 67890123}, {ID: 78901234}, {ID: 89012345}, {ID: 90123456}, {ID: 1234567},
					{ID: 2345678}, {ID: 3456789}, {ID: 4567890}, {ID: 5678901}, {ID: 6789012},
					{ID: 7890123}, {ID: 8901234}, {ID: 9012345}, {ID: 123456}, {ID: 234567},
					{ID: 345678}, {ID: 456789}, {ID: 567890}, {ID: 678901}, {ID: 789012},
					{ID: 890123}, {ID: 901234}, {ID: 12345}, {ID: 23456}, {ID: 34567},
					{ID: 45678}, {ID: 56789}, {ID: 67890}, {ID: 78901}, {ID: 89012},
					{ID: 90123}, {ID: 1234}, {ID: 2345}, {ID: 3456}, {ID: 4567},
					{ID: 5678}, {ID: 6789}, {ID: 7890}, {ID: 8901}, {ID: 9012},
					{ID: 123}, {ID: 234}, {ID: 345}, {ID: 456}, {ID: 567},
					{ID: 678}, {ID: 789}, {ID: 890}, {ID: 901}, {ID: 12},
					{ID: 23}, {ID: 45}, {ID: 56}, {ID: 67}, {ID: 78},
					{ID: 89}, {ID: 90}, {ID: 98}, {ID: 87}, {ID: 76},
					{ID: 65}, {ID: 54}, {ID: 43}, {ID: 32}, {ID: 21},
					{ID: 10}, {ID: 987}, {ID: 876}, {ID: 765}, {ID: 654},
					{ID: 543}, {ID: 432}, {ID: 321}, {ID: 210}, {ID: 109},
					{ID: 9876}, {ID: 8765}, {ID: 7654}, {ID: 6543}, {ID: 5432},
					{ID: 4321}, {ID: 3210}, {ID: 2109}, {ID: 1098}, {ID: 98765},
					{ID: 87654}, {ID: 76543}, {ID: 65432}, {ID: 54321}, {ID: 43210},
					{ID: 32109}, {ID: 21098}, {ID: 10987}, {ID: 987654}, {ID: 876543},
					{ID: 765432}, {ID: 654321}, {ID: 543210}, {ID: 432109}, {ID: 321098},
					{ID: 210987}, {ID: 109876}, {ID: 9876543}, {ID: 8765432}, {ID: 7654321},
					{ID: 6543210}, {ID: 5432109}, {ID: 4321098}, {ID: 3210987}, {ID: 2109876},
					{ID: 1098765}, {ID: 98765432}, {ID: 87654321}, {ID: 76543210}, {ID: 65432109},
					{ID: 54321098}, {ID: 43210987}, {ID: 32109876}, {ID: 21098765}, {ID: 10987654},
					{ID: 987654321}, {ID: 876543210}, {ID: 765432109}, {ID: 654321098}, {ID: 543210987},
					{ID: 432109876}, {ID: 321098765}, {ID: 210987654}, {ID: 109876543}, {ID: 9876543210},
					{ID: 8765432109}, {ID: 7654321098}, {ID: 6543210987}, {ID: 5432109876}, {ID: 4321098765},
					{ID: 3210987654}, {ID: 2109876543}, {ID: 1098765432},
				}

				for i, g := range expect {
					expect[i] = sanitize.Guild(g, 1, 1, 1, 1)
				}

				m.GuildsAfter(after, c.limit, expect)

				actual, err := s.GuildsAfter(after, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("nil guilds", func(t *testing.T) {
		m, s := NewSession(t)

		m.GuildsAfter(0, 100, nil)

		actual, err := s.GuildsAfter(0, 100)
		require.NoError(t, err)

		assert.Len(t, actual, 0)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := []discord.Guild{
			{
				ID:   456,
				Name: "abc",
			},
			{
				ID:   789,
				Name: "def",
			},
		}

		m.GuildsAfter(123, 100, expect)

		actual, err := s.GuildsAfter(321, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m, _ := NewSession(t)

		assert.Panics(t, func() {
			m.GuildsAfter(0, 1, []discord.Guild{{}, {}})
		})
	})
}

func TestMocker_LeaveGuild(t *testing.T) {
	m, s := NewSession(t)

	var guildID discord.GuildID = 123

	m.LeaveGuild(guildID)

	err := s.LeaveGuild(guildID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_ModifyGuild(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		data := api.ModifyGuildData{
			Name: "abc",
		}

		expect := sanitize.Guild(discord.Guild{
			ID:   123,
			Name: "abc",
		}, 1, 1, 1, 1)

		m.ModifyGuild(data, expect)

		actual, err := s.ModifyGuild(expect.ID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := sanitize.Guild(discord.Guild{
			ID:   123,
			Name: "abc",
		}, 1, 1, 1, 1)

		m.ModifyGuild(api.ModifyGuildData{
			Name: "abc",
		}, expect)

		actual, err := s.ModifyGuild(expect.ID, api.ModifyGuildData{
			Name: "def",
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteGuild(t *testing.T) {
	m, s := NewSession(t)

	var guildID discord.GuildID = 123

	m.DeleteGuild(guildID)

	err := s.DeleteGuild(guildID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_VoiceRegionsGuild(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.GuildID = 123

		expect := []discord.VoiceRegion{
			{
				ID:   "abc",
				Name: "ABC",
			},
			{
				ID:   "def",
				Name: "DEF",
			},
		}

		m.VoiceRegionsGuild(guildID, expect)

		actual, err := s.VoiceRegionsGuild(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("nil voice regions", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.GuildID = 123

		//noinspection GoPreferNilSlice
		expect := []discord.VoiceRegion{}

		m.VoiceRegionsGuild(guildID, nil)

		actual, err := s.VoiceRegionsGuild(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})
}

func TestMocker_AuditLog(t *testing.T) {
	successCases := []struct {
		name string
		data api.AuditLogData
	}{
		{
			name: "limit = 0",
			data: api.AuditLogData{
				Limit: 0,
			},
		},
		{
			name: "limit above 100",
			data: api.AuditLogData{
				Limit: 101,
			},
		},
		{
			name: "no data",
			data: api.AuditLogData{
				Limit: 50,
			},
		},
		{
			name: "UserID",
			data: api.AuditLogData{
				Limit:  50,
				UserID: 123,
			},
		},
		{
			name: "ActionType",
			data: api.AuditLogData{
				Limit:      50,
				ActionType: discord.EmojiDelete,
			},
		},
		{
			name: "Before",
			data: api.AuditLogData{
				Limit:  50,
				Before: 123,
			},
		},
		{
			name: "multi",
			data: api.AuditLogData{
				Limit:      50,
				UserID:     123,
				ActionType: discord.EmojiDelete,
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var guildID discord.GuildID = 123

				expect := sanitize.AuditLog(discord.AuditLog{
					Users: []discord.User{
						{
							ID:       256827968133791744,
							Username: "Ragnar89",
						},
					},
					Entries: []discord.AuditLogEntry{
						{
							ID:         456,
							UserID:     256827968133791744,
							ActionType: discord.EmojiUpdate,
						},
					},
				})

				m.AuditLog(guildID, c.data, expect)

				actual, err := s.AuditLog(guildID, c.data)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)

				m.Eval()
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var guildID discord.GuildID = 123

		expect := sanitize.AuditLog(discord.AuditLog{
			Users: []discord.User{
				{
					ID:       256827968133791744,
					Username: "Ragnar89",
				},
			},
			Entries: []discord.AuditLogEntry{
				{
					ID:         456,
					UserID:     256827968133791744,
					ActionType: discord.EmojiUpdate,
				},
			},
		})

		m.AuditLog(guildID, api.AuditLogData{
			ActionType: discord.EmojiUpdate,
			Limit:      100,
		}, expect)

		actual, err := s.AuditLog(guildID, api.AuditLogData{
			ActionType: discord.EmojiDelete,
			Limit:      100,
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_Integrations(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.GuildID = 123

		expect := []discord.Integration{
			{
				ID:   456,
				Name: "abc",
			},
			{
				ID:   789,
				Name: "def",
			},
		}

		for i, integration := range expect {
			expect[i] = sanitize.Integration(integration, 1, 1, 1)
		}

		m.Integrations(guildID, expect)

		actual, err := s.Integrations(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("nil integrations", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.GuildID = 123

		//noinspection GoPreferNilSlice
		expect := []discord.Integration{}

		m.Integrations(guildID, nil)

		actual, err := s.Integrations(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})
}

func TestMocker_AttachIntegration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			guildID         discord.GuildID       = 123
			integrationID   discord.IntegrationID = 345
			integrationType                       = discord.Twitch
		)

		m.AttachIntegration(guildID, integrationID, integrationType)

		err := s.AttachIntegration(guildID, integrationID, integrationType)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var guildID discord.GuildID = 123

		m.AttachIntegration(guildID, 345, discord.Twitch)

		err := s.AttachIntegration(guildID, 345, discord.YouTube)
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_ModifyIntegration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			guildID       discord.GuildID       = 123
			integrationID discord.IntegrationID = 345
			data                                = api.ModifyIntegrationData{
				ExpireGracePeriod: option.NullInt,
			}
		)

		m.ModifyIntegration(guildID, integrationID, data)

		err := s.ModifyIntegration(guildID, integrationID, data)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			guildID       discord.GuildID       = 123
			integrationID discord.IntegrationID = 456
		)

		m.ModifyIntegration(guildID, integrationID, api.ModifyIntegrationData{
			ExpireGracePeriod: option.NullInt,
		})

		err := s.ModifyIntegration(guildID, integrationID, api.ModifyIntegrationData{
			ExpireGracePeriod: option.NewNullableInt(12),
		})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_SyncIntegration(t *testing.T) {
	m, s := NewSession(t)

	var (
		guildID       discord.GuildID       = 123
		integrationID discord.IntegrationID = 456
	)

	m.SyncIntegration(guildID, integrationID)

	err := s.SyncIntegration(guildID, integrationID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_GuildWidget(t *testing.T) {
	m, s := NewSession(t)

	var guildID discord.GuildID = 123

	expect := discord.GuildWidget{
		ChannelID: 345,
	}

	m.GuildWidget(guildID, expect)

	actual, err := s.GuildWidget(guildID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_ModifyGuildWidget(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			guildID discord.GuildID = 123
			data                    = api.ModifyGuildWidgetData{
				ChannelID: 345,
			}
			expect = discord.GuildWidget{
				Enabled:   false,
				ChannelID: 345,
			}
		)

		m.ModifyGuildWidget(guildID, data, expect)

		actual, err := s.ModifyGuildWidget(guildID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			guildID discord.GuildID = 123
			expect                  = discord.GuildWidget{
				Enabled:   false,
				ChannelID: 345,
			}
		)

		m.ModifyGuildWidget(guildID, api.ModifyGuildWidgetData{
			ChannelID: 456,
		}, expect)

		actual, err := s.ModifyGuildWidget(guildID, api.ModifyGuildWidgetData{
			ChannelID: 789,
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_GuildVanityURL(t *testing.T) {
	m, s := NewSession(t)

	var guildID discord.GuildID = 123

	expect := sanitize.Invite(discord.Invite{
		Code: "abc",
		InviteMetadata: discord.InviteMetadata{
			Uses: 3,
		},
	}, 1, 1, 1, 1, 1, 1, 1)

	m.GuildVanityURL(guildID, expect)

	actual, err := s.GuildVanityURL(guildID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_GuildImage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			guildID discord.GuildID = 123
			style                   = api.GuildBanner3
		)

		expect := []byte{1, 30, 0, 15, 24}

		reader := bytes.NewBuffer(expect)

		m.GuildImage(guildID, style, reader)

		actualReader, err := s.GuildImage(guildID, style)
		require.NoError(t, err)

		actual, err := ioutil.ReadAll(actualReader)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var guildID discord.GuildID = 123

		expect := []byte{1, 30, 0, 15, 24}

		reader := bytes.NewBuffer(expect)

		m.GuildImage(guildID, api.GuildBanner1, reader)

		actualReader, err := s.GuildImage(guildID, api.GuildBanner2)
		require.NoError(t, err)

		actual, err := ioutil.ReadAll(actualReader)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		assert.True(t, tMock.Failed())
	})
}
