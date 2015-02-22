#!/bin/bash

set -e

[ -z "$IRCD_PORT_6667_TCP_ADDR" ] && echo "Need to link an ircd as a container." && exit 1;

go test ./...
