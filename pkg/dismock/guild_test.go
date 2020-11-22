package dismock

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_CreateGuild(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		data := api.CreateGuildData{Name: "abc"}

		expect := discord.Guild{
			ID:                     123,
			Name:                   "abc",
			OwnerID:                1,
			RulesChannelID:         2,
			PublicUpdatesChannelID: 3,
		}

		m.CreateGuild(data, expect)

		actual, err := s.CreateGuild(data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := discord.Guild{
			ID:                     123,
			Name:                   "abc",
			OwnerID:                1,
			RulesChannelID:         2,
			PublicUpdatesChannelID: 3,
		}

		m.CreateGuild(api.CreateGuildData{Name: "abc"}, expect)

		actual, err := s.CreateGuild(api.CreateGuildData{Name: "def"})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_Guild(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	expect := discord.Guild{
		ID:                     123,
		Name:                   "abc",
		OwnerID:                1,
		RulesChannelID:         2,
		PublicUpdatesChannelID: 3,
	}

	m.Guild(expect)

	actual, err := s.Guild(expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_GuildWithCount(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	expect := discord.Guild{
		ID:                     123,
		Name:                   "abc",
		OwnerID:                1,
		RulesChannelID:         2,
		PublicUpdatesChannelID: 3,
		ApproximateMembers:     3,
		ApproximatePresences:   2,
	}

	m.GuildWithCount(expect)

	actual, err := s.GuildWithCount(expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_GuildPreview(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	expect := discord.GuildPreview{
		ID:   123,
		Name: "abc",
	}

	m.GuildPreview(expect)

	actual, err := s.GuildPreview(expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

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
				defer m.Eval()

				expect := make([]discord.Guild, c.guilds)

				for i := 0; i < c.guilds; i++ {
					expect[i] = discord.Guild{
						ID:                     discord.GuildID(i + 1),
						OwnerID:                1,
						RulesChannelID:         2,
						PublicUpdatesChannelID: 3,
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
		defer m.Eval()

		m.Guilds(100, nil)

		actual, err := s.Guilds(100)
		require.NoError(t, err)

		assert.Len(t, actual, 0)
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
				defer m.Eval()

				var before discord.GuildID = 999999999999

				expect := make([]discord.Guild, c.guilds)

				for i := 0; i < c.guilds; i++ {
					expect[i] = discord.Guild{
						ID:                     discord.GuildID(i + 1),
						OwnerID:                1,
						RulesChannelID:         2,
						PublicUpdatesChannelID: 3,
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
		defer m.Eval()

		m.GuildsBefore(0, 100, nil)

		actual, err := s.GuildsBefore(0, 100)
		require.NoError(t, err)

		assert.Len(t, actual, 0)
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
			},
			{
				ID:                     456,
				Name:                   "def",
				OwnerID:                1,
				RulesChannelID:         2,
				PublicUpdatesChannelID: 3,
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
				defer m.Eval()

				var after discord.GuildID = 123

				expect := make([]discord.Guild, c.guilds)

				for i := 0; i < c.guilds; i++ {
					expect[i] = discord.Guild{
						ID:                     after + 1,
						OwnerID:                1,
						RulesChannelID:         2,
						PublicUpdatesChannelID: 3,
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
		defer m.Eval()

		m.GuildsAfter(0, 100, nil)

		actual, err := s.GuildsAfter(0, 100)
		require.NoError(t, err)

		assert.Len(t, actual, 0)
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
			},
			{
				ID:                     789,
				Name:                   "def",
				OwnerID:                1,
				RulesChannelID:         2,
				PublicUpdatesChannelID: 3,
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
	defer m.Eval()

	var guildID discord.GuildID = 123

	m.LeaveGuild(guildID)

	err := s.LeaveGuild(guildID)
	require.NoError(t, err)
}

func TestMocker_ModifyGuild(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		data := api.ModifyGuildData{Name: "abc"}

		expect := discord.Guild{
			ID:                     123,
			Name:                   "abc",
			OwnerID:                1,
			RulesChannelID:         2,
			PublicUpdatesChannelID: 3,
		}

		m.ModifyGuild(data, expect)

		actual, err := s.ModifyGuild(expect.ID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := discord.Guild{
			ID:                     123,
			Name:                   "abc",
			OwnerID:                1,
			RulesChannelID:         2,
			PublicUpdatesChannelID: 3,
		}

		m.ModifyGuild(api.ModifyGuildData{Name: "abc"}, expect)

		actual, err := s.ModifyGuild(expect.ID, api.ModifyGuildData{Name: "def"})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteGuild(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	var guildID discord.GuildID = 123

	m.DeleteGuild(guildID)

	err := s.DeleteGuild(guildID)
	require.NoError(t, err)
}

func TestMocker_VoiceRegionsGuild(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

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
	})

	t.Run("nil voice regions", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		//noinspection GoPreferNilSlice
		expect := []discord.VoiceRegion{}

		m.VoiceRegionsGuild(guildID, nil)

		actual, err := s.VoiceRegionsGuild(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})
}

func TestMocker_AuditLog(t *testing.T) {
	successCases := []struct {
		name string
		data api.AuditLogData
	}{
		{
			name: "limit = 0",
			data: api.AuditLogData{Limit: 0},
		},
		{
			name: "limit above 100",
			data: api.AuditLogData{Limit: 101},
		},
		{
			name: "no data",
			data: api.AuditLogData{Limit: 50},
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
				defer m.Eval()

				var guildID discord.GuildID = 123

				expect := discord.AuditLog{
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
				}

				m.AuditLog(guildID, c.data, expect)

				actual, err := s.AuditLog(guildID, c.data)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var guildID discord.GuildID = 123

		expect := discord.AuditLog{
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
		}

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
		defer m.Eval()

		var guildID discord.GuildID = 123

		expect := []discord.Integration{
			{
				ID:     456,
				Name:   "abc",
				RoleID: 1,
				User:   discord.User{ID: 2},
			},
			{
				ID:     789,
				Name:   "def",
				RoleID: 1,
				User:   discord.User{ID: 2},
			},
		}

		m.Integrations(guildID, expect)

		actual, err := s.Integrations(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})

	t.Run("nil integrations", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var guildID discord.GuildID = 123

		//noinspection GoPreferNilSlice
		expect := []discord.Integration{}

		m.Integrations(guildID, nil)

		actual, err := s.Integrations(guildID)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})
}

func TestMocker_AttachIntegration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var (
			guildID         discord.GuildID       = 123
			integrationID   discord.IntegrationID = 345
			integrationType                       = discord.Twitch
		)

		m.AttachIntegration(guildID, integrationID, integrationType)

		err := s.AttachIntegration(guildID, integrationID, integrationType)
		require.NoError(t, err)
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
		defer m.Eval()

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
	defer m.Eval()

	var (
		guildID       discord.GuildID       = 123
		integrationID discord.IntegrationID = 456
	)

	m.SyncIntegration(guildID, integrationID)

	err := s.SyncIntegration(guildID, integrationID)
	require.NoError(t, err)
}

func TestMocker_GuildWidgetSettings(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	var guildID discord.GuildID = 123

	expect := discord.GuildWidgetSettings{ChannelID: 345}

	m.GuildWidgetSettings(guildID, expect)

	actual, err := s.GuildWidgetSettings(guildID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_ModifyGuildWidget(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var (
			guildID discord.GuildID = 123
			data                    = api.ModifyGuildWidgetData{
				ChannelID: 345,
			}
			expect = discord.GuildWidgetSettings{
				Enabled:   false,
				ChannelID: 345,
			}
		)

		m.ModifyGuildWidget(guildID, data, expect)

		actual, err := s.ModifyGuildWidget(guildID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			guildID discord.GuildID = 123
			expect                  = discord.GuildWidgetSettings{
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

func TestMocker_GuildWidget(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	var guildID discord.GuildID = 123

	expect := discord.GuildWidget{
		ID:            123,
		Name:          "abc",
		InviteURL:     "def",
		Channels:      []discord.Channel{{ID: 456}},
		Members:       []discord.User{{ID: 789}},
		PresenceCount: 5,
	}

	m.GuildWidget(guildID, expect)

	actual, err := s.GuildWidget(guildID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_GuildVanityInvite(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	var guildID discord.GuildID = 123

	expect := discord.Invite{
		Code:           "abc",
		Channel:        discord.Channel{ID: 456},
		InviteMetadata: discord.InviteMetadata{Uses: 3},
	}

	m.GuildVanityInvite(guildID, expect)

	actual, err := s.GuildVanityInvite(guildID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_GuildWidgetImage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)
		defer m.Eval()

		var (
			guildID discord.GuildID = 123
			style                   = api.GuildBanner3
		)

		expect := []byte{1, 30, 0, 15, 24}

		reader := bytes.NewBuffer(expect)

		m.GuildWidgetImage(guildID, style, reader)

		actualReader, err := s.GuildWidgetImage(guildID, style)
		require.NoError(t, err)

		actual, err := ioutil.ReadAll(actualReader)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var guildID discord.GuildID = 123

		expect := []byte{1, 30, 0, 15, 24}

		reader := bytes.NewBuffer(expect)

		m.GuildWidgetImage(guildID, api.GuildBanner1, reader)

		actualReader, err := s.GuildWidgetImage(guildID, api.GuildBanner2)
		require.NoError(t, err)

		actual, err := ioutil.ReadAll(actualReader)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		assert.True(t, tMock.Failed())
	})
}
