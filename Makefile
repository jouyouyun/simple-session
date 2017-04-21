PREFIX=/usr
GOPATH_DIR=gopath
GOBUILD=go build
PRJ=simple-session

all: build

prepare:
	@mkdir -p out
	@if [ ! -d ${GOPATH_DIR}/src/${PRJ} ]; then \
		mkdir -p ${GOPATH_DIR}/src; \
		ln -sf ${CURDIR}/src ${GOPATH_DIR}/src/${PRJ}; \
	fi

build: out/${PRJ}

out/${PRJ}: prepare
	env GOPATH="${CURDIR}/${GOPATH_DIR}:${GOPATH}" ${GOBUILD} -o $@ ${PRJ}

clean:
	rm -rf ${GOPATH_DIR}
	rm -rf out

rebuild: clean build
