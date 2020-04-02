import time
import os
import subprocess

print("\tPYTHON\t\t\tGO SEQUENTIAL\t\tGO CONCURRENT")

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
    RUNS = 5
    # PYTHON
    p_avg = 0
    for runs in range(RUNS):
        p_tic = time.perf_counter()
        subprocess.run(["python3", "../p0-python/test-main.py", "temp.p"])
        p_toc = time.perf_counter()
        p_avg += p_toc - p_tic
    p_avg /=  RUNS
    # POGO SEQUENTIAL
    g_avg = 0
    for runs in range(RUNS):
        g_tic = time.perf_counter()
        subprocess.run(["./p0-go-sequential", "temp.p"])
        g_toc = time.perf_counter()
        g_avg += g_toc - g_tic
    g_avg /= RUNS
    # POGO CONCURRENT
    c_avg = 0
    for runs in range(RUNS):
        c_tic = time.perf_counter()
        subprocess.run(["./p0-go-concurrent", "temp.p"])
        c_toc = time.perf_counter()
        c_avg += c_toc - c_tic
    c_avg /= RUNS
    print(size, "\t", p_avg, "\t", g_avg, "\t", c_avg, "\t")

os.remove("temp.p")
