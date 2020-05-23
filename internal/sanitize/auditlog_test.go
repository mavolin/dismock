package sanitize

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
)

func TestAuditLog(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.AuditLog
		expect discord.AuditLog
	}{
		{
			name: "none",
			in: discord.AuditLog{
				Webhooks: []discord.Webhook{
					{
						ID: 123,
						User: discord.User{
							ID: 456,
						},
						ChannelID: 1,
					},
				},
				Users: []discord.User{
					{
						ID: 789,
					},
				},
				Entries: []discord.AuditLogEntry{
					{
						ID:     012,
						UserID: 345,
					},
				},
				Integrations: []discord.Integration{
					{
						ID:     678,
						RoleID: 901,
						User: discord.User{
							ID: 234,
						},
					},
				},
			},
			expect: discord.AuditLog{
				Webhooks: []discord.Webhook{
					{
						ID: 123,
						User: discord.User{
							ID: 456,
						},
						ChannelID: 1,
					},
				},
				Users: []discord.User{
					{
						ID: 789,
					},
				},
				Entries: []discord.AuditLogEntry{
					{
						ID:     012,
						UserID: 345,
					},
				},
				Integrations: []discord.Integration{
					{
						ID:     678,
						RoleID: 901,
						User: discord.User{
							ID: 234,
						},
					},
				},
			},
		},
		{
			name: "webhooks",
			in: discord.AuditLog{
				Webhooks: []discord.Webhook{{}},
			},
			expect: discord.AuditLog{
				Webhooks: []discord.Webhook{
					{
						ID: 1,
						User: discord.User{
							ID: 1,
						},
						ChannelID: 1,
					},
				},
			},
		},
		{
			name: "user",
			in: discord.AuditLog{
				Users: []discord.User{{}},
			},
			expect: discord.AuditLog{
				Users: []discord.User{
					{
						ID: 1,
					},
				},
			},
		},
		{
			name: "entries",
			in: discord.AuditLog{
				Entries: []discord.AuditLogEntry{{}},
			},
			expect: discord.AuditLog{
				Entries: []discord.AuditLogEntry{
					{
						ID:     1,
						UserID: 1,
					},
				},
			},
		},
		{
			name: "integrations",
			in: discord.AuditLog{
				Integrations: []discord.Integration{{}},
			},
			expect: discord.AuditLog{
				Integrations: []discord.Integration{
					{
						ID:     1,
						RoleID: 1,
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
			actual := AuditLog(c.in)

			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestAuditLogEntry(t *testing.T) {
	testCases := []struct {
		name   string
		in     discord.AuditLogEntry
		expect discord.AuditLogEntry
	}{
		{
			name: "none",
			in: discord.AuditLogEntry{
				ID:     321,
				UserID: 654,
			},
			expect: discord.AuditLogEntry{
				ID:     321,
				UserID: 654,
			},
		},
		{
			name: "id",
			in: discord.AuditLogEntry{
				UserID: 654,
			},
			expect: discord.AuditLogEntry{
				ID:     123,
				UserID: 654,
			},
		},
		{
			name: "userID",
			in: discord.AuditLogEntry{
				ID: 321,
			},
			expect: discord.AuditLogEntry{
				ID:     321,
				UserID: 456,
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := AuditLogEntry(c.in, 123, 456)

			assert.Equal(t, c.expect, actual)
		})
	}
}
