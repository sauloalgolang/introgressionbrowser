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

def checkBlockOfBlocks(block, blocks, checkPos=True):
    minposition = block["minposition"]
    maxposition = block["maxposition"]
    numsnps = block["numsnps"]

    if checkPos:
        minpositionCalc = min([b["minposition"] for b in blocks])
        assert minposition == minpositionCalc, "minposition mismatch: {} != {}".format(minposition, minpositionCalc)
        print(" min position OK", minposition)

        maxpositionCalc = max([b["maxposition"] for b in blocks])
        assert maxposition == maxpositionCalc, "maxposition mismatch: {} != {}".format(maxposition, maxpositionCalc)
        print(" max position OK", maxposition)

    numsnpsCalc = sum([b["numsnps"] for b in blocks])
    assert numsnps == numsnpsCalc, "numsnps mismatch: {} != {}".format(numsnps, numsnpsCalc)
    print(" num snps OK", numsnps)

    matrix = block["matrix"]["data"]
    sumBlock = sum(matrix)
    sumBlocks = sum([sum(b["matrix"]["data"]) for b in blocks])

    assert sumBlock == sumBlocks, "sumBlock mismatch: {} != {}".format(sumBlock, sumBlocks)

    for p, blockP in enumerate(matrix):
        blocksP = sum([b["matrix"]["data"][p] for b in blocks])
        assert blockP == blocksP, "sumBlock p {} mismatch: {} != {}".format(p, blockP, blocksP)

    return True

def main(prefix):
    basefile = prefix + ".yaml"

    print('reading', basefile)

    with open(basefile, 'rt') as fhd:
        data = load(fhd, Loader=Loader)

        # print(data)

        chromosomesnames = data["chromosomesnames"]

        print(chromosomesnames)

        samples = data["samples"]
        numsamples = data["numsamples"]
        blocksize = data["blocksize"]
        keepemptyblock = data["keepemptyblock"]
        numregisters = data["numregisters"]
        numsnps = data["numsnps"]
        numblocks = data["numblocks"]
        chromosomes = data["chromosomes"]
        block = data["block"]

        numsamplesCalc = len(samples)
        assert numsamples == numsamplesCalc, "numsamples mismatch: {} != {}".format(numsamples, numsamplesCalc)
        print(" num samples OK", numsamples)

        numsnpsCalc = sum([c["numsnps"] for c in chromosomes.values()])
        assert numsnps == numsnpsCalc, "numsnps mismatch: {} != {}".format(numsnps, numsnpsCalc)
        print(" num snps OK", numsnps)

        numblocksCalc = sum([c["numblocks"] for c in chromosomes.values()])
        assert numblocks == numblocksCalc, "numblocks mismatch: {} != {}".format(numblocks, numblocksCalc)
        print(" num blocks OK", numblocks)

        assert checkBlockOfBlocks(block, [b["block"] for b in chromosomes.values()], checkPos=False)
        print(" block OK")

        for chromosomeName, chromosome in chromosomes.items():
            print(chromosomeName)
            
            chromosomeNameC = chromosome["chromosome"]
            
            assert chromosomeName == chromosomeNameC
            
            minpositionC = chromosome["minposition"]
            maxpositionC = chromosome["maxposition"]
            numblocksC = chromosome["numblocks"]
            numsnpsC = chromosome["numsnps"]
            numsamplesC = chromosome["numsamples"]
            blockC = chromosome["block"]
            blocksC = chromosome["blocks"]

            minpositionCCalc = min([b["minposition"] for b in blocksC])
            assert minpositionC == minpositionCCalc, " minposition mismatch: {} != {}".format(minpositionC, minpositionCCalc)
            print(" min position OK", minpositionC)

            maxpositionCCalc = max([b["maxposition"] for b in blocksC])
            assert maxpositionC == maxpositionCCalc, " maxposition mismatch: {} != {}".format(maxpositionC, maxpositionCCalc)
            print(" max position OK", maxpositionC)

            numblocksCCalc = len(blocksC)
            assert numblocksC == numblocksCCalc, " numblocks mismatch: {} != {}".format(numblocksC, numblocksCCalc)
            print(" num blocks OK", numblocksC)

            numsnpsCCalc = sum([b["numsnps"] for b in blocksC])
            assert numsnpsC == numsnpsCCalc, " numsnps mismatch: {} != {}".format(numsnpsC, numsnpsCCalc)
            print(" num snps OK", numsnpsC)

            assert numsamples == numsamplesC, " numsamples mismatch: {} != {}".format(numsamples, numsamplesC)
            print(" num samples OK", numsamplesC)

            assert checkBlockOfBlocks(blockC, [b for b in blocksC])
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
    print("ALL OK")



if __name__ == "__main__":
    main(sys.argv[1])