package common

type PermissionModule string
type PermissionCode int
type OtpFactor string

const (
	SmsOtpFactor   OtpFactor = "SMS_OTP_FACTOR"
	EmailOtpFactor OtpFactor = "EMAIL_OTP_FACTOR"
)

const (
	Authorization   = "Authorization"
	TokenTypeBearer = "Bearer"
	TokenTypeBasic  = "Basic"

	KeyOpsModule     = "ops_module"
	KeyOpsModuleCode = "ops_module_code"

	AttributeAudience        = "audience"
	OriginalSource           = "original"
	SourceOneU               = "one_u"
	AuthenSDKSource          = "authen_sdk"
	ErrorRaw                 = "error_raw"
	Channel                  = "X-Channel"
	NumberSecondPerDay       = 86400
	VietnamPhoneNumberPrefix = "+84"
	ContentType              = "Content-Type"
	ContentJSON              = "application/json"
	SmsConnection            = "sms"
	SMSType                  = "sms"
	EmailType                = "email"
	EmailConnection          = "email"
	Email                    = "email"
	VerifyOTPType            = "Verify_Type"
	VerifyOTPResetPassword   = "reset-password"
	VerifyOTPResendPassword  = "resend-password"
	TypePhoneNumber          = "phone_number"
	TypeEmail                = "email"
	Sub                      = "sub"
	RandomSaltLength         = 36
	PasswordTTLByDay         = 180
	UserID                   = "user_id"
	Tenant                   = "tenant"
	RegexOTP                 = `^[0-9]{6}.*`
	TestPhoneExpire          = 120
	OTPTimeToLive
	IsSetPassword        = "is_set_password"
	Source               = "source"
	ExpiredAt            = "expired_at"
	PasswordLength       = 6
	SecurePasswordLength = 6
	BrandName            = "brand_name"
	Otp                  = "otp"
	PIN                  = "PIN"
	Pin                  = "pin"
	PWD                  = "pwd"
	Sender               = "sender"
	SenderName           = "sender_name"
	AuthorizationCode    = "authorization_code"
	Env                  = "env"
	Link                 = "link1"
	EmailNew             = "email_new"
	EventTypeSignIn      = "signin"
	EventTypeSignUp      = "signup"
	TypePassword         = "password"
	TypeOtpPassword      = "otp"
	TypePinPassword      = "pin"
	TypeSecuredPassword  = "secured_password"
	TypeDynamicPassword  = "dynamic_password"
	TokenID              = "token_id"
	OPS                  = "OPS"

	TypeLogin                   = "login"
	TypeResetPassword           = "reset_password"
	TypeUnlockPasswordPhone     = "phone_number"
	TypeUnlockPasswordEmail     = "email"
	UnlockOTPByEmail            = "email"
	UnlockOTPByPhoneNumber      = "phone_number"
	UserPasswordConnection      = "Username-Password-Authentication"
	UserPasswordConnectionVinID = "vinid-auth"
	LinkAccountByEmail          = "by_email"
)
const (
	JWTOTPPhone         = "sms"
	JWTOTPEmail         = "email"
	JWTSecurityQuestion = "security_question"
	JWTFaceMatching     = "face_matching"
	JWTTokenLevel2      = "token_level2"
	TokenSource         = "token_source"
	TokenType           = "X-token-type"
)

const (
	XUserID            = "X-USER-ID"
	VUserID            = "V-USER-ID"
	RequestIDKey       = "X-REQUEST-ID"
	DeviceIDKey        = "X-DEVICE-ID"
	RFDeviceIDKey      = "X-RF-Device-ID"
	RequestDeviceIDKey = "X-Device-Request-ID"
	DeviceSessionIDKey = "X-Device-Session-ID"
	DeviceOSToken      = "X-Device-OS-Token"
	VersionHeader      = "X-Version"
	DeviceOSKey        = "X-Device-OS"
	DeviceVersionKey   = "X-Device-Version"
	DeviceUUID         = "X-Device-UUID"

	XApiKey         = "X-API-KEY"
	XApiSecret      = "X-API-SECRET"
	VinPayApiKey    = "X-Api-Key"
	VinPayApiSecret = "X-Api-Secret"
	AcceptLanguage  = "Accept-Language"
	Vietnamese      = "vi"
	English         = "en"
)

const (
	ApplicationConfig = "oneid:application:config:%s:%s"
)

const (
	SearchFieldUsername    = "username"
	DeviceTokenGenFailure  = "DEVICE_TOKEN_GENERATION_FAILED"
	EmailErrorTMPL         = "email-error.tmpl"
	SearchFieldPhoneNumber = "user_metadata.phone_number"
	SearchFieldEmail       = "user_metadata.email"
	FullNameKH             = "Khách Hàng"
)

const (
	ApplicationModule     PermissionModule = "GeneralConfiguration_AuthenticationVinid"
	ApplicationUserModule PermissionModule = "GeneralConfiguration_AuthenAppUser"
	AppProviderModule     PermissionModule = "GeneralConfiguration_AppProvider"
)

const (
	ListApplicationCode   PermissionCode = 1000
	AddNewApplicationCode PermissionCode = 1001
	EditApplicationCode   PermissionCode = 1002
	RemoveApplicationCode PermissionCode = 1003

	ListUserCode PermissionCode = 1000
	EditUserCode PermissionCode = 1002

	ViewAppProviderCode PermissionCode = 1000
	EditAppProviderCode PermissionCode = 1002
)

const (
	StatusActive          = "active"
	StatusInactive        = "inactive"
	StatusWaitingToDelete = "waiting_to_delete"
	StatusCancelDelete    = "cancel_delete"
	StatusDeleted         = "deleted"
)

var UserStatusMap = map[string]bool{
	StatusActive:          true,
	StatusInactive:        true,
	StatusWaitingToDelete: true,
	StatusCancelDelete:    true,
	StatusDeleted:         true,
}

const (
	UserMobileType = "MOBILE"
)
