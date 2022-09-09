// SPDX-License-Identifier: Apache-2.0

package planner

import (
	"testing"

	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sambaoperatorv1alpha1 "github.com/samba-in-kubernetes/samba-operator/api/v1alpha1"
	"github.com/samba-in-kubernetes/samba-operator/internal/conf"
	"github.com/samba-in-kubernetes/samba-operator/internal/smbcc"
)

func TestUpdate(t *testing.T) {
	t.Run("simpleUpdate", testSimpleUpdate)
	t.Run("secondShare", testSecondShare)
}

func testSimpleUpdate(t *testing.T) {
	state := smbcc.New()
	assert.Len(t, state.Shares, 0)
	assert.Len(t, state.Globals, 0)

	p := New(InstanceConfiguration{
		SmbShare: &sambaoperatorv1alpha1.SmbShare{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "share1",
				Namespace: "smbshares",
				UID:       "phonyuid1",
			},
			Spec: sambaoperatorv1alpha1.SmbShareSpec{
				ShareName:      "share1",
				ReadOnly:       false,
				Browseable:     true,
				SecurityConfig: "",
				CommonConfig:   "",
				Storage: sambaoperatorv1alpha1.SmbShareStorageSpec{
					Pvc: &sambaoperatorv1alpha1.SmbSharePvcSpec{
						Name: "mydata",
						Path: "share1",
					},
				},
			},
		},
		GlobalConfig: &conf.OperatorConfig{},
	}, state)

	// apply changes
	changed, err := p.Update()
	assert.NoError(t, err)
	assert.True(t, changed)

	// changes already applied
	changed, err = p.Update()
	assert.NoError(t, err)
	assert.False(t, changed)

	assert.Len(t, state.Shares, 1)
	assert.Len(t, state.Globals, 1)
	assert.Contains(t, state.Shares, smbcc.Key("share1"))
}

func testSecondShare(t *testing.T) {
	state := smbcc.New()
	assert.Len(t, state.Shares, 0)
	assert.Len(t, state.Globals, 0)

	p := New(InstanceConfiguration{
		SmbShare: &sambaoperatorv1alpha1.SmbShare{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "share1",
				Namespace: "smbshares",
				UID:       "phonyuid1",
			},
			Spec: sambaoperatorv1alpha1.SmbShareSpec{
				ShareName:      "share1",
				ReadOnly:       false,
				Browseable:     true,
				SecurityConfig: "",
				CommonConfig:   "",
				Storage: sambaoperatorv1alpha1.SmbShareStorageSpec{
					Pvc: &sambaoperatorv1alpha1.SmbSharePvcSpec{
						Name: "mydata",
						Path: "share1",
					},
				},
			},
		},
		GlobalConfig: &conf.OperatorConfig{},
	}, state)

	// apply changes
	changed, err := p.Update()
	assert.NoError(t, err)
	assert.True(t, changed)

	p2 := New(InstanceConfiguration{
		SmbShare: &sambaoperatorv1alpha1.SmbShare{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "share2",
				Namespace: "smbshares",
				UID:       "phonyuid2",
			},
			Spec: sambaoperatorv1alpha1.SmbShareSpec{
				ShareName:      "share2",
				ReadOnly:       false,
				Browseable:     true,
				SecurityConfig: "",
				CommonConfig:   "",
				Storage: sambaoperatorv1alpha1.SmbShareStorageSpec{
					Pvc: &sambaoperatorv1alpha1.SmbSharePvcSpec{
						Name: "mydata",
						Path: "share2",
					},
				},
			},
		},
		GlobalConfig: &conf.OperatorConfig{},
	}, state)
	changed, err = p2.Update()
	assert.NoError(t, err)
	assert.True(t, changed)

	assert.Len(t, state.Shares, 2)
	assert.Len(t, state.Globals, 1)
	assert.Contains(t, state.Shares, smbcc.Key("share1"))
	assert.Contains(t, state.Shares, smbcc.Key("share2"))
}
