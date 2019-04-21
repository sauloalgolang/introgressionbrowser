#!/usr/bin/env python3

import os
import sys

import numpy as np
import pandas as pd

np.set_printoptions(edgeitems=3)
np.core.arrayprint._line_width = 80

fname = "res/output_360_merged_2.50.vcf.gz_summary.bin"
fname = "res/output_360_merged_2.50.vcf.gz_chromosomes.bin"

class ReadIbrowserBinary():
    dt0 = np.dtype([
        ('hasData', bool), 
        ('serial', np.int64),
        ('counterBits', np.int64),
        ('dataLen', np.int64),
        ('sumData', np.uint64)
    ])

    def __init__(self, infile):
        self.infile = infile
        self.fhd = None
        self.isOpen = False
        self.fileSize = os.stat(infile).st_size
        
        with open(self.infile, 'rb') as fhd:
            d = np.fromfile(fhd, dtype=self.dt0, count=1)[0]
            
            # print( "d", d )

            self.counterBits = d["counterBits"]
            self.dataLen = d["dataLen"]

            self.dataFmt = None
            self.dataFmtLen = None
            
            if self.counterBits == 16:
                self.dataFmt = np.int16
                self.dataFmtLen = 2
            elif self.counterBits == 32:
                self.dataFmt = np.int32
                self.dataFmtLen = 4
            elif self.counterBits == 64:
                self.dataFmt = np.int64
                self.dataFmtLen = 8
            else:
                print("unknown counter bits", self.counterBits)
                sys.exit(1)

            self.dataSize     = self.dataLen * self.dataFmtLen
            self.registerSize = 1 + 8 + 8 + 8 + 8 + self.dataSize

            self.dt = np.dtype([
                ('hasData'    , bool     ), #1  1
                ('serial'     , np.int64 ), #8  9
                ('counterBits', np.int64 ), #8 17
                ('dataLen'    , np.int64 ), #8 25
                ('sumData'    , np.uint64), #8 33
                ('data', self.dataFmt, self.dataLen)
            ])

            assert (self.fileSize % self.registerSize) == 0, "wrong file size: fileSize {} % registerSize {} != 0. {}".format(self.fileSize, self.registerSize, self.fileSize % self.registerSize)
            
            self.numRegisters = int(self.fileSize / self.registerSize) - 1
            self.offset = 0
            self.memmap = np.memmap(infile, dtype=self.dt, mode='r')

    def open(self):
        if self.isOpen:
            raise Exception("opening already closed file")

        self.fhd = open(self.infile, 'rb')
        self.isOpen = True

    def close(self):
        if not self.isOpen:
            raise Exception("closing closed file")

        self.fhd.close()
        self.isOpen = False

    def reset(self):
        if not self.isOpen:
            raise Exception("resetting closed file")
        
        self.fhd.seek(0)

    def __enter__(self):
        if not self.isOpen:
            self.open()
        
        return self

    def __exit__(self, tp, value, traceback):
        if self.isOpen:
            self.close()

    def __iter__(self):
        if not self.isOpen:
            self.open()

        for i in range(self.numRegisters):
            yield self.getRegister(i)

    def getRegisters(self, regs):
        for regp in regs:
            yield self.getRegister(regp)

    def getRegistersRange(self, regFrom, regTo):
        for regp in range(regFrom, regTo+1):
            yield self.getRegister(regp)

    def _registerLoopToDataframe(self, func, args, kwargs):
        matrix = None

        for reg in func(*args, **kwargs):
            if matrix is None:
                matrix = pd.DataFrame({reg['serial']: reg['data']})
            else:
                matrix[reg['serial']] = reg['data']

        return matrix

    def getRegistersAsDataFrame(self, regs):
        return self._registerLoopToDataframe(self.getRegisters, [regs], {})

    def getRegistersRangeAsDataFrame(self, regFrom, regTo):
        return self._registerLoopToDataframe(self.getRegistersRange, [regFrom, regTo], {})

    def getRegister(self, regnum):
        if not self.isOpen:
            raise Exception("cannot get register in closed file")

        if regnum > self.numRegisters:
            raise Exception("register {} > {} number of registers".format(regnum, self.numRegisters))

        offset = regnum * self.registerSize

        self.fhd.seek(offset)
        
        e = np.fromfile(self.fhd, dtype=self.dt, count=1)
        
        if not e["hasData"]:
            return None
        
        assert len(e) == 1

        r = e[0]

        serial = r['serial']
        r['data'] = np.cumsum(r['data'])
        f = r['data']
        s = np.sum(f)
        
        assert s == r['sumData']

        return r

    def __str__(self):
        return ("infile {infile}"+
                " fileSize {fileSize}"+
                " numRegisters {numRegisters}"+
                " isOpen {isOpen}"+
                " counterBits {counterBits}"+
                " dataLen {dataLen}"+
                " dataFmt {dataFmt}"+
                " dataFmtLen {dataFmtLen}"+
                " dataSize {dataSize}"+
                " registerSize {registerSize}"+
                " offset {offset}"
        ).format(**{
            "infile": self.infile,
            "fileSize": self.fileSize,
            "numRegisters": self.numRegisters,
            "isOpen": self.isOpen,
            "counterBits": self.counterBits,
            "dataLen": self.dataLen,
            "dataFmt": self.dataFmt,
            "dataFmtLen": self.dataFmtLen,
            "dataSize": self.dataSize,
            "registerSize": self.registerSize,
            "dt": self.dt,
            "offset": self.offset
        })
        
    def __repr__(self):
        return str(self)

if __name__ == "__main__":
    with ReadIbrowserBinary(fname) as ibb:
        print("ibb", ibb)

        matrix = ibb.getRegister(0)

        print("matrix.shape", matrix['data'].shape)

    ibb = ReadIbrowserBinary(fname)
    
    print(ibb)

    print("loop 1")
    for reg in ibb:
        print(str(reg)[:50])

    print("loop 2")
    for reg in ibb:
        print(str(reg)[:50])

    print("getRegisters 2,3,5")
    gr1 = ibb.getRegisters([2,3,5])
    print("\n".join([ str(reg)[:50] for reg in gr1 ]))

    print("getRegistersRange 2-2")
    gr2 = ibb.getRegistersRange(2, 2)
    print("\n".join([ str(reg)[:50] for reg in gr2 ]))

    print("getRegistersRange 10-12")
    gr3 = ibb.getRegistersRange(10, 12)
    print("\n".join([ str(reg)[:50] for reg in gr3 ]))

    print("getRegistersAsDataFrame 2,3,5")
    gr4 = ibb.getRegistersAsDataFrame([2,3,5])
    print(gr4)

    print("getRegistersRangeAsDataFrame 10-12")
    gr5 = ibb.getRegistersRangeAsDataFrame(10, 12)
    print(gr5)
    
    print("getRegister 19")
    try:
        gr6 = ibb.getRegister(19)
        print(gr6)
    except:
        print("raised exception. OK")
        pass

    ibb.reset()

    print("memmap")

    for i in range(ibb.numRegisters):
        # reg = ibb.getRegister(i)
        reg = ibb.memmap[i]
        print(str(reg)[:50])
        print("  ", reg['data'][:10])