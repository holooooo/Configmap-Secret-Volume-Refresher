# Copyright 2021 Holooooo.
# Use of this source code is governed by the WTFPL
# license that can be found in the LICENSE file.

FROM alpine

WORKDIR /home
COPY ./out/csvr /usr/local/csvr
RUN chmod +x /usr/local/csvr
ENTRYPOINT [ "/usr/local/csvr" ]