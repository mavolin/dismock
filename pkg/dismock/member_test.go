package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_Member(t *testing.T) {
	m, s := NewSession(t)

	var guildID discord.GuildID = 123

	expect := discord.Member{User: discord.User{ID: 456}}

	m.Member(guildID, expect)

	actual, err := s.Member(guildID, expect.User.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

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

func TestMocker_AddMember(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			guildID discord.GuildID = 123
			data                    = api.AddMemberData{Token: "abc"}
		)

		expect := discord.Member{User: discord.User{ID: 345}}

		m.AddMember(guildID, data, expect)

		actual, err := s.AddMember(guildID, expect.User.ID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)
		m, s := NewSession(tMock)

		var guildID discord.GuildID = 123

		expect := discord.Member{User: discord.User{ID: 345}}

		m.AddMember(guildID, api.AddMemberData{Token: "abc"}, expect)

		actual, err := s.AddMember(guildID, expect.User.ID, api.AddMemberData{
			Token: "cba",
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_ModifyMember(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			guildID discord.GuildID = 123
			userID  discord.UserID  = 456

			data = api.ModifyMemberData{Nick: option.NewString("abc")}
		)

		m.ModifyMember(guildID, userID, data)

		err := s.ModifyMember(guildID, userID, data)
		require.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)
		m, s := NewSession(tMock)

		var (
			guildID discord.GuildID = 123
			userID  discord.UserID  = 465
		)

		m.ModifyMember(guildID, userID, api.ModifyMemberData{
			Nick: option.NewString("abc"),
		})

		err := s.ModifyMember(guildID, userID, api.ModifyMemberData{
			Nick: option.NewString("cba"),
		})
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_PruneCount(t *testing.T) {
	successCases := []struct {
		name string
		data api.PruneCountData
	}{
		{
			name: "days = 0",
			data: api.PruneCountData{Days: 0},
		},
		{
			name: "nil roles",
			data: api.PruneCountData{
				Days:          5,
				IncludedRoles: nil,
			},
		},
		{
			name: "empty roles",
			data: api.PruneCountData{
				Days:          5,
				IncludedRoles: []discord.RoleID{},
			},
		},
		{
			name: "roles",
			data: api.PruneCountData{
				Days:          5,
				IncludedRoles: []discord.RoleID{123, 456},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var (
					guildID discord.GuildID = 123
					expect  uint            = 25
				)

				m.PruneCount(guildID, c.data, expect)

				actual, err := s.PruneCount(guildID, c.data)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)
		m, s := NewSession(tMock)

		var (
			guildID discord.GuildID = 123
			expect  uint            = 25
		)

		m.PruneCount(guildID, api.PruneCountData{
			Days:          5,
			IncludedRoles: []discord.RoleID{},
		}, expect)

		actual, err := s.PruneCount(guildID, api.PruneCountData{
			Days:          7,
			IncludedRoles: []discord.RoleID{},
		})
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_Prune(t *testing.T) {
	successCases := []struct {
		name string
		data api.PruneData
	}{
		{
			name: "days = 0",
			data: api.PruneData{Days: 0},
		},
		{
			name: "nil roles",
			data: api.PruneData{
				Days:          5,
				IncludedRoles: nil,
			},
		},
		{
			name: "empty roles",
			data: api.PruneData{
				Days:          5,
				IncludedRoles: []discord.RoleID{},
			},
		},
		{
			name: "roles",
			data: api.PruneData{
				Days:          5,
				IncludedRoles: []discord.RoleID{123, 456},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var (
					guildID discord.GuildID = 123
					expect  uint            = 25
				)

				m.Prune(guildID, c.data, expect)

				actual, err := s.Prune(guildID, c.data)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)
		m, s := NewSession(tMock)

		var (
			guildID discord.GuildID = 123
			expect  uint            = 25
		)

		m.Prune(guildID, api.PruneData{
			Days:          5,
			ReturnCount:   true,
			IncludedRoles: []discord.RoleID{},
		}, expect)

		actual, err := s.Prune(guildID, api.PruneData{
			Days:          7,
			ReturnCount:   true,
			IncludedRoles: []discord.RoleID{},
		})
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_Kick(t *testing.T) {
	m, s := NewSession(t)

	var (
		guildID discord.GuildID = 123
		userID  discord.UserID  = 456
	)

	m.Kick(guildID, userID)

	err := s.Kick(guildID, userID)
	require.NoError(t, err)
}

func TestMocker_Bans(t *testing.T) {
	m, s := NewSession(t)

	var guildID discord.GuildID = 123

	expect := []discord.Ban{
		{User: discord.User{ID: 123}},
		{User: discord.User{ID: 456}},
	}

	m.Bans(guildID, expect)

	actual, err := s.Bans(guildID)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)
}

func TestMocker_GetBan(t *testing.T) {
	m, s := NewSession(t)

	var guildID discord.GuildID = 123

	expect := discord.Ban{User: discord.User{ID: 123}}

	m.GetBan(guildID, expect)

	actual, err := s.GetBan(guildID, expect.User.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)
}

func TestMocker_Ban(t *testing.T) {
	successCases := []struct {
		name string
		data api.BanData
	}{
		{
			name: "deleteDays > 7",
			data: api.BanData{DeleteDays: option.NewUint(8)},
		},
		{
			name: "deleteDays",
			data: api.BanData{DeleteDays: option.NewUint(5)},
		},
		{
			name: "reason",
			data: api.BanData{
				DeleteDays: option.NewUint(8),
				Reason:     option.NewString("abc"),
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				var (
					guildID discord.GuildID = 123
					userID  discord.UserID  = 456
				)

				m.Ban(guildID, userID, c.data)

				err := s.Ban(guildID, userID, c.data)
				require.NoError(t, err)
			})
		}
	})
}

func TestMocker_Unban(t *testing.T) {
	m, s := NewSession(t)

	var (
		guildID discord.GuildID = 123
		userID  discord.UserID  = 456
	)

	m.Unban(guildID, userID)

	err := s.Unban(guildID, userID)
	require.NoError(t, err)
}
