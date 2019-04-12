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

def checkBlockOfBlocks(prefix, chromosomeName, blocknames, block, blockMatrix, checkPos=True):
  minposition = block["minposition"]
  maxposition = block["maxposition"]
  numsnps = block["numsnps"]

  blocks = blocksKeys(prefix, chromosomeName, blocknames, ["minposition", "maxposition", "numsnps"])

  if checkPos:
    minpositionCalc = min(blocks["minposition"])
    assert minposition == minpositionCalc, "minposition mismatch: {} != {}".format(minposition, minpositionCalc)
    print(" min position OK", minposition)

    maxpositionCalc = max(blocks["maxposition"])
    assert maxposition == maxpositionCalc, "maxposition mismatch: {} != {}".format(maxposition, maxpositionCalc)
    print(" max position OK", maxposition)

  numsnpsCalc = sum(blocks["numsnps"])
  assert numsnps == numsnpsCalc, "numsnps mismatch: {} != {}".format(numsnps, numsnpsCalc)
  print(" num snps OK", numsnps)

  blockMatrixC = sumMatrices(prefix, chromosomeName, blocknames)
  matrixDiff = subtractMatrices(blockMatrix["data"], blockMatrixC)
  assert sum(matrixDiff) == 0
  print(" matrix OK")

  return True

def readChromosomes(prefix, chromosomesnames):
  return {c: readChromosome(prefix, c) for c in chromosomesnames}

def readChromosome(prefix, chromosomeName):
  basefile = prefix + "." + chromosomeName + ".yaml"
  fhd = open(basefile, 'rt')
  data = load(fhd, Loader=Loader)
  fhd.close()
  return data

def readBlockIB(prefix):
  return readBlock(prefix, 0, False)

def readBlockMatrixIB(prefix):
  return readBlockMatrix(prefix, 0, False)

def readBlockChrom(prefix, chromosomeName):
  return readBlock(prefix + "." + chromosomeName, 0, False)

def readBlockMatrixChrom(prefix, chromosomeName):
  return readBlockMatrix(prefix + "." + chromosomeName, 0, False)

def readBlocksChrom(prefix, chromosomeName, pos):
  return readBlock(prefix + "." + chromosomeName, pos, True)

def readBlocksMatrixChrom(prefix, chromosomeName, pos):
  return readBlockMatrix(prefix + "." + chromosomeName, pos, True)

def readBlock(prefix, pos, isblocks=False):
  basefile = prefix + "_block{}.{:012d}.yaml".format("s" if isblocks else "", pos)
  fhd = open(basefile, 'rt')
  data = load(fhd, Loader=Loader)
  fhd.close()
  return data

def readBlockMatrix(prefix, pos, isblocks=False):
  basefile = prefix + "_block{}.{:012d}_matrix.yaml".format("s" if isblocks else "",pos)
  fhd = open(basefile, 'rt')
  data = load(fhd, Loader=Loader)
  fhd.close()
  return data

def blocksKey(prefix, chromosomeName, blocknames, key):
  return [ readBlocksChrom(prefix, chromosomeName, pos)[key] for pos in blocknames.keys() ]

def blocksKeys(prefix, chromosomeName, blocknames, keys):
  res = {}
  for pos in blocknames.keys():
    block = readBlocksChrom(prefix, chromosomeName, pos)
    for key in keys:
      l = res.get(key, [])
      l.append(block[key])
      res[key] = l
  return res

def subtractMatrices(blockMatrix, blockMatrixC):
  for i in range(len(blockMatrix)):
    blockMatrix[i] -= blockMatrixC[i]
  return blockMatrix

def sumMatrices(prefix, chromosomeName, blocknames):
  res = None
  for pos in blocknames.keys():
    matrix = readBlocksMatrixChrom(prefix, chromosomeName, pos)
    if res is None:
      res = matrix["data"]
    else:
      for i in range(len(res)):
        res[i] += matrix["data"][i]
  return res


def main(prefix):
    basefile = prefix + ".yaml"

    print('reading', basefile)

    with open(basefile, 'rt') as fhd:
        data = load(fhd, Loader=Loader)

        # print(data)

        chromosomesnames = data["chromosomesnames"]

        print("chromosomesnames", chromosomesnames)

        samples = data["samples"]
        numsamples = data["numsamples"]
        blocksize = data["blocksize"]
        keepemptyblock = data["keepemptyblock"]
        numregisters = data["numregisters"]
        numsnps = data["numsnps"]
        numblocks = data["numblocks"]

        chromosomes = readChromosomes(prefix, chromosomesnames)

        numsamplesCalc = len(samples)
        assert numsamples == numsamplesCalc, "numsamples mismatch: {} != {}".format(numsamples, numsamplesCalc)
        print(" num samples OK", numsamples)

        numsnpsCalc = sum([c["numsnps"] for c in chromosomes.values()])
        assert numsnps == numsnpsCalc, "numsnps mismatch: {} != {}".format(numsnps, numsnpsCalc)
        print(" num snps OK", numsnps)

        numblocksCalc = sum([c["numblocks"] for c in chromosomes.values()])
        assert numblocks == numblocksCalc, "numblocks mismatch: {} != {}".format(numblocks, numblocksCalc)
        print(" num blocks OK", numblocks)

        block = readBlockIB(prefix)
        blockMatrix = readBlockMatrixIB(prefix)

        for chromosomeName, chromosome in chromosomes.items():
            print(chromosomeName)
            # print(chromosome)
            
            chromosomeNameC = chromosome["chromosomename"]
            
            assert chromosomeName == chromosomeNameC
            
            minpositionC = chromosome["minposition"]
            maxpositionC = chromosome["maxposition"]
            numblocksC = chromosome["numblocks"]
            numsnpsC = chromosome["numsnps"]
            numsamplesC = chromosome["numsamples"]

            blockC = readBlockChrom(prefix, chromosomeName)
            blockMatrixC = readBlockMatrixChrom(prefix, chromosomeName)
            
            blocknames = chromosome["blocknames"]

            block["numsnps"] -= blockC["numsnps"]
            blockMatrix["data"] = [blockMatrix["data"][p] - blockMatrixC["data"][p] for p in range(len(blockMatrix["data"]))]
            numsnps -= numsnpsC
            numblocks -= len(blocknames)

            minpositionCCalc = min(blocksKey(prefix, chromosomeName, blocknames, "minposition"))
            assert minpositionC == minpositionCCalc, " minposition mismatch: {} != {}".format(minpositionC, minpositionCCalc)
            print(" min position OK", minpositionC)

            maxpositionCCalc = max(blocksKey(prefix, chromosomeName, blocknames, "maxposition"))
            assert maxpositionC == maxpositionCCalc, " maxposition mismatch: {} != {}".format(maxpositionC, maxpositionCCalc)
            print(" max position OK", maxpositionC)

            numblocksCCalc = len(blocknames)
            assert numblocksC == numblocksCCalc, " numblocks mismatch: {} != {}".format(numblocksC, numblocksCCalc)
            print(" num blocks OK", numblocksC)

            numsnpsCCalc = sum(blocksKey(prefix, chromosomeName, blocknames, "numsnps"))
            assert numsnpsC == numsnpsCCalc, " numsnps mismatch: {} != {}".format(numsnpsC, numsnpsCCalc)
            print(" num snps OK", numsnpsC)

            assert numsamples == numsamplesC, " numsamples mismatch: {} != {}".format(numsamples, numsamplesC)
            print(" num samples OK", numsamplesC)

            assert checkBlockOfBlocks(prefix, chromosomeName, blocknames, blockC, blockMatrixC)
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
    
    assert numsnps == 0
    assert block["numsnps"] == 0
    print("numsnps OK")

    assert numblocks == 0
    print("numblocks OK")

    assert sum(blockMatrix["data"]) == 0
    print("matrix OK")
    
    print("ALL OK")



if __name__ == "__main__":
    main(sys.argv[1])