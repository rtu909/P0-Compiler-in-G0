import sys
import pywasm
import subprocess
from parser import compileFile

# Get and validate input
file_name = sys.argv[1]
assert file_name[-2:] == ".p" and file_name[:-2]

# # Compile the file
compileFile(file_name)

