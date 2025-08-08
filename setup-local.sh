#!/bin/sh

echo "[ ℹ  ] Updating packages"
go get -u all

echo "[ ℹ  ] Cleaning go mod"
go mod tidy

echo "[ ✅ ] Setup completed"
