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

DEBUG                         = True
DEFAULT_BIN_SIZE              = 250_000
DEFAULT_COUNTER_TYPE_MATRIX   = np.uint16
DEFAULT_COUNTER_TYPE_PAIRWISE = np.uint32
DEFAULT_THREADS               = mp.cpu_count()

MatrixType           = typing.OrderedDict[typing.Tuple[str,str],int]
ChromosomeMatrixType = typing.OrderedDict[int, typing.List[int]]
BinSnpsType          = typing.OrderedDict[int, int]
BinPairwiseCountType = typing.OrderedDict[int, typing.List[int]]
BinAlignmentType     = typing.OrderedDict[int, typing.List[str]]
BinAlignmentTypeInt  = typing.OrderedDict[int, typing.List[typing.List[str]]]
TriangleIndexType    = typing.OrderedDict[typing.Tuple[int,int], int]
IUPACType            = typing.Dict[typing.FrozenSet[str], str]

"""
real    78m37.755s
user    76m10.670s
sys     00m09.592s

real    219m41.370s
user    217m45.398s
sys     000m15.839s
"""

"""
calculate distance

https://stackoverflow.com/questions/36250729/how-to-convert-triangle-matrix-to-square-in-numpy
"""



class Chromosome():
    def __init__(self,
            vcf_name: str,
            bin_width: int,
            chromosome_order: int,
            chromosome_name: str,
            type_matrix_counter   = DEFAULT_COUNTER_TYPE_MATRIX,
            type_pairwise_counter = DEFAULT_COUNTER_TYPE_PAIRWISE
        ):
        self.type_matrix_counter         = type_matrix_counter
        self.type_pairwise_counter       = type_pairwise_counter
        self.type_matrix_counterMaxVal   = np.iinfo(self.type_matrix_counter).max
        self.type_pairwise_counterMaxVal = np.iinfo(self.type_pairwise_counter).max

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

        self.matrixNp     = None
        self.totalsNp     = None
        self.pairwiNp     = None

    @property
    def filename(self) -> str:
        return f"{self.vcf_name}.{self.bin_width}.{self.chromosome_order:06d}.{self.chromosome_name}.npz"

    @property
    def exists(self) -> bool:
        return os.path.exists(self.filename)

    def addFromVcf(self,
            chromosome_snps          : int,
            chromosome_matrix        : ChromosomeMatrixType,
            matrix_size              : int,
            bin_alignment            : BinAlignmentType,
            bin_snps                 : BinSnpsType,
            bin_pairwise_count       : BinPairwiseCountType,
            chromosome_first_position: int,
            chromosome_last_position : int,
            sample_names             : typing.List[str]
        ):
        print(f"adding chromosome data: {self.chromosome_name}")

        self.matrix_size                   = matrix_size

        bin_names                          = list(chromosome_matrix.keys())
        self.bin_min                       = min(bin_names)
        self.bin_max                       = max(bin_names)
        self.bin_count                     = len(bin_names)

        self.chromosome_snps               = chromosome_snps
        self.sample_count                  = len(sample_names)
        self.sample_names                  = sample_names
        self.chromosome_first_position     = chromosome_first_position
        self.chromosome_last_position      = chromosome_last_position

        self.matrixNp                     = np.zeros((self.bin_max, self.matrix_size), self.type_matrix_counter)
        self.totalsNp                     = np.zeros( self.bin_max, np.uint64)
        self.pairwiNp                     = np.zeros((self.bin_max, self.sample_count), self.type_pairwise_counter)
        self.alignmentNp                  = None
        if bin_alignment:
            self.alignmentNp              = np.zeros((self.bin_max, self.sample_count), np.unicode_)
        else:
            self.alignmentNp              = np.zeros(0, np.unicode_)

        for binNum in range(self.bin_max):
            # print("binNum", binNum)
            chromosome_matrix_bin  = chromosome_matrix .get(binNum, [0 ] * self.matrix_size)
            bin_snp_bin            = bin_snps          .get(binNum, 0)
            bin_pairwise_count_bin = bin_pairwise_count.get(binNum, [0 ] * self.sample_count)
            bin_alignment_bin      = None
            if bin_alignment:
                bin_alignment_bin  = bin_alignment     .get(binNum, [""] * self.sample_count)
            
            # print("bin_pairwise_count_bin", bin_pairwise_count_bin)
            assert not any([v > self.type_matrix_counterMaxVal   for v in chromosome_matrix_bin ]), f"value went over the maximum value ({self.type_matrix_counterMaxVal  }) for container {self.type_matrix_counter  }: {[v for v in chromosome_matrix_bin  if v > self.type_matrix_counterMaxVal  ]}"
            assert not any([v > self.type_pairwise_counterMaxVal for v in bin_pairwise_count_bin]), f"value went over the maximum value ({self.type_pairwise_counterMaxVal}) for container {self.type_pairwise_counter}: {[v for v in bin_pairwise_count_bin if v > self.type_pairwise_counterMaxVal]}"
            
            # binData             = [maxVal if v > maxVal else v for v in binData]
            self.matrixNp[binNum,:]    = chromosome_matrix_bin
            self.totalsNp[binNum]      = bin_snp_bin
            self.pairwiNp[binNum,:]    = bin_pairwise_count_bin
            if bin_alignment:
                self.alignmentNp[binNum,:] = bin_alignment_bin

    def save(self):
        print(f"saving numpy array:           {self.filename}")
        print(f"  vcf_name                    {self.vcf_name}")
        print(f"  chromosome_name             {self.chromosome_name}")
        print(f"  chromosome_snps             {self.chromosome_snps}")
        print(f"  chromosome_order            {self.chromosome_order}")
        print(f"  matrix_size                 {self.matrix_size}")
        print(f"  bin_count                   {self.bin_count}")
        print(f"  bin_min                     {self.bin_min}")
        print(f"  bin_max                     {self.bin_max}")
        print(f"  bin_width                   {self.bin_width}")
        print(f"  sample_count                {self.sample_count}")
        print(f"  chromosome_first_position   {self.chromosome_first_position}")
        print(f"  chromosome_last_position    {self.chromosome_last_position}")
        print(f"  type_matrix_counter         {self.type_matrix_counter}")
        print(f"  type_pairwise_counter       {self.type_pairwise_counter}")
        print(f"  type_matrix_counterMaxVal   {self.type_matrix_counterMaxVal}")
        print(f"  type_pairwise_counterMaxVal {self.type_pairwise_counterMaxVal}")
        print(f"  matrixNp                    {self.matrixNp.shape}")
        print(f"  totalsNp                    {self.totalsNp.shape}")
        print(f"  pairwiNp                    {self.pairwiNp.shape}")
        print(f"  alignmentNp                 {self.alignmentNp.shape}")

        type_matrix_counterName   = self.type_matrix_counter.__name__
        type_pairwise_counterName = self.type_pairwise_counter.__name__

        sample_namesNp           = np.array(self.sample_names, np.unicode_)
        info_namesNp             = np.array(["matrix_size"   , "bin_count"         , "bin_min"                , "bin_max"                   , "bin_width"   , "chromosome_snps"   , "sample_count"   , "chromosome_order"   , "chromosome_first_position"   , "chromosome_last_position"   , "type_matrix_counterMaxVal"   , "type_pairwise_counterMaxVal"   ], np.unicode_)
        info_valuesNp            = np.array([self.matrix_size, self.bin_count      , self.bin_min             , self.bin_max                , self.bin_width, self.chromosome_snps, self.sample_count, self.chromosome_order, self.chromosome_first_position, self.chromosome_last_position, self.type_matrix_counterMaxVal, self.type_pairwise_counterMaxVal], np.int64   )
        meta_namesNp             = np.array(["vcf_name"      , "chromosome_name"   , "type_matrix_counterName", "type_pairwise_counterName"], np.unicode_)
        meta_valuesNp            = np.array([self.vcf_name   , self.chromosome_name, type_matrix_counterName  , type_pairwise_counterName  ], np.unicode_)

        np.savez_compressed(self.filename,
            countMatrix  = self.matrixNp,
            countTotals  = self.totalsNp,
            countPairw   = self.pairwiNp,
            alignments   = self.alignmentNp,
            sample_names = sample_namesNp,
            info_names   = info_namesNp,
            info_values  = info_valuesNp,
            meta_names   = meta_namesNp,
            meta_values  = meta_valuesNp
        )

    def load(self):
        print(f"loading numpy array:          {self.filename}")
        data                             = np.load(self.filename, mmap_mode='r', allow_pickle=False)

        self.matrixNp                    = data['countMatrix']
        self.totalsNp                    = data['countTotals']
        self.pairwiNp                    = data['countPairw']
        self.alignmentNp                 = data['alignments']
        
        sample_namesNp                   = data['sample_names']
        
        info_namesNp                     = data['info_names']
        info_valuesNp                    = data['info_values']

        meta_namesNp                     = data['meta_names']
        meta_valuesNp                    = data['meta_values']

        sample_names                     = sample_namesNp.tolist()

        info_names                       = info_namesNp.tolist()
        info_values                      = info_valuesNp.tolist()
        info_values                      = [int(v) for v in info_values]
        info_dict                        = {info_names[k]: info_values[k]  for k in range(len(info_names))}
        # print(info_dict)

        meta_names                       = meta_namesNp.tolist()
        meta_values                      = meta_valuesNp.tolist()
        meta_dict                        = {meta_names[k]: meta_values[k]  for k in range(len(meta_names))}
        # print(meta_dict)

        self.vcf_name                    = meta_dict["vcf_name"]
        self.type_matrix_counter         = self.matrixNp.dtype

        self.matrix_size                 = info_dict["matrix_size"]
        assert self.matrix_size == self.matrixNp.shape[1]

        self.bin_count                   = info_dict["bin_count"]
        self.bin_min                     = info_dict["bin_min"]
        self.bin_max                     = info_dict["bin_max"]
        assert self.bin_max == self.matrixNp.shape[0]
        assert self.bin_max == self.pairwiNp.shape[0]

        self.bin_width                   = info_dict["bin_width"]

        self.chromosome_snps             = info_dict["chromosome_snps"]
        self.chromosome_name             = meta_dict["chromosome_name"]
        self.chromosome_order            = info_dict["chromosome_order"]

        self.chromosome_first_position   = info_dict["chromosome_first_position"]
        self.chromosome_last_position    = info_dict["chromosome_last_position"]
        assert (self.chromosome_first_position // self.bin_width) == self.bin_min
        assert (self.chromosome_last_position  // self.bin_width) == self.bin_max

        type_matrix_counterName          = meta_dict["type_matrix_counterName"]
        type_pairwise_counterName        = meta_dict["type_pairwise_counterName"]        
        self.type_matrix_counter         = getattr(np, type_matrix_counterName)
        self.type_pairwise_counter       = getattr(np, type_pairwise_counterName)

        self.type_matrix_counterMaxVal   = info_dict["type_matrix_counterMaxVal"]
        self.type_pairwise_counterMaxVal = info_dict["type_pairwise_counterMaxVal"]
        assert np.iinfo(self.type_matrix_counter).max   == self.type_matrix_counterMaxVal
        assert np.iinfo(self.type_pairwise_counter).max == self.type_pairwise_counterMaxVal

        self.sample_count                = info_dict["sample_count"]
        self.sample_names                = sample_names
        assert len(self.sample_names) == self.sample_count
        assert self.sample_count == self.pairwiNp.shape[1]
        if self.alignmentNp.shape[0] != 0:
            assert self.sample_count == self.alignmentNp.shape[1]

        print(f"  vcf_name                    {self.vcf_name}")
        print(f"  chromosome_name             {self.chromosome_name}")
        print(f"  chromosome_snps             {self.chromosome_snps}")
        print(f"  chromosome_order            {self.chromosome_order}")
        print(f"  matrix_size                 {self.matrix_size}")
        print(f"  bin_count                   {self.bin_count}")
        print(f"  bin_min                     {self.bin_min}")
        print(f"  bin_max                     {self.bin_max}")
        print(f"  bin_width                   {self.bin_width}")
        print(f"  sample_count                {self.sample_count}")
        print(f"  chromosome_first_position   {self.chromosome_first_position}")
        print(f"  chromosome_last_position    {self.chromosome_last_position}")
        print(f"  type_matrix_counter         {self.type_matrix_counter}")
        print(f"  type_pairwise_counter       {self.type_pairwise_counter}")
        print(f"  type_matrix_counterMaxVal   {self.type_matrix_counterMaxVal}")
        print(f"  type_pairwise_counterMaxVal {self.type_pairwise_counterMaxVal}")
        print(f"  matrixNp                    {self.matrixNp.shape}")
        print(f"  totalsNp                    {self.totalsNp.shape}")
        print(f"  pairwiNp                    {self.pairwiNp.shape}")
        print(f"  alignmentNp                 {self.alignmentNp.shape}")


class Genome():
    def __init__(self,
            vcf_name: str,
            bin_width: int        = DEFAULT_BIN_SIZE,
            type_matrix_counter   = DEFAULT_COUNTER_TYPE_MATRIX,
            type_pairwise_counter = DEFAULT_COUNTER_TYPE_PAIRWISE,
            save_alignment       : bool =False,
            diff_matrix: MatrixType=None,
            IUPAC: IUPACType=None
        ):
        self.vcf_name             : str = vcf_name
        self.bin_width            : int = bin_width
        self.type_matrix_counter        = type_matrix_counter
        self.type_pairwise_counter      = type_pairwise_counter

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
                type_pairwise_counter = self.type_pairwise_counter
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
                    matrix_size  = calculateMatrixSize(sample_count)
                    indexes      = triangleToIndex(sample_count)

                    print("sample_names", sample_names)
                    print("num samples ", sample_count)
                    print("matrix_size ", matrix_size)
                    # print("indexes    ", indexes)
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

            save_alignment  : bool,
            IUPAC           : IUPACType,
            diff_matrix     : MatrixType
        ) -> typing.Tuple[int, int]:


        bgzip = BGzip(vcf_name)

        bin_snps                 : BinSnpsType = OrderedDict()
        bin_pairwise_count       : BinPairwiseCountType = OrderedDict()
        bin_alignment            : BinAlignmentTypeInt = OrderedDict()
        chromosome_matrix        : MatrixType = OrderedDict()
        chromosome_first_position: int = 0
        chromosome_last_position : int = 0
        chromosome_snps          : int = 0

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

            if DEBUG:
                if len(chromosome_matrix) > 3:
                    break

            if binNum not in chromosome_matrix:
                print(f"  New bin: {chrom} :: {pos:12,d} => {binNum:12,d}")
                chromosome_matrix[binNum]  = [0] * matrix_size
                bin_snps[binNum]           = 0
                bin_pairwise_count[binNum] = [0] * sample_count
                bin_alignment[binNum]      = None
                if save_alignment:
                    bin_alignment[binNum]  = [[] for _ in range(sample_count)]

            samples                   = [s.split(";")[0]            for s in samples]
            samples                   = [s if len(s) == 3 else None for s in samples]
            # samples = [s.replace("|", "").replace("/", "") if s is not None else None for s in samples]
            samples                   = [tuple([int(i) for i in s.replace("/", "|").split("|")]) if s is not None else None for s in samples]
            vals                      = chromosome_matrix[binNum]
            paiw                      = bin_pairwise_count[binNum]
            aling                     = None
            if save_alignment:
                aling                 = bin_alignment[binNum]
            chromosome_snps          += 1
            bin_snps[binNum]         += 1
            chromosome_last_position  = pos
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

                    pairind           = indexes[(sample2num,sample1num)]
                    vals[pairind]    += value
                    paiw[sample1num] += value
                    paiw[sample2num] += value

        Genome._processVcf_save_chrom_data(
            vcf_name                     = vcf_name,
            bin_width                    = bin_width,
            sample_names                 = sample_names,
            sample_count                 = sample_count,

            chromosome_name              = chromosome_name,
            chromosome_order             = chromosome_order,
            chromosome_snps              = chromosome_snps,
            chromosome_matrix            = chromosome_matrix,
            chromosome_first_position    = chromosome_first_position,
            chromosome_last_position     = chromosome_last_position,
            bin_alignment                = bin_alignment if save_alignment else None,
            bin_snps                     = bin_snps,
            bin_pairwise_count           = bin_pairwise_count,
            matrix_size                  = matrix_size,

            type_matrix_counter          = type_matrix_counter,
            type_pairwise_counter        = type_pairwise_counter
        )

        chromosome_bins = len(chromosome_matrix)

        return chromosome_bins, chromosome_snps

    @staticmethod
    def _processVcf_save_chrom_data(
            vcf_name                 : str,
            bin_width                : int,
            sample_names             : typing.List[str],
            sample_count             : int,
            chromosome_name          : str,
            chromosome_order         : int,
            chromosome_snps          : int,
            chromosome_matrix        : ChromosomeMatrixType,
            chromosome_first_position: int,
            chromosome_last_position : int,
            bin_alignment            : BinAlignmentTypeInt,
            bin_snps                 : BinSnpsType,
            bin_pairwise_count       : BinPairwiseCountType, 
            matrix_size              : int,
            type_matrix_counter,
            type_pairwise_counter
        ):

        if bin_alignment:
            for binNum, samples in bin_alignment.items():
                # print("binNum", binNum)
                for sample_num in range(len(samples)):
                    samples[sample_num] = "".join(samples[sample_num])
                    # print(" sample_num", sample_num, samples[sample_num])
            # bin_alignment[binNum][sample1num].append(nuc)

        # self.chromosome_names.append(chromosome_name)
        # self.chromosome_count += 1

        chromosome = Chromosome(
            vcf_name              = vcf_name,
            bin_width             = bin_width,
            chromosome_order      = chromosome_order,
            chromosome_name       = chromosome_name,
            type_matrix_counter   = type_matrix_counter,
            type_pairwise_counter = type_pairwise_counter
        )

        chromosome.addFromVcf(
            chromosome_snps           = chromosome_snps,
            chromosome_matrix         = chromosome_matrix,
            matrix_size               = matrix_size,
            bin_alignment             = bin_alignment,
            bin_snps                  = bin_snps,
            bin_pairwise_count        = bin_pairwise_count,
            chromosome_first_position = chromosome_first_position,
            chromosome_last_position  = chromosome_last_position,
            sample_names              = sample_names
        )

        chromosome.save()

    def _load_db(self):
        if self.complete:
            return

        print(f"loading numpy array:          {self.filename}")
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
        type_matrix_counterName    = meta_dict["type_matrix_counterName"]
        type_pairwise_counterName  = meta_dict["type_pairwise_counterName"]        
        self.type_matrix_counter   = getattr(np, type_matrix_counterName)
        self.type_pairwise_counter = getattr(np, type_pairwise_counterName)

        self.bin_width             = info_dict["bin_width"]
        self.chromosome_count      = info_dict["chromosome_count"]
        self.sample_count          = info_dict["sample_count"]
        self.genome_bins           = info_dict["genome_bins"]
        self.genome_snps           = info_dict["genome_snps"]

        assert self.chromosome_count == len(self.chromosome_names)
        assert self.sample_count     == len(self.sample_names)

        print(f"  vcf_name                    {self.vcf_name}")
        print(f"  bin_width                   {self.bin_width}")
        print(f"  chromosome_count            {self.chromosome_count}")
        print(f"  sample_count                {self.sample_count}")
        print(f"  type_matrix_counter         {self.type_matrix_counter}")
        print(f"  type_pairwise_counter       {self.type_pairwise_counter}")
        print(f"  genome_bins                 {self.genome_bins}")
        print(f"  genome_snps                 {self.genome_snps}")

        chromosome_snps = 0
        chromosome_bins = 0
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
            
            chromosome_snps += chromosome.chromosome_snps
            chromosome_bins += chromosome.bin_count
            
            self._chromosomes.append(chromosome)
        
        assert self.genome_bins == chromosome_bins
        assert self.genome_snps == chromosome_snps

    def save(self):
        print(f"saving numpy array:           {self.filename}")
        print(f"  vcf_name                    {self.vcf_name}")
        print(f"  bin_width                   {self.bin_width}")
        print(f"  chromosome_count            {self.chromosome_count}")
        print(f"  sample_count                {self.sample_count}")
        print(f"  type_matrix_counter         {self.type_matrix_counter}")
        print(f"  type_pairwise_counter       {self.type_pairwise_counter}")
        print(f"  genome_bins                 {self.genome_bins}")
        print(f"  genome_snps                 {self.genome_snps}")

        type_matrix_counterName   = self.type_matrix_counter.__name__
        type_pairwise_counterName = self.type_pairwise_counter.__name__

        sample_namesNp            = np.array(self.sample_names    , np.unicode_)
        chromosome_namesNp        = np.array(self.chromosome_names, np.unicode_)

        info_namesNp              = np.array(["bin_width"   , "chromosome_count"   , "sample_count"   , "genome_bins"   , "genome_snps"   ], np.unicode_)
        info_valuesNp             = np.array([self.bin_width, self.chromosome_count, self.sample_count, self.genome_bins, self.genome_snps], np.int64   )

        meta_namesNp              = np.array(["vcf_name"   , "type_matrix_counterName", "type_pairwise_counterName"], np.unicode_)
        meta_valuesNp             = np.array([self.vcf_name, type_matrix_counterName  , type_pairwise_counterName  ], np.unicode_)

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

            current_compressed_offset  : int,
            current_uncompressed_offset: int,
            current_compressed_size    : int,
            current_uncompressed_size  : int,
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
            print(f"NEW CHROMOSOME chrom_name       {chrom_name:>15s}")
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

    @property
    def chromosomes(self) -> typing.List[str]:
        chroms = [k for k in self._data.keys()]
        chroms.sort(key=lambda x: self._data[x]["entry_num"])
        return chroms

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
                        print('WRONG CHROMOSOME TAIL', chrom, line)
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

def triangleToIndex(size: int) -> TriangleIndexType:
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
    # b = np.tril_indices(size)
    # (array([0, 1, 1, 2, 2, 2, 3, 3, 3, 3]), array([0, 0, 1, 0, 1, 2, 0, 1, 2, 3]))
    # >>> b = np.tril_indices(4, -1)
    # (array([1, 2, 2, 3, 3, 3]), array([0, 0, 1, 0, 1, 2]))
    # >>> a = np.arange(16).reshape(4, 4)
    # >>> a[b]
    # array([ 4,  8,  9, 12, 13, 14])
    # >>> bp = list(zip(b[0].tolist(), b[1].tolist()))
    # >>> bp
    # [(1, 0), (2, 0), (2, 1), (3, 0), (3, 1), (3, 2)]

    indexes = OrderedDict()
    b       = np.tril_indices(size, -1)
    bp      = list(zip(b[0].tolist(), b[1].tolist()))
    # print("b", b)
    # print("bp", bp)
    
    for bl, bc in enumerate(bp):
        # print("index", bc, bl)
        indexes[bc] = bl

    # for bl in range(size):
    #     indexes[(bl, bl)] = len(bp) + bl
    
    return indexes

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
    genome.load(threads=1 if DEBUG else DEFAULT_THREADS)


if __name__ == "__main__":
    main()
