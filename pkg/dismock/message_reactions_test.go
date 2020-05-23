package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_React(t *testing.T) {
	m, s := NewArikawaSession(t)

	var (
		channelID discord.Snowflake = 123
		messageID discord.Snowflake = 456
		emoji                       = "abc"
	)

	m.React(channelID, messageID, emoji)

	err := s.React(channelID, messageID, emoji)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_Unreact(t *testing.T) {
	m, s := NewArikawaSession(t)

	var (
		channelID discord.Snowflake = 123
		messageID discord.Snowflake = 456
		emoji                       = "abc"
	)

	m.Unreact(channelID, messageID, emoji)

	err := s.Unreact(channelID, messageID, emoji)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_Reactions(t *testing.T) {
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
					channelID discord.Snowflake = 123
					messageID discord.Snowflake = 456
					emoji                       = "abc"
				)

				expect := []discord.User{ // more than 100 entries so multiple requests are mocked
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

				m.Reactions(channelID, messageID, c.limit, emoji, expect)

				actual, err := s.Reactions(channelID, messageID, c.limit, emoji)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("nil users", func(t *testing.T) {
		m, s := NewArikawaSession(t)

		var (
			channelID discord.Snowflake = 123
			messageID discord.Snowflake = 456
			emoji                       = "abc"
		)

		var expect []discord.User

		m.Reactions(channelID, messageID, 100, emoji, expect)

		actual, err := s.Reactions(channelID, messageID, 100, emoji)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m, _ := NewArikawaSession(t)

		assert.Panics(t, func() {
			m.Reactions(123, 456, 1, "abc", []discord.User{{}, {}})
		})
	})
}

func TestMocker_ReactionsBefore(t *testing.T) {
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
					channelID discord.Snowflake = 123
					messageID discord.Snowflake = 456
					emoji                       = "abc"

					before discord.Snowflake = 999999999999
				)

				expect := []discord.User{ // more than 100 entries so multiple requests are mocked
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

				m.ReactionsBefore(channelID, messageID, before, c.limit, emoji, expect)

				actual, err := s.ReactionsBefore(channelID, messageID, before, c.limit, emoji)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("nil users", func(t *testing.T) {
		m, s := NewArikawaSession(t)

		var (
			channelID discord.Snowflake = 123
			messageID discord.Snowflake = 456
			emoji                       = "abc"
		)

		//noinspection GoPreferNilSlice
		var expect []discord.User

		m.ReactionsBefore(channelID, messageID, 0, 100, emoji, expect)

		actual, err := s.ReactionsBefore(channelID, messageID, 0, 100, emoji)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewArikawaSession(tMock)

		var (
			channelID discord.Snowflake = 123
			messageID discord.Snowflake = 456
			emoji                       = "abc"
		)

		expect := []discord.User{
			{
				ID: 123,
			},
			{
				ID: 456,
			},
		}

		m.ReactionsBefore(channelID, messageID, 890, 100, emoji, expect)

		actual, err := s.ReactionsBefore(channelID, messageID, 789, 100, emoji)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m, _ := NewArikawaSession(t)

		assert.Panics(t, func() {
			m.ReactionsBefore(123, 456, 0, 1, "abc", []discord.User{{}, {}})
		})
	})
}

func TestMocker_ReactionsAfter(t *testing.T) {
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
					channelID discord.Snowflake = 123
					messageID discord.Snowflake = 456
					emoji                       = "abc"

					after discord.Snowflake = 123
				)

				expect := []discord.User{ // more than 100 entries so multiple requests are mocked
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

				m.ReactionsAfter(channelID, messageID, after, c.limit, emoji, expect)

				actual, err := s.ReactionsAfter(channelID, messageID, after, c.limit, emoji)
				require.NoError(t, err)

				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})

	t.Run("nil users", func(t *testing.T) {
		m, s := NewArikawaSession(t)

		var (
			channelID discord.Snowflake = 123
			messageID discord.Snowflake = 456
			emoji                       = "abc"
		)

		var expect []discord.User

		m.ReactionsAfter(channelID, messageID, 0, 100, emoji, expect)

		actual, err := s.ReactionsAfter(channelID, messageID, 0, 100, emoji)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewArikawaSession(tMock)

		var (
			channelID discord.Snowflake = 123
			messageID discord.Snowflake = 456
			emoji                       = "abc"
		)

		expect := []discord.User{
			{
				ID: 456,
			},
			{
				ID: 789,
			},
		}

		m.ReactionsAfter(channelID, messageID, 123, 100, emoji, expect)

		actual, err := s.ReactionsAfter(channelID, messageID, 321, 100, emoji)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})

	t.Run("limit smaller than guilds", func(t *testing.T) {
		m, _ := NewArikawaSession(t)

		assert.Panics(t, func() {
			m.ReactionsAfter(123, 456, 0, 1, "abc", []discord.User{{}, {}})
		})
	})
}

func TestMocker_DeleteUserReaction(t *testing.T) {
	m, s := NewArikawaSession(t)

	var (
		channelID discord.Snowflake = 123
		messageID discord.Snowflake = 456
		userID    discord.Snowflake = 789
		emoji                       = "abc"
	)

	m.DeleteUserReaction(channelID, messageID, userID, emoji)

	err := s.DeleteUserReaction(channelID, messageID, userID, emoji)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_DeleteReactions(t *testing.T) {
	m, s := NewArikawaSession(t)

	var (
		channelID discord.Snowflake = 123
		messageID discord.Snowflake = 456
		emoji                       = "abc"
	)

	m.DeleteReactions(channelID, messageID, emoji)

	err := s.DeleteReactions(channelID, messageID, emoji)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_DeleteAllReactions(t *testing.T) {
	m, s := NewArikawaSession(t)

	var (
		channelID discord.Snowflake = 123
		messageID discord.Snowflake = 456
	)

	m.DeleteAllReactions(channelID, messageID)

	err := s.DeleteAllReactions(channelID, messageID)
	require.NoError(t, err)

	m.Eval()
}
