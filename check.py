#!/usr/bin/env python3

import sys
import math

from yaml import load, dump
try:
    from yaml import CLoader as Loader, CDumper as Dumper
except ImportError:
    from yaml import Loader, Dumper

"""
samples:
- TS-1
- TS-2
- TS-100
- TS-101
numsamples: 361
blocksize: 1000000
keepemptyblock: true
numregisters: 1226
numsnps: 1226
numblocks: 13
chromosomesnames:
block:
  blocknumber: 0
  minposition: 0
  maxposition: 0
  numsnps: 1226
  numsamples: 361
  matrix:
  - 14
  - 1
chromosomes:
  SL2.50ch00:
    chromosome: SL2.50ch00
    minposition: 0
    maxposition: 8853
    numblocks: 1
    numsnps: 98
    numsamples: 361
    keepemptyblock: true
    block:
      blocknumber: 0
      minposition: 0
      maxposition: 8853
      numsnps: 98
      numsamples: 361
      matrix:
        - 14
        - 1
    blocks:
    - blocknumber: 0
      minposition: 0
      maxposition: 8853
      numsnps: 98
      numsamples: 361
      matrix:
        - 14
        - 1
    blocknames:
      0: 0
"""

def checkBlockOfBlocks(block, blocks, dataKey, checkPos=True):
    minposition = block["minposition"]
    maxposition = block["maxposition"]
    numsnps     = block["numsnps"]

    blocksminposition = [b["minposition"] for b in blocks]
    blocksmaxposition = [b["maxposition"] for b in blocks]
    blocksnumsnps     = [b["numsnps"    ] for b in blocks]

    numsnpsCalc = sum(blocksnumsnps)
    assert numsnps == numsnpsCalc, "numsnps mismatch: {} != {}".format(numsnps, numsnpsCalc)
    print(" num snps OK", numsnps)

    if checkPos:
        minpositionCalc = min(blocksminposition)
        assert minposition == minpositionCalc, "minposition mismatch: {} != {}".format(minposition, minpositionCalc)
        print(" min position OK", minposition)

        maxpositionCalc = max(blocksmaxposition)
        assert maxposition == maxpositionCalc, "maxposition mismatch: {} != {}".format(maxposition, maxpositionCalc)
        print(" max position OK", maxposition)

    # print(" checking matrix")
    # blockMatrixC = sumMatrices(prefix, chromosomeName, blocknames)
    # matrixDiff = subtractMatrices(blockMatrix[dataKey], blockMatrixC)
    # assert sum(matrixDiff) == 0
    # print("  matrix OK")

    return True

    # bits = data["numbits"]
    # dataKey = None

    # if bits == 16:
    #     dataKey = "data16"
    # elif bits == 32:
    #     dataKey = "data32"
    # elif bits == 64:
    #     dataKey = "data64"
    # else:
    #     sys.stderr.write("unknown bits {}\n".format(bits))
    #     sys.exit(1)

def subtractMatrices(blockMatrix, blockMatrixC):
    assert isinstance(blockMatrix, list), "matrix is not an list: {} {}".format(isinstance(blockMatrix, list), blockMatrix) 
    assert isinstance(blockMatrixC, list), "matrix is not an list: {} {}".format(isinstance(blockMatrixC, list), blockMatrixC)
    assert len(blockMatrix) == len(blockMatrixC), "matrices have different lenghts: {} {}".format(len(blockMatrix), len(blockMatrixC))
    
    for i in range(len(blockMatrix)):
        blockMatrix[i] -= blockMatrixC[i]
  
    return blockMatrix

def sumMatrices(prefix, chromosomeName, blocknames):
    res = None
    reg = 0

    for pos in blocknames.keys():
        matrix, dataKey = readBlocksMatrixChrom(prefix, chromosomeName, pos)

        reg += 1
        print(".", end="")
        
        if reg % 100 == 0:
            print("{:8d}".format(reg))
            sys.stdout.flush()

        if res is None:
            res = matrix[dataKey]
            # print(res)
        
        else:
            for i in range(len(res)):
                res[i] += matrix[dataKey][i]

    print("{:8d}".format(reg))
    print()
    
    return res


def main(prefix):
    basefile = prefix + ".yaml"

    print('reading', basefile)

    dataKey = None

    with open(basefile, 'rt') as fhd:
        data = load(fhd, Loader=Loader)

        # print(data)

        chromosomesnames = data["chromosomesnames"]

        print("chromosomesnames", chromosomesnames)

        samples = data["samples"]
        numsamples = data["numsamples"]
        #blocksize = data["blocksize"]
        #keepemptyblock = data["keepemptyblock"]
        #numregisters = data["numregisters"]
        numsnps = data["numsnps"]
        numblocks = data["numblocks"]

        chromosomes = data["chromosomes"]

        numsamplesCalc = len(samples)
        assert numsamples == numsamplesCalc, "numsamples mismatch: {} != {}".format(numsamples, numsamplesCalc)
        print(" num samples OK", numsamples)

        numsnpsCalc = sum([c["numsnps"] for c in chromosomes.values()])
        assert numsnps == numsnpsCalc, "numsnps mismatch: {} != {}".format(numsnps, numsnpsCalc)
        print(" num snps OK", numsnps)

        numblocksCalc = sum([c["blockmanager"]["numblocks"] for c in chromosomes.values()])
        assert numblocks == numblocksCalc, "numblocks mismatch: {} != {}".format(numblocks, numblocksCalc)
        print(" num blocks OK", numblocks)

        blockmanager = data["blockmanager"]
        blockIndex   = blockmanager["blocknames"]["_whole_genome"]
        block        = blockmanager["blocks"][blockIndex]
        blocks = []

        sumBlocks = numblocks
        sumBlocksSnps = block["numsnps"]
        sumSnps = numsnps

        for chromosomeName, chromosome in chromosomes.items():
            print(chromosomeName)
            # print(chromosome)
            
            chromosomeNameC = chromosome["chromosomename"]
            
            assert chromosomeName == chromosomeNameC
            
            minpositionC  = chromosome["minposition"]
            maxpositionC  = chromosome["maxposition"]
            numsnpsC      = chromosome["numsnps"]
            numsamplesC   = chromosome["numsamples"]

            blockIndexC   = blockmanager["blocknames"][chromosomeName]
            blockC        = blockmanager["blocks"][blockIndexC]
            
            blockmanagerC = chromosome["blockmanager"]
            numblocksC    = blockmanagerC["numblocks"]

            blocksC       = blockmanagerC["blocks"]
            blocks.append(blockC)
            
            blocknames    = blockmanagerC["blocknames"]
            blocknumbers  = blockmanagerC["blocknumbers"]

            sumBlocksSnps -= blockC["numsnps"]
            print("  numsnps SUB", sumBlocksSnps)
            
            # blockMatrix[dataKey] = [blockMatrix[dataKey][p] - blockMatrixC[dataKey][p] for p in range(len(blockMatrix[dataKey]))]
            
            numsnps   -= numsnpsC
            numblocks -= len(blocknumbers)
            
            minpositionCCalc = min([c["minposition"] for c in blocksC])
            assert minpositionC == minpositionCCalc, " minposition mismatch: {} != {} - {}".format(minpositionC, minpositionCCalc, [c["minposition"] for c in blocksC])
            print(" min position OK", minpositionC)

            maxpositionCCalc = max([c["maxposition"] for c in blocksC])
            assert maxpositionC == maxpositionCCalc, " maxposition mismatch: {} != {}".format(maxpositionC, maxpositionCCalc)
            print(" max position OK", maxpositionC)

            numblocksCCalc = len(blocknumbers)
            assert numblocksC == numblocksCCalc, " numblocks mismatch: {} != {}".format(numblocksC, numblocksCCalc)
            print(" num blocks OK", numblocksC)

            numsnpsCCalc = sum([c["numsnps"] for c in blocksC])
            assert numsnpsC == numsnpsCCalc, " numsnps mismatch: {} != {}".format(numsnpsC, numsnpsCCalc)
            print(" num snps OK", numsnpsC)

            assert numsamples == numsamplesC, " numsamples mismatch: {} != {}".format(numsamples, numsamplesC)
            print(" num samples OK", numsamplesC)

            assert checkBlockOfBlocks(blockC, blocksC, dataKey)
            print(" block OK")

            # block:
            #     blocknumber: 0
            #     minposition: 0
            #     maxposition: 8853
            #     numsnps: 98
            #     numsamples: 361
            #     matrix:
            #         - 14
            #         - 1
            # blocks:
            #     blocknumber: 0
            #     minposition: 0
            #     maxposition: 8853
            #     numsnps: 98
            #     numsamples: 361
            #     matrix:
            #         - 14
            #         - 1
    
    assert numblocks == 0, "number of blocks did not lower to zero: {} - {}".format(numblocks, sumBlocks)
    print("numblocks OK", sumBlocks)

    assert numsnps == 0, "number of SNPS did not lower to zero: {} - {}".format(numsnps, sumSnps)
    print("numsnps OK", sumSnps)

    assert sumBlocksSnps == 0, "number of block numsnps did not lower to zero: {}".format(sumBlocksSnps)
    print("block numsnps OK", sumBlocksSnps)

    assert checkBlockOfBlocks(block, blocks, dataKey, checkPos=False)
    print(" block OK")

    # assert sum(blockMatrix[dataKey]) == 0, "matrix sum did not lower to zero: {} - {}".format(numblocks, sum(blockMatrix[dataKey]))
    # print("matrix OK")
    
    print("ALL OK")

if __name__ == "__main__":
    main(sys.argv[1])
