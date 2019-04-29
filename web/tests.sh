#!/bin/bash

set -xeu

#bin/ibrowser web -vv &

curl http://127.0.0.1:8000
curl http://127.0.0.1:8000/

curl http://127.0.0.1:8000/data
curl http://127.0.0.1:8000/data/

curl http://127.0.0.1:8000/api
curl http://127.0.0.1:8000/api/

curl -X POST http://127.0.0.1:8000/api/update

curl http://127.0.0.1:8000/api/databases

curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz
curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/
curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/summary
curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/summary/matrix
curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/summary/matrix/table

curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/chromosomes

curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/chromosomes/SL2.50ch02
curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/chromosomes/SL2.50ch02/summary
curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/chromosomes/SL2.50ch02/summary/matrix
curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/chromosomes/SL2.50ch02/summary/matrix/table

curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/chromosomes/SL2.50ch02/blocks
curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/chromosomes/SL2.50ch02/blocks/0
curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/chromosomes/SL2.50ch02/blocks/0/matrix
curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/chromosomes/SL2.50ch02/blocks/0/matrix/table
curl http://127.0.0.1:8000/api/databases/output_360_merged_2.50.vcf.gz/chromosomes/SL2.50ch02/blocks/1/matrix/table

curl http://127.0.0.1:8000/api/plots/output_360_merged_2.50.vcf.gz/SL2.50ch02/TS-111
