#!/bin/bash

docker build -t elemental-ircd-test .

docker run --rm -it --link hub:ircd elemental-ircd-test
