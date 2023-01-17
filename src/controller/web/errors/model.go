package errors

type ErrorResponse struct {
	Error string `json:"error"`
}

const (
	ErrExpiredRefresh = "refresh is expired"
	ErrExpiredAccess  = "access is expired"
)

const (
	ErrIncorrectAccess         = "incorrect access"
	ErrIncorrectRefresh        = "incorrect refresh"
	ErrIncorrectPasswordLength = "incorrect password length"
	ErrIncorrectLogin          = "incorrect login"
	ErrIncorrectLoginLength    = "incorrect login length"
	ErrIncorrectId             = "incorrect id"
	ErrIncorrectContent        = "incorrect content"
	ErrIncorrectBody           = "incorrect body"
	ErrIncorrectLimit          = "incorrect limit"
	ErrIncorrectOffset         = "incorrect offset"
	ErrIncorrectFeedType       = "incorrect feed type"
)

const (
	ErrClosedWall     = "wall is closed"
	ErrClosedComments = "comments is closed"
)

const (
	ErrNotFoundUser          = "user not found"
	ErrNotFoundInvite        = "invite not found"
	ErrNotFoundRefresh       = "refresh not found"
	ErrNotFoundAccess        = "access not found"
	ErrNotFoundMessage       = "messages not found"
	ErrNotFoundSubscribers   = "subscribers not found"
	ErrNotFoundSubscriptions = "subscriptions not found"
	ErrNotFoundReply         = "reply not found"
	ErrNotFoundWall          = "wall not found"
	ErrNotFoundImage         = "image not found"
	ErrNotFoundFeed          = "feed not found"
)

const (
	ErrInternalError            = "internal error"
	ErrBadCaptcha               = "bad captcha solution"
	ErrRegisterInviteOnly       = "invite only"
	ErrInviteIsAlreadyActivated = "invite is already activated"
	ErrAccountIsAlreadyExists   = "account is already exists"
	ErrLimitValueErr            = "minimum limit 1, maximum 30"
	ErrNotEnoughPermissions     = "not enough permissions"
	ErrLoginFailed              = "failed to log in"
	ErrCantSubscribeYourself    = "you can't subscribe to yourself"
	ErrImageSizeIsTooBig        = "image size is too big"
)
