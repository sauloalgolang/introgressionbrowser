#!/usr/bin/env python3

import numpy as np

fname = "../res/output_360_merged_2.50.vcf.gz_summary.bin"
fname = "../res/output_360_merged_2.50.vcf.gz_chromosomes.bin"

dt = np.dtype([
    ('hasData', bool), 
    ('serial', np.int64),
    ('counterBits', np.int64),
    ('dataLen', np.int64),
    ('sumData', np.uint64)
])

d = np.fromfile(fname, dtype=dt, count=1)

print( "d", d )

dt = np.dtype([
    ('hasData', bool), 
    ('serial', np.int64),
    ('counterBits', np.int64),
    ('dataLen', np.int64),
    ('sumData', np.uint64),
    ('data', np.int32, d["dataLen"])
])

e = np.fromfile(fname, dtype=dt) #, count=2)

# print( "e", e )

for r in e:
    r['data'] = np.cumsum(r['data'])
    f = r['data']
    s = np.sum(f)
    print("r", r)
    assert s == r['sumData']

