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

curl http://127.0.0.1:8000/api/database

curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz
curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/
curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/summary
curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/summary/matrix
curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/summary/matrix/table

curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/chromosome

curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/chromosome/SL2.50ch02
curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/chromosome/SL2.50ch02/summary
curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/chromosome/SL2.50ch02/summary/matrix
curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/chromosome/SL2.50ch02/summary/matrix/table

curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/chromosome/SL2.50ch02/block/0
curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/chromosome/SL2.50ch02/block/0/matrix
curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/chromosome/SL2.50ch02/block/0/matrix/table
curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/chromosome/SL2.50ch02/block/1/matrix/table

#curl http://127.0.0.1:8000/api/database/output_360_merged_2.50.vcf.gz/chromosome/SL2.50ch02/block/2/matrix/table
