// Copyright 2021 Holooooo.
// Use of this source code is governed by the WTFPL
// license that can be found in the LICENSE file.

package cache

import (
	"sync"

	corev1 "k8s.io/api/core/v1"
)

type Cache struct {
	mx sync.RWMutex

	resourceType resourceType

	// key is Resources namespace/name
	// value is all pods that depend on the resources
	influenceCache map[string]map[string]*corev1.Pod
	podCache       map[string]*corev1.Pod
	resourcesCache map[string]resource
}

type resource interface {
	GetName() string
	GetNamespace() string
}

type resourceType string

var (
	ConfigmapType resourceType = "configmap"
	SecretType    resourceType = "secret"
)

func NewConfigmapCache(resList []corev1.ConfigMap, podList []corev1.Pod) *Cache {
	cache := &Cache{
		resourceType:   ConfigmapType,
		influenceCache: make(map[string]map[string]*corev1.Pod),
		podCache:       make(map[string]*corev1.Pod),
		resourcesCache: make(map[string]resource),
	}
	for _, res := range resList {
		cache.PushRes(&res)
	}
	for _, pod := range podList {
		cache.PushPod(&pod)
	}
	return cache
}

func NewSecretCache(resList []corev1.Secret, podList []corev1.Pod) *Cache {
	cache := &Cache{
		resourceType:   SecretType,
		influenceCache: make(map[string]map[string]*corev1.Pod),
		podCache:       make(map[string]*corev1.Pod),
		resourcesCache: make(map[string]resource),
	}
	for _, res := range resList {
		cache.PushRes(&res)
	}
	for _, pod := range podList {
		cache.PushPod(&pod)
	}
	return cache
}

func (c *Cache) GetPods(obj resource) []*corev1.Pod {
	c.mx.RLock()
	defer c.mx.RUnlock()

	var res = make([]*corev1.Pod, 0)
	for _, pod := range c.influenceCache[format(obj.GetName(), obj.GetNamespace())] {
		res = append(res, pod)
	}
	return res
}

func (c *Cache) PushRes(obj resource) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.resourcesCache[format(obj.GetName(), obj.GetNamespace())] = obj
}

func (c *Cache) PushPod(obj *corev1.Pod) {
	c.mx.Lock()
	defer c.mx.Unlock()

	objFmt := format(obj.Name, obj.Namespace)
	oldpod := c.podCache[objFmt]
	newObjRes := c.getAllRes(obj)
	addRes, delRes := make(map[string]resource), make(map[string]resource)
	if oldpod != nil {
		oldRes := c.getAllRes(oldpod)
		for k, v := range newObjRes {
			if oldRes[k] != nil {
				continue
			}
			addRes[k] = v
		}
		for k, v := range oldRes {
			if newObjRes[k] != nil {
				continue
			}
			delRes[k] = v
		}
	} else {
		addRes = newObjRes
	}

	for key, val := range addRes {
		tempCache := c.influenceCache[key]
		if tempCache == nil {
			tempCache = make(map[string]*corev1.Pod)
		}
		tempCache[objFmt] = obj
		resfmt := format(val.GetName(), val.GetNamespace())
		c.influenceCache[resfmt] = tempCache
	}

	for key, val := range delRes {
		tempCache := c.influenceCache[key]
		if tempCache == nil {
			continue
		}
		delete(tempCache, objFmt)
		c.influenceCache[format(val.GetName(), val.GetNamespace())] = tempCache
	}

	c.podCache[objFmt] = obj
}

func (c *Cache) DeleteRes(obj resource) {
	c.mx.Lock()
	defer c.mx.Unlock()

	objFmt := format(obj.GetName(), obj.GetNamespace())
	delete(c.resourcesCache, objFmt)
	delete(c.influenceCache, objFmt)
}

func (c *Cache) DeletePod(obj *corev1.Pod) {
	c.mx.Lock()
	defer c.mx.Unlock()

	resMap := c.getAllRes(obj)
	for _, res := range resMap {
		delete(c.influenceCache[format(res.GetName(), res.GetNamespace())], format(obj.Name, obj.Namespace))
	}
	delete(c.podCache, format(obj.Name, obj.Namespace))
}

func (c *Cache) getAllRes(obj *corev1.Pod) map[string]resource {
	var res = make(map[string]resource)
	for _, v := range obj.Spec.Volumes {
		configFmt := ""
		switch c.resourceType {
		case ConfigmapType:
			if v.ConfigMap == nil {
				continue
			}
			configFmt = format(v.ConfigMap.Name, obj.Namespace)
		case SecretType:
			if v.Secret == nil {
				continue
			}
			configFmt = format(v.Secret.SecretName, obj.Namespace)
		}
		if newRes, ok := c.resourcesCache[configFmt]; ok {
			res[configFmt] = newRes
		}
	}
	return res
}
