#!/bin/sh
BUILD_NAME=ottp-feedback
BUILD_DIR=./release

[ "$1" == "clean" ] && {
	echo "CleaningUp..."
	[ -d $BUILD_DIR/ ] && rm -rf $BUILD_DIR
	exit
}

build_bin() {
	echo "== Building: $GOOS-$GOARCH"
	[ "$GOOS" == "windows" ] && GOEXT=.exe || unset GOEXT
	FN=${BUILD_NAME}_${GOOS}-${GOARCH}
	go build -ldflags "-s -w -X main.depl_ver=$DEPL_VER" -o $BUILD_DIR/${FN}${GOEXT} $* && {
		echo "   OK"
		echo "   └ Packing TAR..."
		tar cf $BUILD_DIR/${FN}_${DEPL_VER}.tar --remove-files -C $BUILD_DIR/ ${FN}${GOEXT}
		tar Af $BUILD_DIR/${FN}_${DEPL_VER}.tar $BUILD_DIR/files_${DEPL_VER}.tar
		gzip -f $BUILD_DIR/${FN}_${DEPL_VER}.tar
	}
}

# Get versions 
GIT_COMM="`git log -n1 --pretty='%h'`"
GIT_TAG="`git describe --exact-match --tags $GIT_COMM || echo dev`"

DEPL_VER="${GIT_TAG}-${GIT_COMM}"
echo "[i] Version: $DEPL_VER"

[ ! -d $BUILD_DIR/ ] && mkdir -p $BUILD_DIR

# Static pack: create
tar cf $BUILD_DIR/files_${DEPL_VER}.tar -C $BUILD_DIR/../files/ --exclude='config.hjson' .
tar rf $BUILD_DIR/files_${DEPL_VER}.tar -C $BUILD_DIR/../ README.md

# Build
GOOS=windows GOARCH=amd64  build_bin
GOOS=linux   GOARCH=amd64  build_bin
GOOS=linux   GOARCH=arm64  build_bin
GOOS=linux   GOARCH=arm    build_bin
GOOS=linux   GOARCH=mipsle build_bin

# Static pack: remove
rm $BUILD_DIR/files_${DEPL_VER}.tar
