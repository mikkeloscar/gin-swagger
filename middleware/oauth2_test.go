package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	ginoauth2 "github.com/zalando/gin-oauth2"
)

func TestScopesAuth(t *testing.T) {
	for _, tc := range []struct {
		msg       string
		scopes    []string
		container *ginoauth2.TokenContainer
		accepted  bool
	}{
		{
			msg:    "all scopes in the token gets accepted",
			scopes: []string{"a", "b"},
			container: &ginoauth2.TokenContainer{
				Scopes: map[string]interface{}{
					"a": nil,
					"b": nil,
				},
			},
			accepted: true,
		},
		{
			msg:    "too many scopes in token gets accepted",
			scopes: []string{"a"},
			container: &ginoauth2.TokenContainer{
				Scopes: map[string]interface{}{
					"a": nil,
					"b": nil,
				},
			},
			accepted: true,
		},
		{
			msg:    "missing scope in token does not get accepted",
			scopes: []string{"a", "b", "c"},
			container: &ginoauth2.TokenContainer{
				Scopes: map[string]interface{}{
					"a": nil,
					"b": nil,
				},
			},
			accepted: false,
		},
	} {
		t.Run(tc.msg, func(t *testing.T) {
			fn := ScopesAuth(tc.scopes...)
			if tc.accepted != fn(tc.container, &gin.Context{Keys: make(map[string]interface{})}) {
				t.Errorf("expected accepted: %t, got %t", tc.accepted, fn(tc.container, nil))
			}
		})
	}
}
