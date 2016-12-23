RM=rm -f
RMFORCE=rm -rf
MKDIR = mkdir -p
DESTDIR=$(SR_CODE_BASE)/snaproute/src/out/bin
SR_PLATFORMD_CODE_BASE=$(SR_CODE_BASE)/snaproute/src/infra/platformd
PARAMSDIR = $(DESTDIR)/params
PLATFORMDIR = $(PARAMSDIR)/platformd
GENERATED_IPC=$(SR_CODE_BASE)/generated/src
IPC_GEN_CMD=thrift
SRCS=main.go
IPC_SRCS=rpc/platformd.thrift
COMP_NAME=platformd
GOLDFLAGS=-r /opt/flexswitch/sharedlib


ifeq ($(PLATFORM_BUILD_TARGET), )
PLATFORM_BUILD_TARGET = dummy
endif

#CGO CFLAGS
CGO_CFLAGS=-I$(SR_PLATFORMD_CODE_BASE)/pluginManager/pluginCommon/
export CGO_CFLAGS+=-I$(SR_PLATFORMD_CODE_BASE)/pluginManager/onlp/inc
export CGO_CFLAGS+=-I$(SR_PLATFORMD_CODE_BASE)/pluginManager/openBMCVoyager
export CGO_LDFLAGS+=-ldl -Wl,--unresolved-symbols=ignore-in-shared-libs
export CGO_LDFLAGS+=-lm

.PHONY: all

all:ipc exe
ipc:
	$(IPC_GEN_CMD) -r --gen go -out $(GENERATED_IPC) $(IPC_SRCS)

exe: $(SRCS)
	@bash makeSymLinks.sh $(PLATFORM_BUILD_TARGET)
	go build -o $(DESTDIR)/$(COMP_NAME) -ldflags="$(GOLDFLAGS)" $(SRCS)

guard:
ifndef SR_CODE_BASE
	$(error SR_CODE_BASE is not set)
endif

install:
	$(MKDIR) $(PLATFORMDIR)
	install params/openBMC_platformd.conf $(PLATFORMDIR)/openBMC_platformd.conf
	install params/onlp_platformd.conf $(PLATFORMDIR)/onlp_platformd.conf
	install params/dummy_platformd.conf $(PLATFORMDIR)/dummy_platformd.conf
	install params/openBMCVoyager_platformd.conf $(PLATFORMDIR)/openBMCVoyager_platformd.conf
	install params/openBMCOpen19_platform.conf $(PLATFORMDIR)/openBMCOpen19_platform.conf
	install params/dummy_platformd.conf $(PLATFORMDIR)/dummy_platformd.conf

clean:guard
	$(RM) $(DESTDIR)/$(COMP_NAME)
	$(RMFORCE) $(GENERATED_IPC)/$(COMP_NAME)
