## 打包项目
## 依赖前端项目: pubot-web -> npm run build

target := bin

build_args := -ldflags="-w -s" -trimpath

.PHONY: clean build

# 设置静态编译
CLIB_ENABLE := $(shell go env CGO_ENABLED)
ifeq ($(CLIB_ENABLE),1)
	go env -w CGO_ENABLED=0
endif

# 设置平台
GO_PLATFORM := $(shell go env GOOS)
ifneq ($(GO_PLATFORM), linux)
    go env -w GOOS=linux
endif
# 设置架构
GO_ARCH := $(shell go env GOARCH)
ifneq ($(GO_ARCH), amd64)
    go env -w GOARCH=amd64
endif

build:
	@if [ ! -d $(target) ]; then  \
	  mkdir $(target);  \
	fi
	go build $(build_args) -o $(target)
	@cp config.yaml bin/config.yaml

clean:
	rm -fr $(target)