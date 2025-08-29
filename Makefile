### 构建项目

bin := pubot

.PHONY: clean, build, test

# 设置静态编译
CLIB_ENABLE := $(go env CGO_ENABLED)
ifeq ($(CLIB_ENABLE),1)
	go env -w CGO_ENABLED=0
endif

build:
	go build -C cmd/pubot -trimpath -ldflags="-s -w"  -o $(bin)

clean:
	rm -f cmd/pubot/$(bin)

test:
	@go test -v -count=1 -skip="TestUserGet|TestUserDelete|TestUserUpdate|TestUserList|TestUserCreate" ./...
	