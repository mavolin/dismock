package dismock

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/dismock/v2/internal/mockutil"
)

const maxFetchGuilds = 100

// CreateGuild mocks a CreateGuild request.
func (m *Mocker) CreateGuild(d api.CreateGuildData, g discord.Guild) {
	m.MockAPI("CreateGuild", http.MethodPost, "/guilds",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &d)
			mockutil.WriteJSON(t, w, g)
		})
}

// Guild mocks a Guild request.
//
// The ID field of the passed discord.Guild must be set.
func (m *Mocker) Guild(g discord.Guild) {
	m.MockAPI("Guild", http.MethodGet, "/guilds/"+g.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, g)
		})
}

// GuildWithCount mocks a GuildWithCount request.
//
// The ID field of the passed discord.Guild must be set.
func (m *Mocker) GuildWithCount(g discord.Guild) {
	m.MockAPI("GuildWithCount", http.MethodGet, "/guilds/"+g.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckQuery(t, r.URL.Query(), url.Values{
				"with_counts": {"true"},
			})
			mockutil.WriteJSON(t, w, g)
		})
}

// GuildPreview mocks a GuildPreview request.
func (m *Mocker) GuildPreview(p discord.GuildPreview) {
	m.MockAPI("GuildPreview", http.MethodGet, "/guilds/"+p.ID.String()+"/preview",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, p)
		})
}

// Guilds mocks a Guilds request.
//
// This method will sanitize Guilds.ID, Guilds.OwnerID, Guilds.Emojis.ID and
// Guilds.Roles.ID.
func (m *Mocker) Guilds(limit uint, g []discord.Guild) {
	if g == nil {
		g = []discord.Guild{}
	}

	if len(g) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent guilds (%d vs. %d)", len(g), limit))
	}

	var after discord.GuildID

	for i := 0; i <= len(g)/maxFetchGuilds; i++ {
		var (
			from = uint(i) * maxFetchGuilds
			to   = uint(math.Min(float64(from+maxFetchGuilds), float64(len(g))))

			fetch = to - from // we use this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or mexFetchGuild, depending on which is smaller, instead.
			if fetch < maxFetchGuilds {
				fetch = uint(math.Min(float64(limit), float64(maxFetchGuilds)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchGuilds
			fetch = maxFetchGuilds
		}

		m.guildsRange(0, after, fmt.Sprintf("Guilds #%d", i+1), fetch, g[from:to])

		if to-from < maxFetchGuilds {
			break
		}

		after = g[to-1].ID
	}
}

// GuildsBefore mocks a GuildsBefore request.
//
// This method will sanitize Guilds.ID, Guilds.OwnerID, Guilds.Emojis.ID and
// Guilds.Roles.ID.
func (m *Mocker) GuildsBefore(before discord.GuildID, limit uint, g []discord.Guild) {
	if g == nil {
		g = []discord.Guild{}
	}

	if len(g) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent guilds (%d vs. %d)", len(g), limit))
	}

	req := len(g)/maxFetchGuilds + 1

	from := uint(math.Min(float64(uint(req)*maxFetchGuilds), float64(len(g))))

	for i := req; i > 0; i-- {
		no := req - i + 1

		to := from
		from = uint(math.Max(float64(0), float64(int(to-maxFetchGuilds))))

		fetch := to - from // we use this as the sent limit

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or mexFetchGuild, depending on which is smaller, instead.
			if fetch < maxFetchGuilds {
				fetch = uint(math.Min(float64(limit), float64(maxFetchGuilds)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchGuilds
			fetch = maxFetchGuilds
		}

		m.guildsRange(before, 0, fmt.Sprintf("GuildsBefore #%d", no), fetch, g[from:to])

		if to-from < maxFetchGuilds {
			break
		}

		before = g[from].ID
	}
}

// GuildsAfter mocks a GuildsAfter request.
//
// This method will sanitize Guilds.ID, Guilds.OwnerID, Guilds.Emojis.ID and
// Guilds.Roles.ID.WithToken
func (m *Mocker) GuildsAfter(after discord.GuildID, limit uint, g []discord.Guild) {
	if g == nil {
		g = []discord.Guild{}
	}

	if len(g) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent guilds (%d vs. %d)", len(g), limit))
	}

	for i := 0; i <= len(g)/maxFetchGuilds; i++ {
		var (
			from = uint(i) * maxFetchGuilds
			to   = uint(math.Min(float64(from+maxFetchGuilds), float64(len(g))))

			fetch = to - from // we use this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or maxFetchGuilds, depending on which is smaller, instead.
			if fetch < maxFetchGuilds {
				fetch = uint(math.Min(float64(limit), float64(maxFetchGuilds)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchGuilds
			fetch = maxFetchGuilds
		}

		m.guildsRange(0, after, fmt.Sprintf("GuildsAfter #%d", i+1), fetch, g[from:to])

		if to-from < maxFetchGuilds {
			break
		}

		after = g[to-1].ID
	}
}

// guildsRange mocks a single request to the GET /guilds endpoint.
func (m *Mocker) guildsRange(before, after discord.GuildID, name string, limit uint, g []discord.Guild) {
	m.MockAPI(name, http.MethodGet, "/users/@me/guilds",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"limit": {strconv.FormatUint(uint64(limit), 10)},
			}

			if after != 0 {
				expect["after"] = []string{after.String()}
			}

			if before != 0 {
				expect["before"] = []string{before.String()}
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			mockutil.WriteJSON(t, w, g)
		})
}

// LeaveGuild mocks a LeaveGuild request.
func (m *Mocker) LeaveGuild(id discord.GuildID) {
	m.MockAPI("LeaveGuild", http.MethodDelete, "/users/@me/guilds/"+id.String(), nil)
}

// ModifyGuild mocks a ModifyGuild request.
//
// The ID field of the passed discord.Guild must be set.
func (m *Mocker) ModifyGuild(d api.ModifyGuildData, g discord.Guild) {
	m.MockAPI("ModifyGuild", http.MethodPatch, "/guilds/"+g.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &d)
			mockutil.WriteJSON(t, w, g)
		})
}

// DeleteGuild mocks a DeleteGuild request.
func (m *Mocker) DeleteGuild(id discord.GuildID) {
	m.MockAPI("DeleteGuild", http.MethodDelete, "/guilds/"+id.String(), nil)
}

// VoiceRegionsGuild mocks a VoiceRegionsGuild request.
func (m *Mocker) VoiceRegionsGuild(guildID discord.GuildID, vr []discord.VoiceRegion) {
	if vr == nil {
		vr = []discord.VoiceRegion{}
	}

	m.MockAPI("VoiceRegionsGuild", http.MethodGet, "/guilds/"+guildID.String()+"/regions",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, vr)
		})
}

// AuditLog mocks a AuditLog request.
func (m *Mocker) AuditLog(guildID discord.GuildID, d api.AuditLogData, al discord.AuditLog) {
	switch {
	case d.Limit == 0:
		d.Limit = 50
	case d.Limit > 100:
		d.Limit = 100
	}

	m.MockAPI("AuditLog", http.MethodGet, "/guilds/"+guildID.String()+"/audit-logs",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"limit": {strconv.Itoa(int(d.Limit))},
			}

			if d.UserID != 0 {
				expect["user_id"] = []string{d.UserID.String()}
			}

			if d.ActionType != 0 {
				expect["action_type"] = []string{strconv.FormatUint(uint64(d.ActionType), 10)}
			}

			if d.Before != 0 {
				expect["before"] = []string{d.Before.String()}
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			mockutil.WriteJSON(t, w, al)
		})
}

// Integrations mocks a Integrations request.
func (m *Mocker) Integrations(guildID discord.GuildID, integrations []discord.Integration) {
	if integrations == nil {
		integrations = []discord.Integration{}
	}

	m.MockAPI("Integrations", http.MethodGet, "/guilds/"+guildID.String()+"/integrations",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, integrations)
		})
}

type attachIntegrationPayload struct {
	Type discord.Service       `json:"type"`
	ID   discord.IntegrationID `json:"id"`
}

// AttachIntegration mocks a AttachIntegration request.
func (m *Mocker) AttachIntegration(
	guildID discord.GuildID, integrationID discord.IntegrationID, integrationType discord.Service,
) {
	m.MockAPI("AttachIntegration", http.MethodPost, "/guilds/"+guildID.String()+"/integrations",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := &attachIntegrationPayload{
				Type: integrationType,
				ID:   integrationID,
			}

			mockutil.CheckJSON(t, r.Body, expect)
			w.WriteHeader(http.StatusNoContent)
		})
}

// ModifyIntegration mocks a ModifyIntegration request.
func (m *Mocker) ModifyIntegration(
	guildID discord.GuildID, integrationID discord.IntegrationID, d api.ModifyIntegrationData,
) {
	m.MockAPI("ModifyIntegration", http.MethodPatch,
		"/guilds/"+guildID.String()+"/integrations/"+integrationID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// SyncIntegration mocks a SyncIntegration request.
func (m *Mocker) SyncIntegration(guildID discord.GuildID, integrationID discord.IntegrationID) {
	m.MockAPI("SyncIntegration", http.MethodPost,
		"/guilds/"+guildID.String()+"/integrations/"+integrationID.String()+"/sync", nil)
}

// GuildWidgetSettings mocks a GuildWidgetSettings request.
func (m *Mocker) GuildWidgetSettings(guildID discord.GuildID, s discord.GuildWidgetSettings) {
	m.MockAPI("GuildWidgetSettings", http.MethodGet, "/guilds/"+guildID.String()+"/widget",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, s)
		})
}

// ModifyGuildWidget mocks a ModifyGuildWidget request.
func (m *Mocker) ModifyGuildWidget(
	guildID discord.GuildID, d api.ModifyGuildWidgetData, s discord.GuildWidgetSettings,
) {
	m.MockAPI("ModifyGuildWidget", http.MethodPatch, "/guilds/"+guildID.String()+"/widget",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &d)
			mockutil.WriteJSON(t, w, s)
		})
}

// GuildWidget mocks a GuildWidget request.
func (m *Mocker) GuildWidget(guildID discord.GuildID, widget discord.GuildWidget) {
	m.MockAPI("GuildWidgetSettings", http.MethodGet, "/guilds/"+guildID.String()+"/widget.json",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, widget)
		})
}

// GuildVanityURL mocks a GuildVanityURL request.
func (m *Mocker) GuildVanityInvite(guildID discord.GuildID, i discord.Invite) {
	m.MockAPI("GuildVanityURL", http.MethodGet, "/guilds/"+guildID.String()+"/vanity-url",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, i)
		})
}

// GuildWidgetImage mocks a GuildWidgetImage request.
func (m *Mocker) GuildWidgetImage(guildID discord.GuildID, style api.GuildWidgetImageStyle, img io.Reader) {
	m.MockAPI("GuildWidgetImage", http.MethodGet, "/guilds/"+guildID.String()+"/widget.png",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckQuery(t, r.URL.Query(), url.Values{
				"style": {string(style)},
			})

			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}
