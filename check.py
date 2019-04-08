#!/usr/bin/env python3

import sys

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

"""
"""
$ ./tri.py
dimension 8
[ 0,  0,  0][ 0,  1,  1][ 0,  2,  2][ 0,  3,  3][ 0,  4,  4][ 0,  5,  5][ 0,  6,  6][ 0,  7,  7]
[ 1,  0,  8][ 1,  1,  9][ 1,  2, 10][ 1,  3, 11][ 1,  4, 12][ 1,  5, 13][ 1,  6, 14][ 1,  7, 15]
[ 2,  0, 16][ 2,  1, 17][ 2,  2, 18][ 2,  3, 19][ 2,  4, 20][ 2,  5, 21][ 2,  6, 22][ 2,  7, 23]
[ 3,  0, 24][ 3,  1, 25][ 3,  2, 26][ 3,  3, 27][ 3,  4, 28][ 3,  5, 29][ 3,  6, 30][ 3,  7, 31]
[ 4,  0, 32][ 4,  1, 33][ 4,  2, 34][ 4,  3, 35][ 4,  4, 36][ 4,  5, 37][ 4,  6, 38][ 4,  7, 39]
[ 5,  0, 40][ 5,  1, 41][ 5,  2, 42][ 5,  3, 43][ 5,  4, 44][ 5,  5, 45][ 5,  6, 46][ 5,  7, 47]
[ 6,  0, 48][ 6,  1, 49][ 6,  2, 50][ 6,  3, 51][ 6,  4, 52][ 6,  5, 53][ 6,  6, 54][ 6,  7, 55]
[ 7,  0, 56][ 7,  1, 57][ 7,  2, 58][ 7,  3, 59][ 7,  4, 60][ 7,  5, 61][ 7,  6, 62][ 7,  7, 63]

[  ,   ,   ][ 0,  1,  0][ 0,  2,  1][ 0,  3,  2][ 0,  4,  3][ 0,  5,  4][ 0,  6,  5][ 0,  7,  6]
[  ,   ,   ][  ,   ,   ][ 1,  2,  7][ 1,  3,  8][ 1,  4,  9][ 1,  5, 10][ 1,  6, 11][ 1,  7, 12]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][ 2,  3, 13][ 2,  4, 14][ 2,  5, 15][ 2,  6, 16][ 2,  7, 17]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 3,  4, 18][ 3,  5, 19][ 3,  6, 20][ 3,  7, 21]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 4,  5, 22][ 4,  6, 23][ 4,  7, 24]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 5,  6, 25][ 5,  7, 26]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 6,  7, 27]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ]
triSize 28
[(0, (0, 1)), (1, (0, 2)), (2, (0, 3)), (3, (0, 4)), (4, (0, 5)), (5, (0, 6)), (6, (0, 7)), (7, (1, 2)), (8, (1, 3)), (9, (1, 4)), (10, (1, 5)), (11, (1, 6)), (12, (1, 7)), (13, (2, 3)), (14, (2, 4)), (15, (2, 5)), (16, (2, 6)), (17, (2, 7)), (18, (3, 4)), (19, (3, 5)), (20, (3, 6)), (21, (3, 7)), (22, (4, 5)), (23, (4, 6)), (24, (4, 7)), (25, (5, 6)), (26, (5, 7)), (27, (6, 7))]
[  ,   ,   ][ 0,  1,  0][ 0,  2,  1][ 0,  3,  2][ 0,  4,  3][ 0,  5,  4][ 0,  6,  5][ 0,  7,  6]
[  ,   ,   ][  ,   ,   ][ 1,  2,  7][ 1,  3,  8][ 1,  4,  9][ 1,  5, 10][ 1,  6, 11][ 1,  7, 12]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][ 2,  3, 13][ 2,  4, 14][ 2,  5, 15][ 2,  6, 16][ 2,  7, 17]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 3,  4, 18][ 3,  5, 19][ 3,  6, 20][ 3,  7, 21]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 4,  5, 22][ 4,  6, 23][ 4,  7, 24]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 5,  6, 25][ 5,  7, 26]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 6,  7, 27]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ]

$ ./tri.py
dimension 9
[ 0,  0,  0][ 0,  1,  1][ 0,  2,  2][ 0,  3,  3][ 0,  4,  4][ 0,  5,  5][ 0,  6,  6][ 0,  7,  7][ 0,  8,  8]
[ 1,  0,  9][ 1,  1, 10][ 1,  2, 11][ 1,  3, 12][ 1,  4, 13][ 1,  5, 14][ 1,  6, 15][ 1,  7, 16][ 1,  8, 17]
[ 2,  0, 18][ 2,  1, 19][ 2,  2, 20][ 2,  3, 21][ 2,  4, 22][ 2,  5, 23][ 2,  6, 24][ 2,  7, 25][ 2,  8, 26]
[ 3,  0, 27][ 3,  1, 28][ 3,  2, 29][ 3,  3, 30][ 3,  4, 31][ 3,  5, 32][ 3,  6, 33][ 3,  7, 34][ 3,  8, 35]
[ 4,  0, 36][ 4,  1, 37][ 4,  2, 38][ 4,  3, 39][ 4,  4, 40][ 4,  5, 41][ 4,  6, 42][ 4,  7, 43][ 4,  8, 44]
[ 5,  0, 45][ 5,  1, 46][ 5,  2, 47][ 5,  3, 48][ 5,  4, 49][ 5,  5, 50][ 5,  6, 51][ 5,  7, 52][ 5,  8, 53]
[ 6,  0, 54][ 6,  1, 55][ 6,  2, 56][ 6,  3, 57][ 6,  4, 58][ 6,  5, 59][ 6,  6, 60][ 6,  7, 61][ 6,  8, 62]
[ 7,  0, 63][ 7,  1, 64][ 7,  2, 65][ 7,  3, 66][ 7,  4, 67][ 7,  5, 68][ 7,  6, 69][ 7,  7, 70][ 7,  8, 71]
[ 8,  0, 72][ 8,  1, 73][ 8,  2, 74][ 8,  3, 75][ 8,  4, 76][ 8,  5, 77][ 8,  6, 78][ 8,  7, 79][ 8,  8, 80]

[  ,   ,   ][ 0,  1,  0][ 0,  2,  1][ 0,  3,  2][ 0,  4,  3][ 0,  5,  4][ 0,  6,  5][ 0,  7,  6][ 0,  8,  7]
[  ,   ,   ][  ,   ,   ][ 1,  2,  8][ 1,  3,  9][ 1,  4, 10][ 1,  5, 11][ 1,  6, 12][ 1,  7, 13][ 1,  8, 14]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][ 2,  3, 15][ 2,  4, 16][ 2,  5, 17][ 2,  6, 18][ 2,  7, 19][ 2,  8, 20]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 3,  4, 21][ 3,  5, 22][ 3,  6, 23][ 3,  7, 24][ 3,  8, 25]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 4,  5, 26][ 4,  6, 27][ 4,  7, 28][ 4,  8, 29]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 5,  6, 30][ 5,  7, 31][ 5,  8, 32]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 6,  7, 33][ 6,  8, 34]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 7,  8, 35]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ]
triSize 36
[(0, (0, 1)), (1, (0, 2)), (2, (0, 3)), (3, (0, 4)), (4, (0, 5)), (5, (0, 6)), (6, (0, 7)), (7, (0, 8)), (8, (1, 2)), (9, (1, 3)), (10, (1, 4)), (11, (1, 5)), (12, (1, 6)), (13, (1, 7)), (14, (1, 8)), (15, (2, 3)), (16, (2, 4)), (17, (2, 5)), (18, (2, 6)), (19, (2, 7)), (20, (2, 8)), (21, (3, 4)), (22, (3, 5)), (23, (3, 6)), (24, (3, 7)), (25, (3, 8)), (26, (4, 5)), (27, (4, 6)), (28, (4, 7)), (29, (4, 8)), (30, (5, 6)), (31, (5, 7)), (32, (5, 8)), (33, (6, 7)), (34, (6, 8)), (35, (7, 8))]
[  ,   ,   ][ 0,  1,  0][ 0,  2,  1][ 0,  3,  2][ 0,  4,  3][ 0,  5,  4][ 0,  6,  5][ 0,  7,  6][ 0,  8,  7]
[  ,   ,   ][  ,   ,   ][ 1,  2,  8][ 1,  3,  9][ 1,  4, 10][ 1,  5, 11][ 1,  6, 12][ 1,  7, 13][ 1,  8, 14]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][ 2,  3, 15][ 2,  4, 16][ 2,  5, 17][ 2,  6, 18][ 2,  7, 19][ 2,  8, 20]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 3,  4, 21][ 3,  5, 22][ 3,  6, 23][ 3,  7, 24][ 3,  8, 25]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 4,  5, 26][ 4,  6, 27][ 4,  7, 28][ 4,  8, 29]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 5,  6, 30][ 5,  7, 31][ 5,  8, 32]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 6,  7, 33][ 6,  8, 34]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][ 7,  8, 35]
[  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ][  ,   ,   ]
"""

# https://stackoverflow.com/questions/27086195/linear-index-upper-triangular-matrix
def triPos(i,j,n):
	k = (n*(n-1)/2) - (n-i)*((n-i)-1)/2 + j - i - 1
	return int(k)

def triCoord(k, n):
	i = n - 2 - math.floor(math.sqrt(-8*k + 4*n*(n-1)-7)/2.0 - 0.5)
	j = k + i + 1 - n*(n-1)/2 + (n-i)*((n-i)-1)/2
	return int(i), int(j)

def main():
	dimension = 8
	print("dimension", dimension)

	matrix = []
	for d in range(dimension):
		matrix.append([])
		for e in range(dimension):
			matrix[d].append((d,e))

	c = 0
	for d in range(dimension):
		for e in range(dimension):
			print("[{:2d}, {:2d}, {:2d}]".format(d, e, c), end="")
			c += 1
		print()

	print()

	c = 0
	for d in range(dimension):
		for e in range(dimension):
			if d < e:
				print("[{:2d}, {:2d}, {:2d}]".format(d, e, c), end="")
				c += 1
			else:
				print("[  ,   ,   ]", end="")
		print()

	triSize = int(dimension * (dimension-1) / 2)
	print("triSize", triSize)
	tri = [(p,triCoord(p, dimension)) for p in range(triSize)]
	print(tri)

	for d in range(dimension):
		for e in range(dimension):
			if d < e:
				p = triPos(d, e, dimension)
				cd, ce = triCoord(p, dimension)
				assert d == cd
				assert e == ce
				assert p == tri[p][0]
				assert d == tri[p][1][0]
				assert e == tri[p][1][1]
				print("[{:2d}, {:2d}, {:2d}]".format(cd, ce, p), end="")
			else:
				print("[  ,   ,   ]", end="")
		print()

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




if __name__ == "__main__":
    main(sys.argv[1])