#!/usr/bin/env python3

import os
import sys

from collections import OrderedDict

import numpy as np
import pandas as pd
import pyarrow as pa
import pyarrow.parquet as pq

from yaml import load, dump
try:
    from yaml import CLoader as Loader, CDumper as Dumper
except ImportError:
    from yaml import Loader, Dumper

np.set_printoptions(edgeitems=3)
np.core.arrayprint._line_width = 80

DEBUG = False

def readIbrowserBinary(infile, colNames=None):
    dtI = np.dtype([
        ('hasData', bool),
        ('serial', np.int64),
        ('counterBits', np.int64),
        ('dataLen', np.int64),
        ('sumData', np.uint64)
    ])

    # type RegisterHeader struct {
    #     HasData     bool
    #     Serial      uint64
    #     CounterBits uint64
    #     DataLen     uint64
    #     SumData     uint64
    # }

    fileSize = os.stat(infile).st_size
    fileSizeI = os.stat(infile + ".idx").st_size

    with open(infile + ".idx", 'rb') as fhd:
        d = np.fromfile(fhd, dtype=dtI, count=1)[0]

        counterBits = d["counterBits"]
        dataLen = d["dataLen"]
        dataFmt = None
        dataFmtLen = None

        if counterBits == 16:
            dataFmt = np.uint16
            dataFmtLen = 2
        elif counterBits == 32:
            dataFmt = np.uint32
            dataFmtLen = 4
        elif counterBits == 64:
            dataFmt = np.uint64
            dataFmtLen = 8
        else:
            print("unknown counter bits", counterBits)
            sys.exit(1)

        dataSize = dataLen * dataFmtLen
        registerSize = dataSize
        registerSizeI = 1 + 8 + 8 + 8 + 8

        dt = np.dtype([
            ('data', dataFmt, dataLen)
        ])

        assert (fileSizeI % registerSizeI) == 0, "wrong index file size: fileSize {} % registerSize {} != 0. {}".format(
            fileSizeI, registerSizeI, fileSizeI % registerSizeI)
        assert (fileSize % registerSize) == 0, "wrong data  file size: fileSize {} % registerSize {} != 0. {}".format(
            fileSize, registerSize, fileSize % registerSize)

        numRegistersI = int(fileSizeI / registerSizeI) - 1
        numRegisters = int(fileSize / registerSize) - 1

        assert numRegisters == numRegistersI, "mismatch in number of registers {} != {}".format(
            numRegisters, numRegistersI)

    memmapi = np.memmap(infile + ".idx", dtype=dtI, mode='r')[:-1]
    memmap = np.memmap(infile, dtype=dt, mode='r')[:-1]

    if colNames is None:
        colNames = [i for i in range(numRegisters)]

    # matrixi = pd.DataFrame(OrderedDict([(key, memmapi[key]) for key in [
    #                        "sumData"]]), index=memmapi['serial'], copy=False).T
    # matrix = pd.DataFrame(OrderedDict(
    #     [(colNames[serial], reg['data']) for serial, reg in enumerate(memmap)]), copy=False)

    matrixi = pd.DataFrame(memmapi["sumData"], columns=["sumData"], copy=False).T
    matrix  = pd.DataFrame(memmap["data"], copy=False).T

    matrixi.columns = colNames
    
    return numRegisters, matrixi, matrix


def checkMatrixIndex(matrixi, matrix):
    matrixS = matrix.sum(axis=0).T
    matrixD = matrixS - matrixi
    diff = matrixD.sum(axis=1)[0]

    # print("matrixi\n", matrixi)
    # print("matrixS\n", matrixS)
    # print("matrixD\n", matrixD)
    print("  diff ", diff)
    
    assert diff == 0, "calculated sum differs from sum: {} != 0".format(diff)


def checkZero(zero, others):
    othersS = others.sum(axis=1)
    othersD = othersS - zero
    zeroD = othersD.sum()

    # print("matrix \n", matrix)
    # print("zero   \n", zero)
    # print("others \n", others)
    # print("othersS\n", othersS)
    # print("othersD\n", othersD)
    print("  zeroD", zeroD)
    
    assert zeroD == 0, "calculated zerp sum differs from sum: {} != 0".format(zeroD)


def checkSummary(matrixi, matrix):
    # print("matrixi\n", matrixi)
    # print("matrix\n", matrix)
    checkMatrixIndex(matrixi, matrix)
    zero = matrix.iloc[:, 0]
    others = matrix.iloc[:, 1:]
    checkZero(zero, others)


def checkChromosome(matrixi, matrix, matrixR):
    # print("cMatrixi\n", matrixi)
    # print("cMatrix\n", matrix)
    # print("matrixR\n", matrixR)
    checkMatrixIndex(matrixi, matrix)
    checkZero(matrixR, matrix)


def main(prefix):
    basefile    = prefix + ".yaml"
    outfileh5   = prefix + ".h5"
    outfilepq   = prefix + ".parquet"
    summary     = prefix + "_summary.bin"
    chromosomes = prefix + "_chromosomes_{}.bin"

    print("Loading YAML: ", basefile)
    data = load(open(basefile, 'rt'), Loader=Loader)

    # output_360_merged_2.50.vcf.gz.yaml
    # output_360_merged_2.50.vcf.gz_summary.bin
    # output_360_merged_2.50.vcf.gz_chromosomes_SL2.50ch00.bin

    blockNames = [e['chromosomename'] for e in data['blockmanager']['blocks']]
    # numblocks  = data['blockmanager']['blocks'][numblocks]
    # print('blockNames', blockNames)

    print("Loading summary: ", summary)
    _, matrixi, matrix = readIbrowserBinary(summary, colNames=blockNames)

    print("checking summary")
    # checkSummary(matrixi, matrix)

    chromosomesDfs = OrderedDict()
    cdf = None
    for blockPos, blockName in enumerate(blockNames):
        if blockPos == 0:
            continue

        filename = chromosomes.format(blockName)
        print("Reading chromosome: ", blockName)
        print(" blockName    ", blockName)
        print(" filename     ", filename)

        cNumRegisters, cMatrixi, cMatrix = readIbrowserBinary(filename)
        
        cMatrix.columns = pd.MultiIndex.from_product([[blockName], cMatrix.columns])
        # print('cMatrix', cMatrix)

        print("  cNumRegisters", cNumRegisters)

        numblocks = data['chromosomes'][blockName]['blockmanager']['numblocks']
        print("  numblocks    ", numblocks)

        assert numblocks == cNumRegisters

        if DEBUG:
            print(" Checking chromosome: ", blockName)
            # checkChromosome(cMatrixi, cMatrix, matrix[blockName])
            # chromosomesDfs[blockName] = cMatrix
            pass

        if cdf is None:
            cdf = cMatrix
        else:
            cdf = cdf.merge(cMatrix, left_index=True, right_index=True, validate="one_to_one", copy=False)

    # print(cdf)

    # print("Merging chromosomes")
    # cdf = pd.concat(chromosomesDfs, axis=1, copy=False)
    # print(cdf)
    
    basefileKey = basefile.replace('/', '_').replace('.', '_')
    print("Saving chromosomes to hdf5:", outfileh5)
    cdf.to_hdf(outfileh5, key=basefileKey, mode="w", format="fixed", append=False)#, data_columns=None)

    print("Saving chromosomes to parquet:", outfilepq)
    table = pa.Table.from_pandas(cdf, preserve_index=True, nthreads=4)
    pq.write_table(table, outfilepq, version="2.0")

    # cdf.to_parquet(outfilepq, index=True)

    if DEBUG or True:
        print("Loading chromosomes from hdf5:", outfileh5)
        cdfh5 = pd.read_hdf(outfileh5, key=basefileKey)

        assert cdf.equals(cdfh5), "hdf5 data differ"

        print("Loading chromosomes from parquet:", outfilepq)
        # cdfpq = pd.read_parquet(outfilepq)
        cdfpq = pq.read_table(outfilepq).to_pandas()

        # print(cdf)
        # print(cdfh5)
        # print(cdfpq)
        assert cdf.equals(cdfpq), "parquet data differ"


if __name__ == "__main__":
    # prefix = sys.argv[1]
    # prefix = "res/output_360_merged_2.50.vcf.gz"
    prefix = "res/all/output_360_merged_2.50.vcf.gz"
    main(prefix)
