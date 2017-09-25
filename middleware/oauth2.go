package middleware

import (
	"github.com/gin-gonic/gin"
	ginoauth2 "github.com/zalando/gin-oauth2"
)

// ScopesAuth is an AccessCheckFunction that gives access if the token includes
// all of the specified scopes.
func ScopesAuth(scopes ...string) ginoauth2.AccessCheckFunction {
	// convert scopes slice to set.
	authScopes := make(map[string]struct{})
	for _, scope := range scopes {
		authScopes[scope] = struct{}{}
	}

	return func(tc *ginoauth2.TokenContainer, ctx *gin.Context) bool {
		for scope := range authScopes {
			value, ok := tc.Scopes[scope]
			if !ok {
				return false
			}
			ctx.Set(scope, value)
		}
		// set uid and realm
		// ctx.Set("uid")
		return true
	}
}

// User defines a user with UID and Realm.
type User struct {
	UID   string
	Realm string
}

// GetUser gets user (uid and realm) from a gin context.
func GetUser(ctx *gin.Context) User {
	user := User{}
	uid, ok := ctx.Get("uid")
	if !ok {
		return user
	}
	user.UID = uid.(string)

	realm, ok := ctx.Get("realm")
	if !ok {
		return user
	}
	user.Realm = realm.(string)

	return user
}
