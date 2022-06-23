package kube_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/weaveworks/weave-gitops/pkg/kube"
	"github.com/weaveworks/weave-gitops/pkg/server/auth"
	"k8s.io/client-go/rest"
)

var _ kube.ConfigGetter = (*kube.ImpersonatingConfigGetter)(nil)

func TestImpersonatingConfigGetterPrincipalInContext(t *testing.T) {
	g := kube.NewImpersonatingConfigGetter(&rest.Config{})
	ctx := auth.WithPrincipal(context.TODO(), &auth.UserPrincipal{ID: "user@example.com"})

	cfg := g.Config(ctx)

	want := &rest.Config{
		Impersonate: rest.ImpersonationConfig{
			UserName: "user@example.com",
		},
	}
	if diff := cmp.Diff(want, cfg); diff != "" {
		t.Fatalf("incorrect client config:\n%s", diff)
	}
}

func TestImpersonatingConfigGetterPrincipalInContextWithToken(t *testing.T) {
	g := kube.NewImpersonatingConfigGetter(&rest.Config{})
	ctx := auth.WithPrincipal(context.TODO(), &auth.UserPrincipal{ID: "user@example.com", Token: "test-token"})

	cfg := g.Config(ctx)

	want := &rest.Config{
		BearerToken: "test-token",
	}
	if diff := cmp.Diff(want, cfg); diff != "" {
		t.Fatalf("incorrect client config:\n%s", diff)
	}
}

func TestImpersonatingConfigGetterPrincipalInContextWithGroups(t *testing.T) {
	g := kube.NewImpersonatingConfigGetter(&rest.Config{})
	ctx := auth.WithPrincipal(context.TODO(), &auth.UserPrincipal{ID: "user@example.com", Groups: []string{"test-group"}})

	cfg := g.Config(ctx)

	want := &rest.Config{
		Impersonate: rest.ImpersonationConfig{
			UserName: "user@example.com",
			Groups:   []string{"test-group"},
		},
	}
	if diff := cmp.Diff(want, cfg); diff != "" {
		t.Fatalf("incorrect client config:\n%s", diff)
	}
}
