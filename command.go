package multiplexer

import (
	"git.randomchars.net/freenitori/log"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

// NoCommandMatched is called when no command is matched.
var NoCommandMatched = func(context *Context) {
	context.SendMessage("Command not found.")
}

// CommandHandler represents the handler function of a Route.
type CommandHandler func(*Context)

// Route represents a command route.
type Route struct {
	Pattern       string
	AliasPatterns []string
	Description   string
	Category      *CommandCategory
	Handler       CommandHandler
}

// CommandCategory represents a category of Route.
type CommandCategory struct {
	Routes      []*Route
	Title       string
	Description string
}

// MatchRoute fuzzy matches a message to a route.
func (mux *Multiplexer) MatchRoute(message string) (*Route, []string) {
	fields := strings.Fields(message)
	if len(fields) == 0 {
		return nil, nil
	}

	var route *Route
	var similarityRating int
	var fieldIndex int

	for fieldIndex, fieldIter := range fields {
		for _, routeIter := range mux.Routes {
			if routeIter.Pattern == fieldIter {
				return routeIter, fields[fieldIndex:]
			}
			for _, aliasPattern := range routeIter.AliasPatterns {
				if aliasPattern == fieldIter {
					return routeIter, fields[fieldIndex:]
				}
			}
			if strings.HasPrefix(routeIter.Pattern, fieldIter) {
				if len(fieldIter) > similarityRating {
					route = routeIter
					similarityRating = len(fieldIter)
				}
			}
		}
	}
	return route, fields[fieldIndex:]
}

func (mux *Multiplexer) handleMessageCommand(session *discordgo.Session, create *discordgo.MessageCreate) {

	// Ignore self and bot messages
	if create.Author.ID == session.State.User.ID || create.Author.Bot {
		return
	}

	// Make Context
	context := mux.NewContextMessage(session, create.Message, create)
	if context == nil {
		return
	}

	// Call not targeted hooks and return
	if !context.IsTargeted {
		go func() {
			for _, hook := range mux.NotTargeted {
				hook(context)
			}
		}()
		return
	}

	// Log the processed message
	var hostName string
	if context.IsPrivate {
		hostName = "Private Messages"
	} else {
		hostName = "\"" + context.Guild.Name + "\""
	}
	log.Infof("(Shard %s) \"%s\"@%s > %s",
		strconv.Itoa(session.ShardID),
		context.User.Username+"#"+context.User.Discriminator,
		hostName,
		context.Message.Content)

	// Figure out the route of the message
	if !(context.HasMention && !context.HasLeadingMention) {
		route, fields := mux.MatchRoute(context.Text)
		if route != nil {
			context.Fields = fields
			route.Handler(context)
			return
		}
	}

	NoCommandMatched(context)
}
