package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	sealedSecretName = "sealedsecret"
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

	se.Logger.Debugf("pod env updated %s", sealedSecretName)
	se.UpdateEnvVar(mpod, sealedSecretName)

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

	return mpod, nil
}

func (se injectEnv) UpdateEnvVar(pod *corev1.Pod, specificEnvName string) {
	for i, container := range pod.Spec.Containers {
		for j, envVar := range container.Env {
			if envVar.Name == specificEnvName {
				//TODO get secret from etcd
				pod.Spec.Containers[i].Env[j].Value = "sealedsecret"
			}
		}
	}
	for i, container := range pod.Spec.InitContainers {
		for j, envVar := range container.Env {
			if envVar.Name == specificEnvName {
				//TODO get secret from etcd
				pod.Spec.InitContainers[i].Env[j].Value = "sealedsecret"
			}
		}
	}
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
