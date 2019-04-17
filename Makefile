GOROOT=$(shell go env GOROOT)
$(info GOROOT  $(GOROOT))


ifndef FORMAT
FORMAT=yaml
endif

$(info FORMAT  $(FORMAT))


ifndef OUTFILE
OUTFILE=res/output
endif

$(info OUTFILE $(OUTFILE))


CGO_CFLAGS="-g -O3"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O3"
CGO_FFLAGS="-g -O3"
CGO_LDFLAGS="-g -O3"


GIT_COMMIT_HASH=$(shell git log --pretty=format:'%H' -n 1 | sed "s/\"/'/g")
GIT_COMMIT_AUTHOR=$(shell git log --pretty=format:'%an (%aE) %ai' -n 1 | sed "s/\"/'/g")
GIT_COMMIT_COMMITER=$(shell git log --pretty=format:'%cn (%cE) %ci' -n 1 | sed "s/\"/'/g")
GIT_COMMIT_NOTES=$(shell git log --pretty=format:'%N' -n 1 | sed "s/\"/'/g")
GIT_COMMIT_TITLE=$(shell git log --pretty=format:'%s' -n 1 | sed "s/\"/'/g")
GIT_STATUS=$(shell bash -c 'git diff-index --quiet HEAD; if [ "$$?" == "1" ]; then echo "dirty"; else echo "clean"; fi')
GIT_DIFF=$(shell bash -c 'git diff | md5sum --text | cut -f1 -d" "')
GO_VERSION=$(shell go version)
# TIMESTAMP=$(shell date +"%Y-%m-%d_%H-%M-%S")

# VERSION=$(GIT_COMMIT) - Status $(GIT_STATUS) - Diff $(GIT_DIFF)
#Timestamp $(TIMESTAMP)

$(info GIT_COMMIT_HASH     $(GIT_COMMIT_HASH))
$(info GIT_COMMIT_AUTHOR   $(GIT_COMMIT_AUTHOR))
$(info GIT_COMMIT_COMMITER $(GIT_COMMIT_COMMITER))
$(info GIT_COMMIT_NOTES    $(GIT_COMMIT_NOTES))
$(info GIT_COMMIT_TITLE    $(GIT_COMMIT_TITLE))
$(info GIT_STATUS          $(GIT_STATUS))
$(info GIT_DIFF            $(GIT_DIFF))
$(info GO_VERSION          $(GO_VERSION))
# $(info TIMESTAMP           $(TIMESTAMP))
# $(info VERSION             $(VERSION))
$(info )

LDFLAGS=-X 'main.IBROWSER_GIT_STATUS=$(GIT_STATUS)'
LDFLAGS+=-X 'main.IBROWSER_GIT_DIFF=$(GIT_DIFF)'
LDFLAGS+=-X 'main.IBROWSER_GIT_COMMIT_HASH=$(GIT_COMMIT_HASH)'
LDFLAGS+=-X 'main.IBROWSER_GIT_COMMIT_AUTHOR=$(GIT_COMMIT_AUTHOR)'
LDFLAGS+=-X 'main.IBROWSER_GIT_COMMIT_COMMITER=$(GIT_COMMIT_COMMITER)'
LDFLAGS+=-X 'main.IBROWSER_GIT_COMMIT_NOTES=$(GIT_COMMIT_NOTES)'
LDFLAGS+=-X 'main.IBROWSER_GIT_COMMIT_TITLE=$(GIT_COMMIT_TITLE)'
LDFLAGS+=-X 'main.IBROWSER_GO_VERSION=$(GO_VERSION)'

# $(info $(LDFLAGS))
# $(info )

.PHONY: help

help:
	@echo " bin"
	@echo " ibrowser"
	@echo " ibrowser.wasm"
	@echo " httpserver"
	@echo ""
	@echo " examples"
	@echo " serve"
	@echo ""
	@echo " get"
	@echo " requirements"
	@echo " wasm_exec.js"
	@echo ""
	@echo " run150"
	@echo " run360"
	@echo ""
	@echo " prof"

.PHONY: ibrowser ibrowser.wasm httpserver bin version

bin: ibrowser ibrowser.wasm httpserver

ibrowser: bin/ibrowser

ibrowser.exe: bin/ibrowser.exe

ibrowser.wasm: bin/ibrowser.wasm

httpserver: bin/httpserver

version:
	@echo "LDFLAGS $(LDFLAGS)"

bin/ibrowser: version */*.go
	cd main/ && GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -v -p 4 -o ../$@ .
	md5sum $@
	bin/ibrowser --version

bin/ibrowser.exe: version */*.go
	cd main/ && GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -v -p 4 -o ../$@ .
	@#bin/ibrowser.exe --version
	md5sum $@

bin/ibrowser.wasm: ibrowser */*.go
	cd main/ && GOOS=js GOARCH=wasm go build -ldflags="$(LDFLAGS) -s -w" -v -p 4 -o ../$@ .
	md5sum $@

bin/httpserver: opt/httpserver/httpserver.go
	cd opt/httpserver/ && go build -v -p 4 -o ../../$@ .



.PHONY: serve

serve: bin/httpserver wasm_exec.js
	bin/httpserver



.PHONY: requirements get

requirements: get httpserver wasm_exec.js

get:
	pip3 install -r requirements.txt
	go mod tidy -v
	cat go.mod

wasm_exec.js:
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" .



.PHONY: examples

examples: data/150_VCFs_2.50.tar.gz data/360_merged_2.50.vcf.gz

data/150_VCFs_2.50.tar.gz:
	wget --no-clobber https://s3.eu-central-1.amazonaws.com/saulo.ibrowser/360_merged_2.50.vcf.gz -O data/360_merged_2.50.vcf.gz.tmp && mv data/360_merged_2.50.vcf.gz.tmp data/360_merged_2.50.vcf.gz

data/360_merged_2.50.vcf.gz:
	wget --no-clobber https://s3.eu-central-1.amazonaws.com/saulo.ibrowser/150_VCFs_2.50.tar.gz   -O data/150_VCFs_2.50.tar.gz.tmp   && mv data/150_VCFs_2.50.tar.gz.tmp   data/150_VCFs_2.50.tar.gz



.PHONY: clean run150 run360

clean:
	rm -v $(OUTFILE)*.yaml   || true
	rm -v $(OUTFILE)*.bson   || true
	rm -v $(OUTFILE)*.bin    || true
	rm -v $(OUTFILE)*.gob    || true
	rm -v $(OUTFILE)*.gz     || true
	rm -v $(OUTFILE)*.snappy || true

run150: clean ibrowser data/150_VCFs_2.50.tar.gz
	time bin/ibrowser -format $(FORMAT) -outfile $(OUTFILE)_150_VCFs_2.50.tar.gz data/150_VCFs_2.50.tar.gz

run360: clean ibrowser data/360_merged_2.50.vcf.gz
	time bin/ibrowser -format $(FORMAT) -outfile $(OUTFILE)_360_merged_2.50.vcf.gz data/360_merged_2.50.vcf.gz



.PHONY: test test_load test_save

test: test_load test_save

test_save: clean ibrowser data/360_merged_2.50.vcf.gz
	time bin/ibrowser -format $(FORMAT) -outfile $(OUTFILE)_360_merged_2.50.vcf.gz -debugMaxRegisterChrom 1000 -threads 4 -check data/360_merged_2.50.vcf.gz

test_load: test_save
	bin/ibrowser -load -check res/output_360_merged_2.50.vcf.gz



.PHONY: clean run150 run360 prof prof_run

prof: prof_run ibrowser.cpu.prof
	go tool pprof -tree bin/ibrowser ibrowser.cpu.prof
	go tool pprof -tree bin/ibrowser ibrowser.mem.prof

ibrowser.cpu.prof: clean
	rm -v ibrowser.cpu.prof ibrowser.mem.prof || true
	time bin/ibrowser -cpuprofile ibrowser.cpu.prof -memprofile ibrowser.mem.prof -format $(FORMAT) data/360_merged_2.50.vcf.gz

prof_run: clean ibrowser data/360_merged_2.50.vcf.gz

check:
	ls -la
	ls -la res/
	ls -la bin/
	python3 ./check.py $(OUTFILE)_360_merged_2.50.vcf.gz
