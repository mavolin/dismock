package sanitize

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
)

func TestInvite(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.Invite
		expect discord.Invite
	}{
		{
			name: "none",
			in: discord.Invite{
				Guild: &discord.Guild{
					ID:                     321,
					OwnerID:                654,
					RulesChannelID:         987,
					PublicUpdatesChannelID: 210,
				},
				Channel: discord.Channel{
					ID: 543,
				},
				Inviter: &discord.User{
					ID: 876,
				},
				Target: &discord.User{
					ID: 109,
				},
			},
			expect: discord.Invite{
				Guild: &discord.Guild{
					ID:                     321,
					OwnerID:                654,
					RulesChannelID:         987,
					PublicUpdatesChannelID: 210,
				},
				Channel: discord.Channel{
					ID: 543,
				},
				Inviter: &discord.User{
					ID: 876,
				},
				Target: &discord.User{
					ID: 109,
				},
			},
		},
		{
			name: "guild",
			in: discord.Invite{
				Guild: new(discord.Guild),
				Channel: discord.Channel{
					ID: 543,
				},
				Inviter: &discord.User{
					ID: 876,
				},
				Target: &discord.User{
					ID: 109,
				},
			},
			expect: discord.Invite{
				Guild: &discord.Guild{
					ID:                     123,
					OwnerID:                456,
					RulesChannelID:         789,
					PublicUpdatesChannelID: 012,
				},
				Channel: discord.Channel{
					ID: 543,
				},
				Inviter: &discord.User{
					ID: 876,
				},
				Target: &discord.User{
					ID: 109,
				},
			},
		},
		{
			name: "channel",
			in: discord.Invite{
				Guild: &discord.Guild{
					ID:                     321,
					OwnerID:                654,
					RulesChannelID:         987,
					PublicUpdatesChannelID: 210,
				},
				Inviter: &discord.User{
					ID: 876,
				},
				Target: &discord.User{
					ID: 109,
				},
			},
			expect: discord.Invite{
				Guild: &discord.Guild{
					ID:                     321,
					OwnerID:                654,
					RulesChannelID:         987,
					PublicUpdatesChannelID: 210,
				},
				Channel: discord.Channel{
					ID: 345,
				},
				Inviter: &discord.User{
					ID: 876,
				},
				Target: &discord.User{
					ID: 109,
				},
			},
		},
		{
			name: "inviter",
			in: discord.Invite{
				Guild: &discord.Guild{
					ID:                     321,
					OwnerID:                654,
					RulesChannelID:         987,
					PublicUpdatesChannelID: 210,
				},
				Channel: discord.Channel{
					ID: 543,
				},
				Inviter: new(discord.User),
				Target: &discord.User{
					ID: 109,
				},
			},
			expect: discord.Invite{
				Guild: &discord.Guild{
					ID:                     321,
					OwnerID:                654,
					RulesChannelID:         987,
					PublicUpdatesChannelID: 210,
				},
				Channel: discord.Channel{
					ID: 543,
				},
				Inviter: &discord.User{
					ID: 678,
				},
				Target: &discord.User{
					ID: 109,
				},
			},
		},
		{
			name: "target",
			in: discord.Invite{
				Guild: &discord.Guild{
					ID:                     321,
					OwnerID:                654,
					RulesChannelID:         987,
					PublicUpdatesChannelID: 210,
				},
				Channel: discord.Channel{
					ID: 543,
				},
				Inviter: &discord.User{
					ID: 876,
				},
				Target: new(discord.User),
			},
			expect: discord.Invite{
				Guild: &discord.Guild{
					ID:                     321,
					OwnerID:                654,
					RulesChannelID:         987,
					PublicUpdatesChannelID: 210,
				},
				Channel: discord.Channel{
					ID: 543,
				},
				Inviter: &discord.User{
					ID: 876,
				},
				Target: &discord.User{
					ID: 901,
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Invite(c.in, 123, 456, 789, 012, 345, 678, 901)

			assert.Equal(t, c.expect, actual)
		})
	}
}
