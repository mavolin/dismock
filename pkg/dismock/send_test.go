package dismock

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/dismock/internal/sanitize"
)

func TestMocker_SendMessageComplex(t *testing.T) {
	successCases := []struct {
		name string
		data api.SendMessageData
	}{
		{
			name: "no files",
			data: api.SendMessageData{
				Content: "abc",
			},
		},
		{
			name: "with file",
			data: api.SendMessageData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("def"),
					},
				},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := NewSession(t)

				expect := sanitize.Message(discord.Message{
					ChannelID: 123,
				}, 1, 1, 1)

				cp := c.data

				cp.Files = make([]api.SendMessageFile, len(c.data.Files))
				copy(cp.Files, c.data.Files) // the readers of the file will be consumed twice

				// the files are copied now, but the reader for them may be a pointer and wasn't
				// deep copied. therefore we create two readers using the data from the original
				// reader
				for i, f := range c.data.Files {
					b, err := ioutil.ReadAll(f.Reader)
					require.NoError(t, err)

					cp.Files[i].Reader = bytes.NewBuffer(b)
					c.data.Files[i].Reader = bytes.NewBuffer(b)
				}

				m.SendMessageComplex(c.data, expect)

				actual, err := s.SendMessageComplex(expect.ChannelID, cp)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)

				m.Eval()
			})
		}
	})

	failureCases := []struct {
		name  string
		data1 api.SendMessageData
		data2 api.SendMessageData
	}{
		{
			name: "different content",
			data1: api.SendMessageData{
				Content: "abc",
			},
			data2: api.SendMessageData{
				Content: "cba",
			},
		},
		{
			name: "different file",
			data1: api.SendMessageData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("def"),
					},
				},
			},
			data2: api.SendMessageData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("fed"),
					},
				},
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				tMock := new(testing.T)

				m, s := NewSession(tMock)

				expect := sanitize.Message(discord.Message{
					ChannelID: 123,
				}, 1, 1, 1)

				m.SendMessageComplex(c.data1, expect)

				actual, err := s.SendMessageComplex(expect.ChannelID, c.data2)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)
				assert.True(t, tMock.Failed())
			})
		}
	})
}

func TestMocker_ExecuteWebhook(t *testing.T) {
	successCases := []struct {
		name string
		data api.ExecuteWebhookData
	}{
		{
			name: "no files",
			data: api.ExecuteWebhookData{
				Content: "abc",
			},
		},
		{
			name: "with file",
			data: api.ExecuteWebhookData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("def"),
					},
				},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m := New(t)

				var (
					webhookID discord.WebhookID = 123
					token                       = "abc"
				)

				cp := c.data

				cp.Files = make([]api.SendMessageFile, len(c.data.Files))
				copy(cp.Files, c.data.Files) // the readers of the file will be consumed twice

				// the files are copied now, but the reader for them may be a pointer and wasn't
				// deep copied. therefore we create two readers using the data from the original
				// reader
				for i, f := range c.data.Files {
					b, err := ioutil.ReadAll(f.Reader)
					require.NoError(t, err)

					cp.Files[i].Reader = bytes.NewBuffer(b)
					c.data.Files[i].Reader = bytes.NewBuffer(b)
				}

				m.ExecuteWebhook(webhookID, token, c.data)

				err := webhook.Execute(webhookID, token, cp)
				require.NoError(t, err)

				m.Eval()
			})
		}
	})

	failureCases := []struct {
		name  string
		data1 api.ExecuteWebhookData
		data2 api.ExecuteWebhookData
	}{
		{
			name: "different content",
			data1: api.ExecuteWebhookData{
				Content: "abc",
			},
			data2: api.ExecuteWebhookData{
				Content: "cba",
			},
		},
		{
			name: "different file",
			data1: api.ExecuteWebhookData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("def"),
					},
				},
			},
			data2: api.ExecuteWebhookData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("fed"),
					},
				},
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				tMock := new(testing.T)

				m := New(tMock)

				var (
					webhookID discord.WebhookID = 123
					token                       = "abc"
				)

				m.ExecuteWebhook(webhookID, token, c.data1)

				err := webhook.Execute(webhookID, token, c.data2)
				require.NoError(t, err)

				assert.True(t, tMock.Failed())
			})
		}
	})
}

func TestMocker_ExecuteWebhookAndWait(t *testing.T) {
	successCases := []struct {
		name string
		data api.ExecuteWebhookData
	}{
		{
			name: "no files",
			data: api.ExecuteWebhookData{
				Content: "abc",
			},
		},
		{
			name: "with file",
			data: api.ExecuteWebhookData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("def"),
					},
				},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m := New(t)

				var (
					webhookID discord.WebhookID = 123
					token                       = "abc"
				)

				expect := sanitize.Message(discord.Message{
					ChannelID: 123,
				}, 1, 1, 1)

				cp := c.data

				cp.Files = make([]api.SendMessageFile, len(c.data.Files))
				copy(cp.Files, c.data.Files) // the readers of the file will be consumed twice

				// the files are copied now, but the reader for them may be a pointer and wasn't
				// deep copied. therefore we create two readers using the data from the original
				// reader
				for i, f := range c.data.Files {
					b, err := ioutil.ReadAll(f.Reader)
					require.NoError(t, err)

					cp.Files[i].Reader = bytes.NewBuffer(b)
					c.data.Files[i].Reader = bytes.NewBuffer(b)
				}

				m.ExecuteWebhookAndWait(webhookID, token, c.data, expect)

				actual, err := webhook.ExecuteAndWait(webhookID, token, cp)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)

				m.Eval()
			})
		}
	})

	failureCases := []struct {
		name  string
		data1 api.ExecuteWebhookData
		data2 api.ExecuteWebhookData
	}{
		{
			name: "different content",
			data1: api.ExecuteWebhookData{
				Content: "abc",
			},
			data2: api.ExecuteWebhookData{
				Content: "cba",
			},
		},
		{
			name: "different file",
			data1: api.ExecuteWebhookData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("def"),
					},
				},
			},
			data2: api.ExecuteWebhookData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("fed"),
					},
				},
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				tMock := new(testing.T)

				m := New(tMock)

				var (
					webhookID discord.WebhookID = 123
					token                       = "abc"
				)

				expect := sanitize.Message(discord.Message{
					ChannelID: 123,
				}, 1, 1, 1)

				m.ExecuteWebhookAndWait(webhookID, token, c.data1, expect)

				actual, err := webhook.ExecuteAndWait(webhookID, token, c.data2)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)
				assert.True(t, tMock.Failed())
			})
		}
	})
}
