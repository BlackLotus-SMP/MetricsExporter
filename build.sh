
#!/bin/bash

DIR=$(dirname "$0")/build
mkdir -p "$DIR"
cd "$DIR" || exit

sum="sha1sum"

export GO111MODULE=on
echo "Setting GO111MODULE to" $GO111MODULE

if ! hash sha1sum 2>/dev/null; then
	if ! hash shasum 2>/dev/null; then
		echo "I can't see 'sha1sum' or 'shasum'"
		echo "Please install one of them!"
		exit
	fi
	# shellcheck disable=SC2034
	sum="shasum"
fi

UPX=false
if hash upx 2>/dev/null; then
	UPX=true
fi

VERSION=$(date -u +%Y%m%d)
LDFLAGS="-X main.VERSION=$VERSION -s -w"
GCFLAGS=""

# AMD64
OSES=(linux darwin windows)
for os in "${OSES[@]}"; do
	suffix=""
	if [ "$os" == "windows" ]
	then
		suffix=".exe"
	fi
	env CGO_ENABLED=0 GOOS="$os" GOARCH=amd64 go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o "${os}"_amd64"${suffix}" ..
done

# 386
OSES=(linux windows)
for os in "${OSES[@]}"; do
	suffix=""
	if [ "$os" == "windows" ]
	then
		suffix=".exe"
	fi
	env CGO_ENABLED=0 GOOS="$os" GOARCH=386 go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o "${os}"_386"${suffix}" ..
done

#Apple M1 device
env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o darwin_arm64 ..

# ARM
ARMS=(5 6 7)
for v in "${ARMS[@]}"; do
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM="$v" go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o linux_arm"${v}" ..
done

# ARM64
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o linux_arm64 ..

# MIPS32LE
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o linux_mipsle ..
env CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o linux_mips ..

if $UPX; then upx -9 *;fi