package dismock

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/dismock/v2/internal/mockutil"
)

const maxFetchMembers = 1000

// Member mocks a Member request.
//
// The User.ID field of the passed member must be set.
func (m *Mocker) Member(guildID discord.GuildID, member discord.Member) {
	m.MockAPI("Member", http.MethodGet, "/guilds/"+guildID.String()+"/members/"+member.User.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, member)
		})
}

// Members mocks a Members request.
func (m *Mocker) Members(guildID discord.GuildID, limit uint, members []discord.Member) {
	if members == nil {
		members = []discord.Member{}
	}

	if len(members) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent Members (%d vs. %d)", len(members), limit))
	}

	var after discord.UserID

	for i := 0; i <= len(members)/maxFetchMembers; i++ {
		var (
			from = uint(i) * maxFetchMembers
			to   = uint(math.Min(float64(from+maxFetchMembers), float64(len(members))))

			fetch = to - from // we use this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or maxFetchMembers, depending on which is smaller, instead.
			if fetch < maxFetchMembers {
				fetch = uint(math.Min(float64(limit), float64(maxFetchMembers)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchMembers
			fetch = maxFetchMembers
		}

		m.membersAfter(guildID, after, fmt.Sprintf("Members #%d", i+1), fetch, members[from:to])

		if to-from < maxFetchMembers {
			break
		}

		after = members[to-1].User.ID
	}
}

// MembersAfter mocks a MembersAfter request.
func (m *Mocker) MembersAfter(guildID discord.GuildID, after discord.UserID, limit uint, members []discord.Member) {
	if members == nil {
		members = []discord.Member{}
	}

	if len(members) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent Members (%d vs. %d)", len(members), limit))
	}

	for i := 0; i <= len(members)/maxFetchMembers; i++ {
		var (
			from = uint(i) * maxFetchMembers
			to   = uint(math.Min(float64(from+maxFetchMembers), float64(len(members))))

			fetch = to - from // we use this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// use either limit or maxFetchMembers, depending on which is smaller, instead.
			if fetch < maxFetchMembers {
				fetch = uint(math.Min(float64(limit), float64(maxFetchMembers)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should use
			// maxFetchMembers
			fetch = maxFetchMembers
		}

		m.membersAfter(guildID, after, fmt.Sprintf("MembersAfter #%d", i+1), fetch, members[from:to])

		if to-from < maxFetchMembers {
			break
		}

		after = members[to-1].User.ID
	}
}

// membersAfter mocks a single request to the GET /Members endpoint.
func (m *Mocker) membersAfter(
	guildID discord.GuildID, after discord.UserID, name string, limit uint, g []discord.Member,
) {
	m.MockAPI(name, http.MethodGet, "/guilds/"+guildID.String()+"/members",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"limit": {strconv.FormatUint(uint64(limit), 10)},
			}

			if after != 0 {
				expect["after"] = []string{after.String()}
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			mockutil.WriteJSON(t, w, g)
		})
}

// AddMember mocks a AddMember request.
//
// The User.ID field of the passed discord.Member must be set.
func (m *Mocker) AddMember(guildID discord.GuildID, d api.AddMemberData, member discord.Member) {
	m.MockAPI("AddMember", http.MethodPut, "/guilds/"+guildID.String()+"/members/"+member.User.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &d)
			mockutil.WriteJSON(t, w, member)
		})
}

// ModifyMember mocks a ModifyMember request.
func (m *Mocker) ModifyMember(guildID discord.GuildID, userID discord.UserID, d api.ModifyMemberData) {
	m.MockAPI("ModifyMember", http.MethodPatch, "/guilds/"+guildID.String()+"/members/"+userID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

type pruneBody struct {
	Pruned uint `json:"pruned"`
}

// PruneCount mocks a PruneCount request.
func (m *Mocker) PruneCount(guildID discord.GuildID, d api.PruneCountData, pruned uint) {
	if d.Days == 0 {
		d.Days = 7
	}

	m.MockAPI("PruneCount", http.MethodGet, "/guilds/"+guildID.String()+"/prune",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"days": {strconv.Itoa(int(d.Days))},
			}

			if len(d.IncludedRoles) > 0 {
				expect["include_roles"] = make([]string, len(d.IncludedRoles))

				for i, r := range d.IncludedRoles {
					expect["include_roles"][i] = r.String()
				}
			}

			resp := pruneBody{
				Pruned: pruned,
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			mockutil.WriteJSON(t, w, resp)
		})
}

// Prune mocks a Prune request.
func (m *Mocker) Prune(guildID discord.GuildID, d api.PruneData, pruned uint) {
	if d.Days == 0 {
		d.Days = 7
	}

	m.MockAPI("Prune", http.MethodPost, "/guilds/"+guildID.String()+"/prune",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"days":                {strconv.Itoa(int(d.Days))},
				"compute_prune_count": {strconv.FormatBool(d.ReturnCount)},
			}

			if len(d.IncludedRoles) > 0 {
				expect["include_roles"] = make([]string, len(d.IncludedRoles))

				for i, r := range d.IncludedRoles {
					expect["include_roles"][i] = r.String()
				}
			}

			resp := pruneBody{
				Pruned: pruned,
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			mockutil.WriteJSON(t, w, resp)
		})
}

// Kick mocks a Kick request.
func (m *Mocker) Kick(guildID discord.GuildID, userID discord.UserID) {
	m.MockAPI("Kick", http.MethodDelete, "/guilds/"+guildID.String()+"/members/"+userID.String(), nil)
}

// Bans mocks a Bans request.
func (m *Mocker) Bans(guildID discord.GuildID, b []discord.Ban) {
	m.MockAPI("Bans", http.MethodGet, "/guilds/"+guildID.String()+"/bans",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, b)
		})
}

// GetBan mocks a GetBan request.
//
// The User.ID field of the passed discord.Ban must be set.
func (m *Mocker) GetBan(guildID discord.GuildID, b discord.Ban) {
	m.MockAPI("GetBan", http.MethodGet, "/guilds/"+guildID.String()+"/bans/"+b.User.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, b)
		})
}

// Ban mocks a Ban request.
func (m *Mocker) Ban(guildID discord.GuildID, userID discord.UserID, d api.BanData) {
	if *d.DeleteDays > 7 {
		*d.DeleteDays = 7
	}

	m.MockAPI("Ban", http.MethodPut, "/guilds/"+guildID.String()+"/bans/"+userID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := make(url.Values)

			if d.DeleteDays != nil {
				expect["delete_message_days"] = []string{strconv.Itoa(int(*d.DeleteDays))}
			}

			if d.Reason != nil {
				expect["reason"] = []string{*d.Reason}
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			w.WriteHeader(http.StatusNoContent)
		})
}

// Unban mocks a Unban request.
func (m *Mocker) Unban(guildID discord.GuildID, userID discord.UserID) {
	m.MockAPI("Unban", http.MethodDelete, "/guilds/"+guildID.String()+"/bans/"+userID.String(), nil)
}
