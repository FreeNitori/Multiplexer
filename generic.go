package multiplexer

import "github.com/bwmarrin/discordgo"

// Event handler that fires when ready
func (mux *Multiplexer) onReady(session *discordgo.Session, ready *discordgo.Ready) {
	go func() {
		for _, hook := range mux.Ready {
			hook(&Context{
				Multiplexer: mux,
				User:        session.State.User,
				Session:     session,
				Event:       ready,
			})
		}
	}()
	return
}

// Event handler that fires when a guild member is added
func (mux *Multiplexer) onGuildMemberAdd(session *discordgo.Session, add *discordgo.GuildMemberAdd) {
	go func() {
		for _, hook := range mux.GuildMemberAdd {
			guild := GetGuild(session, add.GuildID)
			if guild == nil {
				return
			}
			hook(&Context{
				Multiplexer: mux,
				Member:      add.Member,
				User:        add.Member.User,
				Session:     session,
				Guild:       guild,
				Event:       add,
			})
		}
	}()
	return
}

// Event handler that fires when a guild member is removed
func (mux *Multiplexer) onGuildMemberRemove(session *discordgo.Session, remove *discordgo.GuildMemberRemove) {
	go func() {
		for _, hook := range mux.GuildMemberRemove {
			guild := GetGuild(session, remove.GuildID)
			if guild == nil {
				return
			}
			hook(&Context{
				Multiplexer: mux,
				Member:      remove.Member,
				User:        remove.Member.User,
				Session:     session,
				Guild:       guild,
				Event:       remove,
			})
		}
	}()
	return
}

// Event handler that fires when a guild is deleted
func (mux *Multiplexer) onGuildDelete(session *discordgo.Session, delete *discordgo.GuildDelete) {
	go func() {
		for _, hook := range mux.GuildDelete {
			hook(&Context{
				Multiplexer: mux,
				Session:     session,
				Guild:       delete.Guild,
				Event:       delete,
			})
		}
	}()
	return
}

// Event handler that fires when a message is created
func (mux *Multiplexer) onMessageCreate(session *discordgo.Session, create *discordgo.MessageCreate) {
	go func() {
		for _, hook := range mux.MessageCreate {
			context := mux.NewContextMessage(session, create.Message, create)
			if context == nil {
				return
			}
			hook(context)
		}
	}()
	return
}

// Event handler that fires when a message is deleted
func (mux *Multiplexer) onMessageDelete(session *discordgo.Session, delete *discordgo.MessageDelete) {
	go func() {
		for _, hook := range mux.MessageDelete {
			hook(&Context{
				Multiplexer: mux,
				Message:     delete.Message,
				Session:     session,
				Event:       delete,
			})
		}
	}()
	return
}

// Event handler that fires when a message is updated
func (mux *Multiplexer) onMessageUpdate(session *discordgo.Session, update *discordgo.MessageUpdate) {
	go func() {
		for _, hook := range mux.MessageUpdate {
			hook(&Context{
				Multiplexer: mux,
				Message:     update.Message,
				Session:     session,
				Event:       update,
			})
		}
	}()
	return
}

// Event handler that fires when a reaction is added
func (mux *Multiplexer) onMessageReactionAdd(session *discordgo.Session, add *discordgo.MessageReactionAdd) {
	go func() {
		for _, hook := range mux.MessageReactionAdd {
			message, err := session.ChannelMessage(add.ChannelID, add.MessageID)
			if err != nil {
				log.Errorf("Unable to get message %s from channel %s, %s", add.MessageID, add.ChannelID, err)
				return
			}
			context := mux.NewContextMessage(session, message, add)
			if context == nil {
				return
			}
			hook(context)
		}
	}()
	return
}

// Event handler that fires when a reaction is removed
func (mux *Multiplexer) onMessageReactionRemove(session *discordgo.Session, remove *discordgo.MessageReactionRemove) {
	go func() {
		for _, hook := range mux.MessageReactionRemove {
			message, err := session.ChannelMessage(remove.ChannelID, remove.MessageID)
			if err != nil {
				log.Errorf("Unable to get message %s from channel %s, %s", remove.MessageID, remove.ChannelID, err)
				return
			}
			context := mux.NewContextMessage(session, message, remove)
			if context == nil {
				return
			}
			hook(context)
		}
	}()
}
