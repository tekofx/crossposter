package merrors

type MError struct {
	Code    MErrorCode `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
}
type MErrorCode int

func New(code MErrorCode, message string) *MError {
	return &MError{
		Code:    code,
		Message: message,
	}
}

const (
	// 0 --> 999 | SYSTEM UNEXPECTED ERRORS
	UnexpectedErrorCode                 MErrorCode = 0
	DatabaseErrorCode                   MErrorCode = 1
	NotImplementedErrorCode             MErrorCode = 2
	NothingChangedErrorCode             MErrorCode = 3
	CannotGenerateAuthTokenErrorCode    MErrorCode = 4
	CannotCreateValidationCodeErrorCode MErrorCode = 5
	CannotReadFileErrorCode             MErrorCode = 6
	CannotCreateRequestErrorCode        MErrorCode = 7
	DoRequestErrorCode                  MErrorCode = 8
	CannotConvertToIntErrorCode         MErrorCode = 9
	CannotDownloadFileErrorCode         MErrorCode = 10
	CannotCreateFileErrorCode           MErrorCode = 11
	CannotSaveFileErrorCode             MErrorCode = 12

	// 1000 -> 1999 | VALIDATION ERRORS
	InvalidRequestErrorCode   MErrorCode = 1000
	ParseJSONErrorCode        MErrorCode = 1001
	ReadResponseBodyErrorCode MErrorCode = 1002

	// 2000 -> 2099 | BLUESKY ERRORS
	BskyAuthRequestErrorCode MErrorCode = 2000
	BskyAuthErrorCode        MErrorCode = 2001
	BskyAuthDecodeErrorCode  MErrorCode = 2002
	BskyUploadBlobErrorCode  MErrorCode = 2003
	BskyPostRequestErrorCode MErrorCode = 2004
	BskyPostErrorCode        MErrorCode = 2005

	// 2100 -> 2199 | TWITTER ERRORS
	TwitterClientCreationErrorCode          MErrorCode = 2100
	TwitterInitializeMediaErrorCode         MErrorCode = 2101
	TwitterCannotPostTextErrorCode          MErrorCode = 2102
	TwitterCannotAppendMediaUploadErrorCode MErrorCode = 2103
	TwitterCannotFinalizeInputErrorCode     MErrorCode = 2104
	TwitterCannotCreatePostErrorCode        MErrorCode = 2105
	TwitterServiceUnavailableErrorCode      MErrorCode = 2106

	// 2200 -> 2299 | INSTAGRAM ERRORS
	InstagramUploadImageErrorCode        MErrorCode = 2200
	InstagramInvalidAccessTokenErrorCode MErrorCode = 2201

	// 2400 -> 2499 | TELEGRAM BOT ERRORS
	TelegramCannotSendMessageToOwnerErrorCode   MErrorCode = 2400
	TelegramCannotSendMediaGroupErrorCode       MErrorCode = 2401
	TelegramCannotSendMessageToChannelErrorCode MErrorCode = 2402
	TelegramArgumentNotProvidedErrorCode        MErrorCode = 2403

	// 4000 -> 4999 | DATABASE ERRORS
	NotFoundErrorCode   MErrorCode = 4000
	UpdatePostErrorCode MErrorCode = 4001
	RemovePostErrorCode MErrorCode = 4002

	// 5000 -> 5999 | AUTHORITATION ERRORS
	UnauthorizedErrorCode MErrorCode = 5000

	// 6000 -> 7999 | PERMISSION ERRORS
	AccessDeniedErrorCode          MErrorCode = 6000
	NotEnoughtPermissionsErrorCode MErrorCode = 6001
)

const (
	AccessDeniedMessage string = "access denied"

	CannotConnectToDatabaseMessage string = "cannot connect to database"
	DatabaseConnectionEmptyMessage string = "database connection cannot be empty"
	ServiceIDEmptyMessage          string = "service id cannot be empty"
	RegisteredDomainsEmptyMessage  string = "registered domains cannot be empty"
	SecretEmptyMessage             string = "secret cannot be empty"

	TokenEmptyMessage   string = "token cannot be empty"
	TokenInvalidMessage string = "invalid token"

	FileTooLargeMessage string = "%s is too long; the maximum size is %dMB"

	PasswordEmptyMessage                string = "password cannot be empty"
	PasswordShortMessage                string = "password is short"
	PasswordNoNumericMessage            string = "password must contain at least one numeric character"
	PasswordNoSpecialCharacterMessage   string = "password must contain at least one special character"
	PasswordNoLowercaseCharacterMessage string = "password must contain at least one uppercase character"
	PasswordNoUppercaseCharacterMessage string = "password must contain at least one lowercase character"

	EmailEmptyMessage string = "email cannot be empty"

	UserEmptyMessage               string = "user cannot be empty"
	UserNotFoundMessage            string = "user not found"
	UserIDNegativeMessage          string = "user id must be positive"
	UserCannotDeleteMessage        string = "cannot delete user"
	UserCannotUpdateMessage        string = "cannot update user"
	UserAlreadyExistsMessage       string = "user already exists"
	UserProfilePictureEmptyMessage string = "no profile picture provided"

	DeviceEmptyMessage          string = "device cannot be empty"
	DeviceAddressEmptyMessage   string = "device address cannot be empty"
	DeviceUserAgentEmptyMessage string = "device user agent cannot be empty"
)
