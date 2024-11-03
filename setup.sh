#!/bin/sh

GHE_CONFIG_DIR
REPO="https://github.com/GMkonan/ghe.git"

error() {
	echo -e "\033[0;31m[ERROR]\033[0m $1"
	exit 1
}

git clone "$REPO"
cd ghe

if ! command -v go &> /dev/null; then
	error "Go is not installed"
fi

go build .

go install

echo -e "\033[0;32m Setup finished!"
