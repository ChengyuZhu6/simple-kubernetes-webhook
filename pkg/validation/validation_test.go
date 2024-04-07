package validation

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	pkgUtils "github.com/slackhq/simple-kubernetes-webhook/pkg/utils"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidatePod(t *testing.T) {
	k8sClient, err := pkgUtils.NewK8SClient()
	if err != nil {
		logrus.WithError(err).Fatalf("error creating k8s client")
	}
	v := NewValidator(k8sClient, logger())

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "lifespan",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lifespan",
				Image: "busybox",
			}},
		},
	}

	val, err := v.ValidatePod(pod)
	assert.Nil(t, err)
	assert.True(t, val.Valid)
}

func logger() *logrus.Entry {
	mute := logrus.StandardLogger()
	mute.Out = ioutil.Discard
	return mute.WithField("logger", "test")
}
