package constants

import "errors"

const (
	CodeFieldRequired            = 100
	CodeSearchNotAvailable       = 101
	CodeInvalidProductLink       = 102
	CodeMerchantNotActive        = 103
	CodeMerchantNotSupported     = 104
	CodeAlphaCharsOnly           = 105
	CodeInvalidEmailFormat       = 106
	CodeInsufficientLength       = 107
	CodeEmailAlreadyRegistered   = 108
	CodeInvalidMobile            = 109
	CodeMobileAlreadyInUse       = 110
	CodeInvalidEmailOrPassword   = 111
	CodeBasketNotFound           = 112
	CodeItemNotFound             = 113
	CodeItemNotFoundInBasket     = 114
	CodeQtyMustBeGreaterThanZero = 115
)

var ErrBasketNotFound = errors.New("basket not found")
var ErrItemNotFound = errors.New("item not found")
var ErrItemNotFoundInBasket = errors.New("item not found in basket")
