package multiplexer

// InvalidArgument is the message sent when the user passes an invalid argument.
const InvalidArgument = "Invalid argument."

// ErrorOccurred is the message sent when the event handler catches an error.
const ErrorOccurred = "Something went wrong and I am very confused! Please try again!"

// GuildOnly is the message sent when a guild-only command is issued in private.
const GuildOnly = "This command can only be issued from a guild."

// FeatureDisabled is the message sent when a feature requested by the user is disabled.
const FeatureDisabled = "This feature is currently disabled."

// AdminOnly is the message sent when an unprivileged user invokes an admin-only request.
const AdminOnly = "This command is only available to system administrators!"

// OperatorOnly is the message sent when an unprivileged user invokes an operator-only request.
const OperatorOnly = "This command is only available to operators!"

// PermissionDenied is the message sent when the user invokes a request without sufficient permission.
const PermissionDenied = "You are not allowed to issue this command!"

// MissingUser is the message sent when a specified user does not exist.
const MissingUser = "Specified user does not exist."

// LackingPermission is the message sent when lacking permission for an operation.
const LackingPermission = "Lacking permission to perform specified action."

// KappaColor is the primary color of the kappa.
const KappaColor = 0x3492c4
