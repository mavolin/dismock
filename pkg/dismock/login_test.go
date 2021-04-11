package dismock

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			email    = "abc@def.ghi"
			password = "jkl"
		)

		expect := api.LoginResponse{Token: "mno"}

		m.Login(email, password, expect)

		actual, err := s.Login(email, password)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)
		m, s := NewSession(tMock)

		expect := api.LoginResponse{Token: "mno"}

		m.Login("cba@fed.ihg", "jkl", expect)

		actual, err := s.Login("abc@def.ghi", "jkl")
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}

func TestMocker_TOTP(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := NewSession(t)

		var (
			code   = "abc"
			ticket = "def"
		)

		expect := api.LoginResponse{Token: "ghi"}

		m.TOTP(code, ticket, expect)

		actual, err := s.TOTP(code, ticket)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, s := NewSession(tMock)

		expect := api.LoginResponse{Token: "ghi"}

		m.TOTP("abc", "def", expect)

		actual, err := s.TOTP("cba", "def")
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
		assert.True(t, tMock.Failed())
	})
}
