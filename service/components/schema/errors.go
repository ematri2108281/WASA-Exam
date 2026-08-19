package schema

import "errors"

var (
	ErrUserDoesNotExist            = errors.New("user does not exist")
	ErrUsernameTaken               = errors.New("username already exists")
	ErrConversationDoesNotExist    = errors.New("conversation does not exist")
	ErrMessageDoesNotExist         = errors.New("message does not exist")
	ErrReactionDoesNotExist        = errors.New("reaction does not exist")
	ErrGroupDoesNotExist           = errors.New("group does not exist")
	ErrGroupNotFound               = errors.New("group not found")
	ErrUnauthorized                = errors.New("authentication required")
	ErrUnauthorizedToDeleteMessage = errors.New("unauthorized to delete message")
	ErrNotAParticipant             = errors.New("user is not a participant of this conversation")
	ErrNameTooShort                = errors.New("name is too short (min 3 characters)")
	ErrNameTooLong                 = errors.New("name is too long (max 20 characters)")
	ErrInvalidGroupName            = errors.New("invalid group name")
	ErrInvalidReactionType         = errors.New("invalid reaction type")
	ErrReactionAlreadyPresent      = errors.New("user has already reacted with this type")
	ErrUserAlreadyInGroup          = errors.New("user is already in this group")
	ErrUserNotInGroup              = errors.New("user is not in this group")
	ErrInternal                    = errors.New("internal server error")
	ErrDatabase                    = errors.New("database operation failed")
)
