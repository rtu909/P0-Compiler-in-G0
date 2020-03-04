# P0 Compiler

The original P0 compiler from the course notes, with the following changes:

 * The WASM generation for the `read()` built-in function has been fixed.
 * A new built-in, `print()`, has been added for printing characters directly
   to the output.
 * Only the WASM generator is currently available
 * The file `main.py` has been added. It compiles the input `.p` file to `.wat`,
   assembles with `wat2wasm`, then executes it in the Python WASM VM

