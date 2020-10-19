package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/dismock/internal/sanitize"
)

func TestMocker_Messages(t *testing.T) {
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

				var channelID discord.ChannelID = 123

				expect := []discord.Message{ // more than 100 entries so multiple requests are mocked
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

				for i, msg := range expect {
					expect[i] = sanitize.Message(msg, 1, channelID, 1)
				}

				m.Messages(channelID, c.limit, expect)

				actual, err := s.Messages(channelID, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("nil messages", func(t *testing.T) {
		m, s := NewSession(t)

		var channelID discord.ChannelID = 123

		m.Messages(channelID, 100, nil)

		actual, err := s.Messages(channelID, 100)
		require.NoError(t, err)

		assert.Len(t, actual, 0)

		m.Eval()
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

		var (
			channelID discord.ChannelID = 123
			around    discord.MessageID = 456
		)

		expect := []discord.Message{
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
		}

		for i, msg := range expect {
			expect[i] = sanitize.Message(msg, 1, channelID, 1)
		}

		m.MessagesAround(channelID, around, 100, expect)

		actual, err := s.MessagesAround(channelID, around, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	limitCases := []struct {
		name  string
		limit uint
	}{
		{
			name:  "limit 0",
			limit: 0,
		},
		{
			name:  "limit > 100",
			limit: 101,
		},
	}

	for _, c := range limitCases {
		t.Run(c.name, func(t *testing.T) {
			m, s := NewSession(t)

			var (
				channelID discord.ChannelID = 123
				around    discord.MessageID = 456
			)

			expect := []discord.Message{
				{ID: 123}, {ID: 234}, {ID: 345}, {ID: 456}, {ID: 567},
				{ID: 678}, {ID: 789}, {ID: 890}, {ID: 901}, {ID: 12},
				{ID: 23}, {ID: 45}, {ID: 56}, {ID: 67}, {ID: 78},
				{ID: 89}, {ID: 90}, {ID: 98}, {ID: 87}, {ID: 76},
				{ID: 65}, {ID: 54}, {ID: 43}, {ID: 32}, {ID: 21},
				{ID: 10}, {ID: 987}, {ID: 876}, {ID: 765}, {ID: 654},
				{ID: 543}, {ID: 432}, {ID: 321}, {ID: 210}, {ID: 109},
			}

			for i, msg := range expect {
				expect[i] = sanitize.Message(msg, 1, channelID, 1)
			}

			m.MessagesAround(channelID, around, c.limit, expect)

			actual, err := s.MessagesAround(channelID, around, c.limit)
			require.NoError(t, err)

			assert.Equal(t, expect, actual)

			m.Eval()
		})
	}

	t.Run("nil messages", func(t *testing.T) {
		m, s := NewSession(t)

		var channelID discord.ChannelID = 123

		//noinspection GoPreferNilSlice
		expect := []discord.Message{}

		m.MessagesAround(channelID, 0, 100, nil)

		actual, err := s.MessagesAround(channelID, 0, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var channelID discord.ChannelID = 123

		expect := []discord.Message{
			{
				ID: 123,
			},
			{
				ID: 456,
			},
		}

		for i, msg := range expect {
			expect[i] = sanitize.Message(msg, 1, channelID, 1)
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

				var (
					channelID discord.ChannelID = 123
					before    discord.MessageID = 3
				)

				expect := []discord.Message{ // more than 100 entries so multiple requests are mocked
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

				for i, msg := range expect {
					expect[i] = sanitize.Message(msg, 1, channelID, 1)
				}

				m.MessagesBefore(channelID, before, c.limit, expect)

				actual, err := s.MessagesBefore(channelID, before, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("nil messages", func(t *testing.T) {
		m, s := NewSession(t)

		var channelID discord.ChannelID = 123

		m.MessagesBefore(channelID, 0, 100, nil)

		actual, err := s.MessagesBefore(channelID, 0, 100)
		require.NoError(t, err)

		assert.Len(t, actual, 0)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var channelID discord.ChannelID = 123

		expect := []discord.Message{
			{
				ID: 123,
			},
			{
				ID: 456,
			},
		}

		for i, msg := range expect {
			expect[i] = sanitize.Message(msg, 1, channelID, 1)
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

				var (
					channelID discord.ChannelID = 123
					after     discord.MessageID = 456
				)

				expect := []discord.Message{ // more than 100 entries so multiple requests are mocked
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

				for i, msg := range expect {
					expect[i] = sanitize.Message(msg, 1, channelID, 1)
				}
				m.MessagesAfter(channelID, after, c.limit, expect)

				actual, err := s.MessagesAfter(channelID, after, c.limit)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("nil guilds", func(t *testing.T) {
		m, s := NewSession(t)

		var channelID discord.ChannelID = 123

		var expect []discord.Message

		m.MessagesAfter(channelID, 0, 100, expect)

		actual, err := s.MessagesAfter(channelID, 0, 100)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var channelID discord.ChannelID = 123

		expect := []discord.Message{
			{
				ID: 456,
			},
			{
				ID: 789,
			},
		}

		for i, msg := range expect {
			expect[i] = sanitize.Message(msg, 1, channelID, 1)
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

	expect := sanitize.Message(discord.Message{
		ID:        123,
		ChannelID: 465,
	}, 1, 1, 1)

	m.Message(expect)

	actual, err := s.Message(expect.ChannelID, expect.ID)
	require.NoError(t, err)

	assert.Equal(t, expect, *actual)

	m.Eval()
}

func TestMocker_SendText(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		expect := sanitize.Message(discord.Message{
			ID:        123,
			ChannelID: 456,
			Content:   "abc",
		}, 1, 1, 1)

		m.SendText(expect)

		actual, err := s.SendText(expect.ChannelID, expect.Content)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var channelID discord.ChannelID = 123

		m.SendText(sanitize.Message(discord.Message{
			ChannelID: channelID,
			Content:   "abc",
		}, 1, 1, 1))

		_, err := s.SendText(channelID, "cba")
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_SendEmbed(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		expect := sanitize.Message(discord.Message{
			ID:        123,
			ChannelID: 456,
			Embeds: []discord.Embed{
				{
					Title:       "def",
					Description: "ghi",
				},
			},
		}, 1, 1, 1)

		m.SendEmbed(expect)

		actual, err := s.SendEmbed(expect.ChannelID, expect.Embeds[0])
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
		)

		m.SendEmbed(sanitize.Message(discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds: []discord.Embed{
				{
					Title:       "def",
					Description: "ghi",
				},
			},
		}, 1, 1, 1))

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
				Content:   "abc",
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				expect := sanitize.Message(c.msg, 1, 1, 1)

				var embed *discord.Embed
				if len(expect.Embeds) > 0 {
					embed = &expect.Embeds[0]
				}

				m.SendMessage(embed, expect)

				actual, err := s.SendMessage(expect.ChannelID, expect.Content, embed)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)

				m.Eval()
			})
		}

		t.Run("param embed", func(t *testing.T) {
			m, s := NewSession(t)

			var (
				msg = sanitize.Message(discord.Message{
					ID:        123,
					ChannelID: 456,
					Content:   "abc",
				}, 1, 1, 1)

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

			m.Eval()
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

		m.SendMessage(&embed, sanitize.Message(discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds:    []discord.Embed{embed},
		}, 1, 1, 1))

		_, err := s.SendMessage(channelID, "cba", &embed)
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_EditText(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		expect := sanitize.Message(discord.Message{
			ID:        123,
			ChannelID: 456,
			Content:   "abc",
		}, 1, 1, 1)

		m.EditText(expect)

		actual, err := s.EditText(expect.ChannelID, expect.ID, expect.Content)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
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

		m.EditText(sanitize.Message(discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds:    []discord.Embed{embed},
		}, 1, 1, 1))

		_, err := s.EditText(channelID, messageID, "cba")
		require.NoError(t, err)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_EditEmbed(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		expect := sanitize.Message(discord.Message{
			ID:        123,
			ChannelID: 456,
			Embeds: []discord.Embed{
				{
					Title:       "def",
					Description: "ghi",
				},
			},
		}, 1, 1, 1)

		m.EditEmbed(expect)

		actual, err := s.EditEmbed(expect.ChannelID, expect.ID, expect.Embeds[0])
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var (
			channelID discord.ChannelID = 123
			messageID discord.MessageID = 456
		)

		m.EditEmbed(sanitize.Message(discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds: []discord.Embed{
				{
					Title:       "def",
					Description: "ghi",
				},
			},
		}, 1, 1, 1))

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
				Content:   "abc",
			},
		},
		{
			name:           "don't suppressEmbeds",
			suppressEmbeds: false,
			msg: discord.Message{
				ID:        123,
				ChannelID: 456,
				Content:   "abc",
			},
		},
		{
			name: "with embed",
			msg: discord.Message{
				ID:        123,
				ChannelID: 456,
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
				Content:   "abc",
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				expect := sanitize.Message(c.msg, 1, 1, 1)

				var embed *discord.Embed
				if len(expect.Embeds) > 0 {
					embed = &expect.Embeds[0]
				}

				m.EditMessage(embed, expect, c.suppressEmbeds)

				actual, err := s.EditMessage(expect.ChannelID, expect.ID, expect.Content, embed, c.suppressEmbeds)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)

				m.Eval()
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

		m.EditMessage(&embed, sanitize.Message(discord.Message{
			ID:        messageID,
			ChannelID: channelID,
			Content:   "abc",
			Embeds:    []discord.Embed{embed},
		}, 1, 1, 1), false)

		_, err := s.EditMessage(channelID, messageID, "cba", &embed, false)
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

		expect := sanitize.Message(discord.Message{
			ID:        123,
			ChannelID: 456,
			Content:   "abc",
		}, 1, 1, 1)

		m.EditMessageComplex(data, expect)

		actual, err := s.EditMessageComplex(expect.ChannelID, expect.ID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := sanitize.Message(discord.Message{
			ID:        123,
			ChannelID: 456,
			Content:   "abc",
		}, 1, 1, 1)

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

	var (
		channelID discord.ChannelID = 123
		messageID discord.MessageID = 456
	)

	m.DeleteMessage(channelID, messageID)

	err := s.DeleteMessage(channelID, messageID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_DeleteMessages(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			channelID  discord.ChannelID = 123
			messageIDs                   = []discord.MessageID{456, 789}
		)

		m.DeleteMessages(channelID, messageIDs)

		err := s.DeleteMessages(channelID, messageIDs)
		require.NoError(t, err)

		m.Eval()
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
