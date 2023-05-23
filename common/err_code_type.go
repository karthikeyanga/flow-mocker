package common

type ErrorCodeType string

const (
	ERR_CODE_NONE                           ErrorCodeType = ""
	ERR_CODE_BAD_REQUEST                    ErrorCodeType = "BAD_REQUEST"
	ERR_CODE_DEPRECATED_API_VERSION         ErrorCodeType = "DEPRECATED_API_VERSION"
	ERR_CODE_FORBIDDEN                      ErrorCodeType = "FORBIDDEN"
	ERR_CODE_EXCEPTION                      ErrorCodeType = "500"
	ERR_CODE_OTP_MISMATCH                   ErrorCodeType = "OTP_MISMATCH"
	ERR_CODE_DISPLAYABLE                    ErrorCodeType = "DISPLAY_ERROR"
	ERR_CODE_CLOSE_WINDOW                   ErrorCodeType = "CLOSE_WEBVIEW"
	ERR_CODE_INVALID_STAGE                  ErrorCodeType = "INVALID_STAGE"
	ERR_CODE_INTERNAL_ERROR                 ErrorCodeType = "INTERNAL_ERROR"
	ERR_CODE_VALIDATION                     ErrorCodeType = "VALIDATION_ERROR"
	ERR_CODE_UNKNOWN                        ErrorCodeType = "ERROR_UNKNOWN"
	ERR_DUPLICATE_TRANSACTION               ErrorCodeType = "DUPLICATE_TRANSACTION"
	ERR_CODE_DISPLAYABLE_PAGE               ErrorCodeType = "DISPLAY_ERROR_PAGE"
	ERR_CODE_LOGIN_REQUIRED                 ErrorCodeType = "LOGIN_REQUIRED"
	ERR_CODE_NO_RULE_MATCHED                ErrorCodeType = "NO_RULE_MATCHED"
	ERR_CODE_INVALID_TRANSACTION_REQUEST    ErrorCodeType = "INVALID_TRANSACTION_REQUEST"
	ERR_CODE_INELIGIBLE_TRANSACTION_REQUEST ErrorCodeType = "TXN_NOT_ELIGIBLE"
	ERR_CODE_INVALID_TRANSACTION_STATUS     ErrorCodeType = "INVALID_TRANSACTION_STATUS"
	ERR_CODE_INVALID_REFUND_REQUEST                       = "INVALID_REFUND_REQUEST"
)
