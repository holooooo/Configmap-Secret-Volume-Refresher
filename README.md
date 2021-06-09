<!--
 Copyright 2021 Holooooo.
 Use of this source code is governed by the WTFPL
 license that can be found in the LICENSE file.
-->

# Configmap Secret Volume Refresher(CSVR)

CSVR can make configmap/secret volume be refreshed immediately.

After kuberlet 1.17, the default value of `sync-frequency` changed from 10s to 1min.
That means onfigmap/secret volume may take up to 1min to get new data from configmap/secret.

Fortunately, we can refresh volume immediately by edit an annotation of pod.This will not cause pod restart.
And csvr is a controller to listen pod, secret, configmap changes, and refresh pod after volume needs update.

## Usage

Clone this repo, and apply the `deployment.yaml` and `rbac.yaml`. Then it will listen pod of all the namespaces.

You can also add the following flag to `args`.

| Flag            | Describe                                                                          | Default | Example               |
| --------------- | --------------------------------------------------------------------------------- | ------- | --------------------- |
| pod_selector    | if not nil, controller will filter out pod than not contain label and value       | nil     | holooooo.io/csvr=true |
| cm_selector     | if not nil, controller will filter out configmap than not contain label and value | nil     | holooooo.io/csvr=true |
| secret_selector | if not nil, controller will filter out secret than not contain label and value    | nil     | holooooo.io/csvr=true |
| resync_duration | how often to list all the resouces for keep cache consistent with cluster         | 2min    | 30s                   |

## Warn

It not supported `subpath`.
