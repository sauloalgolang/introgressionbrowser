#!/usr/bin/env python3

import os
import sys
import typing

import numpy as np

import flexx
from   flexx import flx

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
    def metric(self) -> str:
        return self._chromosome.metric



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
    def sample_names(self) -> reader.SampleNamesType:
        return self._chromosome.sample_names
            
    @property
    def sample_count(self) -> int:
        return self._chromosome.sample_count


    @property
    def filename(self) -> str:
        return self._chromosome.filename

    @property
    def matrix(self) -> np.ndarray:
        return self._chromosome.matrix

    @property
    def matrix_dtype(self) -> np.dtype:
        return self._chromosome.matrix_dtype



    def matrix_sample(self, sample_name) -> np.ndarray:
        return self._chromosome.matrix_sample(sample_name, metric=self.metric)

    def matrix_bin(self, binNum: int) -> np.ndarray:
        return self._chromosome.matrix_bin(binNum)

    def matrix_bin_square(self, binNum: int) -> np.ndarray:
        return self._chromosome.matrix_bin_square(binNum)

    def matrix_bin_sample(self, binNum: int, sample_name: str) -> np.ndarray:
        return self._chromosome.matrix_bin_sample(binNum, sample_name)

    def matrix_bin_dist(self, binNum: int) -> np.ndarray:
        return self._chromosome.matrix_bin_dist(binNum, metric=self.metric, dtype=self.matrix_dtype)

    def matrix_bin_dist_square(self, binNum: int) -> np.ndarray:
        return self._chromosome.matrix_bin_dist_square(binNum, self.metric)

    def matrix_bin_dist_sample(self, binNum: int, sample_name: str) -> np.ndarray:
        return self._chromosome.matrix_bin_dist_sample(binNum, sample_name, metric=self.metric)

    def matrix_bin_dist_sample_square(self, binNum: int, sample_name: str) -> np.ndarray:
        return self._chromosome.matrix_bin_dist_sample_square(binNum, sample_name, metric=self.metric)


    # self.matrixNp                     : np.ndarray = None
    # self.matrixNp                     = np.zeros((self.bin_max, self.matrix_size ), self.type_matrix_counter  )

    # self.pairwiNp                     : np.ndarray = None
    # self.pairwiNp                     = np.zeros((self.bin_max, self.sample_count), self.type_pairwise_counter)

    # self.binsnpNp                     : np.ndarray = None
    # self.binsnpNp                     = np.zeros( self.bin_max,                     self.type_pairwise_counter)

    # self.alignmentNp                  : np.ndarray = None
    # self.alignmentNp                  = np.zeros((self.bin_max, self.sample_count), np.unicode_)

    # self.positionNp                   : np.ndarray = None
    # self.positionNp                   = np.zeros((self.bin_max, position_size    ), self.type_positions)


class GenomeController(flx.PyComponent):
    def init(self, genome: reader.Genome):
        self._genome = genome

    @property
    def vcf_name(self) -> str:
        return self._genome.vcf_name

    @property
    def bin_width(self) -> int:
        return self._genome.bin_width

    @property
    def metric(self) -> str:
        return self._genome.metric

    @property
    def sample_names(self) -> reader.SampleNamesType:
        return self._genome.sample_names

    @property
    def sample_count(self) -> int:
        return self._genome.sample_count

    @property
    def chromosome_names(self) -> reader.ChromosomeNamesType:
        return self._genome.chromosome_names

    @property
    def chromosome_count(self) -> int:
        return self._genome.chromosome_count

    @property
    def genome_bins(self) -> int:
        return self._genome.genome_bins

    @property
    def genome_snps(self) -> int:
        return self._genome.genome_snps

    @property
    def filename(self) -> str:
        return self._genome.filename

    def get_chromosome(self, chromosome_name: str):
        return ChromosomeController(self._genome.get_chromosome(chromosome_name))


class MainController(flx.PyComponent):
    genomes : reader.Genomes = None
    # https://flexx.readthedocs.io/en/stable/examples/send_data_src.html

    def init(self):
        self._genomes = MainController.genomes
        self.update()

    @property
    def genomes(self) -> typing.List[str]:
        return self._genomes.genomes

    @property
    def genome(self) -> GenomeController:
        return GenomeController(self._genomes.genome)

    @property
    def chromosomes(self) -> reader.ChromosomeNamesType:
        return self._genomes.chromosomes

    @property
    def chromosome(self) -> ChromosomeController:
        return ChromosomeController(self._genomes.chromosome)

    def genome_info(self, genome_name: str):
        return self._genomes.genome_info(genome_name)

    def bin_widths(self, genome_name: str) -> typing.List[str]:
        return self._genomes.bin_widths(genome_name)

    def bin_width_info(self, genome_name: str, bin_width: int):
        return self._genomes.bin_width_info(genome_name, bin_width)

    def metrics(self, genome_name: str, bin_width: int) -> typing.List[str]:
        return self._genomes.metrics(genome_name, bin_width)

    def metric_info(self, genome_name: str, bin_width: int, metric: str) -> typing.List[str]:
        return self._genomes.metric_info(genome_name, bin_width, metric)

    def update(self):
        self._genomes.update()

    def load_genome(self, genome_name: str, bin_width: int, metric: str) -> GenomeController:
        return self._genomes.load_genome(genome_name, bin_width, metric)

    def load_chromosome(self, genome_name: str, bin_width: int, metric: str, chromosome_name: str) -> ChromosomeController:
        return self._genomes.load_genome(genome_name, bin_width, metric, chromosome_name)



def main():
    folder_name = sys.argv[1]

    MainController.genomes = reader.Genomes(folder_name)

    # https://flexx.readthedocs.io/en/stable/guide/running.html
    # https://flexx.readthedocs.io/en/stable/guide/reactions.html

    flexx.config.hostname           = '0.0.0.0'
    flexx.config.port               = 5000
    flexx.config.log_level          = "DEBUG" if DEBUG else "INFO"
    flexx.config.tornado_debug      = DEBUG
    flexx.config.ws_timeout         = 20
    flexx.config.browser_stacktrace = True
    flexx.config.cookie_secret      = "0123456789"

    app = flx.App(MainController)
    app.serve('')  # Serve at http://domain.com/
    flx.start()  # mainloop will exit when the app is closed

if __name__ == "__main__":
    main()