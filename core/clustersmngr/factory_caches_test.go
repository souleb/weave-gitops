package clustersmngr_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cheshir/ttlcache"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	"github.com/weaveworks/weave-gitops/pkg/server/auth"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestUsersNamespaces(t *testing.T) {
	g := NewGomegaWithT(t)

	un := clustersmngr.UsersNamespaces{Cache: ttlcache.New(1 * time.Second)}

	user := &auth.UserPrincipal{ID: "user-id"}

	ns := v1.Namespace{}
	ns.Name = "ns1"

	clusterName := "cluster-1"

	un.Set(user, clusterName, []v1.Namespace{ns})

	t.Run("namespaces of a single cluster", func(t *testing.T) {
		nss, found := un.Get(user, clusterName)
		g.Expect(found).To(BeTrue())
		g.Expect(nss).To(Equal([]v1.Namespace{ns}))
	})

	t.Run("all namespaces from all", func(t *testing.T) {
		nsMap := un.GetAll(user, []clustersmngr.Cluster{{Name: clusterName}})
		g.Expect(nsMap).To(Equal(map[string][]v1.Namespace{clusterName: {ns}}))
	})
}

func TestClusters(t *testing.T) {
	g := NewGomegaWithT(t)

	cs := clustersmngr.Clusters{}

	c1 := "cluster-1"
	c2 := "cluster-2"
	clusters := []clustersmngr.Cluster{{Name: c1}, {Name: c2}}

	// simulating concurrent access
	go cs.Set(clusters)
	go cs.Set(clusters)

	cs.Set(clusters)

	g.Expect(cs.Get()).To(Equal([]clustersmngr.Cluster{{Name: c1}, {Name: c2}}))

	g.Expect(cs.Hash()).To(Equal(fmt.Sprintf("%s%s", c1, c2)))
}

func TestClustersNamespaces(t *testing.T) {
	testSuite := []struct {
		name       string
		namespaces []v1.Namespace
	}{
		{
			name: "single namespace",
			namespaces: []v1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "a",
					},
				},
			},
		},
		{
			name: "multiple namespaces",
			namespaces: []v1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "a",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "c",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "b",
					},
				},
			},
		},
	}

	for _, tc := range testSuite {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			cs := clustersmngr.ClustersNamespaces{}
			clusterName := "cluster-1"

			// simulating concurrent access
			go cs.Set(clusterName, []v1.Namespace{tc.namespaces[0]})
			go cs.Set(clusterName, []v1.Namespace{tc.namespaces[0]})

			cs.Set(clusterName, []v1.Namespace{tc.namespaces[0]})

			if len(tc.namespaces) == 1 {
				g.Expect(cs.Get(clusterName)).To(Equal([]v1.Namespace{tc.namespaces[0]}))
				cs.Clear()
				g.Expect(cs.Get(clusterName)).To(HaveLen(0))
				return
			}

			for _, ns := range tc.namespaces[1:] {
				cs.AddNamespace(clusterName, ns)
			}

			g.Expect(cs.Get(clusterName)).To(HaveLen(len(tc.namespaces)))

			cs.RemoveNamespace(clusterName, tc.namespaces[0])

			g.Expect(cs.Get(clusterName)).To(HaveLen(len(tc.namespaces) - 1))

			cs.Clear()

			g.Expect(cs.Get(clusterName)).To(HaveLen(0))

		})
	}
}

func TestClusterSet_Set(t *testing.T) {
	cs := clustersmngr.Clusters{}
	cluster1 := newTestCluster("cluster1", "server1")
	cluster2 := newTestCluster("cluster2", "server2")
	cluster3 := newTestCluster("cluster2", "server3")

	clusters := []clustersmngr.Cluster{cluster1, cluster2, cluster3}

	added, removed := cs.Set(clusters)
	if diff := cmp.Diff([]clustersmngr.Cluster{cluster1, cluster2, cluster3}, added); diff != "" {
		t.Fatalf("failed to calculate added:\n%s", diff)
	}

	if diff := cmp.Diff([]clustersmngr.Cluster{}, removed); diff != "" {
		t.Fatalf("failed to calculate removed:\n%s", diff)
	}

	clusters = []clustersmngr.Cluster{cluster1}

	added, removed = cs.Set(clusters)
	if diff := cmp.Diff([]clustersmngr.Cluster{}, added); diff != "" {
		t.Fatalf("failed to calculate added:\n%s", diff)
	}

	if diff := cmp.Diff([]clustersmngr.Cluster{cluster2, cluster3}, removed); diff != "" {
		t.Fatalf("failed to calculate removed:\n%s", diff)
	}
}

func newTestCluster(name, server string) clustersmngr.Cluster {
	return clustersmngr.Cluster{
		Name:   name,
		Server: server,
	}
}
