import argparse
from pathlib import Path

from .convert import convert

parser = argparse.ArgumentParser()
parser.add_argument("input", type=Path, help="input Excel file")
parser.add_argument("output", type=Path, help="output directory")
args = parser.parse_args()
convert(args.input, args.output)
