package handlers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	corev2 "github.com/iden3/go-iden3-core/v2/w3c"
)

var (
	ErrDID        = validation.NewError("validation_is_did", "must be a valid DID")
	ValidationDID = validation.NewStringRuleWithError(isDID, ErrDID)
)

func isDID(value string) bool {
	_, err := corev2.ParseDID(value)
	return err == nil
}
