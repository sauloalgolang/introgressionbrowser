#!/usr/bin/env python3

import os
import sys
import typing

import numpy as np

import flexx
from flexx import flx

import reader

DEBUG       = True


class ChromosomeController(flx.PyComponent):
    def init(self, chromosome: reader.Chromosome):
        self._chromosome = chromosome

    @property
    def vcf_name(self) -> str:
        return self._chromosome.vcf_name

    @property
    def bin_width(self) -> int:
        return self._chromosome.bin_width

    @property
    def chromosome_order(self) -> int:
        return self._chromosome.chromosome_order

    @property
    def chromosome_name(self) -> str:
        return self._chromosome.chromosome_name

    @property
    def matrix_size(self) -> int:
        return self._chromosome.matrix_size

    @property
    def bin_max(self) -> int:
        return self._chromosome.bin_max

    @property
    def bin_min(self) -> int:
        return self._chromosome.bin_min

    @property
    def bin_count(self) -> int:
        return self._chromosome.bin_count

    @property
    def bin_snps_min(self) -> int:
        return self._chromosome.bin_snps_min

    @property
    def bin_snps_max(self) -> int:
        return self._chromosome.bin_snps_max

    @property
    def chromosome_snps(self) -> int:
        return self._chromosome.chromosome_snps

    @property
    def chromosome_first_position(self) -> int:
        return self._chromosome.chromosome_first_position

    @property
    def chromosome_last_position(self) -> int:
        return self._chromosome.chromosome_last_position

    @property
    def sample_names(self) -> typing.List[str]:
        return self._chromosome.sample_names
            
    @property
    def sample_count(self) -> int:
        return self._chromosome.sample_count

    @property
    def matrix(self) -> np.ndarray:
        return self._chromosome.matrix

    def matrix_bin(self, binNum) -> np.ndarray:
        return self._chromosome.matrix_bin(binNum)

    def matrix_bin_matrix(self, binNum) -> np.ndarray:
        return self._chromosome.matrix_bin_matrix(binNum)

    def matrix_bin_matrix_dist(self, binNum, metric='jaccard') -> np.ndarray:
        return self._chromosome.matrix_bin_matrix_dist(binNum, metric=metric)

    def matrix_bin_matrix_dist_sample(self, binNum, sample, metric='jaccard') -> np.ndarray:
        return self._chromosome.matrix_bin_matrix_dist_sample(binNum, sample, metric=metric)

    # self.matrixNp                     : np.ndarray = None
    # self.matrixNp                     = np.zeros((self.bin_max, self.matrix_size ), self.type_matrix_counter  )

    # # self.pairwiNp                   : np.ndarray = None
    # self.pairwiNp                     = np.zeros((self.bin_max, self.sample_count), self.type_pairwise_counter)

    # self.totalsNp                     : np.ndarray = None
    # self.totalsNp                     = np.zeros( self.bin_max,                     self.type_pairwise_counter)

    # self.alignmentNp                  : np.ndarray = None
    # self.alignmentNp                  = np.zeros((self.bin_max, self.sample_count), np.unicode_)

    # self.positionNp                   : np.ndarray = None
    # self.positionNp                   = np.zeros((self.bin_max, position_size    ), self.type_positions)

class GenomeController(flx.PyComponent):
    datahandler : reader.Genome = None
    # https://flexx.readthedocs.io/en/stable/examples/send_data_src.html

    def init(self):
        assert GenomeController.datahandler is not None

    @property
    def vcf_name(self) -> str:
        return GenomeController.datahandler.vcf_name

    @property
    def bin_width(self) -> int:
        return GenomeController.datahandler.bin_width

    @property
    def sample_names(self) -> int:
        return GenomeController.datahandler.sample_names

    @property
    def sample_count(self):
        return GenomeController.datahandler.sample_count

    @property
    def chromosome_names(self):
        return GenomeController.datahandler.chromosome_names

    @property
    def chromosome_count(self):
        return GenomeController.datahandler.chromosome_count

    @property
    def genome_bins(self):
        return GenomeController.datahandler.genome_bins

    @property
    def genome_snps(self):
        return GenomeController.datahandler.genome_snps

    def get_chromosome(self, chromosome_name):
        return ChromosomeController(GenomeController.datahandler.get_chromosome(chromosome_name))


class DataHandler():
    def __init__(self, filename):
        self.filename = filename
        self.genome   = None
        self.load_data()

    def load_data(self):
        self.genome = reader.Genome(self.filename)
        print(self.genome.filename)
        assert self.genome.exists
        print(self.genome)
        self.genome.load()
        assert self.genome.loaded
        assert self.genome.complete


def main():
    filename            = sys.argv[1]

    GenomeController.datahandler = DataHandler(filename)

    # https://flexx.readthedocs.io/en/stable/guide/running.html
    # https://flexx.readthedocs.io/en/stable/guide/reactions.html

    flexx.config.hostname           = '0.0.0.0'
    flexx.config.port               = 5000
    flexx.config.log_level          = "DEBUG" if DEBUG else "INFO"
    flexx.config.tornado_debug      = DEBUG
    flexx.config.ws_timeout         = 20
    flexx.config.browser_stacktrace = True
    flexx.config.cookie_secret      = "0123456789"

    app = flx.App(GenomeController)
    app.serve('')  # Serve at http://domain.com/
    flx.start()  # mainloop will exit when the app is closed

if __name__ == "__main__":
    main()