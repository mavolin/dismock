package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_Member(t *testing.T) {
	m, s := NewArikawaSession(t)

	var guildID discord.Snowflake = 123

	expect := discord.Member{
		User: discord.User{
			ID: 456,
		},
	}

	m.Member(guildID, expect)

	actual, err := s.Member(guildID, expect.User.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_Members(t *testing.T) {
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
				m, s := NewArikawaSession(t)

				var guildID discord.Snowflake = 123

				expect := []discord.Member{ // more than 100 entries so multiple requests are mocked
					{User: discord.User{ID: 1234567890}}, {User: discord.User{ID: 2345678901}},
					{User: discord.User{ID: 3456789012}}, {User: discord.User{ID: 4567890123}},
					{User: discord.User{ID: 5678901234}}, {User: discord.User{ID: 6789012345}},
					{User: discord.User{ID: 7890123456}}, {User: discord.User{ID: 8901234567}},
					{User: discord.User{ID: 9012345678}}, {User: discord.User{ID: 123456789}},
					{User: discord.User{ID: 234567890}}, {User: discord.User{ID: 345678901}},
					{User: discord.User{ID: 456789012}}, {User: discord.User{ID: 567890123}},
					{User: discord.User{ID: 678901234}}, {User: discord.User{ID: 789012345}},
					{User: discord.User{ID: 890123456}}, {User: discord.User{ID: 901234567}},
					{User: discord.User{ID: 12345678}}, {User: discord.User{ID: 23456789}},
					{User: discord.User{ID: 34567890}}, {User: discord.User{ID: 45678901}},
					{User: discord.User{ID: 56789012}}, {User: discord.User{ID: 67890123}},
					{User: discord.User{ID: 78901234}}, {User: discord.User{ID: 89012345}},
					{User: discord.User{ID: 90123456}}, {User: discord.User{ID: 1234567}},
					{User: discord.User{ID: 2345678}}, {User: discord.User{ID: 3456789}},
					{User: discord.User{ID: 4567890}}, {User: discord.User{ID: 5678901}},
					{User: discord.User{ID: 6789012}}, {User: discord.User{ID: 7890123}},
					{User: discord.User{ID: 8901234}}, {User: discord.User{ID: 9012345}},
					{User: discord.User{ID: 123456}}, {User: discord.User{ID: 234567}},
					{User: discord.User{ID: 345678}}, {User: discord.User{ID: 456789}},
					{User: discord.User{ID: 567890}}, {User: discord.User{ID: 678901}},
					{User: discord.User{ID: 789012}}, {User: discord.User{ID: 890123}},
					{User: discord.User{ID: 901234}}, {User: discord.User{ID: 12345}}, {User: discord.User{ID: 23456}},
					{User: discord.User{ID: 34567}}, {User: discord.User{ID: 45678}}, {User: discord.User{ID: 56789}},
					{User: discord.User{ID: 67890}}, {User: discord.User{ID: 78901}}, {User: discord.User{ID: 89012}},
					{User: discord.User{ID: 90123}}, {User: discord.User{ID: 1234}}, {User: discord.User{ID: 2345}},
					{User: discord.User{ID: 3456}}, {User: discord.User{ID: 4567}}, {User: discord.User{ID: 5678}},
					{User: discord.User{ID: 6789}}, {User: discord.User{ID: 7890}}, {User: discord.User{ID: 8901}},
					{User: discord.User{ID: 9012}}, {User: discord.User{ID: 123}}, {User: discord.User{ID: 234}},
					{User: discord.User{ID: 345}}, {User: discord.User{ID: 456}}, {User: discord.User{ID: 567}},
					{User: discord.User{ID: 678}}, {User: discord.User{ID: 789}}, {User: discord.User{ID: 890}},
					{User: discord.User{ID: 901}}, {User: discord.User{ID: 12}}, {User: discord.User{ID: 23}},
					{User: discord.User{ID: 45}}, {User: discord.User{ID: 56}}, {User: discord.User{ID: 67}},
					{User: discord.User{ID: 78}}, {User: discord.User{ID: 89}}, {User: discord.User{ID: 90}},
					{User: discord.User{ID: 98}}, {User: discord.User{ID: 87}}, {User: discord.User{ID: 76}},
					{User: discord.User{ID: 65}}, {User: discord.User{ID: 54}}, {User: discord.User{ID: 43}},
					{User: discord.User{ID: 32}}, {User: discord.User{ID: 21}}, {User: discord.User{ID: 10}},
					{User: discord.User{ID: 987}}, {User: discord.User{ID: 876}}, {User: discord.User{ID: 765}},
					{User: discord.User{ID: 654}}, {User: discord.User{ID: 543}}, {User: discord.User{ID: 432}},
					{User: discord.User{ID: 321}}, {User: discord.User{ID: 210}}, {User: discord.User{ID: 109}},
					{User: discord.User{ID: 9876}}, {User: discord.User{ID: 8765}}, {User: discord.User{ID: 7654}},
					{User: discord.User{ID: 6543}}, {User: discord.User{ID: 5432}}, {User: discord.User{ID: 4321}},
					{User: discord.User{ID: 3210}}, {User: discord.User{ID: 2109}}, {User: discord.User{ID: 1098}},
					{User: discord.User{ID: 98765}}, {User: discord.User{ID: 87654}}, {User: discord.User{ID: 76543}},
					{User: discord.User{ID: 65432}}, {User: discord.User{ID: 54321}}, {User: discord.User{ID: 43210}},
					{User: discord.User{ID: 32109}}, {User: discord.User{ID: 21098}}, {User: discord.User{ID: 10987}},
					{User: discord.User{ID: 987654}}, {User: discord.User{ID: 876543}},
					{User: discord.User{ID: 765432}}, {User: discord.User{ID: 654321}},
					{User: discord.User{ID: 543210}}, {User: discord.User{ID: 432109}},
					{User: discord.User{ID: 321098}}, {User: discord.User{ID: 210987}},
					{User: discord.User{ID: 109876}}, {User: discord.User{ID: 9876543}},
					{User: discord.User{ID: 8765432}}, {User: discord.User{ID: 7654321}},
					{User: discord.User{ID: 6543210}}, {User: discord.User{ID: 5432109}},
					{User: discord.User{ID: 4321098}}, {User: discord.User{ID: 3210987}},
					{User: discord.User{ID: 2109876}}, {User: discord.User{ID: 1098765}},
					{User: discord.User{ID: 98765432}}, {User: discord.User{ID: 87654321}},
					{User: discord.User{ID: 76543210}}, {User: discord.User{ID: 65432109}},
					{User: discord.User{ID: 54321098}}, {User: discord.User{ID: 43210987}},
					{User: discord.User{ID: 32109876}}, {User: discord.User{ID: 21098765}},
					{User: discord.User{ID: 10987654}}, {User: discord.User{ID: 987654321}},
					{User: discord.User{ID: 876543210}}, {User: discord.User{ID: 765432109}},
					{User: discord.User{ID: 654321098}}, {User: discord.User{ID: 543210987}},
					{User: discord.User{ID: 432109876}}, {User: discord.User{ID: 321098765}},
					{User: discord.User{ID: 210987654}}, {User: discord.User{ID: 109876543}},
					{User: discord.User{ID: 9876543210}}, {User: discord.User{ID: 8765432109}},
					{User: discord.User{ID: 7654321098}}, {User: discord.User{ID: 6543210987}},
					{User: discord.User{ID: 5432109876}}, {User: discord.User{ID: 4321098765}},
					{User: discord.User{ID: 3210987654}}, {User: discord.User{ID: 2109876543}},
					{User: discord.User{ID: 1098765432}},
				}

				m.Members(guildID, c.limit, expect)

				actual, err := s.Members(guildID, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("limit smaller than members", func(t *testing.T) {
		m, _ := NewArikawaSession(t)

		assert.Panics(t, func() {
			m.Members(123, 1, []discord.Member{{}, {}})
		})
	})
}

func TestMocker_MembersAfter(t *testing.T) {
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
				m, s := NewArikawaSession(t)

				var (
					guildID discord.Snowflake = 123
					after   discord.Snowflake = 456
				)

				expect := []discord.Member{ // more than 100 entries so multiple requests are mocked
					{User: discord.User{ID: 1234567890}}, {User: discord.User{ID: 2345678901}},
					{User: discord.User{ID: 3456789012}}, {User: discord.User{ID: 4567890123}},
					{User: discord.User{ID: 5678901234}}, {User: discord.User{ID: 6789012345}},
					{User: discord.User{ID: 7890123456}}, {User: discord.User{ID: 8901234567}},
					{User: discord.User{ID: 9012345678}}, {User: discord.User{ID: 123456789}},
					{User: discord.User{ID: 234567890}}, {User: discord.User{ID: 345678901}},
					{User: discord.User{ID: 456789012}}, {User: discord.User{ID: 567890123}},
					{User: discord.User{ID: 678901234}}, {User: discord.User{ID: 789012345}},
					{User: discord.User{ID: 890123456}}, {User: discord.User{ID: 901234567}},
					{User: discord.User{ID: 12345678}}, {User: discord.User{ID: 23456789}},
					{User: discord.User{ID: 34567890}}, {User: discord.User{ID: 45678901}},
					{User: discord.User{ID: 56789012}}, {User: discord.User{ID: 67890123}},
					{User: discord.User{ID: 78901234}}, {User: discord.User{ID: 89012345}},
					{User: discord.User{ID: 90123456}}, {User: discord.User{ID: 1234567}},
					{User: discord.User{ID: 2345678}}, {User: discord.User{ID: 3456789}},
					{User: discord.User{ID: 4567890}}, {User: discord.User{ID: 5678901}},
					{User: discord.User{ID: 6789012}}, {User: discord.User{ID: 7890123}},
					{User: discord.User{ID: 8901234}}, {User: discord.User{ID: 9012345}},
					{User: discord.User{ID: 123456}}, {User: discord.User{ID: 234567}},
					{User: discord.User{ID: 345678}}, {User: discord.User{ID: 456789}},
					{User: discord.User{ID: 567890}}, {User: discord.User{ID: 678901}},
					{User: discord.User{ID: 789012}}, {User: discord.User{ID: 890123}},
					{User: discord.User{ID: 901234}}, {User: discord.User{ID: 12345}}, {User: discord.User{ID: 23456}},
					{User: discord.User{ID: 34567}}, {User: discord.User{ID: 45678}}, {User: discord.User{ID: 56789}},
					{User: discord.User{ID: 67890}}, {User: discord.User{ID: 78901}}, {User: discord.User{ID: 89012}},
					{User: discord.User{ID: 90123}}, {User: discord.User{ID: 1234}}, {User: discord.User{ID: 2345}},
					{User: discord.User{ID: 3456}}, {User: discord.User{ID: 4567}}, {User: discord.User{ID: 5678}},
					{User: discord.User{ID: 6789}}, {User: discord.User{ID: 7890}}, {User: discord.User{ID: 8901}},
					{User: discord.User{ID: 9012}}, {User: discord.User{ID: 123}}, {User: discord.User{ID: 234}},
					{User: discord.User{ID: 345}}, {User: discord.User{ID: 456}}, {User: discord.User{ID: 567}},
					{User: discord.User{ID: 678}}, {User: discord.User{ID: 789}}, {User: discord.User{ID: 890}},
					{User: discord.User{ID: 901}}, {User: discord.User{ID: 12}}, {User: discord.User{ID: 23}},
					{User: discord.User{ID: 45}}, {User: discord.User{ID: 56}}, {User: discord.User{ID: 67}},
					{User: discord.User{ID: 78}}, {User: discord.User{ID: 89}}, {User: discord.User{ID: 90}},
					{User: discord.User{ID: 98}}, {User: discord.User{ID: 87}}, {User: discord.User{ID: 76}},
					{User: discord.User{ID: 65}}, {User: discord.User{ID: 54}}, {User: discord.User{ID: 43}},
					{User: discord.User{ID: 32}}, {User: discord.User{ID: 21}}, {User: discord.User{ID: 10}},
					{User: discord.User{ID: 987}}, {User: discord.User{ID: 876}}, {User: discord.User{ID: 765}},
					{User: discord.User{ID: 654}}, {User: discord.User{ID: 543}}, {User: discord.User{ID: 432}},
					{User: discord.User{ID: 321}}, {User: discord.User{ID: 210}}, {User: discord.User{ID: 109}},
					{User: discord.User{ID: 9876}}, {User: discord.User{ID: 8765}}, {User: discord.User{ID: 7654}},
					{User: discord.User{ID: 6543}}, {User: discord.User{ID: 5432}}, {User: discord.User{ID: 4321}},
					{User: discord.User{ID: 3210}}, {User: discord.User{ID: 2109}}, {User: discord.User{ID: 1098}},
					{User: discord.User{ID: 98765}}, {User: discord.User{ID: 87654}}, {User: discord.User{ID: 76543}},
					{User: discord.User{ID: 65432}}, {User: discord.User{ID: 54321}}, {User: discord.User{ID: 43210}},
					{User: discord.User{ID: 32109}}, {User: discord.User{ID: 21098}}, {User: discord.User{ID: 10987}},
					{User: discord.User{ID: 987654}}, {User: discord.User{ID: 876543}},
					{User: discord.User{ID: 765432}}, {User: discord.User{ID: 654321}},
					{User: discord.User{ID: 543210}}, {User: discord.User{ID: 432109}},
					{User: discord.User{ID: 321098}}, {User: discord.User{ID: 210987}},
					{User: discord.User{ID: 109876}}, {User: discord.User{ID: 9876543}},
					{User: discord.User{ID: 8765432}}, {User: discord.User{ID: 7654321}},
					{User: discord.User{ID: 6543210}}, {User: discord.User{ID: 5432109}},
					{User: discord.User{ID: 4321098}}, {User: discord.User{ID: 3210987}},
					{User: discord.User{ID: 2109876}}, {User: discord.User{ID: 1098765}},
					{User: discord.User{ID: 98765432}}, {User: discord.User{ID: 87654321}},
					{User: discord.User{ID: 76543210}}, {User: discord.User{ID: 65432109}},
					{User: discord.User{ID: 54321098}}, {User: discord.User{ID: 43210987}},
					{User: discord.User{ID: 32109876}}, {User: discord.User{ID: 21098765}},
					{User: discord.User{ID: 10987654}}, {User: discord.User{ID: 987654321}},
					{User: discord.User{ID: 876543210}}, {User: discord.User{ID: 765432109}},
					{User: discord.User{ID: 654321098}}, {User: discord.User{ID: 543210987}},
					{User: discord.User{ID: 432109876}}, {User: discord.User{ID: 321098765}},
					{User: discord.User{ID: 210987654}}, {User: discord.User{ID: 109876543}},
					{User: discord.User{ID: 9876543210}}, {User: discord.User{ID: 8765432109}},
					{User: discord.User{ID: 7654321098}}, {User: discord.User{ID: 6543210987}},
					{User: discord.User{ID: 5432109876}}, {User: discord.User{ID: 4321098765}},
					{User: discord.User{ID: 3210987654}}, {User: discord.User{ID: 2109876543}},
					{User: discord.User{ID: 1098765432}},
				}

				m.MembersAfter(guildID, after, c.limit, expect)

				actual, err := s.MembersAfter(guildID, after, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewArikawaSession(tMock)

		expect := []discord.Member{
			{
				User: discord.User{
					ID: 456,
				},
			},
			{
				User: discord.User{
					ID: 789,
				},
			},
		}

		m.MembersAfter(123, 456, 100, expect)

		actual, err := s.MembersAfter(123, 654, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m, _ := NewArikawaSession(t)

		assert.Panics(t, func() {
			m.MembersAfter(123, 0, 1, []discord.Member{{}, {}})
		})
	})
}

func TestMocker_AddMember(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewArikawaSession(t)

		var (
			guildID discord.Snowflake = 123

			data = api.AddMemberData{
				Token: "abc",
			}
		)

		expect := discord.Member{
			User: discord.User{
				ID: 345,
			},
		}

		m.AddMember(guildID, data, expect)

		actual, err := s.AddMember(guildID, expect.User.ID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewArikawaSession(tMock)

		var guildID discord.Snowflake = 123

		expect := discord.Member{
			User: discord.User{
				ID: 345,
			},
		}

		m.AddMember(guildID, api.AddMemberData{
			Token: "abc",
		}, expect)

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
		m, s := NewArikawaSession(t)

		var (
			guildID discord.Snowflake = 123
			userID  discord.Snowflake = 456

			data = api.ModifyMemberData{
				Nick: option.NewString("abc"),
			}
		)

		m.ModifyMember(guildID, userID, data)

		err := s.ModifyMember(guildID, userID, data)
		require.NoError(t, err)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewArikawaSession(tMock)

		var (
			guildID discord.Snowflake = 123
			userID  discord.Snowflake = 465
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
			data: api.PruneCountData{
				Days: 0,
			},
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
				IncludedRoles: []discord.Snowflake{},
			},
		},
		{
			name: "roles",
			data: api.PruneCountData{
				Days:          5,
				IncludedRoles: []discord.Snowflake{123, 456},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewArikawaSession(t)

				var guildID discord.Snowflake = 123

				var expect uint = 25

				m.PruneCount(guildID, c.data, expect)

				actual, err := s.PruneCount(guildID, c.data)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewArikawaSession(tMock)

		var guildID discord.Snowflake = 123

		var expect uint = 25

		m.PruneCount(guildID, api.PruneCountData{
			Days:          5,
			IncludedRoles: []discord.Snowflake{},
		}, expect)

		actual, err := s.PruneCount(guildID, api.PruneCountData{
			Days:          7,
			IncludedRoles: []discord.Snowflake{},
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
			data: api.PruneData{
				Days: 0,
			},
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
				IncludedRoles: []discord.Snowflake{},
			},
		},
		{
			name: "roles",
			data: api.PruneData{
				Days:          5,
				IncludedRoles: []discord.Snowflake{123, 456},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewArikawaSession(t)

				var guildID discord.Snowflake = 123

				var expect uint = 25

				m.Prune(guildID, c.data, expect)

				actual, err := s.Prune(guildID, c.data)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewArikawaSession(tMock)

		var guildID discord.Snowflake = 123

		var expect uint = 25

		m.Prune(guildID, api.PruneData{
			Days:          5,
			ReturnCount:   true,
			IncludedRoles: []discord.Snowflake{},
		}, expect)

		actual, err := s.Prune(guildID, api.PruneData{
			Days:          7,
			ReturnCount:   true,
			IncludedRoles: []discord.Snowflake{},
		})
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_Kick(t *testing.T) {
	m, s := NewArikawaSession(t)

	var (
		guildID discord.Snowflake = 123
		userID  discord.Snowflake = 456
	)

	m.Kick(guildID, userID)

	err := s.Kick(guildID, userID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_Bans(t *testing.T) {
	m, s := NewArikawaSession(t)

	var guildID discord.Snowflake = 123

	expect := []discord.Ban{
		{
			User: discord.User{
				ID: 123,
			},
		},
		{
			User: discord.User{
				ID: 456,
			},
		},
	}

	m.Bans(guildID, expect)

	actual, err := s.Bans(guildID)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_GetBan(t *testing.T) {
	m, s := NewArikawaSession(t)

	var guildID discord.Snowflake = 123

	expect := discord.Ban{
		User: discord.User{
			ID: 123,
		},
	}

	m.GetBan(guildID, expect)

	actual, err := s.GetBan(guildID, expect.User.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_Ban(t *testing.T) {
	successCases := []struct {
		name string
		data api.BanData
	}{
		{
			name: "deleteDays > 7",
			data: api.BanData{
				DeleteDays: option.NewUint(8),
			},
		},
		{
			name: "deleteDays",
			data: api.BanData{
				DeleteDays: option.NewUint(5),
			},
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
				m, s := NewArikawaSession(t)

				var (
					guildID discord.Snowflake = 123
					userID  discord.Snowflake = 456
				)

				m.Ban(guildID, userID, c.data)

				err := s.Ban(guildID, userID, c.data)
				require.NoError(t, err)

				m.Eval()
			})
		}
	})
}

func TestMocker_Unban(t *testing.T) {
	m, s := NewArikawaSession(t)

	var (
		guildID discord.Snowflake = 123
		userID  discord.Snowflake = 456
	)

	m.Unban(guildID, userID)

	err := s.Unban(guildID, userID)
	require.NoError(t, err)

	m.Eval()
}
