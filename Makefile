GOROOT=$(shell go env GOROOT)

$(info GOROOT $(GOROOT))

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

.PHONY: ibrowser ibrowser.wasm httpserver bin

bin: ibrowser ibrowser.wasm httpserver

ibrowser: bin/ibrowser

ibrowser.wasm: bin/ibrowser.wasm

httpserver: bin/httpserver


bin/ibrowser: */*.go
	cd main/ && go build -v -o ../$@ .

bin/ibrowser.wasm: */*.go
	cd main/ && GOOS=js GOARCH=wasm go build -ldflags="-s -w" -v -o ../$@ .

bin/httpserver: opt/httpserver/httpserver.go
	cd opt/httpserver/ && go build -v -o ../../$@ .



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

examples: 150_VCFs_2.50.tar.gz 360_merged_2.50.vcf.gz

150_VCFs_2.50.tar.gz:
	wget ftp://ftp.solgenomics.net/genomes/tomato_150/150_VCFs_2.50.tar.gz

360_merged_2.50.vcf.gz:
	wget ftp://ftp.solgenomics.net/genomes/tomato_360/360_merged_2.50.vcf.gz



.PHONY: run150 run360 prof

run150: 150_VCFs_2.50.tar.gz
	time bin/ibrowser 150_VCFs_2.50.tar.gz

run360: 360_merged_2.50.vcf.gz
	time bin/ibrowser 360_merged_2.50.vcf.gz

prof:
	rm ibrowser.cpu.prof ibrowser.mem.prof || true
	bin/ibrowser -cpuprofile ibrowser.cpu.prof -memprofile ibrowser.mem.prof 360_merged_2.50.vcf.gz
	go tool pprof -tree bin/ibrowser ibrowser.cpu.prof
	go tool pprof -tree bin/ibrowser ibrowser.mem.prof
