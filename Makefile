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
	@echo ""
	@echo " prof"

.PHONY: ibrowser ibrowser.wasm httpserver bin

bin: ibrowser ibrowser.wasm httpserver

ibrowser: bin/ibrowser

ibrowser.wasm: bin/ibrowser.wasm

httpserver: bin/httpserver


bin/ibrowser: */*.go
	cd main/ && go build -v -o ../$@ .

bin/ibrowser.exe: */*.go
	cd main/ && GOOS=windows GOARCH=amd64 go build -v -o ../$@ .

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



.PHONY: clean run150 run360 prof prof_run

clean:
	rm -v output*.yaml | true

run150: clean ibrowser 150_VCFs_2.50.tar.gz
	time bin/ibrowser 150_VCFs_2.50.tar.gz

run360: clean ibrowser 360_merged_2.50.vcf.gz
	time bin/ibrowser 360_merged_2.50.vcf.gz

prof: prof_run ibrowser.cpu.prof
	go tool pprof -tree bin/ibrowser ibrowser.cpu.prof
	go tool pprof -tree bin/ibrowser ibrowser.mem.prof

ibrowser.cpu.prof: clean
	rm -v ibrowser.cpu.prof ibrowser.mem.prof || true
	time bin/ibrowser -cpuprofile ibrowser.cpu.prof -memprofile ibrowser.mem.prof 360_merged_2.50.vcf.gz

prof_run: clean ibrowser 360_merged_2.50.vcf.gz

check: output.yaml
	grep -v " -" output.yaml
	./check.py output
