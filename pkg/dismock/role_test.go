package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_AddRole(t *testing.T) {
	m, s := NewSession(t)

	var (
		guildID discord.Snowflake = 123
		userID  discord.Snowflake = 456
		roleID  discord.Snowflake = 789
	)

	m.AddRole(guildID, userID, roleID)

	err := s.AddRole(guildID, userID, roleID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_RemoveRole(t *testing.T) {
	m, s := NewSession(t)

	var (
		guildID discord.Snowflake = 123
		userID  discord.Snowflake = 456
		roleID  discord.Snowflake = 789
	)

	m.RemoveRole(guildID, userID, roleID)

	err := s.RemoveRole(guildID, userID, roleID)
	require.NoError(t, err)

	m.Eval()
}

func TestMocker_Roles(t *testing.T) {
	m, s := NewSession(t)

	var guildID discord.Snowflake = 123

	expect := []discord.Role{
		{
			ID: 456,
		},
		{
			ID: 789,
		},
	}

	m.Roles(guildID, expect)

	actual, err := s.Roles(guildID)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_CreateRole(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.Snowflake = 123

		data := api.CreateRoleData{
			Name: "abc",
		}

		expect := discord.Role{
			ID:   456,
			Name: "abc",
		}

		m.CreateRole(guildID, data, expect)

		actual, err := s.CreateRole(guildID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var guildID discord.Snowflake = 123

		expect := discord.Role{
			ID:   456,
			Name: "abc",
		}

		m.CreateRole(guildID, api.CreateRoleData{
			Name: "abc",
		}, expect)

		actual, err := s.CreateRole(guildID, api.CreateRoleData{
			Name: "cba",
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_MoveRole(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.Snowflake = 123

		data := []api.MoveRoleData{
			{
				ID:       456,
				Position: option.NewNullableInt(1),
			},
			{
				ID:       789,
				Position: option.NewNullableInt(0),
			},
		}

		expect := []discord.Role{
			{
				ID:       456,
				Name:     "abc",
				Position: 1,
			},
			{
				ID:       789,
				Name:     "def",
				Position: 0,
			},
		}

		m.MoveRole(guildID, data, expect)

		actual, err := s.MoveRole(guildID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var guildID discord.Snowflake = 123

		expect := []discord.Role{
			{
				ID:       456,
				Name:     "abc",
				Position: 1,
			},
			{
				ID:       789,
				Name:     "def",
				Position: 0,
			},
		}

		m.MoveRole(guildID, []api.MoveRoleData{
			{
				ID:       456,
				Position: option.NewNullableInt(1),
			},
			{
				ID:       789,
				Position: option.NewNullableInt(0),
			},
		}, expect)

		actual, err := s.MoveRole(guildID, []api.MoveRoleData{
			{
				ID:       654,
				Position: option.NewNullableInt(1),
			},
			{
				ID:       987,
				Position: option.NewNullableInt(0),
			},
		})
		require.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_ModifyRole(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var guildID discord.Snowflake = 123

		data := api.ModifyRoleData{
			Name: option.NewNullableString("abc"),
		}

		expect := discord.Role{
			ID:       456,
			Name:     "abc",
			Position: 1,
		}

		m.ModifyRole(guildID, data, expect)

		actual, err := s.ModifyRole(guildID, expect.ID, data)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		var guildID discord.Snowflake = 123

		expect := discord.Role{
			ID:       456,
			Name:     "abc",
			Position: 1,
		}

		m.ModifyRole(guildID, api.ModifyRoleData{
			Name: option.NewNullableString("abc"),
		}, expect)

		actual, err := s.ModifyRole(guildID, expect.ID, api.ModifyRoleData{
			Name: option.NewNullableString("cba"),
		})
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_DeleteRole(t *testing.T) {
	m, s := NewSession(t)

	var (
		guildID discord.Snowflake = 123
		roleID  discord.Snowflake = 456
	)

	m.DeleteRole(guildID, roleID)

	err := s.DeleteRole(guildID, roleID)
	require.NoError(t, err)

	m.Eval()
}
