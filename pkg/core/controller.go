// Copyright 2021 Holooooo.
// Use of this source code is governed by the WTFPL
// license that can be found in the LICENSE file.

package core

import (
	"context"
	"csvr/pkg/cache"
	bootconfig "csvr/pkg/core/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

var (
	clientSet *kubernetes.Clientset

	configmapCache *cache.Cache
	secretCache    *cache.Cache
)

func Run(cs *kubernetes.Clientset, config *bootconfig.Config) error {
	clientSet = cs

	err := initCache(config)
	if err != nil {
		return err
	}

	stopCh := make(chan struct{})
	secretInformer := newSecretInformer(config.SecretSelector, config.ResyncDuration)
	configmapInformer := newConfigMapInformer(config.CmSelector, config.ResyncDuration)
	podInformer := newPodInformer(config.PodSelector, config.ResyncDuration)

	go func() {
		configmapInformer.Run(stopCh)
		tryClose(stopCh)
	}()
	go func() {
		secretInformer.Run(stopCh)
		tryClose(stopCh)
	}()
	go func() {
		podInformer.Run(stopCh)
		tryClose(stopCh)
	}()

	<-stopCh
	return nil
}

func initCache(config *bootconfig.Config) error {
	podList, err := clientSet.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		LabelSelector: config.PodSelector,
	})
	if err != nil {
		return err
	}
	configmapList, err := clientSet.CoreV1().ConfigMaps(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		LabelSelector: config.CmSelector,
	})
	if err != nil {
		return err
	}
	configmapCache = cache.NewConfigmapCache(configmapList.Items, podList.Items)

	secretList, err := clientSet.CoreV1().Secrets(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		LabelSelector: config.SecretSelector,
	})
	if err != nil {
		return err
	}
	secretCache = cache.NewSecretCache(secretList.Items, podList.Items)

	return nil
}

func tryClose(c chan struct{}) {
	select {
	case <-c:
	default:
		close(c)
	}
}
