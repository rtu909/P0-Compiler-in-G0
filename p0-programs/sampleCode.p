program p;
  type T = array [1..10] of integer;
  var x: integer;
  var z: T;
  procedure q(a: integer; b: integer);
    var y: integer;
    begin y := a; write(y); writeln(); {writes 7}
      a := b; write(x); write(a); writeln(); {writes 5, 5}
      b := y; write(b); write(x); writeln(); {writes 7, 5}
      write(a); write(y); writeln(); {writes 5, 7}
      write(z[5]) {writes 5}
    end;
  procedure r(var c: T);
    begin c[x] := x; q(7, c[x]); write(x) {writes 5}
    end;
  begin x := 5; r(z)
  end
