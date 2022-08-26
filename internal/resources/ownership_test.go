// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	yep = true
)

func sampleConfigMap() *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         "samba-operator.samba.org/v1alpha1",
					Kind:               "SmbShare",
					Name:               "foobar",
					UID:                "pretendapplezebra",
					Controller:         &yep,
					BlockOwnerDeletion: &yep,
				},
				{
					APIVersion:         "samba-operator.samba.org/v1alpha1",
					Kind:               "SmbShare",
					Name:               "bazbaz",
					UID:                "pretendalanzork",
					Controller:         nil,
					BlockOwnerDeletion: nil,
				},
			},
		},
		Data: map[string]string{
			"hello": "world",
		},
	}
	return cm
}

func TestSmbShareOwnerRefs(t *testing.T) {
	t.Run("findTwoRefsExact", func(t *testing.T) {
		cm := sampleConfigMap()
		refs, err := smbShareOwnerRefs(cm)
		assert.NoError(t, err)
		assert.Len(t, refs, 2)
	})

	t.Run("findTwoRefsExtra", func(t *testing.T) {
		cm := sampleConfigMap()
		cm.OwnerReferences = append(
			cm.OwnerReferences,
			metav1.OwnerReference{
				APIVersion:         "something.example.org/v1",
				Kind:               "CoolResource",
				Name:               "yikes",
				UID:                "pretendablezipper",
				Controller:         nil,
				BlockOwnerDeletion: nil,
			},
		)

		refs, err := smbShareOwnerRefs(cm)
		assert.NoError(t, err)
		assert.Len(t, refs, 2)
	})

	t.Run("findThreeRefsAdded", func(t *testing.T) {
		cm := sampleConfigMap()
		cm.OwnerReferences = append(
			cm.OwnerReferences,
			metav1.OwnerReference{
				// pretend we're living in the future.
				APIVersion:         "samba-operator.samba.org/v1beta2",
				Kind:               "SmbShare",
				Name:               "abc",
				UID:                "pretendatomiczoo",
				Controller:         nil,
				BlockOwnerDeletion: nil,
			},
		)

		refs, err := smbShareOwnerRefs(cm)
		assert.NoError(t, err)
		assert.Len(t, refs, 3)
	})

	t.Run("badGV", func(t *testing.T) {
		cm := sampleConfigMap()
		cm.OwnerReferences = append(
			cm.OwnerReferences,
			metav1.OwnerReference{
				APIVersion:         "good/bye/my/friend",
				Kind:               "SmbShare",
				Name:               "abc",
				UID:                "pretendatlanticzen",
				Controller:         nil,
				BlockOwnerDeletion: nil,
			},
		)

		refs, err := smbShareOwnerRefs(cm)
		assert.Error(t, err)
		assert.Nil(t, refs)
	})
}

func TestOwnerRefsToNames(t *testing.T) {
	cm := sampleConfigMap()
	names := ownerRefsToNames(cm.GetOwnerReferences(), cm.GetNamespace())
	if assert.Len(t, names, 2) {
		assert.Equal(t, names[0].Name, "foobar")
		assert.Equal(t, names[1].Name, "bazbaz")
	}
}

func TestExcludeOwnerRefs(t *testing.T) {
	cm := sampleConfigMap()
	refs := excludeOwnerRefs(cm.GetOwnerReferences(), "nomatch", "zzz")
	assert.Len(t, refs, 2)

	refs = excludeOwnerRefs(cm.GetOwnerReferences(), "foobar", "pretendapplezebra")
	assert.Len(t, refs, 1)
}
