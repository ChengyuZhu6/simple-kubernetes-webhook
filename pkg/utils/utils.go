// Copyright (c) 2024 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package utils

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubernetesConfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var k8sClient kubernetes.Interface

func NewK8SClient() (kubernetes.Interface, error) {
	if k8sClient != nil {
		return k8sClient, nil
	}
	kubeConfig, err := kubernetesConfig.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get kubernetes config")
	}

	k8sClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kubernetes client")
	}
	return k8sClient, nil
}

func LookForEnvFrom(envFrom []corev1.EnvFromSource, ns string) ([]corev1.EnvVar, error) {
	var envVars []corev1.EnvVar

	for _, ef := range envFrom {
		if ef.SecretRef != nil {
			data, err := GetDataFromSecret(context.TODO(), ef.SecretRef.Name, ns)
			if err != nil {
				if apierrors.IsNotFound(err) && ef.SecretRef.Optional != nil && *ef.SecretRef.Optional {
					continue
				}
				return envVars, errors.Wrapf(err, "failed to get secret %s/%s", ns, ef.SecretRef.Name)
			}
			for key, value := range data {
				envFromSec := corev1.EnvVar{
					Name:  key,
					Value: string(value),
				}
				envVars = append(envVars, envFromSec)
			}
		}
	}
	return envVars, nil
}

func LookForValueFrom(env corev1.EnvVar, ns string) (*corev1.EnvVar, error) {
	if env.ValueFrom.SecretKeyRef != nil {
		data, err := GetDataFromSecret(context.TODO(), env.ValueFrom.SecretKeyRef.Name, ns)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get secret %s/%s", ns, env.ValueFrom.SecretKeyRef.Name)
		}
		fromSecret := corev1.EnvVar{
			Name:  env.Name,
			Value: string(data[env.ValueFrom.SecretKeyRef.Key]),
		}
		return &fromSecret, nil
	}
	return nil, errors.New("no value")
}

func GetDataFromSecret(ctx context.Context, secretName, ns string) (map[string][]byte, error) {
	k8sClient, err := NewK8SClient()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get k8sclient")
	}
	secret, err := k8sClient.CoreV1().Secrets(ns).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get secret %s/%s", ns, secretName)
	}
	return secret.Data, nil
}
