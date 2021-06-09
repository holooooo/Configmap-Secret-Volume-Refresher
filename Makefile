# Copyright 2021 Holooooo.
# Use of this source code is governed by the WTFPL
# license that can be found in the LICENSE file.

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/csvr
package: build
	docker build --platform=linux/amd64 -f package/Dockerfile -t csvr .