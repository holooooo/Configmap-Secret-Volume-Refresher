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

const SecretAnnotation = "holooooo.io/lastest-secret-version"

func newSecretInformer(selector string, resync time.Duration) informerCache.SharedIndexInformer {
	i := informerCache.NewSharedIndexInformer(
		&informerCache.ListWatch{
			ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
				opts.LabelSelector = selector
				return clientSet.CoreV1().Secrets(metav1.NamespaceAll).List(context.TODO(), opts)
			},
			WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
				opts.LabelSelector = selector
				return clientSet.CoreV1().Secrets(metav1.NamespaceAll).Watch(context.TODO(), opts)
			},
		}, &corev1.Secret{}, 0, informerCache.Indexers{},
	)
	i.AddEventHandlerWithResyncPeriod(
		&informerCache.ResourceEventHandlerFuncs{
			AddFunc:    handleSecretAdd,
			UpdateFunc: handleSecretUpdate,
			DeleteFunc: handleSecretDelete,
		},
		resync,
	)
	return i
}

func handleSecretAdd(obj interface{}) {
	secret := obj.(*corev1.Secret)
	secretCache.PushRes(secret)

	podList := secretCache.GetPods(secret)
	if len(podList) > 0 {
		logrus.Infof("secret %v/%v updated, %v pod need update", secret.GetNamespace(), secret.GetName(), len(podList))
	}
	for _, pod := range podList {
		pod.SetResourceVersion("")
		pod.Annotations[SecretAnnotation] = secret.GetResourceVersion()
		_, err := clientSet.CoreV1().Pods(pod.GetNamespace()).Update(context.TODO(), pod, metav1.UpdateOptions{})
		if err != nil {
			logrus.WithError(err).Warnf("pod %v/%v update failed", pod.GetNamespace(), pod.GetName())
		}
	}
}

func handleSecretUpdate(oldObj, newObj interface{}) {
	handleSecretAdd(newObj)
}

func handleSecretDelete(obj interface{}) {
	secret := obj.(*corev1.Secret)
	secretCache.DeleteRes(secret)
}
