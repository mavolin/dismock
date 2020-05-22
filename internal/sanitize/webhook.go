package sanitize

import "github.com/diamondburned/arikawa/discord"

// Webhook sanitizes a Webhook.
//
// This method will sanitize Webhook.ID and Webhook.User.ID.
func Webhook(w discord.Webhook, id, userID discord.Snowflake) discord.Webhook {
	if w.ID <= 0 {
		w.ID = id
	}

	w.User = User(w.User, userID)

	return w
}
