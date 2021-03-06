# To include this protoc/make snippet, copy this to the top of your Makefile:
#
# default: build
#
# genproto.mk:
# 	@docker pull lightstep/gogoprotoc:latest
# 	@-docker rm -v lightstep-get-genproto-mk
# 	@docker create --name lightstep-get-genproto-mk lightstep/gogoprotoc:latest
# 	@docker cp lightstep-get-genproto-mk:/genproto.mk genproto.mk
# 	@docker rm -v lightstep-get-genproto-mk
#
# include genproto.mk

GOLANG = golang
PBUF   = protobuf
GOGO   = gogo

# Use the final GOPATH element, since that's where circleci puts the code (lame!)
ROOT = $(shell echo ${GOPATH} | tr : \\n | tail -1)

PWD     = $(shell pwd)
TMPNAME = tmpgen
TMPDIR  = $(PWD)/$(TMPNAME)

# List of standard protoc options
PROTOC_OPTS = plugins=grpc

# These flags manage mapping the google-standard protobuf types (e.g., Timestamp)
# into the annotated versions supplied with Gogo.  The trailing `,` matters.
GOGO_OPTS = Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,

define protos_to_gogo_targets
$(foreach proto,$(1),$(GOLANG)-$(GOGO)-$(basename $(proto)))
endef

define protos_to_protobuf_targets
$(foreach proto,$(1),$(GOLANG)-$(PBUF)-$(basename $(proto)))
endef

define gen_gogo_target
$(call gen_protoc_target,$(1),$(GOLANG)/$(GOGO)/$(basename $(1))pb/$(basename $(1)).pb.go,$(GOGO),--gogofaster_out=$(GOGO_OPTS)$(PROTOC_OPTS))
endef

define gen_protobuf_target
$(call gen_protoc_target,$(1),$(GOLANG)/$(PBUF)/$(basename $(1))pb/$(basename $(1)).pb.go,$(PBUF),--go_out=$(PROTOC_OPTS))
endef

define protoc_targets_to_link_targets
$(foreach target,$(1),$(target)-link)
endef

define gen_protoc_link
@mkdir -p "$(subst -,/,$(subst -link,,$(2)))pb"
@echo "// DO NOT EDIT; THIS FILE IS AUTOGENERATED FROM ../../../$(1)" > $(subst -,/,$(subst -link,,$(2)))pb/$(1)
@sed -E 's@import "github.com/lightstep/([^/]+)/(.*).proto"@import "github.com/lightstep/\1/$(GOLANG)/$(3)/\2pb/\2.proto"@g' < $(1) >> $(subst -,/,$(subst -link,,$(2)))pb/$(1)
endef

# $(1) = .proto input
# $(2) = .pb.go output
# $(3) = gogo or protobuf
# $(4) = protoc-output spec
#
# Note: the --proto_path include "." below references the
# docker image's $(ROOT)/src.  /input is mapped to the
# host's $(ROOT)/src.
define gen_protoc_target
  @echo compiling $(1) [$(3)]
  @mkdir -p $(TMPDIR)
  @sed -E 's@import "github.com/lightstep/([^/]+)/(.*).proto"@import "github.com/lightstep/\1/$(GOLANG)/$(3)/\2pb/\2.proto"@g' < $(1) > $(TMPDIR)/$(1)
  @docker run --rm \
    -v $(ROOT)/src:/input:ro \
    -v $(TMPDIR):/output \
    lightstep/gogoprotoc:latest \
    protoc \
    -I./github.com/google/googleapis \
    $(4):/output \
    --proto_path=/input:. \
    /input/$(PKG_PREFIX)/$(TMPNAME)/$(1)
  @mkdir -p $(GOLANG)/$(3)/$(basename $(1))pb/$(basename $(1))pbfakes
  @sed 's@package $(basename $(1))pb@package $(basename $(1))pb // import "$(PKG_PREFIX)/golang/$(3)/$(basename $(1))pb"@' < $(TMPDIR)/$(PKG_PREFIX)/$(TMPNAME)/$(basename $(1)).pb.go > $(GOLANG)/$(3)/$(basename $(1))pb/$(basename $(1)).pb.go
  @rm $(TMPDIR)/$(PKG_PREFIX)/$(TMPNAME)/$(basename $(1)).pb.go
  @rm $(TMPDIR)/$(1)
endef

define clean_protoc_targets
  @rm -rf $(foreach target,$(1),$(subst -,/,$(target)pb))
endef

# generate_fake: runs counterfeiter in docker container to generate fake classes
# $(1) output file path
# $(2) input file path
# $(3) class name
define generate_fake
  @docker run --rm \
	-v $(ROOT):/usergo \
	lightstep/gobuild:latest \
	/bin/bash -c "cd /usergo/src/$(PKG_PREFIX) && counterfeiter -o $(1) $(2) $(3)"
endef
