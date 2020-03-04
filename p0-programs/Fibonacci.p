program fib;

  { Used to store the user's input}
  var x: integer;

  { \brief Calculates the fibonnaci series at a given value and prints it
    \param n The value of the series to be calculated }
  procedure fib(n: integer);
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

  { \brief Calculates the factorial of a given value and prints it
    \param n The number to calculate the factorial of }
  procedure fact(n: integer);
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
    end;
  
  { Get input number from user, then output the fib and fact of that number }
  begin
    writeln();
    writeln();
    read(x);
    writeln();
    fib(x);
    fact(x);
    writeln();
    writeln()
  end
