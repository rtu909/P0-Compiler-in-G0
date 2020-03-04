import sys
import pywasm
import subprocess
from parser import compileFile

# Get and validate input
file_name = sys.argv[1]
assert file_name[-2:] == ".p" and file_name[:-2]

# # Compile the file
compileFile(file_name)

# # Assemble the file
subprocess.run(["mv", file_name[:-2] + ".s", file_name[:-2] + ".wat"])
subprocess.run(["wat2wasm", file_name[:-2] + ".wat"])

def rite(i):
    print(i)

def writeln():
    print()

def reed():
    return int(input())

def printify(c):
    print(chr(c), end='')

# Load the assembled binary into a vm to run
# Also load in references to the functions to run
vm = pywasm.load(file_name[:-2] + ".wasm", {'P0lib': {'write': rite, 'writeln': writeln, 'read': reed, 'print': printify}})
