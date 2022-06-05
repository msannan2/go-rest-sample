GO=go

ifneq ($(filter $(verbose), on ON yes YES 1),)
    export VERBOSE=1
endif

ifeq ($(VERBOSE), 1)
    q=# nil
    export LOG_LEVEL=VERBOSE
else
    q=@
    ifndef LOG_LEVEL
        export IOT_LOG_LEVEL=ERROR
    endif
endif

tidy:
	@echo go tidy
	$(q)$(GO) mod tidy > /dev/null

build: tidy cmd/sample/sample

cmd/sample/sample:
	@echo go $@
	$(q)$(GO) build $(GOFLAGS) -o $@ ./cmd/sample

docker_go_sample:
	docker build --rm -f cmd/sample/Dockerfile -t go-rest-sample:1.0 .

run:
	docker run --rm -p 8080:$(PORT) -e APP_PORT=$(PORT) -e PERSISTENT_FILE=$(PERSISTENT_FILE) --name go-rest-sample go-rest-sample:1.0

clean:
	@echo Cleaning...
	$(q)rm -f cmd/sample/sample

