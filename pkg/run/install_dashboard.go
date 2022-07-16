package run

import (
	"context"
	"fmt"

	"github.com/weaveworks/weave-gitops/pkg/kube"
	"github.com/weaveworks/weave-gitops/pkg/logger"
	"github.com/weaveworks/weave-gitops/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func InstallDashboard(log logger.Logger, ctx context.Context, kubeClient *kube.KubeHTTP) error {
	password, err := utils.ReadPasswordFromStdin("Please enter your password to generate your secret: ")
	if err != nil {
		return fmt.Errorf("could not read password: %w", err)
	}

	secret, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	log.Successf("Secret has been generated:")
	fmt.Println(string(secret))

	log.Actionf("Installing GitOps Dashboard ...")

	return nil
}
