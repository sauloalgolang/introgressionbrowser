[![Build Status](https://travis-ci.org/sauloalgolang/introgressionbrowser.svg?branch=master)](https://travis-ci.org/sauloalgolang/introgressionbrowser)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

# Introgression Browser

Attempt of recoding Introgression Browser in Golang

## Commands

### Main

~~~bash
$ bin/ibrowser -h
Usage:
  ibrowser [OPTIONS] <command>

Help Options:
  -h, --help  Show this help message

Available commands:
  load     Load database
  save     Read VCF and save database
  version  Print version and exit
  web      Start web interface
~~~

### Version

~~~bash
$ bin/ibrowser version -h
Usage:
  ibrowser [OPTIONS] version

Print version and exit

Help Options:
  -h, --help      Show this help message
~~~

### Save

~~~bash
$ bin/ibrowser save -h
Usage:
  ibrowser [OPTIONS] save [save-OPTIONS] [Input VCF file]

Read VCF and save database

Help Options:
  -h, --help                               Show this help message

[save command options]
          --blockSize=                     Block size (default: 100000)
          --chromosomes=                   Comma separated list of chromomomes to read
          --continueOnError                Continue reading the file on parsing error
          --counterBits=                   Number of bits (default: 32)
          --keepEmptyBlocks                Keep empty blocks
          --maxSnpPerBlock=                Maximum number of SNPs per block (default: 18446744073709551615)
          --minSnpPerBlock=                Minimum number of SNPs per block (default: 10)
          --outfile=                       Output file prefix (default: res/output)
          --description=                   Description of the database
          --cpuProfile=                    Write cpu profile to file
          --memProfile=                    Write memory profile to file
          --check                          Check for self consistency
          --compression=[none|snappy|gzip] Compression format: none, snappy, gzip (default: none)
          --format=[yaml]                  File format: yaml (default: yaml)
          --threads=                       Number of threads (default: 4)
          --debug                          Print debug information
          --debugFirstOnly                 Read only fist chromosome from each thread
          --debugMaxRegisterThread=        Maximum number of registers to read per thread (default: 0)
          --debugMaxRegisterChrom=         Maximum number of registers to read per chromosome (default: 0)

[save command arguments]
  Input VCF file:                          Input VCF file
~~~

### Load

~~~bash
$ bin/ibrowser load -h
Usage:
  ibrowser [OPTIONS] load [load-OPTIONS] [Input Database Prefix]

Load previously created database

Help Options:
  -h, --help                               Show this help message

[load command options]
          --softLoad                       Do not load matrices
          --cpuProfile=                    Write cpu profile to file
          --memProfile=                    Write memory profile to file
          --check                          Check for self consistency
          --compression=[none|snappy|gzip] Compression format: none, snappy, gzip (default: none)
          --format=[yaml]                  File format: yaml (default: yaml)
          --threads=                       Number of threads (default: 4)

[load command arguments]
  Input Database Prefix:                   Input database prefix
~~~

### Web

~~~bash
$ bin/ibrowser web -h
Usage:
  ibrowser [OPTIONS] web [web-OPTIONS]

Start web interface

Help Options:
  -h, --help             Show this help message

[web command options]
          --host=        Hostname (default: 127.0.0.1)
          --port=        Port (default: 8000)
          --DatabaseDir= Databases folder (default: res/)
          --HttpDir=     Web page folder to be served folder (default: http/)
      -v, --verbose      Show verbose debug information
~~~

# Make

~~~bash
$ make
GOROOT  /home/saulo/bin/go
FORMAT  yaml
OUTFILE res/output
GIT_COMMIT_HASH     20fbd6b6b195622ac781dcb986443b774e941c7b
GIT_COMMIT_AUTHOR   Saulo Aflitos (sauloal@gmail.com) 2019-04-30 16:19:43 +0200
GIT_COMMIT_COMMITER Saulo Aflitos (sauloal@gmail.com) 2019-04-30 16:19:43 +0200
GIT_COMMIT_NOTES
GIT_COMMIT_TITLE    major rewrite. moved matrix to ibrowser; refactor vcf; refactor interfaces; move dumper to matrix
GIT_STATUS          dirty
GIT_DIFF            1060965a47890b7c6f14b47c7406ac17
GO_VERSION          go version go1.12.1 linux/amd64

 build
  ibrowser
  ibrowser.exe
  ibrowser.wasm

 httpserver

 examples
 serve

 get
 requirements
 wasm_exec.js

 run150
 run360

 prof
 prof_run
 version
 check
 clean

 test
  test_save
  test_load
~~~

# Tests

## Reading test

~~~bash
VcfRaw w/o bufreader - w/o register
1500000 SL2.50ch02 6627975
76.03user

VcfRaw w/ bufread - w/o register
1500000 SL2.50ch02 6627975
75.39user

VcfRaw w/ bufread - w/ register global
1400000 SL2.50ch01 97265751
585.53user

VcfRaw w/ bufread - w/ register all
1400000 SL2.50ch01 97265751
1312.73

VcfGo w/ bufread
1700000 SL2.50ch02 23977842
1746.71user
VcfGo w/o bufread
~~~

## Compression Test

~~~bash
3.8M output_360_merged_2.50.vcf.gz_summary.bin
672K output_360_merged_2.50.vcf.gz_summary.bin.1.gz
563K output_360_merged_2.50.vcf.gz_summary.bin.6.gz
522K output_360_merged_2.50.vcf.gz_summary.bin.9.gz
350K output_360_merged_2.50.vcf.gz_summary.bin.7z
1.1M output_360_merged_2.50.vcf.gz_summary.bin.snappy
~~~