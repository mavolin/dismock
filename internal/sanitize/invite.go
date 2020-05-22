package sanitize

import "github.com/diamondburned/arikawa/discord"

// Invite sanitizes an Invite.
//
// This function will sanitize Invite.Guild.ID, Invite.Guild.OwnerID,
// Invite.Guild.Emojis.ID, Invite.Guild.Roles.ID, Invite.Channel.ID,
// Invite.Inviter.ID.
func Invite(i discord.Invite, guildID, ownerID, channelID, inviterID, targetID discord.Snowflake) discord.Invite {
	if i.Guild != nil {
		*i.Guild = Guild(*i.Guild, guildID, ownerID)
	}

	i.Channel = Channel(i.Channel, channelID)

	if i.Inviter != nil {
		*i.Inviter = User(*i.Inviter, inviterID)
	}

	if i.Target != nil {
		*i.Target = User(*i.Target, targetID)
	}

	return i
}
