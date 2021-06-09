// Copyright 2021 Holooooo.
// Use of this source code is governed by the WTFPL
// license that can be found in the LICENSE file.

package config

import "time"

type Config struct {
	PodSelector    string
	CmSelector     string
	SecretSelector string
	ResyncDuration time.Duration
}
