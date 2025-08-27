GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
	ERROR_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *_error.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
	ERROR_PROTO_FILES=$(shell find . api -name *_error.proto)
endif

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
 	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: errors
# generate errors proto
errors:
	protoc --proto_path=. \
             --proto_path=./third_party \
             --go_out=paths=source_relative:. \
             --go-errors_out=paths=source_relative:. \
             $(ERROR_PROTO_FILES)

.PHONY: swag_admin
# swag_admin
swag_admin:
	#go install github.com/swaggo/swag/cmd/swag@latest
	#swag init -g ./apps/admin/service/task_manage.go --parseDependency --parseInternal -o apps/admin/docs
	swag init -g main.go -d ./apps/admin,\
	./framework/models/rule_engine \
 	--parseDependency --parseInternal -o apps/admin/docs
	#swag init -g main.go -d ./apps/admin,\
#	./framework/models/template/interfaces/admin,\
#	./framework/models/config_center,\
#	./internal/payment/interfaces/admin,\
#	./framework/models/rule_engine\
# 	--parseDependency --parseInternal -o apps/admin/docs

.PHONY: swag_app
# swag_app
swag_app:
	#go install github.com/swaggo/swag/cmd/swag@latest
	swag init   -d ./apps/app,\
	./internal/user/interfaces/app,\
	./internal/declaration/interfaces/app,\
	./internal/payment/interfaces/app \
	 --parseDependency --parseInternal -o ./apps/app/docs

.PHONY: wire_admin
# wire_admin
wire_admin:
	cd apps/admin/  && wire

.PHONY: wire_app
# wire_app
wire_app:
	cd apps/app/ && wire


.PHONY: build_admin
# build
build_admin:
	make wire_admin;
	make swag_admin;
	mkdir -p bin/ && GOOS=linux GOARCH=amd64 go build -ldflags='-checklinkname=0 -extldflags "-static" -s -w -X main.Version=$(VERSION)' -o ./bin/f02-admin ./apps/admin
	#mkdir -p bin/ && go build -ldflags '-extldflags "-static" -s -w -X main.Version=$(VERSION)' -o ./bin/f02-admin ./apps/admin
	chmod +x ./bin/f02-admin

.PHONY: build_app
# build
build_app:
	rm -rf ./bin/app;
	make wire_app;
	make swag_app;
	mkdir -p bin/ && GOOS=linux GOARCH=amd64 go build -ldflags='-checklinkname=0 -extldflags "-static" -s -w -X main.Version=$(VERSION)' -o ./bin/f02-app ./apps/app
	#mkdir -p bin/ &&  go build -ldflags '-extldflags "-static" -s -w -X main.Version=$(VERSION)' -o ./bin/f02-app ./apps/app
	chmod +x ./bin/f02-app