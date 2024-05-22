default:
    go build -ldflags "-s -w" ./cmd/toru/

build:
    go build -ldflags "-s -w" ./cmd/toru/

install:
    test -x ./toru
    cp ./toru /usr/bin
