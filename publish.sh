#!/bin/sh

echo "[ ✅ ] Updating Package Version"
GOPROXY=proxy.golang.org go list -m github.com/amitjangid80/go-mongodb-client@v1.1.4
