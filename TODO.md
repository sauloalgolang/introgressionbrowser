TODO LIST
=========

TODO
----

- [ ] Add ibrowser merger
- [ ] Use logging
- [ ] Let user choose distance matrix to use
- [ ] Self check
- [ ] Implement limits in main function
  - [ ] minSnpPerBlock
  - [ ] maxSnpPerBlock
- [ ] Ibrowser per sample stats (?)
- [ ] Try to write to parquet
  - https://github.com/xitongsys/parquet-go
- [ ] Consider TABIX
  - https://github.com/brentp/bix
  - https://github.com/biogo/hts
  - https://github.com/brentp/cgotabix

DONE
----

- [X] Keep chromosomes ordered
- [X] Save parameters in dump
- [X] Replace as many callbacks as possible
- [X] Use snappy for compression
  - https://github.com/golang/snappy
- [x] Rewrite check.py
- [x] Double check distance matrix is working
- [x] Test YAML Encoder/Decoder using streams
  - https://godoc.org/gopkg.in/yaml.v2#Encoder
- [x] Test using int32 in matrix
  - 64 989
  - 32
    - 1M 543 Mb
    - 360
      - 4.3 Gb RAM
      - 2.1G output.yaml.tar
      - 429M output.yaml.tar.gz
- [X] Test compression
  - 100
    - 194K Apr 11 21:31 output.gob
    -  36K Apr 11 21:32 output.gob.gz
    - 2.1M Apr 11 21:32 output.yaml
    -  87K Apr 11 21:32 output.yaml.gz
  - 1M
    - 233M Apr 11 22:23 output.yaml.tar
    -  41M Apr 11 22:23 output.yaml.tar.gz
    -  62M Apr 11 22:39 output.gob.tar
    -  26M Apr 11 22:39 output.gob.tar.gz
