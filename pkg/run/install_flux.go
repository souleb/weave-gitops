package run

import (
	"context"
	"time"

	"github.com/fluxcd/flux2/pkg/manifestgen/install"
	"github.com/weaveworks/weave-gitops/pkg/kube"
	"github.com/weaveworks/weave-gitops/pkg/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func InstallFlux(log logger.Logger, ctx context.Context, kubeClient *kube.KubeHTTP, installOptions install.Options, kubeConfigArgs genericclioptions.RESTClientGetter) error {
	log.Actionf("Installing Flux ...")

	manifests, err := install.Generate(installOptions, "")
	if err != nil {
		log.Failuref("Couldn't generate manifests")
		return err
	}

	content := []byte(manifests.Content)

	applyOutput, err := apply(log, ctx, kubeClient, kubeConfigArgs, content)
	if err != nil {
		log.Failuref("Flux install failed")
		return err
	}

	log.Println(applyOutput)

	return nil
}

func WaitForDeploymentToBeReady(log logger.Logger, kubeClient *kube.KubeHTTP, deploymentName string, namespace string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: namespace,
		},
	}

	if err := wait.ExponentialBackoff(wait.Backoff{
		Duration: 1 * time.Second,
		Factor:   2,
		Jitter:   1,
		Steps:    10,
	}, func() (done bool, err error) {
		d := deployment.DeepCopy()
		if err := kubeClient.Get(ctx, client.ObjectKeyFromObject(d), d); err != nil {
			return false, err
		}
		// Confirm the state we are observing is for the current generation
		if d.Generation != d.Status.ObservedGeneration {
			return false, nil
		}

		if d.Status.ReadyReplicas == d.Status.Replicas {
			return true, nil
		}

		return false, nil
	}); err != nil {
		return err
	}

	return nil
}
