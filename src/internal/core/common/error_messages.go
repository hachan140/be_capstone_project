package common

const (
	ErrMessageInvalidRequest = "Yêu cầu không hợp lệ"

	ErrMessageInvalidName               = "Tên không hợp lệ"
	ErrMessageInvalidEmail              = "Email không hợp lệ"
	ErrMessageInvalidPassword           = "Mật khẩu không hợp lệ"
	ErrMessageEmailExisted              = "Email has already existed"
	ErrMessageUsernameHasAlreadyExisted = "Username has already existed"
	ErrMessageInvalidUsername           = "Invalid Username"
	ErrMessageInvalidUser               = "User không tồn tại"
	ErrMessageUserSocialDoesnotExist    = "Tài khoản user social không tồn tại"
	ErrMessageRefreshTokenNotFound      = "Không tìm thấy refresh token của user"

	//organization
	ErrMessageOrganizationExisted            = "Tổ chức đã tồn tại"
	ErrMessageInvalidOrganizationName        = "Tên tổ chức không hợp lệ"
	ErrMessageOrganizationNotExist           = "Tổ chức không tồn tại"
	ErrMessageCannotAccessToOrganization     = "Không có quyền truy cập vào tổ chức"
	ErrMessageUserAlreadyInOtherOrganization = "User đã là thành viên của tổ chức khác"
)
