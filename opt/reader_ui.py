#!/usr/bin/env python3

import os
import sys
import typing

import numpy as np

import flexx
from   flexx import flx, ui

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



class SelectorController(flx.Widget):
    CSS = """
    .flx-ComboBox {
        background: #9d9;
    }
    .flx-LineEdit {
        border: 2px solid #9d9;
    }
    .flx-ComboBox.hidden {
        visibility: hidden;
    }
    //.flx-ComboBox.combo_sel {
    //    width: 250px !important;
    //    height: 50px !important;
    //}
    """

    genome_name      = flx.StringProp("", settable=True, doc="current genome")
    bin_width        = flx.IntProp(   -1, settable=True, doc="current bin width")
    metric           = flx.StringProp("", settable=True, doc="current metric")
    chromosome_name  = flx.StringProp("", settable=True, doc="current chromosome")

    _combo_genome_names_text     = "Select genome"
    _combo_bin_widths_text       = "Select bin width"
    _combo_metrics_text          = "Select distance metric"
    _combo_chromosome_names_text = "Select chromosome"

    def init(self, mainController: "MainController"):
        self._is_loaded = False

        super().init()

        self.mainController         = mainController

        with flx.VBox():
            with flx.HFix(flex=1) as self.header:
                self.combo_genome_names     = ui.ComboBox(text=self._combo_genome_names_text    , editable=False, options=[], flex=1, css_class="combo_sel")
                self.combo_bin_widths       = ui.ComboBox(text=self._combo_bin_widths_text      , editable=False, options=[], flex=1, css_class="combo_sel")
                self.combo_metrics          = ui.ComboBox(text=self._combo_metrics_text         , editable=False, options=[], flex=1, css_class="combo_sel")
                self.combo_chromosome_names = ui.ComboBox(text=self._combo_chromosome_names_text, editable=False, options=[], flex=1, css_class="combo_sel")
            with flx.HFix(flex=19) as self.body:
                ui.Label(text="label")

        self._is_loaded = True

        self.setHidden(self.combo_bin_widths      , True)
        self.setHidden(self.combo_metrics         , True)
        self.setHidden(self.combo_chromosome_names, True)

        self.mainController.update_genome_names()

    def setHidden(self, box, status):
        if " hidden" in box.css_class:
            if status:
                pass
            else:
                box.set_css_class(box.css_class.replace(" hidden", ""))

        else:
            if status:
                box.set_css_class(box.css_class + " hidden")
            else:
                pass



    @flx.action
    def reset_genome_names(self):
        self.update_genome_names([])
        self.set_genome_names("")
        self.combo_genome_names.set_text(self._combo_genome_names_text)
        self.setHidden(self.combo_genome_names, True)
        self.reset_bin_widths()

    @flx.action
    def reset_bin_widths(self):
        self.update_bin_widths([])
        self.set_bin_width(-1)
        self.combo_bin_widths.set_text(self._combo_bin_widths_text)
        self.setHidden(self.combo_bin_widths, True)
        self.reset_metrics()

    @flx.action
    def reset_metrics(self):
        self.update_metrics([])
        self.set_metric("")
        self.combo_metrics.set_text(self._combo_metrics_text)
        self.setHidden(self.combo_metrics, True)
        self.reset_chromosome_names()

    @flx.action
    def reset_chromosome_names(self):
        self.update_chromosome_names([])
        self.set_chromosome_name("")
        self.combo_chromosome_names.set_text(self._combo_chromosome_names_text)
        self.setHidden(self.combo_chromosome_names, True)



    @flx.action
    def set_genome_name(self, genome_name: str):
        if self._is_loaded and genome_name != "":
            print("SelectorController.set_genome_name", genome_name)
            self._mutate_genome_name(genome_name)
            self.reset_bin_widths()
            self.mainController.update_bin_widths(self.genome_name)
            self.setHidden(self.combo_bin_widths, False)

    @flx.action
    def set_bin_width(self, bin_width: int):
        if self._is_loaded and bin_width != -1 and bin_width != "":
            print("SelectorController.set_bin_width", bin_width)
            self._mutate_bin_width(bin_width)
            self.reset_metrics()
            self.mainController.update_metrics(self.genome_name, self.bin_width)
            self.setHidden(self.combo_metrics, False)

    @flx.action
    def set_metric(self, metric: str):
        if self._is_loaded and metric != "":
            print("SelectorController.set_metric", metric)
            self._mutate_metric(metric)
            self.reset_chromosome_names()
            self.mainController.update_chromosome_names(self.genome_name, self.bin_width, self.metric)
            self.setHidden(self.combo_chromosome_names, False)
            self.load_genome()

    @flx.action
    def set_chromosome_name(self, chromosome_name: str):
        if self._is_loaded and chromosome_name != "":
            print("SelectorController.set_chromosome_name", chromosome_name)
            self._mutate_chromosome_name(chromosome_name)
            self.load_chromosome()



    def load_genome(self):
        if (
            self._is_loaded            and
            self.genome_name     != "" and
            self.bin_width       != -1 and self.bin_width != "" and
            self.metric          != ""
        ):
            print("loading genome", self.genome_name, self.bin_width, self.metric)

            self.mainController.update_genome(self.genome_name, self.bin_width, self.metric)

    def load_chromosome(self):
        if (
            self._is_loaded and
            self.genome_name     != "" and
            self.bin_width       != -1 and self.bin_width != "" and
            self.metric          != "" and
            self.chromosome_name != ""
        ):
            print("loading chromosome", self.genome_name, self.bin_width, self.metric, self.chromosome_name)

            self.mainController.update_chromosome(self.genome_name, self.bin_width, self.metric, self.chromosome_name)



    @flx.action
    def update_genome(self, genome_name: str, bin_width: int, metric: str, genome: GenomeController):
        assert self.genome_name == genome_name
        assert self.bin_width   == bin_width
        assert self.metric      == metric
        print("SelectorController.update_genome", genome_name, bin_width, metric)

    @flx.action
    def update_chromosome(self, genome_name: str, bin_width: int, metric: str, chromosome_name: str, chromosomeL: ChromosomeController):
        assert self.genome_name     == genome_name
        assert self.bin_width       == bin_width
        assert self.metric          == metric
        assert self.chromosome_name == chromosome_name
        print("SelectorController.update_chromosome", genome_name, bin_width, metric, chromosome_name)



    @flx.action
    def update_genome_names(self, genome_names: typing.List[str]):
        print("SelectorController.update_genome_names", len(genome_names), genome_names)
        self.combo_genome_names.set_options(genome_names)

    @flx.action
    def update_bin_widths(self, bin_widths: typing.List[int]):
        print("SelectorController.update_bin_widths", len(bin_widths), bin_widths)
        self.combo_bin_widths.set_options(bin_widths)

    @flx.action
    def update_metrics(self, metrics: typing.List[str]):
        print("SelectorController.update_metrics", len(metrics), metrics)
        self.combo_metrics.set_options(metrics)

    @flx.action
    def update_chromosome_names(self, chromosome_names: typing.List[str]):
        print("SelectorController.update_chromosome_names", len(chromosome_names), chromosome_names)
        self.combo_chromosome_names.set_options(chromosome_names)



    #https://flexx.readthedocs.io/en/v0.8.0/ui/dropdown.html?highlight=dropdown
    #https://flexx.readthedocs.io/en/v0.8.0/guide/reactions.html?highlight=reaction
    @flx.reaction("combo_genome_names.selected_key")
    def reaction_combo_genome_names(self):
        self.set_genome_name(self.combo_genome_names.selected_key)

    @flx.reaction("combo_bin_widths.selected_key")
    def reaction_combo_bin_widths(self):
        self.set_bin_width(self.combo_bin_widths.selected_key)

    @flx.reaction("combo_metrics.selected_key")
    def reaction_combo_metrics(self):
        self.set_metric(self.combo_metrics.selected_key)

    @flx.reaction("combo_chromosome_names.selected_key")
    def reaction_combo_chromosome_names(self):
        self.set_chromosome_name(self.combo_chromosome_names.selected_key)



class MainController(flx.PyComponent):
    _genomes_cls : reader.Genomes = None
    # https://flexx.readthedocs.io/en/stable/examples/send_data_src.html

    def init(self):
        self._genomes = MainController._genomes_cls
        self.update(verbose=True)
        self.selector = SelectorController(self)

    @flx.action
    def update_genome_names(self):
        print(f"MainController.update_genome_names {self.genomes()}")
        self.selector.update_genome_names(self.genomes())

    @flx.action
    def update_bin_widths(self, genome_name: str):
        bin_widths = self.bin_widths(genome_name)
        print(f"MainController.update_bin_widths {bin_widths}")
        self.selector.update_bin_widths(bin_widths)

    @flx.action
    def update_metrics(self, genome_name: str, bin_width: int):
        metrics = self.metrics(genome_name, bin_width)
        print(f"MainController.update_metrics {metrics}")
        self.selector.update_metrics(metrics)

    @flx.action
    def update_chromosome_names(self, genome_name: str, bin_width: int, metric: str):
        chromosome_names = self.chromosome_names(genome_name, bin_width, metric)
        print(f"MainController.update_chromosome_names {chromosome_names}")
        self.selector.update_chromosome_names([x[1] for x in chromosome_names])

    @flx.action
    def update_genome(self, genome_name: str, bin_width: int, metric: str):
        print(f"MainController.update_genome genome_name {genome_name} bin_width {bin_width} metric {metric}")
        genome  = self.load_genome(self, genome_name, bin_width, metric)
        genomec = GenomeController(genome)
        self.selector.update_genome(genome_name, bin_width, metric, genomec)

    @flx.action
    def update_chromosome(self, genome_name: str, bin_width: int, metric: str, chromosome_name: str):
        print(f"MainController.update_chromosome genome_name {genome_name} bin_width {bin_width} metric {metric} chromosome_name {chromosome_name}")
        chromosome  = self.load_chromosome(self, genome_name, bin_width, metric, chromosome_name)
        chromosomec = ChromosomeController(chromosome)
        self.selector.update_chromosome(genome_name, bin_width, metric, chromosome_name, chromosomec)



    # def genome(self) -> GenomeController:
    #     return GenomeController(self._genomes.genome)

    # def chromosomes(self) -> reader.ChromosomeNamesType:
    #     return self._genomes.chromosomes

    # def chromosome(self) -> ChromosomeController:
    #     return ChromosomeController(self._genomes.chromosome)

    def genomes(self):
        return self._genomes.genomes

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

    def chromosome_names(self, genome_name: str, bin_width: int, metric: str) -> typing.List[typing.Tuple[int, str]]:
        return self._genomes.chromosome_names(genome_name, bin_width, metric)

    def update(self, verbose=False):
        self._genomes.update(verbose=verbose)

    def load_genome(self, genome_name: str, bin_width: int, metric: str) -> reader.Genome:
        return self._genomes.load_genome(genome_name, bin_width, metric)

    def load_chromosome(self, genome_name: str, bin_width: int, metric: str, chromosome_name: str) -> reader.Chromosome:
        return self._genomes.load_chromosome(genome_name, bin_width, metric, chromosome_name)



def main():
    folder_name = sys.argv[1]

    MainController._genomes_cls     = reader.Genomes(folder_name)

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