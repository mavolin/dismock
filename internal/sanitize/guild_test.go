package sanitize

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
)

func TestGuild(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.Guild
		expect discord.Guild
	}{
		{
			name: "none",
			in: discord.Guild{
				ID:                     321,
				OwnerID:                654,
				RulesChannelID:         987,
				PublicUpdatesChannelID: 210,
			},
			expect: discord.Guild{
				ID:                     321,
				OwnerID:                654,
				RulesChannelID:         987,
				PublicUpdatesChannelID: 210,
			},
		},
		{
			name: "id",
			in: discord.Guild{
				OwnerID:                654,
				RulesChannelID:         987,
				PublicUpdatesChannelID: 210,
			},
			expect: discord.Guild{
				ID:                     123,
				OwnerID:                654,
				RulesChannelID:         987,
				PublicUpdatesChannelID: 210,
			},
		},
		{
			name: "ownerID",
			in: discord.Guild{
				ID:                     321,
				RulesChannelID:         987,
				PublicUpdatesChannelID: 210,
			},
			expect: discord.Guild{
				ID:                     321,
				OwnerID:                456,
				RulesChannelID:         987,
				PublicUpdatesChannelID: 210,
			},
		},
		{
			name: "ruleChannelID",
			in: discord.Guild{
				ID:                     321,
				OwnerID:                654,
				PublicUpdatesChannelID: 210,
			},
			expect: discord.Guild{
				ID:                     321,
				OwnerID:                654,
				RulesChannelID:         789,
				PublicUpdatesChannelID: 210,
			},
		},
		{
			name: "publicUpdatesChannelID",
			in: discord.Guild{
				ID:             321,
				OwnerID:        654,
				RulesChannelID: 987,
			},
			expect: discord.Guild{
				ID:                     321,
				OwnerID:                654,
				RulesChannelID:         987,
				PublicUpdatesChannelID: 012,
			},
		},
		{
			name: "emojis",
			in: discord.Guild{
				ID:                     321,
				OwnerID:                654,
				Emojis:                 []discord.Emoji{{}},
				RulesChannelID:         987,
				PublicUpdatesChannelID: 210,
			},
			expect: discord.Guild{
				ID:      321,
				OwnerID: 654,
				Emojis: []discord.Emoji{
					{
						ID: 1,
						User: discord.User{
							ID: 1,
						},
					},
				},
				RulesChannelID:         987,
				PublicUpdatesChannelID: 210,
			},
		},
		{
			name: "roles",
			in: discord.Guild{
				ID:                     321,
				OwnerID:                654,
				Roles:                  []discord.Role{{}},
				RulesChannelID:         987,
				PublicUpdatesChannelID: 210,
			},
			expect: discord.Guild{
				ID:      321,
				OwnerID: 654,
				Roles: []discord.Role{
					{
						ID: 1,
					},
				},
				RulesChannelID:         987,
				PublicUpdatesChannelID: 210,
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Guild(c.in, 123, 456, 789, 012)

			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestGuildPreview(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.GuildPreview
		expect discord.GuildPreview
	}{
		{
			name: "none",
			in: discord.GuildPreview{
				ID: 321,
			},
			expect: discord.GuildPreview{
				ID: 321,
			},
		},
		{
			name: "id",
			in:   discord.GuildPreview{},
			expect: discord.GuildPreview{
				ID: 123,
			},
		},
		{
			name: "emojis",
			in: discord.GuildPreview{
				ID:     321,
				Emojis: []discord.Emoji{{}},
			},
			expect: discord.GuildPreview{
				ID: 321,
				Emojis: []discord.Emoji{
					{
						ID: 1,
						User: discord.User{
							ID: 1,
						},
					},
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := GuildPreview(c.in, 123)

			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestRole(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.Role
		expect discord.Role
	}{
		{
			name: "none",
			in: discord.Role{
				ID: 321,
			},
			expect: discord.Role{
				ID: 321,
			},
		},
		{
			name: "id",
			in:   discord.Role{},
			expect: discord.Role{
				ID: 123,
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Role(c.in, 123)

			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestMember(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.Member
		expect discord.Member
	}{
		{
			name: "none",
			in: discord.Member{
				User: discord.User{
					ID: 321,
				},
			},
			expect: discord.Member{
				User: discord.User{
					ID: 321,
				},
			},
		},
		{
			name: "id",
			in:   discord.Member{},
			expect: discord.Member{
				User: discord.User{
					ID: 123,
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Member(c.in, 123)

			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestBan(t *testing.T) {
	var (
		userID discord.Snowflake = 123
		ban                      = discord.Ban{
			User: discord.User{
				ID: 0,
			},
		}
	)

	expect := discord.Ban{
		User: discord.User{
			ID: userID,
		},
	}

	actual := Ban(ban, userID)

	assert.Equal(t, expect, actual)
}

func TestIntegration(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.Integration
		expect discord.Integration
	}{
		{
			name: "none",
			in: discord.Integration{
				ID:     321,
				RoleID: 654,
				User: discord.User{
					ID: 987,
				},
			},
			expect: discord.Integration{
				ID:     321,
				RoleID: 654,
				User: discord.User{
					ID: 987,
				},
			},
		},
		{
			name: "id",
			in: discord.Integration{
				RoleID: 654,
				User: discord.User{
					ID: 987,
				},
			},
			expect: discord.Integration{
				ID:     123,
				RoleID: 654,
				User: discord.User{
					ID: 987,
				},
			},
		},
		{
			name: "ownerID",
			in: discord.Integration{
				ID: 321,
				User: discord.User{
					ID: 987,
				},
			},
			expect: discord.Integration{
				ID:     321,
				RoleID: 456,
				User: discord.User{
					ID: 987,
				},
			},
		},
		{
			name: "none",
			in: discord.Integration{
				ID:     321,
				RoleID: 654,
			},
			expect: discord.Integration{
				ID:     321,
				RoleID: 654,
				User: discord.User{
					ID: 789,
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Integration(c.in, 123, 456, 789)

			assert.Equal(t, c.expect, actual)
		})
	}
}
