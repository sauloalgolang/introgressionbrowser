GOROOT=$(shell go env GOROOT)

$(info GOROOT $(GOROOT))

.PHONY: help
help:
	@echo " examples"
	@echo " get"
	@echo " httpserver"
	@echo " ibrowser"
	@echo " ibrowser.wasm"
	@echo " main"
	@echo " requirements"
	@echo " serve"
	@echo " wasm_exec.js"

main: src/main/main.go
	cd src/main/ && go build -v -o ../../$@ .

ibrowser: src/ibrowser/ibrowser.go
	cd src/ibrowser/ && go build -v -o ../../$@ .

ibrowser.wasm: src/ibrowser/ibrowser.go
	cd src/ibrowser/ && GOOS=js GOARCH=wasm go build -v -o ../../$@ .

httpserver: src/httpserver/httpserver.go
	cd src/httpserver && go build -v -o ../../$@ .



.PHONY: serve

serve: httpserver wasm_exec.js
	./httpserver


.PHONY: requirements get

requirements: get httpserver wasm_exec.js

get:
	go get -v .

wasm_exec.js:
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" .



.PHONY: examples

examples: 150_VCFs_2.50.tar.gz 360_merged_2.50.vcf.gz

150_VCFs_2.50.tar.gz:
	wget ftp://ftp.solgenomics.net/genomes/tomato_150/150_VCFs_2.50.tar.gz

360_merged_2.50.vcf.gz:
	wget ftp://ftp.solgenomics.net/genomes/tomato_360/360_merged_2.50.vcf.gz
