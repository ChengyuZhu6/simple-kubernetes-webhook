package mutation

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slackhq/simple-kubernetes-webhook/pkg/api"
	pkgUtils "github.com/slackhq/simple-kubernetes-webhook/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// injectEnv is a container for the mutation injecting environment vars
type injectEnv struct {
	K8sClient kubernetes.Interface
	Logger    logrus.FieldLogger
}

// injectEnv implements the podMutator interface
var _ podMutator = (*injectEnv)(nil)

// Name returns the struct name
func (se injectEnv) Name() string {
	return "inject_env"
}

// Mutate returns a new mutated pod according to set env rules
func (se injectEnv) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	se.Logger = se.Logger.WithField("mutation", se.Name())
	mpod := pod.DeepCopy()

	err := se.UpdateEnvVar(mpod)
	if err != nil {
		return mpod, errors.Wrap(err, "failed to mutate pod, %v")
	}
	// build out env var slice
	envVars := []corev1.EnvVar{{
		Name:  "KUBE",
		Value: "true",
	}}

	// inject env vars into pod
	for _, envVar := range envVars {
		se.Logger.Debugf("pod env injected %s", envVar)
		se.injectEnvVar(mpod, envVar)
	}
	PrintEnvVar(se.Logger, mpod)
	return mpod, nil
}

func (se injectEnv) UpdateEnvVar(pod *corev1.Pod) error {
	if err := se.UpdateEnvFrom(pod); err != nil {
		return errors.Wrap(err, "failed to update envFrom")
	}
	se.Logger.Info("updateEnvFrom successful")

	if err := se.UpdateValueFrom(pod); err != nil {
		return errors.Wrap(err, "failed to update valueFrom")
	}
	se.Logger.Info("updateValueFrom successful")

	return nil
}

func (se injectEnv) UpdateEnvFrom(pod *corev1.Pod) error {
	ns := pod.GetNamespace()
	for _, containers := range [][]corev1.Container{pod.Spec.InitContainers, pod.Spec.Containers} {
		for i, container := range containers {
			if len(container.EnvFrom) > 0 {
				envFrom, err := pkgUtils.LookForEnvFrom(container.EnvFrom, ns)
				if err != nil {
					return errors.Wrap(err, "failed to look for envFrom")
				}
				containers[i].Env = append(container.Env, envFrom...)
			}
		}
	}
	return nil
}

func (se injectEnv) UpdateValueFrom(pod *corev1.Pod) error {
	ns := pod.GetNamespace()
	if err := se.UpdateContainerEnvValueFrom(pod.Spec.InitContainers, ns); err != nil {
		return err
	}
	if err := se.UpdateContainerEnvValueFrom(pod.Spec.Containers, ns); err != nil {
		return err
	}
	return nil
}

func (se injectEnv) UpdateContainerEnvValueFrom(containers []corev1.Container, ns string) error {
	for i, container := range containers {
		for j, env := range container.Env {
			if env.ValueFrom != nil {
				valueFrom, err := pkgUtils.LookForValueFrom(env, ns)
				if err != nil {
					return errors.Wrap(err, "failed to look for valueFrom")
				}
				if valueFrom == nil {
					continue
				}
				if api.HasSealedSecretsPrefix(valueFrom.Value) {
					c, err := api.CreateCDHClient()
					if err != nil {
						return fmt.Errorf("failed to create cdh client: %w", err)
					}
					defer c.Close()

					unsealedValue, err := c.UnsealEnv(context.TODO(), string(valueFrom.Value))
					valueFrom.Value = unsealedValue
				}
				containers[i].Env[j] = *valueFrom
			}
		}
	}
	return nil
}

// injectEnvVar injects a var in both containers and init containers of a pod
func (se injectEnv) injectEnvVar(pod *corev1.Pod, envVar corev1.EnvVar) {
	for i, container := range pod.Spec.Containers {
		if !HasEnvVar(container, envVar) {
			pod.Spec.Containers[i].Env = append(container.Env, envVar)
		}
	}
	for i, container := range pod.Spec.InitContainers {
		if !HasEnvVar(container, envVar) {
			pod.Spec.InitContainers[i].Env = append(container.Env, envVar)
		}
	}
}

// HasEnvVar returns true if environment variable exists false otherwise
func HasEnvVar(container corev1.Container, checkEnvVar corev1.EnvVar) bool {
	for _, envVar := range container.Env {
		if envVar.Name == checkEnvVar.Name {
			return true
		}
	}
	return false
}

func PrintEnvVar(logger logrus.FieldLogger, pod *corev1.Pod) {
	for _, container := range pod.Spec.Containers {
		for _, envVar := range container.Env {
			logger.Infof("Print Containers EnvVar: %s", envVar)
		}
	}
	for _, container := range pod.Spec.InitContainers {
		for _, envVar := range container.Env {
			logger.Infof("Print InitContainers EnvVar: %s", envVar)
		}
	}

}
