#!/data/data/com.termux/files/usr/bin/bash
set -e

cd "$(dirname "$0")/go"
go build -o ../jqmas-core ./cmd/jqmas-core
echo "✅ jqmas-core compilado en $(pwd)/../jqmas-core"
