package multiplexer

import "github.com/bwmarrin/discordgo"

// Multiplexer represents the event router.
type Multiplexer struct {
	// Prefix is the default command prefix.
	Prefix string

	// Routes is a slice of pointers to command routes.
	Routes []*Route
	// Categories is a slice of pointers to CommandCategory.
	Categories []*CommandCategory

	// EventHandlers is a slice of event handler functions registered to the library directly
	EventHandlers []interface{}

	NotTargeted           []func(context *Context)
	Ready                 []func(context *Context)
	GuildMemberAdd        []func(context *Context)
	GuildMemberRemove     []func(context *Context)
	GuildDelete           []func(context *Context)
	MessageCreate         []func(context *Context)
	MessageDelete         []func(context *Context)
	MessageUpdate         []func(context *Context)
	MessageReactionAdd    []func(context *Context)
	MessageReactionRemove []func(context *Context)
	VoiceStateUpdate      []func(context *Context)

	// Administrator is the privileged administrator user with all privilege overrides and full access to all commands.
	Administrator *discordgo.User
	// Operator is a slice of operator users with all privilege overrides and access to some restricted commands.
	Operator []*discordgo.User
}

// Route registers a route to the router.
func (mux *Multiplexer) Route(route *Route) *Route {
	route.Category.Routes = append(route.Category.Routes, route)
	mux.Routes = append(mux.Routes, route)
	return route
}

func (mux *Multiplexer) SessionRegisterHandlers(session *discordgo.Session) {
	for _, handler := range mux.EventHandlers {
		session.AddHandler(handler)
	}
}

// IsOperator checks of a user is an operator.
func (mux *Multiplexer) IsOperator(id string) bool {
	if id == mux.Administrator.ID {
		return true
	}
	for _, operator := range mux.Operator {
		if id == operator.ID {
			return true
		}
	}
	return false
}

// IsAdministrator checks of a user is the system administrator.
func (mux *Multiplexer) IsAdministrator(id string) bool {
	return id == mux.Administrator.ID
}
