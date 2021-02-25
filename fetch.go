package multiplexer

import "github.com/bwmarrin/discordgo"

// GetGuild fetches guild from cache then API and returns nil if all fails.
func GetGuild(session *discordgo.Session, id string) *discordgo.Guild {
	var err error
	var guild *discordgo.Guild
	if id != "" {
		guild, err = session.State.Guild(id)
		if err != nil {
			// Attempt direct API fetching
			guild, err = session.Guild(id)
			if err != nil {
				log.Errorf("Unable to fetch guild from API or cache, %s", err)
				return nil
			}
			// Attempt caching the channel
			err = session.State.GuildAdd(guild)
			if err != nil {
				log.Warnf("Unable to cache guild fetched from API, %s", err)
			}
		}
	}
	return guild
}

// GetChannel fetches channel from cache then API and returns nil if all fails.
func GetChannel(session *discordgo.Session, id string) *discordgo.Channel {
	var err error
	var channel *discordgo.Channel
	channel, err = session.State.Channel(id)
	if err != nil {
		// Attempt direct API fetching
		channel, err = session.Channel(id)
		if err != nil {
			log.Errorf("Unable to fetch channel from API or cache, %s", err)
			return nil
		}
		// Attempt caching the channel
		err = session.State.ChannelAdd(channel)
		if err != nil {
			log.Warnf("Unable to cache channel fetched from API, %s", err)
		}
	}
	return channel
}
