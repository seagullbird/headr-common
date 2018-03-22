#!/usr/bin/env bash

# Dispatcher
rm -rf dispatch/mock
mkdir dispatch/mock
mockgen -destination=dispatch/mock/mock_dispatch.go -package=mock github.com/seagullbird/headr-common/mq/dispatch Dispatcher

# Receiver
rm -rf receive/mock
mkdir receive/mock
mockgen -destination=receive/mock/mock_receive.go -package=mock github.com/seagullbird/headr-common/mq/receive Receiver
