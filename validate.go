package vanilla

import (
	"errors"

	"github.com/dlclark/regexp2"
	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
)

func (r register) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, validation.By(getMatch("^(?=.{6,16}$)(?![_.])(?!.*[_.]{2})[a-zA-Z0-9._]+(?<![_.])$"))),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.By(getMatch("^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$"))),
	)
}

func (l login) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Username, validation.Required, validation.By(getMatch("^(?=.{6,16}$)(?![_.])(?!.*[_.]{2})[a-zA-Z0-9._]+(?<![_.])$"))),
		validation.Field(&l.Password, validation.Required, validation.By(getMatch("^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$"))),
	)
}

func getMatch(pattern string) validation.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
		re := regexp2.MustCompile(pattern, regexp2.RE2)
		isMatch, err := re.MatchString(s)
		if err != nil {
			return err
		}

		if !isMatch {
			return errors.New("not match")
		}
		return nil
	}
}
