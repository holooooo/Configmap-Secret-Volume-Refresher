// Copyright 2021 Holooooo.
// Use of this source code is governed by the WTFPL
// license that can be found in the LICENSE file.

package cache

import "fmt"

func format(name, namespace string) string {
	return fmt.Sprintf("%v/%v", namespace, name)
}
