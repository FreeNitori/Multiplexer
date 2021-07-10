package multiplexer

import (
	"fmt"
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strings"
)

// New returns a command router.
func New() *Multiplexer {
	mux := &Multiplexer{
		Categories: []*CommandCategory{AudioCategory, ExperienceCategory, ManualsCategory, MediaCategory, ModerationCategory, SystemCategory},
	}
	mux.EventHandlers = []interface{}{
		mux.handleMessageCommand,
		mux.onReady,
		mux.onGuildMemberAdd,
		mux.onGuildMemberRemove,
		mux.onGuildDelete,
		mux.onMessageCreate,
		mux.onMessageDelete,
		mux.onMessageUpdate,
		mux.onMessageReactionAdd,
		mux.onMessageReactionRemove,
		mux.onVoiceStateUpdate,
	}
	return mux
}

// NewCategory returns a new command category
func NewCategory(name string, description string) *CommandCategory {
	cat := &CommandCategory{
		Title:       name,
		Description: description,
	}
	return cat
}

// NewContextMessage returns pointer to Context generated from a message.
func (mux *Multiplexer) NewContextMessage(session *discordgo.Session, message *discordgo.Message, event interface{}) *Context {
	if message.Author.ID == session.State.User.ID {
		return nil
	}

	guild := GetGuild(session, message.GuildID)
	if guild == nil {
		guild = &discordgo.Guild{}
	}

	channel := GetChannel(session, message.ChannelID)
	if channel == nil {
		log.Errorf("Unable to obtain channel when making Context.")
		return nil
	}

	context := &Context{
		Multiplexer: mux,
		User:        message.Author,
		Message:     message,
		Session:     session,
		Guild:       guild,
		Channel:     channel,
		Event:       event,
		Text:        strings.TrimSpace(message.Content),
		Fields:      nil,
		IsPrivate:   channel.Type == discordgo.ChannelTypeDM,
	}

	// Get guild-specific prefix
	guildPrefix := context.Prefix()

	// Look for ping
	for _, mentionedUser := range message.Mentions {
		if mentionedUser.ID == session.State.User.ID {
			context.IsTargeted, context.HasMention = true, true
			mentionRegex := regexp.MustCompile(fmt.Sprintf("<@!?(%s)>", session.State.User.ID))

			// If message have leading mention
			location := mentionRegex.FindStringIndex(context.Text)
			if len(location) == 0 {
				context.HasLeadingMention = true
			} else if location[0] == 0 {
				context.HasLeadingMention = true
			}

			// Remove the mention string
			context.Text = mentionRegex.ReplaceAllString(context.Text, "")

			break
		}
	}

	// Command prefix included or not
	if !context.IsTargeted && len(guildPrefix) > 0 {
		if strings.HasPrefix(context.Text, guildPrefix) {
			context.IsTargeted, context.HasPrefix = true, true
			context.Text = strings.TrimPrefix(context.Text, guildPrefix)
		}
	}

	if !context.IsPrivate {
		context.Member = message.Member
	}
	return context
}
