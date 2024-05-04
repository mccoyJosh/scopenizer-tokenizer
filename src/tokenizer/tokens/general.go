// Package tokens has the implementation for the scope tokens
// and other generalized tokens (keyword and symbol)
//
// This package, as part of the symbol token,
package tokens

import (
	"errors"
	"reflect"
	"tp/src/util"
)

const (
	GENERAL_TOKEN_STRING = "TOKEN"
	SCOPE_TOKEN_STIRNG   = "SCOPE_TOKEN"
	UNKNOWN_SCOPE_STRING = "__UNKNOWN__"
)

// ValidTokenType
// This ensures the items provided to the token list are in fact tokens.
// It is expecting an interface object, hopefully being a token, i.e. Token or ScopeToken
// This will return a bool, true meaning a valid item was provided. False otherwise.
func ValidTokenType(i interface{}) bool {
	return ValidNormalToken(i) || ValidScopeToken(i)
}

// ValidNormalToken
// Returns true if the provided token is a normal (non-scope) token; false otherwise.
func ValidNormalToken(i interface{}) bool {
	return GetTokenType(i) == GENERAL_TOKEN_STRING
}

// ValidScopeToken
// Returns true if the provided token is a scope token; false otherwise.
func ValidScopeToken(i interface{}) bool {
	return GetTokenType(i) == SCOPE_TOKEN_STIRNG
}

// GetTokenType
// This gets the type a token.
// Returns a string representation of the token's type, either being TOKEN, SCOPE_TOKEN, or other, if it
// is not a valid token type
func GetTokenType(i interface{}) string {
	switch i.(type) {
	case Token:
		return GENERAL_TOKEN_STRING
	case ScopeToken:
		return SCOPE_TOKEN_STIRNG
	default:
		util.Error("Unexpected type, not a token", errors.New("Unexpected type provided: "+reflect.TypeOf(i).String()))
		return "other"
	}
}
