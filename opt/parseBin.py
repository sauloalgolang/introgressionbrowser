#!/usr/bin/env python3

import os
import sys

import numpy as np
import pandas as pd

np.set_printoptions(edgeitems=3)
np.core.arrayprint._line_width = 80

fname = "res/output_360_merged_2.50.vcf.gz_summary.bin"
fname = "res/output_360_merged_2.50.vcf.gz_chromosomes.bin"

def readIbrowserBinary(infile):
    dt0 = np.dtype([
        ('hasData', bool), 
        ('serial', np.int64),
        ('counterBits', np.int64),
        ('dataLen', np.int64),
        ('sumData', np.uint64)
    ])

    fileSize = os.stat(infile).st_size
    
    with open(infile, 'rb') as fhd:
        d = np.fromfile(fhd, dtype=dt0, count=1)[0]

        counterBits = d["counterBits"]
        dataLen = d["dataLen"]

        dataFmt = None
        dataFmtLen = None
        
        if counterBits == 16:
            dataFmt = np.int16
            dataFmtLen = 2
        elif counterBits == 32:
            dataFmt = np.int32
            dataFmtLen = 4
        elif counterBits == 64:
            dataFmt = np.int64
            dataFmtLen = 8
        else:
            print("unknown counter bits", counterBits)
            sys.exit(1)

        dataSize     = dataLen * dataFmtLen
        registerSize = 1 + 8 + 8 + 8 + 8 + dataSize

        dt = np.dtype([
            ('hasData'    , bool     ), #1  1
            ('serial'     , np.int64 ), #8  9
            ('counterBits', np.int64 ), #8 17
            ('dataLen'    , np.int64 ), #8 25
            ('sumData'    , np.uint64), #8 33
            ('data', dataFmt, dataLen)
        ])

        assert (fileSize % registerSize) == 0, "wrong file size: fileSize {} % registerSize {} != 0. {}".format(fileSize, registerSize, fileSize % registerSize)
        
        numRegisters = int(fileSize / registerSize) - 1
    
    memmap = np.memmap(infile, dtype=dt, mode='r')

    return numRegisters, memmap


def registerToDataframe(regs):
    matrix = None

    for reg in regs:
        if matrix is None:
            matrix = pd.DataFrame({reg['serial']: reg['data']})
        else:
            matrix[reg['serial']] = reg['data']

    return matrix


if __name__ == "__main__":
    print("memmap")

    numRegisters, memmap = readIbrowserBinary(fname)

    for i in range(numRegisters):
        # reg = ibb.getRegister(i)
        reg = memmap[i]
        print(str(reg)[:50])
        print("  ", reg['data'][:10])

    df = registerToDataframe(memmap[1:5])

    print(df)

    pd.DataFrame(memmap)

    