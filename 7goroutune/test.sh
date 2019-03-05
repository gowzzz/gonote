#!/bin/bash
# install flamegraph scripts
if [ ! -d "/opt/flamegraph" ]; then
    echo "Installing flamegraph (git clone)"
    git clone --depth=1 https://github.com/brendangregg/FlameGraph.git /opt/flamegraph
fi

# install go-torch using docker
if [ ! -f "bin/go-torch" ]; then
    echo "Installing go-torch via docker"
    docker run --net=party --rm=true -it -v $(pwd)/bin:/go/bin golang go get github.com/uber/go-torch
    # or if you have go installed locally: go get github.com/uber/go-torch
fi

PATH="$PATH:/opt/flamegraph"
bin/go-torch -b profile_cpu.out -f profile_cpu.torch.svg