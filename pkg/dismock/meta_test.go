package dismock

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ================================ Channel ================================

func TestMocker_ChannelIconURL(t *testing.T) {
	m := New(t)

	channel := discord.Channel{
		ID:   123,
		Icon: "abc",
	}

	expect := []byte{1, 30, 0, 15, 24}

	img := bytes.NewBuffer(expect)

	m.ChannelIcon(channel.ID, channel.Icon, img)

	resp, err := m.Client.Get(channel.IconURL())
	require.NoError(t, err)

	actual, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_ChannelIconURLWithType(t *testing.T) {
	m := New(t)

	var (
		channel = discord.Channel{
			ID:   123,
			Icon: "abc",
		}
		imgType = discord.JPEGImage
	)

	expect := []byte{1, 30, 0, 15, 24}

	img := bytes.NewBuffer(expect)

	m.ChannelIconWithType(channel.ID, channel.Icon, imgType, img)

	resp, err := m.Client.Get(channel.IconURLWithType(imgType))
	require.NoError(t, err)

	actual, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

// ================================ Emoji ================================

func TestMocker_EmojiURL(t *testing.T) {
	testCases := []struct {
		name     string
		animated bool
	}{
		{
			name:     "animated",
			animated: true,
		},
		{
			name:     "not animated",
			animated: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			m := New(t)

			emoji := discord.Emoji{
				ID:       123,
				Animated: c.animated,
			}

			expect := []byte{1, 30, 0, 15, 24}

			img := bytes.NewBuffer(expect)

			m.EmojiPicture(emoji.ID, emoji.Animated, img)

			resp, err := m.Client.Get(emoji.EmojiURL())
			require.NoError(t, err)

			actual, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, expect, actual)

			m.Eval()
		})
	}
}

func TestMocker_EmojiURLWithType(t *testing.T) {
	testCases := []struct {
		name     string
		imgType  discord.ImageType
		animated bool
	}{
		{
			name:     "auto - animated",
			imgType:  discord.AutoImage,
			animated: true,
		},
		{
			name:     "auto - not animated",
			imgType:  discord.AutoImage,
			animated: false,
		},
		{
			name:     "non-auto",
			imgType:  discord.PNGImage,
			animated: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			m := New(t)

			emoji := discord.Emoji{
				ID:       123,
				Animated: c.animated,
			}

			expect := []byte{1, 30, 0, 15, 24}

			img := bytes.NewBuffer(expect)

			m.EmojiPictureWithType(emoji.ID, emoji.Animated, c.imgType, img)

			resp, err := m.Client.Get(emoji.EmojiURLWithType(c.imgType))
			require.NoError(t, err)

			actual, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, expect, actual)

			m.Eval()
		})
	}
}

// ================================ Guild ================================

func TestMocker_GuildIconURL(t *testing.T) {
	m := New(t)

	guild := discord.Guild{
		ID:   123,
		Icon: "abc",
	}

	expect := []byte{1, 30, 0, 15, 24}

	img := bytes.NewBuffer(expect)

	m.GuildIcon(guild.ID, guild.Icon, img)

	resp, err := m.Client.Get(guild.IconURL())
	require.NoError(t, err)

	actual, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_GuildIconURLWithType(t *testing.T) {
	m := New(t)

	var (
		guild = discord.Guild{
			ID:   123,
			Icon: "abc",
		}
		imgType = discord.JPEGImage
	)

	expect := []byte{1, 30, 0, 15, 24}

	img := bytes.NewBuffer(expect)

	m.GuildIconWithType(guild.ID, guild.Icon, imgType, img)

	resp, err := m.Client.Get(guild.IconURLWithType(imgType))
	require.NoError(t, err)

	actual, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_BannerURL(t *testing.T) {
	m := New(t)

	guild := discord.Guild{
		ID:     123,
		Banner: "abc",
	}

	expect := []byte{1, 30, 0, 15, 24}

	img := bytes.NewBuffer(expect)

	m.Banner(guild.ID, guild.Banner, img)

	resp, err := m.Client.Get(guild.BannerURL())
	require.NoError(t, err)

	actual, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_BannerURLWithType(t *testing.T) {
	m := New(t)

	var (
		guild = discord.Guild{
			ID:     123,
			Banner: "abc",
		}
		imgType = discord.JPEGImage
	)

	expect := []byte{1, 30, 0, 15, 24}

	img := bytes.NewBuffer(expect)

	m.BannerWithType(guild.ID, guild.Banner, imgType, img)

	resp, err := m.Client.Get(guild.BannerURLWithType(imgType))
	require.NoError(t, err)

	actual, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_SplashURL(t *testing.T) {
	m := New(t)

	guild := discord.Guild{
		ID:     123,
		Splash: "abc",
	}

	expect := []byte{1, 30, 0, 15, 24}

	img := bytes.NewBuffer(expect)

	m.Splash(guild.ID, guild.Splash, img)

	resp, err := m.Client.Get(guild.SplashURL())
	require.NoError(t, err)

	actual, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_SplashURLWithType(t *testing.T) {
	m := New(t)

	var (
		guild = discord.Guild{
			ID:     123,
			Splash: "abc",
		}
		imgType = discord.JPEGImage
	)

	expect := []byte{1, 30, 0, 15, 24}

	img := bytes.NewBuffer(expect)

	m.SplashWithType(guild.ID, guild.Splash, imgType, img)

	resp, err := m.Client.Get(guild.SplashURLWithType(imgType))
	require.NoError(t, err)

	actual, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_DiscoverySplashURL(t *testing.T) {
	m := New(t)

	guild := discord.GuildPreview{
		ID:              123,
		DiscoverySplash: "abc",
	}

	expect := []byte{1, 30, 0, 15, 24}

	img := bytes.NewBuffer(expect)

	m.DiscoverySplash(guild.ID, guild.DiscoverySplash, img)

	resp, err := m.Client.Get(guild.DiscoverySplashURL())
	require.NoError(t, err)

	actual, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestMocker_DiscoverySplashURLWithType(t *testing.T) {
	m := New(t)

	var (
		guild = discord.Guild{
			ID:              123,
			DiscoverySplash: "abc",
		}
		imgType = discord.JPEGImage
	)

	expect := []byte{1, 30, 0, 15, 24}

	img := bytes.NewBuffer(expect)

	m.DiscoverySplashWithType(guild.ID, guild.DiscoverySplash, imgType, img)

	resp, err := m.Client.Get(guild.DiscoverySplashURLWithType(imgType))
	require.NoError(t, err)

	actual, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	m.Eval()
}

func TestFormatImageType(t *testing.T) {
	testCases := []struct {
		name            string
		formatName      string
		imgType         discord.ImageType
		expectExtension string
	}{
		{
			name:            "auto gif",
			formatName:      "a_abc",
			imgType:         discord.AutoImage,
			expectExtension: ".gif",
		},
		{
			name:            "auto png",
			formatName:      "abc",
			imgType:         discord.AutoImage,
			expectExtension: ".png",
		},
		{
			name:            "manual",
			formatName:      "abc",
			imgType:         discord.WebPImage,
			expectExtension: ".webp",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			format := formatImageType(c.formatName, c.imgType)

			assert.True(t, strings.HasSuffix(format, c.expectExtension))
		})
	}
}
