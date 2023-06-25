//go:build mage
// +build mage

package main

import (
	"github.com/knabben/tutorial-istio-sec/3-istio-gw/steps"
	"github.com/magefile/mage/mg"
	"os"
)

const (
	CLUSTER_NAME = "ambient"
	namespace    = "default"

	ISTIO_CONFIG = "istio.yaml"
)

var (
	serviceMesh steps.ServiceMeshI
)

func init() {
	serviceMesh = steps.NewServiceMesh()
}

type SM mg.Namespace

// Install installs resources into the cluster
func (SM) Install() error {
	if os.Getenv("INSTALL_KIND") != "" {
		// Install kind with MetalLB enabled.
		if err := serviceMesh.InstallKind(CLUSTER_NAME, "kind.yaml", true); err != nil {
			return err
		}
	}
	// Install Istio
	if err := serviceMesh.InstallIstio(ISTIO_CONFIG, namespace); err != nil {
		return err
	}

	return nil
}

// Deploy creates the pre-defined topology for tests
func (SM) Deploy() error {
	if err := serviceMesh.DeployApplication(namespace); err != nil {
		return err
	}

	return nil
}

// Policies create a VirtualService and define application Authorization files
func (SM) Policies() error {
	if err := serviceMesh.ApplyPolicies(namespace); err != nil {
		return err
	}

	return nil
}

// Delete cleans up resources from cluster
func (SM) Delete() error {
	if os.Getenv("INSTALL_KIND") != "" {
		if err := serviceMesh.DeleteKind(CLUSTER_NAME); err != nil {
			return err
		}
	} else {
		if err := serviceMesh.DeleteIstio(); err != nil {
			return err
		}
	}

	return nil
}
