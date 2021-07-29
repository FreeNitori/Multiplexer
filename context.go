package multiplexer

import (
	"errors"
	"git.randomchars.net/freenitori/embedutil"
	"git.randomchars.net/freenitori/log"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

// ErrUserNotFound represents the error returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// GetPrefix is the function used to get a prefix of a guild.
var GetPrefix = func(context *Context) string {
	return context.Multiplexer.Prefix
}

// Context carries an event's information.
type Context struct {
	Multiplexer       *Multiplexer
	User              *discordgo.User
	Member            *discordgo.Member
	Message           *discordgo.Message
	Session           *discordgo.Session
	Guild             *discordgo.Guild
	Channel           *discordgo.Channel
	Event             interface{}
	Text              string
	Fields            []string
	IsPrivate         bool
	IsTargeted        bool
	HasPrefix         bool
	HasMention        bool
	HasLeadingMention bool
}

var numericalRegex = regexp.MustCompile("[^0-9]+")

// NumericalRegex returns a compiled regular expression that matches only numbers.
func (Context) NumericalRegex() *regexp.Regexp {
	return numericalRegex
}

// SendMessage sends a text message in the current channel and returns the message.
func (context *Context) SendMessage(message string) *discordgo.Message {
	permissions, err := context.Session.State.UserChannelPermissions(context.Session.State.User.ID, context.Message.ChannelID)
	if !(err == nil && (permissions&discordgo.PermissionSendMessages == discordgo.PermissionSendMessages)) {
		return nil
	}

	resultMessage, err := context.Session.ChannelMessageSend(context.Message.ChannelID, message)
	if err != nil {
		log.Errorf("Error while sending message to guild %s, %s", context.Message.GuildID, err)
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			ErrorOccurred)
		return nil
	}
	return resultMessage
}

// SendEmbed sends an embedutil message in the current channel and returns the message.
func (context *Context) SendEmbed(message string, embed embedutil.Embed) *discordgo.Message {
	var err error
	permissions, err := context.Session.State.UserChannelPermissions(context.Session.State.User.ID, context.Message.ChannelID)
	if !(err == nil && (permissions&discordgo.PermissionSendMessages == discordgo.PermissionSendMessages)) {
		return nil
	}

	var resultMessage *discordgo.Message
	if message == "" {
		resultMessage, err = context.Session.ChannelMessageSendEmbed(context.Message.ChannelID, embed.MessageEmbed)
	} else {
		resultMessage, err = context.Session.ChannelMessageSendComplex(context.Message.ChannelID, &discordgo.MessageSend{
			Content:         message,
			Embed:           embed.MessageEmbed,
			TTS:             false,
			Files:           nil,
			AllowedMentions: nil,
			File:            nil,
		})
	}
	if err != nil {
		log.Errorf("Error while sending embedutil to guild %s, %s", context.Message.GuildID, err)
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			ErrorOccurred)
		return nil
	}
	return resultMessage
}

// HandleError handles a returned error and send the information of it if in debug mode.
func (context *Context) HandleError(err error) bool {
	if err != nil {
		log.Errorf("Error occurred while executing command, %s", err)
		context.SendMessage(ErrorOccurred)
		if log.GetLevel() == logrus.DebugLevel {
			context.SendMessage(err.Error())
		}
		return false
	}
	return true
}

// HasPermission checks a user for a permission.
func (context *Context) HasPermission(permission int) bool {
	// Override check for operators and system administrators
	if context.User.ID == context.Multiplexer.Administrator.ID {
		return true
	}
	for _, user := range context.Multiplexer.Operator {
		if context.User.ID == user.ID {
			return true
		}
	}

	// Check against the user
	permissions, err := context.Session.State.UserChannelPermissions(context.User.ID, context.Message.ChannelID)
	return err == nil && (int(permissions)&permission == permission)
}

// IsOperator checks of a user is an operator.
func (context *Context) IsOperator() bool {
	return context.Multiplexer.IsOperator(context.User.ID)
}

// IsAdministrator checks of a user is the system administrator.
func (context *Context) IsAdministrator() bool {
	return context.Multiplexer.IsAdministrator(context.User.ID)
}

// GetMember gets a member from a string representing it.
func (context *Context) GetMember(query string) *discordgo.Member {
	// Guild only function
	if context.IsPrivate {
		return nil
	}

	// Check if it's a mention or the string is numerical
	_, err := strconv.Atoi(query)
	if strings.HasPrefix(query, "<@") && strings.HasSuffix(query, ">") || err == nil {
		// Strip off the mention thingy
		userID := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(query, "<@"), "!"), ">")
		// Length of a real snowflake after stripping off stuff
		if len(userID) == 18 {
			for _, member := range context.Guild.Members {
				if member.User.ID == userID {
					return member
				}
			}
		}
	} else {
		// Find as username or nickname
		for _, member := range context.Guild.Members {
			if member.User.Username == query || member.Nick == query {
				return member
			}
		}
	}
	return nil
}

// GetChannel gets a channel from a string representing it.
func (context *Context) GetChannel(query string) *discordgo.Channel {
	// Guild only function
	if context.IsPrivate {
		return nil
	}

	// Check if it's a mention or the string is numerical
	_, err := strconv.Atoi(query)
	if strings.HasPrefix(query, "<#") && strings.HasSuffix(query, ">") || err == nil {
		// Strip off the mention thingy
		channelID := numericalRegex.ReplaceAllString(query, "")
		// Length of a real snowflake after stripping off stuff
		if len(channelID) == 18 {
			for _, channel := range context.Guild.Channels {
				if channel.ID == channelID {
					return channel
				}
			}
		}
	} else {
		// Find as channel name
		for _, channel := range context.Guild.Channels {
			if channel.Name == query {
				return channel
			}
		}
	}
	return nil
}

// GetRole gets a channel from a string representing it.
func (context *Context) GetRole(query string) *discordgo.Role {
	// Guild only function
	if context.IsPrivate {
		return nil
	}

	// Check if it's a mention or the string is numerical
	_, err := strconv.Atoi(query)
	if strings.HasPrefix(query, "<@&") && strings.HasSuffix(query, ">") || err == nil {
		// Strip off the mention thingy
		roleID := numericalRegex.ReplaceAllString(query, "")
		// Length of a real snowflake after stripping off stuff
		if len(roleID) == 18 {
			for _, role := range context.Guild.Roles {
				if role.ID == roleID {
					return role
				}
			}
		}
	} else {
		// Find as channel name
		for _, role := range context.Guild.Roles {
			if role.Name == query {
				return role
			}
		}
	}
	return nil
}

// StitchFields stitches together fields of the message.
func (context *Context) StitchFields(start int) string {
	if len(context.Fields) <= start {
		return ""
	}
	message := context.Fields[start]
	for i := start + 1; i < len(context.Fields); i++ {
		message += " " + context.Fields[i]
	}
	return message
}

// Prefix returns the command prefix of a context.
func (context *Context) Prefix() string {
	if context.IsPrivate {
		return context.Multiplexer.Prefix
	} else {
		return GetPrefix(context)
	}
}

// GetVoiceState returns the voice state of a user if found.
func (context *Context) GetVoiceState() (*discordgo.VoiceState, bool) {
	if context.IsPrivate {
		return nil, false
	}
	for _, voiceState := range context.Guild.VoiceStates {
		if voiceState.UserID == context.User.ID {
			return voiceState, true
		}
	}
	return nil, false
}

// MakeVoiceConnection returns the voice connection to a user's voice channel if join-able.
func (context *Context) MakeVoiceConnection() (*discordgo.VoiceConnection, error) {
	if context.IsPrivate {
		return nil, nil
	}
	voiceState, ok := context.GetVoiceState()
	if !ok {
		return nil, nil
	}
	return context.Session.ChannelVoiceJoin(voiceState.GuildID, voiceState.ChannelID, false, true)
}

// Ban creates a ban on the specified user.
func (context *Context) Ban(query string) error {
	// If Nitori has permission
	permissions, err := context.Session.State.UserChannelPermissions(context.Session.State.User.ID, context.Message.ChannelID)
	if !(err == nil && (permissions&discordgo.PermissionBanMembers == discordgo.PermissionBanMembers)) {
		return discordgo.ErrUnauthorized
	}

	// Check if it's a mention or the string is numerical
	_, err = strconv.Atoi(query)
	if strings.HasPrefix(query, "<@") && strings.HasSuffix(query, ">") || err == nil {
		// Strip off the mention thingy
		userID := context.NumericalRegex().ReplaceAllString(query, "")
		// Length of a real snowflake after stripping off stuff
		if len(userID) == 18 {
			err = context.Session.GuildBanCreate(context.Guild.ID, userID, 0)
			return err
		}
	} else {
		member := context.GetMember(query)
		if member == nil {
			return ErrUserNotFound
		}
		err = context.Session.GuildBanCreate(context.Guild.ID, member.User.ID, 0)
		return err
	}
	return ErrUserNotFound
}
