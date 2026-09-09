<div align="center" style="margin-bottom: 15px;">
  <img src="./assets/golang_QIF(2).png" width="400" height="auto">
</div>

# QPACK Interop Format

![GitHub repo size](https://img.shields.io/github/repo-size/NguyenHien-8/qpack-interop-format)
[![GitHub](https://img.shields.io/github/license/NguyenHien-8/qpack-interop-format)](https://github.com/NguyenHien-8/qpack-interop-format/blob/master/LICENSE)

This repository is for storing QIFs used in QPACK development and interop experiments.

## Encoder Input

A QPACK encoder produces a file whose format is described in [QPACK Offline Interop](https://github.com/quicwg/base-drafts/wiki/QPACK-Offline-Interop). Outputs produced by various encoders are stored in `encoded/` directory.

The encoder output server as the decoder input.  The idea of the interop is to take a QIF file, encode it using an encoder from one distribution and decode it using a decoder from another distribution.

## Converting HARs to QIFs

A HAR is easy to manufacture using any of the major browsers. On the other hand, a HAR is inconvenient for the kind of end-to-end testing that the interop suggests: they are difficult to generate and compare.  A HAR can be converted to a QIF format using `bin/har2qif.pl` program.

## Comparing QIFs

QIFs are text files and can be compared using standard UNIX tools such as `diff(1)`.

Because both input and output QIF files can contain comments and, more importantly, because the decoding process may output header lists in a different order, `bin/sort-qif.pl` is provided.  Sample use:

```bash
sh$ encode source.qif > encoded.bin
sh$ decode encoded.bin > result.qif
sh$ grep -v ^# source.qif > canonical-input.qif
sh$ sort-qif.pl --strip-comments result.qif > canonical-output.qif
sh$ diff canonical-input.qif canonical-output.qif
```