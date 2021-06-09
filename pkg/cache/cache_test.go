// Copyright 2021 Holooooo.
// Use of this source code is governed by the WTFPL
// license that can be found in the LICENSE file.

package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewCache(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-1",
			Namespace: "test",
		},
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-1",
			Namespace: "test",
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "test-1",
							},
						},
					},
				},
			},
		},
	}
	cache := NewConfigmapCache([]corev1.ConfigMap{*cm}, []corev1.Pod{*pod})
	cache.PushRes(cm)
	cache.PushPod(pod)
	podlist := cache.GetPods(cm)
	assert.Equal(t, 1, len(podlist))

}

func TestCache_DeleteRes(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-1",
			Namespace: "test",
		},
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-1",
			Namespace: "test",
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "test-1",
							},
						},
					},
				},
			},
		},
	}
	cache := NewConfigmapCache([]corev1.ConfigMap{*cm}, []corev1.Pod{*pod})
	cache.PushRes(cm)
	cache.PushPod(pod)
	podlist := cache.GetPods(cm)
	assert.Equal(t, 1, len(podlist))

	cache.DeleteRes(cm)
	podlist = cache.GetPods(cm)
	assert.Equal(t, 0, len(podlist))
}

func TestCache_DeletePod(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-1",
			Namespace: "test",
		},
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-1",
			Namespace: "test",
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "test-1",
							},
						},
					},
				},
			},
		},
	}

	cache := NewConfigmapCache([]corev1.ConfigMap{*cm}, []corev1.Pod{*pod})
	podlist := cache.GetPods(cm)
	assert.Equal(t, 1, len(podlist))

	cache.DeletePod(pod)
	podlist = cache.GetPods(cm)
	assert.Equal(t, 0, len(podlist))
}
