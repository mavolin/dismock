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

// React mocks a React request.
func (m *Mocker) React(channelID discord.ChannelID, messageID discord.MessageID, e discord.APIEmoji) {
	m.MockAPI("React", http.MethodPut,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+e.PathString()+"/@me", nil)
}

// Unreact mocks a Unreact request.
func (m *Mocker) Unreact(channelID discord.ChannelID, messageID discord.MessageID, e discord.APIEmoji) {
	m.MockAPI("Unreact", http.MethodDelete,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+e.PathString()+"/@me", nil)
}

// Reactions mocks a Reactions request.
func (m *Mocker) Reactions(
	channelID discord.ChannelID, messageID discord.MessageID, limit uint, e discord.APIEmoji, u []discord.User,
) {
	if u == nil {
		u = []discord.User{}
	}

	if len(u) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent users (%d vs. %d)", len(u), limit))
	}

	var after discord.UserID

	for i := 0; i <= len(u)/api.MaxMessageReactionFetchLimit; i++ {
		var (
			from = uint(i) * api.MaxMessageReactionFetchLimit
			to   = uint(math.Min(float64(from+api.MaxMessageReactionFetchLimit), float64(len(u))))

			fetch = to - from // we msg this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// msg either limit or api.MaxMessageReactionFetchLimit, depending on which is smaller, instead.
			if fetch < api.MaxMessageReactionFetchLimit {
				fetch = uint(math.Min(float64(limit), float64(api.MaxMessageReactionFetchLimit)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should msg
			// api.MaxMessageReactionFetchLimit
			fetch = api.MaxMessageReactionFetchLimit
		}

		m.reactionsRange(channelID, messageID, 0, after, fmt.Sprintf("Reactions #%d", i+1), fetch, e, u[from:to])

		if to-from < api.MaxMessageReactionFetchLimit {
			break
		}

		after = u[to-1].ID
	}
}

// ReactionsBefore mocks a ReactionsBefore request.
func (m *Mocker) ReactionsBefore(
	channelID discord.ChannelID, messageID discord.MessageID, before discord.UserID, limit uint, e discord.APIEmoji,
	u []discord.User,
) {
	if u == nil {
		u = []discord.User{}
	}

	if len(u) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent users (%d vs. %d)", len(u), limit))
	}

	req := len(u)/api.MaxMessageReactionFetchLimit + 1

	from := uint(math.Min(float64(uint(req)*api.MaxMessageReactionFetchLimit), float64(len(u))))

	for i := req; i > 0; i-- {
		no := req - i + 1

		to := from
		from = uint(math.Max(float64(0), float64(int(to-api.MaxMessageReactionFetchLimit))))

		fetch := to - from // we msg this as the sent limit

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// msg either limit or api.MaxMessageReactionFetchLimit, depending on which is smaller, instead.
			if fetch < api.MaxMessageReactionFetchLimit {
				fetch = uint(math.Min(float64(limit), float64(api.MaxMessageReactionFetchLimit)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should msg
			// api.MaxMessageReactionFetchLimit
			fetch = api.MaxMessageReactionFetchLimit
		}

		m.reactionsRange(channelID, messageID, before, 0, fmt.Sprintf("ReactionsBefore #%d", no), fetch, e, u[from:to])

		if to-from < api.MaxMessageReactionFetchLimit {
			break
		}

		before = u[from].ID
	}
}

// ReactionsAfter mocks a ReactionsAfter request.
func (m *Mocker) ReactionsAfter(
	channelID discord.ChannelID, messageID discord.MessageID, after discord.UserID, limit uint, e discord.APIEmoji,
	u []discord.User,
) {
	if u == nil {
		u = []discord.User{}
	}

	if len(u) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent users (%d vs. %d)", len(u), limit))
	}

	for i := 0; i <= len(u)/api.MaxMessageReactionFetchLimit; i++ {
		var (
			from = uint(i) * api.MaxMessageReactionFetchLimit
			to   = uint(math.Min(float64(from+api.MaxMessageReactionFetchLimit), float64(len(u))))

			fetch = to - from // we msg this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// msg either limit or api.MaxMessageReactionFetchLimit, depending on which is smaller, instead.
			if fetch < api.MaxMessageReactionFetchLimit {
				fetch = uint(math.Min(float64(limit), float64(api.MaxMessageReactionFetchLimit)))
			}

			limit -= fetch
		} else {
			// this means there is no limit, hence we should msg
			// api.MaxMessageReactionFetchLimit
			fetch = api.MaxMessageReactionFetchLimit
		}

		m.reactionsRange(channelID, messageID, 0, after, fmt.Sprintf("ReactionsAfter #%d", i+1), fetch, e, u[from:to])

		if to-from < api.MaxMessageReactionFetchLimit {
			break
		}

		after = u[to-1].ID
	}
}

// reactionsRange mocks a single request to the GET /reactions endpoint.
func (m *Mocker) reactionsRange(
	channelID discord.ChannelID, messageID discord.MessageID, before, after discord.UserID, name string, limit uint,
	e discord.APIEmoji, u []discord.User,
) {
	m.MockAPI(name, http.MethodGet,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+e.PathString(),
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
			mockutil.WriteJSON(t, w, u)
		})
}

// DeleteUserReaction mocks a DeleteUserReaction request.
func (m *Mocker) DeleteUserReaction(
	channelID discord.ChannelID, messageID discord.MessageID, userID discord.UserID, e discord.APIEmoji,
) {
	user := "@me"
	if userID > 0 {
		user = userID.String()
	}

	m.MockAPI("DeleteUserReaction", http.MethodDelete,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+e.PathString()+"/"+user, nil)
}

// DeleteReactions mocks a DeleteReactions request.
func (m *Mocker) DeleteReactions(channelID discord.ChannelID, messageID discord.MessageID, e discord.APIEmoji) {
	m.MockAPI("DeleteReactions", http.MethodDelete,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+e.PathString(), nil)
}

// DeleteAllReactions mocks a DeleteAllReactions request.
func (m *Mocker) DeleteAllReactions(channelID discord.ChannelID, messageID discord.MessageID) {
	m.MockAPI("DeleteAllReactions", http.MethodDelete,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions", nil)
}
