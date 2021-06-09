// Copyright 2021 Holooooo.
// Use of this source code is governed by the WTFPL
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
	"time"

	"csvr/pkg/core"
	bootconfig "csvr/pkg/core/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	root = &cobra.Command{
		Use:   "csvr",
		Short: "CSVR can make configmap/secret volume be refreshed immediately",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := rest.InClusterConfig()
			if err != nil {
				logrus.WithError(err).Fatalf("Failed to load clientset config")
			}
			cs, err := kubernetes.NewForConfig(cfg)
			if err != nil {
				logrus.WithError(err).Fatalf("Failed to load clientset config")
			}

			err = core.Run(cs, config)
			if err != nil {
				logrus.WithError(err).Fatalf("Failed to run")
			}
		},
	}

	config = &bootconfig.Config{}
)

func init() {
	cobra.OnInitialize(initConfig)
	root.PersistentFlags().StringVar(&config.PodSelector, "pod_selector", "",
		"if not nil, controller will filter out pod than not contain label and value")
	root.PersistentFlags().StringVar(&config.CmSelector, "cm_selector", "",
		"if not nil, controller will filter out configmap than not contain label and value")
	root.PersistentFlags().StringVar(&config.SecretSelector, "secret_selector", "",
		"if not nil, controller will filter out secret than not contain label and value ")
	root.PersistentFlags().DurationVar(&config.ResyncDuration, "resync_duration", 2*time.Minute,
		"how often to list all the resouces for keep cache consistent with cluster")
}

func initConfig() {
	invalidSelector := ""
	if !isLabelSelector(config.SecretSelector) {
		invalidSelector = fmt.Sprintf("secret_selector %v", config.SecretSelector)
	}
	if !isLabelSelector(config.CmSelector) {
		invalidSelector = fmt.Sprintf("cm_selector %v", config.CmSelector)
	}
	if !isLabelSelector(config.PodSelector) {
		invalidSelector = fmt.Sprintf("pod_selector %v", config.PodSelector)
	}
	if len(invalidSelector) > 0 {
		logrus.Fatalf("%v is invalid, please check input", invalidSelector)
	}
}

func isLabelSelector(label string) bool {
	if len(label) == 0 {
		return true
	}
	labelSplit := strings.Split(label, "=")
	if len(labelSplit) != 2 {
		return false
	}
	return false
}

func main() {
	logrus.Info("Starting CSVR Controller")
	err := root.Execute()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to Start, exited")
	}
}
