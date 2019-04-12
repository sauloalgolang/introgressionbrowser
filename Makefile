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

BREAKAT=100
$(info BREAKAT $(BREAKAT))

GIT_COMMIT_HASH=$(shell git log --pretty=format:'%H' -n 1 | sed "s/\"/'/g"))
GIT_COMMIT_AUTHOR=$(shell git log --pretty=format:'%an (%aE) %ai' -n 1 | sed "s/\"/'/g")
GIT_COMMIT_COMMITER=$(shell git log --pretty=format:'%cn (%cE) %ci' -n 1 | sed "s/\"/'/g")
GIT_COMMIT_NOTES=$(shell git log --pretty=format:'%N' -n 1 | sed "s/\"/'/g")
GIT_COMMIT_TITLE=$(shell git log --pretty=format:'%s' -n 1 | sed "s/\"/'/g")
GIT_STATUS=$(shell bash -c 'git diff-index --quiet HEAD; if [ "$$?" == "1" ]; then echo "dirty"; else echo "clean"; fi')
GIT_DIFF=$(shell bash -c 'git diff | md5sum --text | cut -f1 -d" "')
TIMESTAMP=$(shell date +"%Y-%m-%d_%H-%M-%S")

# VERSION=$(GIT_COMMIT) - Status $(GIT_STATUS) - Diff $(GIT_DIFF)
#Timestamp $(TIMESTAMP)

$(info GIT_COMMIT_HASH     $(GIT_COMMIT_HASH))
$(info GIT_COMMIT_AUTHOR   $(GIT_COMMIT_AUTHOR))
$(info GIT_COMMIT_COMMITER $(GIT_COMMIT_COMMITER))
$(info GIT_COMMIT_NOTES    $(GIT_COMMIT_NOTES))
$(info GIT_COMMIT_TITLE    $(GIT_COMMIT_TITLE))
$(info GIT_STATUS          $(GIT_STATUS))
$(info GIT_DIFF            $(GIT_DIFF))
$(info TIMESTAMP           $(TIMESTAMP))
# $(info VERSION             $(VERSION))
$(info )

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
	@echo 'package main\n' > main/commit.go
	@echo 'const IBROWSER_GIT_COMMIT_HASH     = "$(GIT_COMMIT_HASH)"' >> main/commit.go
	@echo 'const IBROWSER_GIT_COMMIT_AUTHOR   = "$(GIT_COMMIT_AUTHOR)"' >> main/commit.go
	@echo 'const IBROWSER_GIT_COMMIT_COMMITER = "$(GIT_COMMIT_COMMITER)"' >> main/commit.go
	@echo 'const IBROWSER_GIT_COMMIT_NOTES    = "$(GIT_COMMIT_NOTES)"' >> main/commit.go
	@echo 'const IBROWSER_GIT_COMMIT_TITLE    = "$(GIT_COMMIT_TITLE)"' >> main/commit.go
	@echo 'const IBROWSER_GIT_STATUS          = "$(GIT_STATUS)"' >> main/commit.go
	@echo 'var   IBROWSER_GIT_DIFF            = ""' >> main/commit.go
	@echo '\n' >> main/commit.go
	cat main/commit.go
	#@echo 'package main\n' > main/commit_diff.go
	# @echo 'var   IBROWSER_GIT_DIFF            = ""' >> main/commit_diff.go
	# cat main/commit_diff.go

bin/ibrowser: version */*.go
	cd main/ && go build -v -p 4 -o ../$@ .
	# cd main/ && go build -v -race -ldflags "-X main.BREAKAT=$(BREAKAT)" -p 4 -o ../$@ .
	bin/ibrowser --version

bin/ibrowser.exe: version */*.go
	cd main/ && GOOS=windows GOARCH=amd64 go build -v -p 4 -o ../$@ .
	#bin/ibrowser.exe --version

bin/ibrowser.wasm: ibrowser */*.go
	cd main/ && GOOS=js GOARCH=wasm go build -ldflags="-s -w" -v -p 4 -o ../$@ .

bin/httpserver: opt/httpserver/httpserver.go
	cd opt/httpserver/ && go build -v -p 4 -o ../../$@ .



.PHONY: serve

serve: bin/httpserver wasm_exec.js
	bin/httpserver



.PHONY: requirements get

requirements: get httpserver wasm_exec.js

get:
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



.PHONY: clean run150 run360 prof prof_run

clean:
	rm -v $(OUTFILE)*.yaml | true
	rm -v $(OUTFILE)*.bson | true
	rm -v $(OUTFILE)*.bin  | true
	rm -v $(OUTFILE)*.gob  | true

run150: clean ibrowser data/150_VCFs_2.50.tar.gz
	time bin/ibrowser -format $(FORMAT) -outfile $(OUTFILE)_150_VCFs_2.50.tar.gz data/150_VCFs_2.50.tar.gz

run360: clean ibrowser data/360_merged_2.50.vcf.gz
	time bin/ibrowser -format $(FORMAT) -outfile $(OUTFILE)_360_merged_2.50.vcf.gz data/360_merged_2.50.vcf.gz

prof: prof_run ibrowser.cpu.prof
	go tool pprof -tree bin/ibrowser ibrowser.cpu.prof
	go tool pprof -tree bin/ibrowser ibrowser.mem.prof

ibrowser.cpu.prof: clean
	rm -v ibrowser.cpu.prof ibrowser.mem.prof || true
	time bin/ibrowser -cpuprofile ibrowser.cpu.prof -memprofile ibrowser.mem.prof -format $(FORMAT) data/360_merged_2.50.vcf.gz

prof_run: clean ibrowser data/360_merged_2.50.vcf.gz

check: output.yaml
	grep -v " -" output.yaml
	./check.py output
