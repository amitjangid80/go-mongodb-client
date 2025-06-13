#!/bin/sh

echo "[ ℹ  ] Updating packages"
go mod download

echo "[ ℹ  ] Cleaning go mod"
go mod tidy

echo "[ ✅ ] Setup completed"
