package validators

import (
	"fmt"
	"regexp"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"github.com/xManan/fusion-cart/internal/app/models"
	"github.com/xManan/fusion-cart/internal/app/types"
)

type ValidationErr struct {
	Message string
	Code    int
}

func ValidateRegistration(u *models.UnverifiedUser) *ValidationErr {
	if u.FirstName == "" {
		return &ValidationErr{"firstName is required", constants.CodeFieldRequired}
	}
	re := regexp.MustCompile(`^[a-zA-Z]+$`)
	if !re.MatchString(u.FirstName) {
		return &ValidationErr{"firstName must contain alphabetic characters only", constants.CodeAlphaCharsOnly}
	}
	if u.LastName == "" {
		return &ValidationErr{"lastName is required", constants.CodeFieldRequired}
	}
	if !re.MatchString(u.LastName) {
		return &ValidationErr{"lastName must contain alphabetic characters only", constants.CodeAlphaCharsOnly}
	}
	if u.Email == "" {
		return &ValidationErr{"email is required", constants.CodeFieldRequired}
	}
	re = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(u.Email) {
		return &ValidationErr{"email format is invalid", constants.CodeInvalidEmailFormat}
	}
	if u.Mobile != "" {
		re = regexp.MustCompile(`^[0-9]+$`)
		if !re.MatchString(u.Mobile) || len(u.Mobile) != 10 {
			return &ValidationErr{"mobile must be a 10-digit number", constants.CodeInvalidMobile}
		}
	}
	if u.Password == "" {
		return &ValidationErr{"password is required", constants.CodeFieldRequired}
	}
	requiredPasswordLength := 8
	if len(u.Password) < requiredPasswordLength {
		return &ValidationErr{fmt.Sprintf("password must be atleast %d characters long", requiredPasswordLength), constants.CodeInsufficientLength}
	}
	return nil
}

func ValidateLogin(body *types.LoginRequestBody) *ValidationErr {
	if body.Email == "" {
		return &ValidationErr{"email is required", constants.CodeFieldRequired}
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(body.Email) {
		return &ValidationErr{"email format is invalid", constants.CodeInvalidEmailFormat}
	}
	if body.Password == "" {
		return &ValidationErr{"password is required", constants.CodeFieldRequired}
	}
	return nil
}
