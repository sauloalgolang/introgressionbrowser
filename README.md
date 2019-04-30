[![Build Status](https://travis-ci.org/sauloalgolang/introgressionbrowser.svg?branch=master)](https://travis-ci.org/sauloalgolang/introgressionbrowser)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

# ibrowser

attempt of ibrowser

~~~
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

## Compression

~~~
3.8M output_360_merged_2.50.vcf.gz_summary.bin
672K output_360_merged_2.50.vcf.gz_summary.bin.1.gz
563K output_360_merged_2.50.vcf.gz_summary.bin.6.gz
522K output_360_merged_2.50.vcf.gz_summary.bin.9.gz
350K output_360_merged_2.50.vcf.gz_summary.bin.7z
1.1M output_360_merged_2.50.vcf.gz_summary.bin.snappy
~~~