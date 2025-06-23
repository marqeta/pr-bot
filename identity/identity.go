package identity

import (
	"errors"

	prbot "github.com/marqeta/pr-bot"
)

type Validator interface {
	ValidateIdentity(identity *CallerIdentity) error
}

// Default validator: allows anything with a non-empty ARN
type AllowAllValidator struct {
	Config *prbot.Config
}

func (v *AllowAllValidator) ValidateIdentity(identity *CallerIdentity) error {
	if identity == nil || !stringInSlice(identity.Arn, v.Config.Identity.AllowedCallerArns) || !stringInSlice(identity.Account, v.Config.Identity.AllowedCallerAccounts) {
		return errors.New("identity missing or invalid")
	}
	return nil
}

func stringInSlice(target string, list []string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}
