program longProgram;
  type a = array [1..3] of integer;
  type b = array [1..4] of integer;
  type c = array [1..5] of integer;
  type d = array [1..6] of integer;
  type e = array [1..7] of integer;
  type f = array [1..8] of integer;
  type g = array [1..9] of integer;
  type h = array [1..10] of integer;
  type i = array [1..11] of integer;
  type j = array [1..12] of integer;
  type k = array [1..13] of integer;
  type l = array [1..14] of integer;
  var v1: a;
  var v2: b;
  var v3: c;
  var v4: d;
  var v5: e;
  var v6: f;
  var v7: g;
  var v8: h;
  var v9: i;
  var v10: j;
  var v11: k;
  var v12: l;
  procedure p1(var par: a; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 3 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  procedure p2(var par: b; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 4 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  procedure p3(var par: c; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 5 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  procedure p4(var par: d; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 6 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  procedure p5(var par: e; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 7 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  procedure p6(var par: f; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 8 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  procedure p7(var par: g; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 9 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  procedure p8(var par: h; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 10 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  procedure p9(var par: i; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 11 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  procedure p10(var par: j; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 12 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  procedure p11(var par: k; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 13 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  procedure p12(var par: l; i: integer);
    var count: integer;
    begin
      count := 1;
      while count <= 14 do
        begin
          par[count] := i;
          count := count + 1
        end;
      write(par[1]);
      writeln()
    end;
  begin
    p1(v1, 5);
    p2(v2, 6);
    p3(v3, 7);
    p4(v4, 8);
    p5(v5, 9);
    p6(v6, 10)
  end
