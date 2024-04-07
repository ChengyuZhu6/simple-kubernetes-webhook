// Copyright (c) 2024 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package utils

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	kubernetesConfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var k8sClient kubernetes.Interface

func newK8SClient() (kubernetes.Interface, error) {
	if k8sClient != nil {
		return k8sClient, nil
	}
	kubeConfig, err := kubernetesConfig.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get kubernetes config")
	}

	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kubernetes client")
	}
	return client, nil
}
