package mutation

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/wI2L/jsondiff"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// Mutator is a container for mutation
type Mutator struct {
	K8sClient kubernetes.Interface
	Logger    *logrus.Entry
}

// NewMutator returns an initialised instance of Mutator
func NewMutator(k8sClient kubernetes.Interface, logger *logrus.Entry) *Mutator {
	return &Mutator{K8sClient: k8sClient, Logger: logger}
}

// podMutators is an interface used to group functions mutating pods
type podMutator interface {
	Mutate(*corev1.Pod) (*corev1.Pod, error)
	Name() string
}

// MutatePodPatch returns a json patch containing all the mutations needed for
// a given pod
func (m *Mutator) MutatePodPatch(pod *corev1.Pod) ([]byte, error) {
	var podName string
	if pod.Name != "" {
		podName = pod.Name
	} else {
		if pod.ObjectMeta.GenerateName != "" {
			podName = pod.ObjectMeta.GenerateName
		}
	}
	log := logrus.WithField("pod_name", podName)

	// list of all mutations to be applied to the pod
	mutations := []podMutator{
		minLifespanTolerations{K8sClient: m.K8sClient, Logger: log},
		injectEnv{K8sClient: m.K8sClient, Logger: log},
	}

	mpod := pod.DeepCopy()

	// apply all mutations
	for _, m := range mutations {
		var err error
		mpod, err = m.Mutate(mpod)
		if err != nil {
			return nil, err
		}
	}

	// generate json patch
	patch, err := jsondiff.Compare(pod, mpod)
	if err != nil {
		return nil, err
	}

	patchb, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	return patchb, nil
}
