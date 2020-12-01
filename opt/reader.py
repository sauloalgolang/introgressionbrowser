#!/usr/bin/env python3

import os
import sys
import typing
import json
import gzip
import struct
import multiprocessing as mp

import numpy as np
from collections import OrderedDict

"""
    python requirements:
        numpy
        flexx

    system requirements:
        bgzip

    optional requirements:
        pypy
"""


"""
    1 thread  - no alignment
    real    220m11.411s
    user    436m20.040s
    sys       0m41.524s

    4 threads - no alignment
    real     80m24.037s
    user    373m12.301s
    sys       0m40.070s

    6 threads - no alignment
    real     81m47.819s
    user    498m14.096s
    sys       1m27.114s

    8 threads - no alignment
    real     57m19.298s
    user    398m 2.372s
    sys       3m 0.994s

    ===========
    6 threads - no alignment - 20 bins
    real     6m12.305s
    user    39m 2.290s
    sys      0m 4.687s

    6 threads - w/ alignment - 20 bins
    real     7m18.622s
    user    45m43.042s
    sys      0m 5.512s

"""

"""
    1023M 360_merged_2.50.vcf.gz
    5.3M 360_merged_2.50.vcf.gz.250000.000000.SL2.50ch00.npz
    23M 360_merged_2.50.vcf.gz.250000.000001.SL2.50ch01.npz
    12M 360_merged_2.50.vcf.gz.250000.000002.SL2.50ch02.npz
    15M 360_merged_2.50.vcf.gz.250000.000003.SL2.50ch03.npz
    18M 360_merged_2.50.vcf.gz.250000.000004.SL2.50ch04.npz
    16M 360_merged_2.50.vcf.gz.250000.000005.SL2.50ch05.npz
    11M 360_merged_2.50.vcf.gz.250000.000006.SL2.50ch06.npz
    14M 360_merged_2.50.vcf.gz.250000.000007.SL2.50ch07.npz
    15M 360_merged_2.50.vcf.gz.250000.000008.SL2.50ch08.npz
    15M 360_merged_2.50.vcf.gz.250000.000009.SL2.50ch09.npz
    16M 360_merged_2.50.vcf.gz.250000.000010.SL2.50ch10.npz
    15M 360_merged_2.50.vcf.gz.250000.000011.SL2.50ch11.npz
    19M 360_merged_2.50.vcf.gz.250000.000012.SL2.50ch12.npz
    4.0K 360_merged_2.50.vcf.gz.250000.npz
    496K 360_merged_2.50.vcf.gz.csi
    2.8M 360_merged_2.50.vcf.gz.gzi
    8.0K 360_merged_2.50.vcf.gz.gzj
"""

DEBUG                         = True
DEBUG_MAX_BIN                 = 2
DEFAULT_BIN_SIZE              = 250_001
DEFAULT_COUNTER_TYPE_MATRIX   = np.uint16
DEFAULT_COUNTER_TYPE_PAIRWISE = np.uint32
DEFAULT_POSITIONS_TYPE        = np.uint32
DEFAULT_THREADS               = mp.cpu_count()

MatrixType           = typing.OrderedDict[typing.Tuple[str,str],int]
ChromosomeMatrixType = typing.OrderedDict[int, typing.List[int]]
BinSnpsType          = typing.OrderedDict[int, int]
BinPairwiseCountType = typing.OrderedDict[int, typing.List[int]]
BinPositionType      = typing.OrderedDict[int, typing.List[int]]
BinPositionTypeInt   = typing.OrderedDict[int, typing.List[typing.List[int]]]
BinAlignmentType     = typing.OrderedDict[int, typing.List[str]]
BinAlignmentTypeInt  = typing.OrderedDict[int, typing.List[typing.List[str]]]
TriangleIndexType    = typing.OrderedDict[typing.Tuple[int,int], int]
IUPACType            = typing.Dict[typing.FrozenSet[str], str]





class Chromosome():
    """
        import reader
        chrom = reader.Chromosome("../data/360_merged_2.50.vcf.gz", 250_000, 0, "SL2.50ch00")
        chrom.filename
        chrom.exists
        chrom.load()
        chrom.matrixNp
        # array([
        # [ 58,  18,  18, ...,   0,   0,   0],
        # [146,  44,  44, ...,   0,   0,   0],
        # [ 92,  18,  16, ...,   0,   0,   0],
        # ...,
        # [ 12,  28,  24, ...,   0,   0,   0],
        # [ 10,  14,  10, ...,   0,   0,   0],
        # [ 68,  78,  64, ...,   0,   0,   0]], dtype=uint16)
        chrom.matrixNp[0,:]
        # array([58, 18, 24, ..., 12,  0,  0], dtype=uint16)
        chrom.matrixNp.shape
        # (87, 64980)
        chrom.sample_count
        # 361
        chrom.matrixNp.sum(axis=1)
        # array([14135926, 23254563, 19677883, 15874453,  8422700, 17916464,
        #        13698021, 16151376,  7881003, 16836865, 19910092, 17204299,
        #        16263240,  7178839, 14871309,  9396853,  5977816, 10596943,
        #        15257031, 20200738, 16195559, 11704718,  8208319, 18487166,
        #        33227260, 19990183,  9529424,  8905346,  6964558,  4755606,
        #         7747795,   511011,  4415965,  5857671,  6814355,  6168579,
        #         1453627,  1744796,  4442980,  2088842,  5274180,  2468293,
        #         1546854,  2285387,  1936388,  1377816,  2966110,  2667990,
        #         4795700,  3008045,  4955744,  4958694,  6057739,  5587246,
        #         8511794,  7176901,  4027879,  4136236,  3431414,  4764401,
        #         4960412,  4165828,  5062155,  3305293,  5670359,  5311667,
        #         6011455,  6375255,  5160923,  4773837,  5889080,  5730958,
        #         6715606,  6523644,  5566348,  6531633,  6160864,  5562834,
        #         5614820,  3962912,  4594549,  3875737,  4870790,  4368941,
        #         4920943,  4689962,  5360525], dtype=uint64)
        m = reader.triangleToMatrix(chrom.sample_count, chrom.matrixNp[0,:])
        m
        #
        #array([
        # [ 0, 58, 18, ..., 16, 32,  0],
        # [58,  0, 18, ..., 16, 34,  0],
        # [18, 18,  0, ..., 14, 16,  0],
        # ...,
        # [16, 16, 14, ...,  0, 12,  0],
        # [32, 34, 16, ..., 12,  0,  0],
        # [ 0,  0,  0, ...,  0,  0,  0]], dtype=uint16)
        #
        m.shape
        #
        #(361, 361)
        #
    """
    def __init__(self,
            vcf_name: str,
            bin_width: int,
            chromosome_order: int,
            chromosome_name: str,
            type_matrix_counter   = DEFAULT_COUNTER_TYPE_MATRIX,
            type_pairwise_counter = DEFAULT_COUNTER_TYPE_PAIRWISE,
            type_positions        = DEFAULT_POSITIONS_TYPE
        ):
        self.type_matrix_counter          = type_matrix_counter
        self.type_pairwise_counter        = type_pairwise_counter
        self.type_positions               = type_positions
        self.type_matrix_counter_max_val  : int = np.iinfo(self.type_matrix_counter).max
        self.type_pairwise_counter_max_val: int = np.iinfo(self.type_pairwise_counter).max
        self.type_positions_max_val       : int = np.iinfo(self.type_positions).max

        self.vcf_name                     : str = vcf_name
        self.bin_width                    : int = bin_width
        self.chromosome_order             : int = chromosome_order
        self.chromosome_name              : str = chromosome_name

        self.matrix_size                  : int = None
        self.bin_max                      : int = None
        self.bin_min                      : int = None
        self.bin_count                    : int = None
        self.chromosome_snps              : int = None
        self.chromosome_first_position    : int = None
        self.chromosome_last_position     : int = None

        self.sample_names                 : typing.List[str] = None
        self.sample_count                 : int = None

        self.matrixNp                     : np.ndarray = None
        self.totalsNp                     : np.ndarray = None
        self.pairwiNp                     : np.ndarray = None
        self.alignmentNp                  : np.ndarray = None
        self.positionNp                   : np.ndarray = None

    @property
    def filename(self) -> str:
        return f"{self.vcf_name}.{self.bin_width}.{self.chromosome_order:06d}.{self.chromosome_name}.npz"

    @property
    def exists(self) -> bool:
        return os.path.exists(self.filename)

    def addFromVcf(self,
            matrix_size              : int,
            chromosome_snps          : int,
            chromosome_matrix        : ChromosomeMatrixType,
            bin_alignment            : BinAlignmentType,
            bin_positions            : BinPositionType,
            bin_snps                 : BinSnpsType,
            bin_pairwise_count       : BinPairwiseCountType,
            chromosome_first_position: int,
            chromosome_last_position : int,
            sample_names             : typing.List[str]
        ):
        print(f"adding chromosome data: {self.chromosome_name}")

        self.matrix_size                  = matrix_size

        bin_names                         = list(chromosome_matrix.keys())

        self.bin_min                      = min(bin_names)
        self.bin_max                      = max(bin_names)
        self.bin_count                    = len(bin_names)

        self.chromosome_snps              = chromosome_snps
        self.sample_count                 = len(sample_names)
        self.sample_names                 = sample_names
        self.chromosome_first_position    = chromosome_first_position
        self.chromosome_last_position     = chromosome_last_position

        self.matrixNp                     = np.zeros((self.bin_max, self.matrix_size ), self.type_matrix_counter  )
        self.pairwiNp                     = np.zeros((self.bin_max, self.sample_count), self.type_pairwise_counter)
        self.totalsNp                     = np.zeros( self.bin_max,                     self.type_pairwise_counter)

        self.alignmentNp                  = None
        self.positionNp                   = None
        position_size                     = None
        if bin_alignment:
            self.alignmentNp              = np.zeros((self.bin_max, self.sample_count), np.unicode_)
            # position_size                 = bin_positions[list(bin_positions.keys())[0]].shape[0]
            position_size = max([v for v in bin_snps.values()])
            # print( "position_size", position_size)
            self.positionNp               = np.zeros((self.bin_max, position_size    ), self.type_positions)
        else:
            self.alignmentNp              = np.zeros(0, np.unicode_        )
            self.positionNp               = np.zeros(0, self.type_positions)

        for binNum in range(self.bin_max):
            # print("  binNum", binNum)
            chromosome_matrix_bin  = chromosome_matrix .get(binNum, np.zeros(self.matrix_size , self.type_matrix_counter  ))
            bin_pairwise_count_bin = bin_pairwise_count.get(binNum, np.zeros(self.sample_count, self.type_pairwise_counter))
            bin_snps_bin           = bin_snps          .get(binNum, 0)
            
            bin_alignment_bin      = None
            bin_positions_bin      = None
            if bin_alignment:
                # print("   creating alignment")
                bin_alignment_bin  = bin_alignment     .get(binNum, np.zeros(self.sample_count, np.unicode_        ))
                bin_positions_bin  = bin_positions     .get(binNum, np.zeros(position_size    , self.type_positions))
            
            # print("bin_pairwise_count_bin", bin_pairwise_count_bin)
            assert not any([v > self.type_matrix_counter_max_val   for v in chromosome_matrix_bin ]), f"value went over the maximum value ({self.type_matrix_counter_max_val  }) for container {self.type_matrix_counter  }: {[v for v in chromosome_matrix_bin  if v > self.type_matrix_counter_max_val  ]}"
            assert not any([v > self.type_pairwise_counter_max_val for v in bin_pairwise_count_bin]), f"value went over the maximum value ({self.type_pairwise_counter_max_val}) for container {self.type_pairwise_counter}: {[v for v in bin_pairwise_count_bin if v > self.type_pairwise_counter_max_val]}"
            
            # binData             = [maxVal if v > maxVal else v for v in binData]
            self.matrixNp[binNum,:]    = chromosome_matrix_bin
            self.pairwiNp[binNum,:]    = bin_pairwise_count_bin
            self.totalsNp[binNum  ]    = bin_snps_bin
            if bin_alignment:
                # print("   adding alignment")
                self.alignmentNp[binNum,:] = bin_alignment_bin
                self.positionNp [binNum,:] = bin_positions_bin
        
        print(f"chromosome data added: {self.chromosome_name}")

    def _get_infos(self):
        return (
            ["matrix_size"   , "bin_count"         , "bin_min"                , "bin_max"                   , "bin_width"   , "chromosome_snps"   , "sample_count"   , "chromosome_order"   , "chromosome_first_position"   , "chromosome_last_position"   , "type_matrix_counter_max_val"   , "type_pairwise_counter_max_val"   , "type_positions_max_val"   ],
            [self.matrix_size, self.bin_count      , self.bin_min             , self.bin_max                , self.bin_width, self.chromosome_snps, self.sample_count, self.chromosome_order, self.chromosome_first_position, self.chromosome_last_position, self.type_matrix_counter_max_val, self.type_pairwise_counter_max_val, self.type_positions_max_val]
        )

    def _get_meta(self):
        type_matrix_counter_name   = self.type_matrix_counter.__name__
        type_pairwise_counter_name = self.type_pairwise_counter.__name__
        type_positions_name        = self.type_positions.__name__
        return (
            ["vcf_name"      , "chromosome_name"   , "type_matrix_counter_name", "type_pairwise_counter_name", "type_positions_name"],
            [self.vcf_name   , self.chromosome_name, type_matrix_counter_name  , type_pairwise_counter_name  , type_positions_name  ]
        )

    def __repr__(self):
        return str(self)

    def __str__(self):
        res = []
        for k in [
            "vcf_name",
            "chromosome_name",
            "chromosome_snps",
            "chromosome_order",
            "chromosome_first_position",
            "chromosome_last_position",
            "bin_count",
            "bin_max",
            "bin_min",
            "bin_width",
            "matrix_size",
            "sample_names",
            "sample_count",
            "type_matrix_counter",
            "type_pairwise_counter",
            "type_positions",
            "type_matrix_counter_max_val",
            "type_pairwise_counter_max_val",
            "type_positions_max_val",
            "matrixNp",
            "totalsNp",
            "pairwiNp",
            "alignmentNp",
            "positionNp"
        ]:

            v = getattr(self, k)
            s = None
            if   isinstance(v, int):
                s = f"{v:,d}"
            elif isinstance(v, str):
                s = f"{v:s}"
            elif isinstance(v, list):
                s = f"{len(v):,d}"
            elif isinstance(v, np.ndarray):
                s = f"{str(v.shape):s}"
            else:
                s = f"{str(v):s}"
            res.append(f"  {k:.<30s}{s:.>30s}")
        return "\n".join(res)

    def save(self):
        print(f"{'saving numpy array:':.<32s}{self.filename:.>30s}")
        print(self)

        info_names, info_vals      = self._get_infos()
        meta_names, meta_vals      = self._get_meta()

        sample_namesNp             = np.array(self.sample_names, np.unicode_)
        info_namesNp               = np.array(info_names, np.unicode_)
        info_valuesNp              = np.array(info_vals , np.int64   )
        meta_namesNp               = np.array(meta_names, np.unicode_)
        meta_valuesNp              = np.array(meta_vals , np.unicode_)

        np.savez_compressed(self.filename,
            countMatrix  = self.matrixNp,
            countTotals  = self.totalsNp,
            countPairw   = self.pairwiNp,
            alignments   = self.alignmentNp,
            positions    = self.positionNp,
            sample_names = sample_namesNp,
            info_names   = info_namesNp,
            info_values  = info_valuesNp,
            meta_names   = meta_namesNp,
            meta_values  = meta_valuesNp
        )

    def load(self):
        print(f"{'loading numpy array:':.<32s}{self.filename:.>30s}")
        data                               = np.load(self.filename, mmap_mode='r', allow_pickle=False)

        self.matrixNp                      = data['countMatrix']
        self.totalsNp                      = data['countTotals']
        self.pairwiNp                      = data['countPairw']
        self.alignmentNp                   = data['alignments']
        self.positionNp                    = data['positions']
        
        sample_namesNp                     = data['sample_names']
        
        info_namesNp                       = data['info_names']
        info_valuesNp                      = data['info_values']

        meta_namesNp                       = data['meta_names']
        meta_valuesNp                      = data['meta_values']

        sample_names                       = sample_namesNp.tolist()

        info_names                         = info_namesNp.tolist()
        info_values                        = info_valuesNp.tolist()
        info_values                        = [int(v) for v in info_values]
        info_dict                          = {info_names[k]: info_values[k]  for k in range(len(info_names))}
        # print(info_dict)

        meta_names                         = meta_namesNp.tolist()
        meta_values                        = meta_valuesNp.tolist()
        meta_dict                          = {meta_names[k]: meta_values[k]  for k in range(len(meta_names))}
        # print(meta_dict)

        vcf_name                           = meta_dict["vcf_name"]
        assert vcf_name == self.vcf_name
        self.vcf_name                      = vcf_name

        self.matrix_size                   = info_dict["matrix_size"]
        assert self.matrix_size == self.matrixNp.shape[1]

        self.bin_count                     = info_dict["bin_count"]
        self.bin_min                       = info_dict["bin_min"]
        self.bin_max                       = info_dict["bin_max"]
        assert self.bin_max == self.matrixNp.shape[0]
        assert self.bin_max == self.pairwiNp.shape[0]

        bin_width                          = info_dict["bin_width"]
        assert bin_width == self.bin_width
        self.bin_width = bin_width

        self.chromosome_snps               = info_dict["chromosome_snps"]
        chromosome_name                    = meta_dict["chromosome_name"]
        chromosome_order                   = info_dict["chromosome_order"]
        assert self.chromosome_name  == chromosome_name
        assert self.chromosome_order == chromosome_order
        self.chromosome_name               = chromosome_name
        self.chromosome_order              = chromosome_order

        self.chromosome_first_position     = info_dict["chromosome_first_position"]
        self.chromosome_last_position      = info_dict["chromosome_last_position"]
        assert (self.chromosome_first_position // self.bin_width) == self.bin_min
        assert (self.chromosome_last_position  // self.bin_width) == self.bin_max

        type_matrix_counter_name           = meta_dict["type_matrix_counter_name"]
        type_pairwise_counter_name         = meta_dict["type_pairwise_counter_name"]
        type_positions_name                = meta_dict["type_positions_name"]

        self.type_matrix_counter           = getattr(np, type_matrix_counter_name)
        self.type_pairwise_counter         = getattr(np, type_pairwise_counter_name)
        self.type_positions                = getattr(np, type_positions_name)

        assert self.type_matrix_counter    == self.matrixNp.dtype
        assert self.type_pairwise_counter  == self.pairwiNp.dtype
        assert self.type_positions         == self.positionNp.dtype

        self.type_matrix_counter_max_val   = info_dict["type_matrix_counter_max_val"]
        self.type_pairwise_counter_max_val = info_dict["type_pairwise_counter_max_val"]
        self.type_positions_max_val        = info_dict["type_positions_max_val"]
        assert np.iinfo(self.type_matrix_counter  ).max == self.type_matrix_counter_max_val
        assert np.iinfo(self.type_pairwise_counter).max == self.type_pairwise_counter_max_val
        assert np.iinfo(self.type_positions       ).max == self.type_positions_max_val

        self.sample_count                  = info_dict["sample_count"]
        self.sample_names                  = sample_names
        assert len(self.sample_names) == self.sample_count
        assert self.sample_count      == self.pairwiNp.shape[1]
        if self.alignmentNp.shape[0] != 0:
            assert self.sample_count  == self.alignmentNp.shape[1]
        
        print(self)


class Genome():
    def __init__(self,
            vcf_name: str,
            bin_width: int        = DEFAULT_BIN_SIZE,
            type_matrix_counter   = DEFAULT_COUNTER_TYPE_MATRIX,
            type_pairwise_counter = DEFAULT_COUNTER_TYPE_PAIRWISE,
            type_positions        = DEFAULT_POSITIONS_TYPE,
            save_alignment       : bool =False,
            diff_matrix          : MatrixType=None,
            IUPAC                : IUPACType=None
        ):
        self.vcf_name             : str = vcf_name
        self.bin_width            : int = bin_width
        self.type_matrix_counter        = type_matrix_counter
        self.type_pairwise_counter      = type_pairwise_counter
        self.type_positions             = type_positions

        self.sample_names         : typing.List[str] = None
        self.sample_count         : int = None

        self.chromosome_names     : typing.List[str] = None
        self.chromosome_count     : int = None

        self.genome_bins          : int = None
        self.genome_snps          : int = None

        self._chromosomes         : typing.List[Chromosome] = None
        self._save_alignment      : bool = save_alignment
        self._diff_matrix         : MatrixType=diff_matrix
        self._IUPAC               : IUPACType=IUPAC

    @property
    def filename(self) -> str:
        return f"{self.vcf_name}.{self.bin_width}.npz"

    @property
    def exists(self) -> bool:
        return os.path.exists(self.filename)

    @property
    def loaded(self) -> bool:
        if not self.exists:
            return False
        
        if self._chromosomes is None:
            return False

        if len(self._chromosomes) == 0:
            return False

        return True

    @property
    def complete(self) -> bool:
        if not self.loaded:
            return False
        
        for chromosome in self._chromosomes:
            if not chromosome.exists:
                return False

        return True

    def _processVcf(self, threads=DEFAULT_THREADS):
        self._chromosomes     = []
        self.chromosome_names = []
        self.chromosome_count = 0
        self.genome_bins      = 0
        self.genome_snps      = 0

        if self._IUPAC is None:
            self._IUPAC = genIUPAC()

        if self._diff_matrix is None:
            self._diff_matrix = genDiffMatrix()

        assert os.path.exists(self.vcf_name)

        sample_names, sample_count, matrix_size, indexes = Genome._processVcf_read_header(self.vcf_name)
        self.sample_names = sample_names
        self.sample_count = sample_count

        bgzip = BGzip(self.vcf_name)

        self.chromosome_names = bgzip.chromosomes
        self.chromosome_count = len(self.chromosome_names)

        mp.set_start_method('spawn')

        results = []
        with mp.Pool(processes=threads) as pool:
            for chromosome_order, chromosome_name in enumerate(self.chromosome_names):
                chromosome = Chromosome(
                    vcf_name              = self.vcf_name,
                    bin_width             = self.bin_width,
                    chromosome_order      = chromosome_order,
                    chromosome_name       = chromosome_name,
                )
                
                if chromosome.exists:
                    print(f"chromosome {chromosome_name} already exists")
                    continue

                print(f"reading chromosome {chromosome_name} from vcf")
                res = pool.apply_async(
                    Genome._processVcf_read_chrom,
                    [],
                    {
                        "chromosome_order"      : chromosome_order,
                        "chromosome_name"       : chromosome_name,

                        "matrix_size"           : matrix_size,
                        "indexes"               : indexes,

                        "vcf_name"              : self.vcf_name,
                        "bin_width"             : self.bin_width,
                        "sample_names"          : self.sample_names,
                        "sample_count"          : self.sample_count,

                        "type_matrix_counter"   : self.type_matrix_counter,
                        "type_pairwise_counter" : self.type_pairwise_counter,
                        "type_positions"        : self.type_positions,

                        "save_alignment"        : self._save_alignment,
                        "IUPAC"                 : self._IUPAC,
                        "diff_matrix"           : self._diff_matrix
                    }
                )
                results.append([False, res, chromosome_order, chromosome_name])
            
            while not all([r[0] for r in results]):
                for resnum, (finished, res, chromosome_order, chromosome_name) in enumerate(results):
                    if finished:
                        continue
                    
                    if not res.ready():
                        continue

                    try:
                        chromosome_bins, chromosome_snps = res.get(timeout=1)
                    except mp.TimeoutError:
                        print(f"waiting for {chromosome_name}")
                        continue

                    print(f"got results from {chromosome_name}")

                    results[resnum][0] = True
                    self.genome_bins += chromosome_bins
                    self.genome_snps += chromosome_snps

        for resnum, (finished, res, chromosome_order, chromosome_name) in enumerate(results):
            print(f"loading {chromosome_name}")
            chromosome = Chromosome(
                vcf_name              = self.vcf_name,
                bin_width             = self.bin_width,
                chromosome_order      = chromosome_order,
                chromosome_name       = chromosome_name,
                type_matrix_counter   = self.type_matrix_counter,
                type_pairwise_counter = self.type_pairwise_counter,
                type_positions        = self.type_positions
            )
            if not chromosome.exists:
                raise IOError(f"chromosome database does not exists: {chromosome.filename}")
            chromosome.load()
            self._chromosomes.append(chromosome)
        print("all chromosomes loaded")

    @staticmethod
    def _processVcf_read_header(vcf_name: str) -> typing.Tuple[typing.List[str], int, int, TriangleIndexType]:
        with openFile(vcf_name, 'rt') as fhd:
            for line in fhd:
                line = line.strip()
                
                if len(line) <= 2:
                    continue

                if line[:2] == "##":
                    # print("header", line)
                    continue
                
                if line[0] == "#":
                    # print("title", line)
                    cols         = line[1:].split("\t")
                    sample_names = cols[9:]
                    sample_count = len(sample_names)
                    # matrix_size  = calculateMatrixSize(sample_count)
                    _, matrix_size, indexes = triangleToIndex(sample_count)

                    print( "sample_names", ",".join(sample_names))
                    print(f"num samples  {sample_count:12,d}")
                    print(f"matrix_size  {matrix_size :12,d}")
                    print(f"indexes      {len(indexes):12,d}")
                    return sample_names, sample_count, matrix_size, indexes

                else:
                    raise ValueError("data before header error", line)

    @staticmethod
    def _processVcf_read_chrom(
            chromosome_order: int,
            chromosome_name : str,

            matrix_size     : int,
            indexes         : TriangleIndexType,
            
            vcf_name        : str,
            bin_width       : int,
            sample_names    : typing.List[str],
            sample_count    : int,

            type_matrix_counter,
            type_pairwise_counter,
            type_positions,

            save_alignment  : bool,
            IUPAC           : IUPACType,
            diff_matrix     : MatrixType
        ) -> typing.Tuple[int, int]:

        print(f"reading {chromosome_name}")

        bgzip = BGzip(vcf_name)

        bin_snps                 : BinSnpsType          = OrderedDict()
        bin_pairwise_count       : BinPairwiseCountType = OrderedDict()
        bin_alignment            : BinAlignmentTypeInt  = OrderedDict()
        bin_positions            : BinPositionTypeInt   = OrderedDict()
        chromosome_matrix        : MatrixType           = OrderedDict()
        chromosome_first_position: int = 0
        chromosome_last_position : int = 0
        chromosome_snps          : int = 0

        lastBinNum = None
        for line in bgzip.get_chromosome(chromosome_name):
            # print(line)
            cols       = line.split("\t")
            chrom      = cols[0]
            ref        = cols[3]
            alts       = cols[4].split(",")
            opts       = [ref] + alts
            pos        = int(cols[1])
            samples    = cols[9:]
            binNum     = pos // bin_width
            assert len(samples) == sample_count
            
            if len(ref) != 1:
                print('H')
                continue

            if any([len(a) > 1 for a in alts]):
                print('h')
                continue

            if binNum not in chromosome_matrix:
                if DEBUG:
                    if len(chromosome_matrix) > DEBUG_MAX_BIN:
                        break

                print(f"  New bin: {chrom} :: {pos:12,d} => {binNum:12,d}")

                if lastBinNum is not None:
                    chromosome_matrix[ lastBinNum] = np.array(chromosome_matrix [lastBinNum], type_matrix_counter  )
                    bin_pairwise_count[lastBinNum] = np.array(bin_pairwise_count[lastBinNum], type_pairwise_counter)
                    bin_snps          [lastBinNum] = np.array(bin_snps          [lastBinNum], type_pairwise_counter)

                    if save_alignment:
                        bin_alignment[lastBinNum]  = ["".join(b) for b in bin_alignment[lastBinNum]]
                        # bin_positions[lastBinNum]  = [p for p in bin_positions  [lastBinNum] if p != -1]
                        bin_positions[lastBinNum]  = np.array(bin_positions     [lastBinNum], type_positions       )

                lastBinNum                 = binNum
                chromosome_matrix[binNum]  = [0] * matrix_size
                bin_pairwise_count[binNum] = [0] * sample_count
                bin_snps[binNum]           =  0

                bin_alignment[binNum]      = None
                if save_alignment:
                    bin_alignment[binNum]  = [[] for _ in range(sample_count)]
                    bin_positions[binNum]  = []

            samples                   = [s.split(";")[0]            for s in samples]
            samples                   = [s if len(s) == 3 else None for s in samples]
            # samples = [s.replace("|", "").replace("/", "") if s is not None else None for s in samples]
            samples                   = [tuple([int(i) for i in s.replace("/", "|").split("|")]) if s is not None else None for s in samples]
            vals                      = chromosome_matrix[binNum]
            paiw                      = bin_pairwise_count[binNum]
            chromosome_snps          += 1
            bin_snps[binNum]         += 1
            chromosome_last_position  = pos

            aling                     = None
            if save_alignment:
                aling                 = bin_alignment[binNum]
                bin_positions[binNum].append(pos)

            if chromosome_first_position == 0:
                chromosome_first_position = pos

            for sample1num in range(sample_count):
                sample1 = samples[sample1num]

                if sample1 is None:
                    if save_alignment:
                        aling[sample1num].append('N')
                    continue
                elif save_alignment:
                    alts1 = frozenset([opts[s] for s in sample1])
                    nuc   = IUPAC[alts1]
                    aling[sample1num].append(nuc)

                for sample2num in range(sample1num+1,sample_count):
                    sample2 = samples[sample2num]

                    if sample2 is None:
                        continue

                    k = (sample1, sample2) if sample1 <= sample2 else (sample2, sample1)

                    value = diff_matrix.get(k, None)
                    if value is None:
                        print(line)
                        raise ValueError("multiallelic", k)

                    pairind           = indexes[(sample1num,sample2num)]
                    vals[pairind   ] += value
                    paiw[sample1num] += value
                    paiw[sample2num] += value

        print(f"cleaning {chromosome_name}")

        chromosome_matrix[ lastBinNum] = np.array(chromosome_matrix [lastBinNum], type_matrix_counter  )
        bin_pairwise_count[lastBinNum] = np.array(bin_pairwise_count[lastBinNum], type_pairwise_counter)
        bin_snps          [lastBinNum] = np.array(bin_snps          [lastBinNum], type_pairwise_counter)

        if save_alignment:
            bin_alignment[lastBinNum]  = ["".join(b) for b in bin_alignment[lastBinNum]]
            # bin_positions[lastBinNum]  = [p for p in bin_positions  [lastBinNum] if p != -1]
            bin_positions[lastBinNum]  = np.array(bin_positions     [lastBinNum], type_positions       )
            position_max_size          = max([p.shape[0] for p in bin_positions.values()])

            for binNum, binval in bin_positions.items():
                binz                   = np.zeros(position_max_size, type_positions)
                binz[:binval.shape[0]] = binval
                bin_positions[binNum]  = binz

        Genome._processVcf_save_chrom_data(
            vcf_name                     = vcf_name,
            bin_width                    = bin_width,
            sample_names                 = sample_names,
            sample_count                 = sample_count,
            matrix_size                  = matrix_size,

            chromosome_name              = chromosome_name,
            chromosome_order             = chromosome_order,
            chromosome_snps              = chromosome_snps,
            chromosome_matrix            = chromosome_matrix,
            chromosome_first_position    = chromosome_first_position,
            chromosome_last_position     = chromosome_last_position,
            
            bin_alignment                = bin_alignment if save_alignment else None,
            bin_positions                = bin_positions if save_alignment else None,
            bin_snps                     = bin_snps,
            bin_pairwise_count           = bin_pairwise_count,
            
            type_matrix_counter          = type_matrix_counter,
            type_pairwise_counter        = type_pairwise_counter,
            type_positions               = type_positions
        )

        chromosome_bins = len(chromosome_matrix)

        print(f"returning {chromosome_name}")

        return chromosome_bins, chromosome_snps

    @staticmethod
    def _processVcf_save_chrom_data(
            vcf_name                 : str,
            bin_width                : int,
            sample_names             : typing.List[str],
            sample_count             : int,
            matrix_size              : int,

            chromosome_name          : str,
            chromosome_order         : int,
            chromosome_snps          : int,
            chromosome_matrix        : ChromosomeMatrixType,
            chromosome_first_position: int,
            chromosome_last_position : int,
            
            bin_alignment            : BinAlignmentType,
            bin_positions            : BinPositionType,
            bin_snps                 : BinSnpsType,
            bin_pairwise_count       : BinPairwiseCountType, 

            type_matrix_counter,
            type_pairwise_counter,
            type_positions
        ):

        print(f"creating {chromosome_name}")

        # self.chromosome_names.append(chromosome_name)
        # self.chromosome_count += 1

        chromosome = Chromosome(
            vcf_name                  = vcf_name,
            bin_width                 = bin_width,
            chromosome_order          = chromosome_order,
            chromosome_name           = chromosome_name,
            type_matrix_counter       = type_matrix_counter,
            type_pairwise_counter     = type_pairwise_counter,
            type_positions            = type_positions
        )

        chromosome.addFromVcf(
            matrix_size               = matrix_size,
            chromosome_snps           = chromosome_snps,
            chromosome_matrix         = chromosome_matrix,
            bin_alignment             = bin_alignment,
            bin_positions             = bin_positions,
            bin_snps                  = bin_snps,
            bin_pairwise_count        = bin_pairwise_count,
            chromosome_first_position = chromosome_first_position,
            chromosome_last_position  = chromosome_last_position,
            sample_names              = sample_names
        )

        print(f"saving {chromosome_name}")

        chromosome.save()

    def __repr__(self):
        return str(self)

    def __str__(self):
        res = []
        for k in [
            "vcf_name",
            "bin_width",
            "chromosome_names",
            "chromosome_count",
            "genome_bins",
            "genome_snps",
            "sample_names",
            "sample_count",
            "type_matrix_counter",
            "type_pairwise_counter",
            "type_positions",
        ]:

            v = getattr(self, k)
            s = None
            if   isinstance(v, int):
                s = f"{v:,d}"
            elif isinstance(v, str):
                s = f"{v:s}"
            elif isinstance(v, list):
                s = f"{len(v):,d}"
            elif isinstance(v, np.ndarray):
                s = f"{str(v.shape):s}"
            else:
                s = f"{str(v):s}"
            res.append(f"  {k:.<30s}{s:.>30s}")
        return "\n".join(res)

    def _load_db(self):
        if self.complete:
            return

        print(f"{'loading numpy array:':.<32s}{self.filename:.>30s}")
        # filename           = f"{self.vcf_name}.{self.chromosome_order:06d}.{self.chromosome_name}.npz"
        data                       = np.load(self.filename, mmap_mode='r', allow_pickle=False)
        
        self.sample_names          = data["sample_names"].tolist()
        self.chromosome_names      = data["chromosome_names"].tolist()

        info_namesNp               = data['info_names']
        info_valuesNp              = data['info_values']

        meta_namesNp               = data['meta_names']
        meta_valuesNp              = data['meta_values']

        info_names                 = info_namesNp.tolist()
        info_values                = info_valuesNp.tolist()
        info_values                = [int(v) for v in info_values]
        info_dict                  = {info_names[k]: info_values[k]  for k in range(len(info_names))}
        # print(info_dict)

        meta_names                 = meta_namesNp.tolist()
        meta_values                = meta_valuesNp.tolist()
        meta_dict                  = {meta_names[k]: meta_values[k]  for k in range(len(meta_names))}
        # print(meta_dict)

        self.vcf_name              = meta_dict["vcf_name"]
        type_matrix_counter_name   = meta_dict["type_matrix_counter_name"]
        type_pairwise_counter_name = meta_dict["type_pairwise_counter_name"]
        type_positions_name        = meta_dict["type_positions_name"]
        self.type_matrix_counter   = getattr(np, type_matrix_counter_name)
        self.type_pairwise_counter = getattr(np, type_pairwise_counter_name)
        self.type_positions        = getattr(np, type_positions_name)

        self.bin_width             = info_dict["bin_width"]
        self.chromosome_count      = info_dict["chromosome_count"]
        self.sample_count          = info_dict["sample_count"]
        self.genome_bins           = info_dict["genome_bins"]
        self.genome_snps           = info_dict["genome_snps"]

        assert self.chromosome_count == len(self.chromosome_names)
        assert self.sample_count     == len(self.sample_names)

        print(self)

        chromosome_snps   = 0
        chromosome_bins   = 0
        self._chromosomes = []
        for chromosome_order, chromosome_name in enumerate(self.chromosome_names):
            chromosome = Chromosome(
                vcf_name              = self.vcf_name,
                bin_width             = self.bin_width,
                chromosome_order      = chromosome_order,
                chromosome_name       = chromosome_name,
                type_matrix_counter   = self.type_matrix_counter,
                type_pairwise_counter = self.type_pairwise_counter
            )
            
            if not chromosome.exists:
                raise IOError(f"chromosome database file {chromosome.filename} does not exists")

            chromosome.load()
            
            assert chromosome.bin_width             == self.bin_width
            assert chromosome.sample_names          == self.sample_names
            assert chromosome.vcf_name              == self.vcf_name
            assert chromosome.bin_width             == self.bin_width
            assert chromosome.chromosome_name       == chromosome_name
            assert chromosome.chromosome_order      == chromosome_order
            assert chromosome.type_matrix_counter   == self.type_matrix_counter
            assert chromosome.type_pairwise_counter == self.type_pairwise_counter
            assert chromosome.type_positions        == self.type_positions
            
            chromosome_snps += chromosome.chromosome_snps
            chromosome_bins += chromosome.bin_count
            
            self._chromosomes.append(chromosome)
        
        assert self.genome_bins == chromosome_bins
        assert self.genome_snps == chromosome_snps

    def save(self):
        print(f"{'saving numpy array:':.<32s}{self.filename:.>30s}")
        print(self)

        type_matrix_counter_name   = self.type_matrix_counter.__name__
        type_pairwise_counter_name = self.type_pairwise_counter.__name__
        type_positions_name        = self.type_positions.__name__

        sample_namesNp             = np.array(self.sample_names    , np.unicode_)
        chromosome_namesNp         = np.array(self.chromosome_names, np.unicode_)

        info_namesNp               = np.array(["bin_width"   , "chromosome_count"   , "sample_count"   , "genome_bins"   , "genome_snps"   ], np.unicode_)
        info_valuesNp              = np.array([self.bin_width, self.chromosome_count, self.sample_count, self.genome_bins, self.genome_snps], np.int64   )

        meta_namesNp               = np.array(["vcf_name"   , "type_matrix_counter_name", "type_pairwise_counter_name", "type_positions_name"], np.unicode_)
        meta_valuesNp              = np.array([self.vcf_name, type_matrix_counter_name  , type_pairwise_counter_name  , type_positions_name  ], np.unicode_)

        np.savez_compressed(self.filename,
            sample_names     = sample_namesNp,
            chromosome_names = chromosome_namesNp,
            info_names       = info_namesNp,
            info_values      = info_valuesNp,
            meta_names       = meta_namesNp,
            meta_values      = meta_valuesNp
        )

    def load(self, threads=DEFAULT_THREADS):
        if self.exists:
            print(f"database exists {self.filename}")
            self._load_db()

        else:
            print(f"database does not exists {self.filename}. reading vcf")
            self._processVcf(threads=threads)

            self.save()

            if not self.complete:
                raise IOError("not complete. not able to load the database")


class BGzip():
    """
        http://www.htslib.org/doc/bgzip.html#:~:text=GZI%20FORMAT,in%20the%20uncompressed%20data%20stream.

        The index format is a binary file listing pairs of compressed and
        uncompressed offsets in a BGZF file. Each compressed offset points to
        the start of a BGZF block. The uncompressed offset is the corresponding
        location in the uncompressed data stream.

        All values are stored as little-endian 64-bit unsigned integers.
        
        The file contents are:
            uint64_t number_entries

        followed by number_entries pairs of:
            uint64_t compressed_offset
            uint64_t uncompressed_offset
    """

    """
        rm 360_merged_2.50.vcf.gz.gzi
        rm 360_merged_2.50.vcf.gz.gzj

        import reader
        b = reader.BGzip("../data/360_merged_2.50.vcf.gz")
        b.chromosomes
        for line in b.get_chromosome('SL2.50ch00'):
            print(line)
        for line in b.get_chromosome('SL2.50ch11'):
            print(line)
    """

    def __init__(self, gzip_file: str):
        self.gzip_file = gzip_file
        self.gzi_file  = gzip_file + ".gzi"
        self.gzj_file  = gzip_file + ".gzj"

        assert os.path.exists(self.gzip_file)

        if not os.path.exists(self.gzi_file):
            print(f"index file {self.gzi_file} does not exists. creating")
            print(os.system(f"bgzip -r {self.gzip_file}"))
            print(f"index file created")

        self._data = OrderedDict()

        if os.path.exists(self.gzj_file):
            self._load()
        else:
            print(f"index file {self.gzj_file} does not exists. creating")
            self._parse_gzi()
            self._save()

    @property
    def chromosomes(self) -> typing.List[str]:
        chroms = [k for k in self._data.keys()]
        chroms.sort(key=lambda x: self._data[x]["entry_num"])
        return chroms

    def _save(self):
        print(f"saving gzj to {self.gzj_file}")
        json.dump(self._data, open(self.gzj_file, 'wt'), indent=1)

    def _load(self):
        print(f"reading gzj from {self.gzj_file}")
        self._data = json.load(open(self.gzj_file, 'rt'))

    def _parse_gzi(self):
        print(f"reading gzi from {self.gzi_file}")

        with open(self.gzip_file, 'rb') as gzip_fhd:
            with open(self.gzi_file, 'rb') as fhd:
                number_entries_fmt  = "<Q"
                number_entries_size = struct.calcsize(number_entries_fmt)
                
                offsets_fmt         = "<QQ"
                offsets_size        = struct.calcsize(offsets_fmt)
                
                number_entries      = struct.unpack(number_entries_fmt, fhd.read(number_entries_size))[0]
                print("number_entries", number_entries)
                
                previous_compressed_offset   = 0
                previous_uncompressed_offset = 0
                current_compressed_offset    = 0
                current_uncompressed_offset  = 0
                for entry_num in range(number_entries):
                    next_compressed_offset, next_uncompressed_offset = struct.unpack(offsets_fmt, fhd.read(offsets_size))
                    previous_compressed_size   = current_compressed_offset   - previous_compressed_offset
                    previous_uncompressed_size = current_uncompressed_offset - previous_uncompressed_offset
                    current_compressed_size    = next_compressed_offset      - current_compressed_offset
                    current_uncompressed_size  = next_uncompressed_offset    - current_uncompressed_offset

                    # print("current_compressed_offset", current_compressed_offset, "compressed_offset", compressed_offset, "current_uncompressed_offset", current_uncompressed_offset, "uncompressed_offset", uncompressed_offset, "compressed_size", compressed_size, "uncompressed_size", uncompressed_size)
                    
                    self._get_first_contig(
                        gzip_fhd,
                        entry_num,
                        
                        previous_compressed_offset,
                        previous_uncompressed_offset,
                        previous_compressed_size,
                        previous_uncompressed_size,

                        current_compressed_offset,
                        current_uncompressed_offset,
                        current_compressed_size,
                        current_uncompressed_size,
                    )
                    
                    previous_compressed_offset   = current_compressed_offset
                    previous_uncompressed_offset = current_uncompressed_offset
                    current_compressed_offset    = next_compressed_offset
                    current_uncompressed_offset  = next_uncompressed_offset
    
    def _get_first_contig(self,
            gzip_fhd : typing.IO,
            entry_num: int,
            
            previous_compressed_offset  : int,
            previous_uncompressed_offset: int,
            previous_compressed_size    : int,
            previous_uncompressed_size  : int,

            current_compressed_offset   : int,
            current_uncompressed_offset : int,
            current_compressed_size     : int,
            current_uncompressed_size   : int,
        ):

        gzip_fhd.seek(current_compressed_offset)

        compressed_block        = gzip_fhd.read(current_compressed_size)
        uncompressed_block      = gzip.decompress(compressed_block)
        assert len(uncompressed_block) == current_uncompressed_size
        uncompressed_block_text = uncompressed_block.decode()
        
        firstNewLine = 0
        first_tab     = 0
        chrom_name    = None
        while True:
            if firstNewLine >= len(uncompressed_block_text):
                raise ValueError(f"firstNewLine {firstNewLine} >= {len(uncompressed_block_text)} len(uncompressed_block_text)")

            firstNewLine  = uncompressed_block_text.find("\n", firstNewLine+1)
            if firstNewLine == -1:
                print("no firstNewLine")
                continue
            
            secondNewLine = uncompressed_block_text.find("\n", firstNewLine+1)
            if secondNewLine == -1:
                print("no secondNewLine")
                firstNewLine += 1
                continue
            
            firstLine     = uncompressed_block_text[firstNewLine+1:secondNewLine]
            if len(firstLine) == 0:
                print("no firstLine", firstNewLine, secondNewLine)
                firstNewLine += 1
                continue
            
            if firstLine[0] == "#":
                print("skipping header", firstLine)
                firstNewLine += 1
                continue
            
            first_tab      = firstLine.find("\t")
            if first_tab == -1:
                print("no first_tab", firstLine)
                continue
            
            chrom_name     = firstLine[:first_tab]
            if len(chrom_name) == 0:
                print("no chrom_name", firstNewLine, first_tab, firstLine)
                chrom_name = None
                continue

            break

        if chrom_name is None:
            raise ValueError("No chromosome name found")

        if chrom_name not in self._data:
            print(f"NEW CHROMOSOME")
            print(f"  chrom_name                   {chrom_name:>15s}")
            print(f"  entry_num                    {entry_num:15,d}")

            print(f"  previous_compressed_offset   {previous_compressed_offset:15,d}")
            print(f"  previous_uncompressed_offset {previous_uncompressed_offset:15,d}")
            print(f"  previous_compressed_size     {previous_compressed_size:15,d}")
            print(f"  previous_uncompressed_size   {previous_uncompressed_size:15,d}")

            print(f"  current_compressed_offset    {current_compressed_offset:15,d}")
            print(f"  current_uncompressed_offset  {current_uncompressed_offset:15,d}")
            print(f"  current_compressed_size      {current_compressed_size:15,d}")
            print(f"  current_uncompressed_size    {current_uncompressed_size:15,d}")

            self._data[chrom_name] = {
                "entry_num"                      : entry_num,
                
                "previous_compressed_offset"     : previous_compressed_offset,
                "previous_uncompressed_offset"   : previous_uncompressed_offset,
                "previous_compressed_size"       : previous_compressed_size,
                "previous_uncompressed_size"     : previous_uncompressed_size,

                "current_compressed_offset"      : current_compressed_offset,
                "current_uncompressed_offset"    : current_uncompressed_offset,
                "current_compressed_size"        : current_compressed_size,
                "current_uncompressed_size"      : current_uncompressed_size,
            }

    def get_chromosome(self, chrom_name: str) -> typing.Generator[str, None, None]:
        if chrom_name not in self._data:
            raise ValueError(f"chromosome {chrom_name} does not exists: {list(self._data.keys())}")
        
        chromosome_matrix = self._data[chrom_name]
        file_pos          = chromosome_matrix["previous_compressed_offset"]
        found_chrom       = False

        with open(self.gzip_file, 'rb') as fhd:
            fhd.seek(file_pos)
            ghd = gzip.open(fhd, 'rt')

            for linenum, line in enumerate(ghd):
                if len(line) == 0:
                    continue

                if line[0] == "#":
                    continue

                if file_pos != 0: #first chromosome
                    if linenum == 0: #possible incomplete line
                        # print('INCOMPLETE LINE', line)
                        continue

                first_tab = line.find("\t")
                if first_tab == -1:
                    print('NO TAB', line)
                    continue

                chrom = line[:first_tab]
                if chrom != chrom_name:
                    if found_chrom:
                        # print('WRONG CHROMOSOME TAIL', chrom, line)
                        break
                    else:
                        # print('WRONG CHROMOSOME HEAD', chrom, line)
                        continue

                found_chrom = True
                yield line

                # break



def openFile(filename, mode):
    if filename.endswith('.gz'):
        return gzip.open(filename, mode)
    else:
        return open(filename, mode)

def calculateMatrixSize(sample_count: int) -> int:
    """
        2 - 1   2 - 3
        - -     1 -
        1 -     2 3

        3 - 3   3 - 6
        - - -   1 - -
        1 - -   2 3 -
        2 3 -   4 5 6

        4 - 6    4 - 10
        - - - -  1 - - -
        1 - - -  2 3 - -
        2 3 - -  4 5 6 -
        4 5 6 -  7 8 9 10

        5 - 10      5 - 15
        - - - - -   1  -  -  -  -
        1 - - - -   2  3  -  -  -
        2 3 - - -   4  5  6  -  -
        4 5 6 - -   7  8  9 10  -
        7 8 9 10 - 11 12 13 14 15

        f = lambda sample_count: sum([x for x in range(sample_count)])
        f = lambda sample_count: sum([x for x in range(sample_count)]) + sample_count
    """
    return sum([x for x in range(sample_count)])# + sample_count

def triangleToIndex(size: int) -> typing.Tuple[int, TriangleIndexType]:
    """
        1
        2 3
        4 5 6
        7 8 9 10

        1,1 =  1

        2,1 =  2
        2,2 =  3

        3,1 =  4
        3,2 =  5
        3,3 =  6

        4,1 =  7
        4,2 =  8
        4,3 =  9
        4,4 = 10

        ========

        creating a lower triangle first and a list of counts in the end

        for 4 samples:
        1: count sample 1    -                    -                   -
        2: 1,2               3: count sample 2    -                   -
        4: 1,3               5: 2,3               6: count sample 3   -
        7: 1,4               8: 2,4               9: 3,4              10: count sample 4

        -                    -                    -          -
        1: 1,2               -                    -          -
        2: 1,3               3: 2,3               -          -
        4: 1,4               5: 2,4               6: 3,4     -
        7    8    9    10
        [ cs1, cs2, cs3, cs4 ]

        1   2   3   4   5   6   7   8   9   10
        1,2 1,2 2,3 1,4 2,4 3,4 cs1 cs2 cs3 cs4
    """

    """
        b = np.tril_indices(size)
        # (array([0, 1, 1, 2, 2, 2, 3, 3, 3, 3]), array([0, 0, 1, 0, 1, 2, 0, 1, 2, 3]))
        b = np.tril_indices(4, -1)
        # (array([1, 2, 2, 3, 3, 3]), array([0, 0, 1, 0, 1, 2]))
        a = np.arange(16).reshape(4, 4)
        a[b]
        # array([ 4,  8,  9, 12, 13, 14])
        bp = list(zip(b[0].tolist(), b[1].tolist()))
        bp
        # [(1, 0), (2, 0), (2, 1), (3, 0), (3, 1), (3, 2)]
    """

    indexes = OrderedDict()
    # b       = np.tril_indices(size, -1)
    b       = np.triu_indices(size,  1)
    l       = len(b[0].tolist())
    bp      = list(zip(b[0].tolist(), b[1].tolist()))
    # print("b", b)
    # print("bp", bp)
    
    for bl, bc in enumerate(bp):
        # print("index", bc, bl)
        indexes[bc] = bl

    # for bl in range(size):
    #     indexes[(bl, bl)] = len(bp) + bl
    
    return b, l, indexes

def triangleToMatrix(size, tri_array):
    """
        ```python
        a = np.array([[1,2,3],[4,5,6],[7,8,9]])
        a
        #
        #array([[1, 2, 3],
        #       [4, 5, 6],
        #       [7, 8, 9]])
        #
        size = a.shape[0]
        i, j = np.tril_indices(size, -1)
        i
        #
        #array([1, 2, 2])
        #
        j
        #
        #array([0, 0, 1])
        #
        a[i,j]
        #
        #array([4, 7, 8])
        #
        b = a[i,j]
        b
        #
        #array([4, 7, 8])
        #
        M = np.zeros([size,size], a.dtype)
        M
        #
        #array([[0, 0, 0],
        #       [0, 0, 0],
        #       [0, 0, 0]])
        #
        M[i,j] = b
        M[j,i] = b
        M
        #
        #array([[0, 4, 7],
        #       [4, 0, 8],
        #       [7, 8, 0]])
        #
        ```
    """

    """
        ```python
        import reader
        import numpy as np
        (i,j), l, od = reader.triangleToIndex(5)
        i
        j
        l
        od
        # OrderedDict([((0, 1), 0), ((0, 2), 1), ((0, 3), 2), ((0, 4), 3), ((1, 2), 4), ((1, 3), 5), ((1, 4), 6), ((2, 3), 7), ((2, 4), 8), ((3, 4), 9)])
        for i in range(5):
            for j in range(i+1, 5):
                od[(i,j)]

        # 0
        # 1
        # 2
        # 3
        # 4
        # 5
        # 6
        # 7
        # 8
        # 9
        ln = np.array([0,1,2,3,4,5,6,7,8,9])
        M = reader.triangleToMatrix(5, ln)
        M
        # array(
        # [[0, 0, 1, 2, 3],
        # [0, 0, 4, 5, 6],
        # [1, 4, 0, 7, 8],
        # [2, 5, 7, 0, 9],
        # [3, 6, 8, 9, 0]])
        ```
    """
    (i,j), l, _ = triangleToIndex(size)
    assert len(tri_array) == l

    M           = np.zeros((size, size), tri_array.dtype)

    M[i, j] = tri_array
    M[j, i] = tri_array
    
    return M

def genDiffMatrix(alphabet: typing.List[str]=[0,1,2,3]) -> MatrixType:
    # diff_matrixSymetricalHomoExtra = {
    #     # HomoHomo = 3
    #     # HeteHete = 2
    #     # HomoHete = 1
    #     '00': { '00': 3, '01': 1, '02': 1, '10': 1, '11': 0, '20': 1, '22': 0 },
    #     '01': {          '01': 2, '02': 1, '10': 2, '11': 1, '20': 1, '22': 0 },
    #     '02': {                   '02': 2, '10': 2, '11': 0, '20': 2, '22': 1 },
    #     '10': {                            '10': 2, '11': 1, '20': 1, '22': 0 },
    #     '11': {                                     '11': 0, '20': 0, '22': 0 },
    #     '20': {                                              '20': 2, '22': 1 },
    #     '22': {                                                       '22': 3 }
    # }
    diff_matrixSymetricalHomoExtra = OrderedDict()

    for n1v in alphabet:
        for n2v in alphabet:
            # nv = n1v + n2v
            nv = (n1v,n2v)
            for o1v in alphabet:
                for o2v in alphabet:
                    # ov = o1v + o2v
                    ov = (o1v,o2v)
                    # k  = (nv,ov) if nv <= ov else (ov,nv)
                    # k  = f"{n1v}|{n2v}|{o1v}|{o2v}" if nv <= ov else f"{o1v}|{o2v}|{n1v}|{n2v}"
                    k = (nv,ov) if nv <= ov else (ov,nv)
                    if k not in diff_matrixSymetricalHomoExtra:
                        val = 0
                        if n1v == o1v:
                            val += 1
                        if n1v == o2v:
                            val += 1
                        if n2v == o1v:
                            val += 1
                        if n2v == o2v:
                            val += 1
                        diff_matrixSymetricalHomoExtra[k] = val
        # print(diff_matrixSymetricalHomoExtra)

    # diff_matrixSymetricalHomoExtra_simple = {}
    # for k1, v1 in diff_matrixSymetricalHomoExtra.items():
    #     for k2, v2 in v1.items():
    #         diff_matrixSymetricalHomoExtra_simple[(k1,k2)] = v2

    return diff_matrixSymetricalHomoExtra

def genIUPAC() -> IUPACType:
    """
        IUPAC nucleotide code	Base
        A	Adenine
        C	Cytosine
        G	Guanine
        T (or U)	Thymine (or Uracil)
        
        M	A or C
        R	A or G
        W	A or T
        S	C or G
        Y	C or T
        K	G or T

        B	C or G or T
        D	A or G or T
        H	A or C or T
        V	A or C or G
        
        N	any base
        
        . or -	gap
    """

    _IUPAC = [
        [('A', 'A'), 'A'],
        [('C', 'C'), 'C'],
        [('G', 'G'), 'G'],
        [('T', 'T'), 'T'],

        [('A', 'C'), 'M'],
        [('A', 'G'), 'R'],
        [('A', 'T'), 'W'],
        [('C', 'G'), 'S'],
        [('C', 'T'), 'Y'],
        [('G', 'T'), 'K'],
    ]

    IUPAC  = {}
    for pair, val in _IUPAC:
        IUPAC[frozenset(pair)] = val

    return IUPAC

def main():
    genome = Genome(
        sys.argv[1],
        bin_width             = DEFAULT_BIN_SIZE,
        type_matrix_counter   = DEFAULT_COUNTER_TYPE_MATRIX,
        type_pairwise_counter = DEFAULT_COUNTER_TYPE_PAIRWISE,
        save_alignment        = False,
        diff_matrix           = genDiffMatrix(alphabet=[0,1,2,3]),
        IUPAC                 = genIUPAC()
    )
    genome.load(threads=6)
    # genome.load(threads=DEFAULT_THREADS if not DEBUG else 1)


if __name__ == "__main__":
    main()
