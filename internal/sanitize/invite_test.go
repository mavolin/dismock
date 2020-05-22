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
					ID:      321,
					OwnerID: 654,
				},
				Channel: discord.Channel{
					ID: 987,
				},
				Inviter: &discord.User{
					ID: 210,
				},
				Target: &discord.User{
					ID: 543,
				},
			},
			expect: discord.Invite{
				Guild: &discord.Guild{
					ID:      321,
					OwnerID: 654,
				},
				Channel: discord.Channel{
					ID: 987,
				},
				Inviter: &discord.User{
					ID: 210,
				},
				Target: &discord.User{
					ID: 543,
				},
			},
		},
		{
			name: "guild",
			in: discord.Invite{
				Guild: new(discord.Guild),
				Channel: discord.Channel{
					ID: 987,
				},
				Inviter: &discord.User{
					ID: 210,
				},
				Target: &discord.User{
					ID: 543,
				},
			},
			expect: discord.Invite{
				Guild: &discord.Guild{
					ID:      123,
					OwnerID: 456,
				},
				Channel: discord.Channel{
					ID: 987,
				},
				Inviter: &discord.User{
					ID: 210,
				},
				Target: &discord.User{
					ID: 543,
				},
			},
		},
		{
			name: "channel",
			in: discord.Invite{
				Guild: &discord.Guild{
					ID:      321,
					OwnerID: 654,
				},
				Inviter: &discord.User{
					ID: 210,
				},
				Target: &discord.User{
					ID: 543,
				},
			},
			expect: discord.Invite{
				Guild: &discord.Guild{
					ID:      321,
					OwnerID: 654,
				},
				Channel: discord.Channel{
					ID: 789,
				},
				Inviter: &discord.User{
					ID: 210,
				},
				Target: &discord.User{
					ID: 543,
				},
			},
		},
		{
			name: "inviter",
			in: discord.Invite{
				Guild: &discord.Guild{
					ID:      321,
					OwnerID: 654,
				},
				Channel: discord.Channel{
					ID: 987,
				},
				Inviter: new(discord.User),
				Target: &discord.User{
					ID: 543,
				},
			},
			expect: discord.Invite{
				Guild: &discord.Guild{
					ID:      321,
					OwnerID: 654,
				},
				Channel: discord.Channel{
					ID: 987,
				},
				Inviter: &discord.User{
					ID: 012,
				},
				Target: &discord.User{
					ID: 543,
				},
			},
		},
		{
			name: "target",
			in: discord.Invite{
				Guild: &discord.Guild{
					ID:      321,
					OwnerID: 654,
				},
				Channel: discord.Channel{
					ID: 987,
				},
				Inviter: &discord.User{
					ID: 210,
				},
				Target: &discord.User{
					ID: 345,
				},
			},
			expect: discord.Invite{
				Guild: &discord.Guild{
					ID:      321,
					OwnerID: 654,
				},
				Channel: discord.Channel{
					ID: 987,
				},
				Inviter: &discord.User{
					ID: 210,
				},
				Target: &discord.User{
					ID: 345,
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Invite(c.in, 123, 456, 789, 012, 345)

			assert.Equal(t, c.expect, actual)
		})
	}
}
