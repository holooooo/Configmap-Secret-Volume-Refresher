// Copyright 2021 Holooooo.
// Use of this source code is governed by the WTFPL
// license that can be found in the LICENSE file.

package core

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	informerCache "k8s.io/client-go/tools/cache"
)

func newPodInformer(selector string, resync time.Duration) informerCache.SharedIndexInformer {
	i := informerCache.NewSharedIndexInformer(
		&informerCache.ListWatch{
			ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
				opts.LabelSelector = selector
				return clientSet.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), opts)
			},
			WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
				opts.LabelSelector = selector
				return clientSet.CoreV1().Pods(metav1.NamespaceAll).Watch(context.TODO(), opts)
			},
		}, &corev1.Pod{}, 0, informerCache.Indexers{},
	)
	i.AddEventHandlerWithResyncPeriod(
		&informerCache.ResourceEventHandlerFuncs{
			AddFunc:    handlePodAdd,
			UpdateFunc: handlePodUpdate,
			DeleteFunc: handlePodDelete,
		},
		resync,
	)
	return i
}

func handlePodAdd(obj interface{}) {
	pod := obj.(*corev1.Pod)
	configmapCache.PushPod(pod)
	secretCache.PushPod(pod)
}

func handlePodUpdate(oldObj, newObj interface{}) {
	handlePodAdd(newObj)
}

func handlePodDelete(obj interface{}) {
	pod := obj.(*corev1.Pod)
	configmapCache.DeletePod(pod)
	secretCache.DeletePod(pod)
}
