default:
    go build -ldflags "-s -w" ./cmd/toru/

build:
    go build -ldflags "-s -w" ./cmd/toru/

android:
    GOARCH=arm64 GOOS=android go build -ldflags "-s -w" ./cmd/toru/

mac:
    GOARCH=arm64 GOOS=darwin go build -ldflags "-s -w" ./cmd/toru/

openbsd:
    GOARCH=amd64 GOOS=openbsd go build -ldflags "-s -w" ./cmd/toru/

windows:
    GOARCH=amd64 GOOS=windows go build -ldflags "-s -w" ./cmd/toru/

install:
    test -x ./toru
    cp ./toru /usr/bin
