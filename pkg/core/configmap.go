// Copyright 2021 Holooooo.
// Use of this source code is governed by the WTFPL
// license that can be found in the LICENSE file.

package core

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	informerCache "k8s.io/client-go/tools/cache"
)

const ConfigMapAnnotation = "holooooo.io/lastest-configmap-version"

func newConfigMapInformer(selector string, resync time.Duration) informerCache.SharedIndexInformer {
	i := informerCache.NewSharedIndexInformer(
		&informerCache.ListWatch{
			ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
				opts.LabelSelector = selector
				return clientSet.CoreV1().ConfigMaps(metav1.NamespaceAll).List(context.TODO(), opts)
			},
			WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
				opts.LabelSelector = selector
				return clientSet.CoreV1().ConfigMaps(metav1.NamespaceAll).Watch(context.TODO(), opts)
			},
		}, &corev1.ConfigMap{}, 0, informerCache.Indexers{},
	)
	i.AddEventHandlerWithResyncPeriod(
		&informerCache.ResourceEventHandlerFuncs{
			AddFunc:    handleConfigMapAdd,
			UpdateFunc: handleConfigMapUpdate,
			DeleteFunc: handleConfigMapDelete,
		},
		resync,
	)
	return i
}

func handleConfigMapAdd(obj interface{}) {
	configmap := obj.(*corev1.ConfigMap)
	configmapCache.PushRes(configmap)

	podList := configmapCache.GetPods(configmap)
	if len(podList) > 0 {
		logrus.Infof("configmap %v/%v updated, %v pod need update", configmap.GetNamespace(), configmap.GetName(), len(podList))
	}
	for _, pod := range podList {
		pod.SetResourceVersion("")
		pod.Annotations[ConfigMapAnnotation] = configmap.GetResourceVersion()
		_, err := clientSet.CoreV1().Pods(pod.GetNamespace()).Update(context.TODO(), pod, metav1.UpdateOptions{})
		if err != nil {
			logrus.WithError(err).Warnf("pod %v/%v update failed", pod.GetNamespace(), pod.GetName())
		}
	}
}

func handleConfigMapUpdate(oldObj, newObj interface{}) {
	handleConfigMapAdd(newObj)
}

func handleConfigMapDelete(obj interface{}) {
	configmap := obj.(*corev1.ConfigMap)
	configmapCache.DeleteRes(configmap)
}
