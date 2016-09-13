package common

import (
	"testing"
	"context"
)

func TestFromContext(t *testing.T) {
	ctx := NewContext(context.Background(), &AppClaims{Username: "foobar"})
	appClaims, ok := FromContext(ctx)
	if ok != true {
		t.Fatal("Could not get claims from context.")
	}
	if appClaims.Username != "foobar" {
		t.Fatal("Context did not contain the claims object.")
	}
}