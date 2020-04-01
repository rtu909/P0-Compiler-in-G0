import time
import os
import parser

print("\tPYTHON\tGO SEQUENTIAL\tGO CONCURRENT")

# Increase the size of the program from 100 to 1000 lines of code
for size in range(1, 11):
    with open("temp.p", "w") as f:
        f.write("""program fib;\n{ Used to store the user's input}\n  var x: integer;\n""")
        for redef in range (1, 10 * size + 1):
            f.write("""
  { \\brief Calculates the fibonnaci series at a given value and prints it
    \\param n The value of the series to be calculated }
  procedure fib""" + str(redef) + """(n: integer);
    var prev: integer;
    var current: integer;
    var next: integer;
    var i: integer;
    begin
      prev := 0;
      current := 1;
      next := prev + current;
      i := 0;
      while i < n do
      begin
        prev := current;
        current := next;
        next := prev + current;
        i := i + 1
      end;
      write(prev)
    end;

  { \\brief Calculates the factorial of a given value and prints it
    \\param n The number to calculate the factorial of }
  procedure fact""" + str(redef) + """(n: integer);
    var product: integer;
    var i: integer;
    begin
      product := 1;
      i := 1;
      while i <= n do
      begin
        product := product * i;
        i := i + 1
      end;
      write(product)
    end;""")
        f.write("""
  { Get input number from user, then output the fib and fact of that number }
  begin
    writeln();
    writeln();
    read(x);
    writeln();
    fib1(x);
    fact1(x);
    writeln();
    writeln()
  end\n""")
    # PYTHON
    p_tic = time.perf_counter()
    parser.compileFile("temp.p")
    p_toc = time.perf_counter()
    # POGO SEQUENTIAL
    g_tic = time.perf_counter()
    # TODO:
    g_toc = time.perf_counter()
    # POGO CONCURRENT
    c_tic = time.perf_counter()
    # TODO:
    c_toc = time.perf_counter()
    print(size, "\t", (p_toc - p_tic), "\t", (g_toc - g_tic), "\t", (c_toc - c_tic), "\t")

os.remove("temp.p")
