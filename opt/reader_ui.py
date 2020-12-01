#!/usr/bin/env python3

import os
import sys
import typing

import numpy as np

import flexx
from flexx import flx

import reader

DEBUG       = True

datahandler = None

class Controller(flx.PyComponent):
    # https://flexx.readthedocs.io/en/stable/examples/send_data_src.html

    def init(self):
        pass

class DataHandler():
    def __init__(self, filename):
        self.filename = filename
        self.genome   = None
        self.load_data()

    def load_data(self):
        self.genome = reader.Genome(self.filename)
        self.genome.load()

def main():
    filename    = sys.argv[0]

    global datahandler
    datahandler = DataHandler(filename)
    # https://flexx.readthedocs.io/en/stable/guide/running.html
    # https://flexx.readthedocs.io/en/stable/guide/reactions.html

    flexx.config.hostname           = '0.0.0.0'
    flexx.config.port               = 5000
    flexx.config.log_level          = "DEBUG" if DEBUG else "INFO"
    flexx.config.tornado_debug      = DEBUG
    flexx.config.ws_timeout         = 20
    flexx.config.browser_stacktrace = True
    flexx.config.cookie_secret      = "0123456789"

    app = flx.App(Controller)
    app.serve('')  # Serve at http://domain.com/
    flx.start()  # mainloop will exit when the app is closed

if __name__ == "__main__":
    main()