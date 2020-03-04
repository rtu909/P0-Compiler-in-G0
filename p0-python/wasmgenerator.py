from scanner import TIMES, DIV, MOD, AND, PLUS, MINUS, OR, EQ, NE, LT, GT, LE, \
     GE, NOT, mark
from symboltable import indent, Var, Ref, Const, Type, Proc, StdProc, Int, Bool, Array, Record

def genProgStart():
    global curlev, memsize, asm
    curlev, memsize = 0, 0
    asm = ['(module',
           '(import "P0lib" "write" (func $write (param i32)))',
           '(import "P0lib" "writeln" (func $writeln))',
           '(import "P0lib" "print" (func $print (param i32)))',
           '(import "P0lib" "read" (func $read (result i32)))']

def genBool(b):
    # b is Bool
    b.size = 1; return b

def genInt(i):
    # i is Int
    i.size = 4; return i

def genRec(r):
    # r is Record
    s = 0
    for f in r.fields:
        f.offset, s = s, s + f.tp.size
    r.size = s
    return r

def genArray(a: Array):
    # a is Array
    a.size = a.length * a.base.size
    return a

def genGlobalVars(sc, start):
    global memsize
    for i in range(start, len(sc)):
        if type(sc[i]) == Var:
            if sc[i].tp in (Int, Bool):
                asm.append('(global $' + sc[i].name + ' (mut i32) i32.const 0)')
            elif type(sc[i].tp) in (Array, Record):
                sc[i].lev, sc[i].adr, memsize = -2, memsize, memsize + sc[i].tp.size
            else: mark('WASM: type?')
    
def genLocalVars(sc, start):
    for i in range(start, len(sc)):
        if type(sc[i]) == Var:
            if sc[i].tp in (Int, Bool):
                asm.append('(local $' + sc[i].name + ' i32)')
            elif type(sc[i].tp) in (Array, Record):
                mark('WASM: no local arrays, records')
            else: mark('WASM: type?')
    return None

def loadItem(x):
    if type(x) == Var:
        if x.lev == 0: asm.append('global.get $' + x.name) # global Var
        elif x.lev == curlev: asm.append('local.get $' + x.name) # local Var
        elif x.lev == -2: # memory Var
            asm.append('i32.const ' + str(x.adr))
            asm.append('i32.load')
        elif x.lev != -1: mark('WASM: var level!') # already on stack if lev == -1
    elif type(x) == Ref:
        if x.lev == -1: asm.append('i32.load')
        elif x.lev == curlev:
            asm.append('local.get $' + x.name)
            asm.append('i32.load')
        else: mark('WASM: ref level!')
    elif type(x) == Const: asm.append('i32.const ' + str(x.val))

def genVar(x):
    # x is Var, Ref
    if 0 < x.lev < curlev: mark('WASM: level!')
    if type(x) == Ref:
        y = Ref(x.tp); y.lev, y.name = x.lev, x.name
        # if type(x.tp) in (Array, Record):
        #    if x.lev > 0: y.name = x.name 
    elif type(x) == Var:
        y = Var(x.tp); y.lev, y.name = x.lev, x.name
        # if x.lev >= 0: y.name = x.name
        if x.lev == -2: y.adr = x.adr
    return y

def genConst(x):
    # x is Const
    return x

def genUnaryOp(op, x):
    loadItem(x)
    if op == MINUS:
        asm.append('i32.const -1')
        asm.append('i32.mul')
        x = Var(Int); x.lev = -1
    elif op == NOT:
        asm.append('i32.eqz')
        x = Var(Bool); x.lev = -1
    elif op == AND:
        asm.append('if (result i32)')
        x = Var(Bool); x.lev = -1
    elif op == OR:
        asm.append('if (result i32)')
        asm.append('i32.const 1')
        asm.append('else')
        x = Var(Bool); x.lev = -1
    else: mark('WASM: unary operator?')
    return x

def genBinaryOp(op, x, y):
    if op in (PLUS, MINUS, TIMES, DIV, MOD):
        loadItem(x); loadItem(y)
        asm.append('i32.add' if op == PLUS else \
                   'i32.sub' if op == MINUS else \
                   'i32.mul' if op == TIMES else \
                   'i32.div_s' if op == DIV else \
                   'i32.rem_s' if op == MOD else '?')
        x = Var(Int); x.lev = -1
    elif op == AND:
        loadItem(y) # x is already on the stack
        asm.append('else')
        asm.append('i32.const 0')
        asm.append('end')
        x = Var(Bool); x.lev = -1
    elif op == OR:
        loadItem(y) # x is already on the stack
        asm.append('end')
        x = Var(Bool); x.lev = -1
    else: assert False
    return x

def genRelation(op, x, y):
    loadItem(x); loadItem(y)
    asm.append('i32.eq' if op == EQ else \
               'i32.ne' if op == NE else \
               'i32.lt_s' if op ==  LT else \
               'i32.gt_s' if op == GT else \
               'i32.le_s' if op == LE else \
               'i32.ge_s' if op == GE else '?')
    x = Var(Bool); x.lev = -1
    return x

def genSelect(x, f):
    # x.f, assuming x.tp is Record and x is global Var, local Ref, stack Ref
    # and f is Field
    if type(x) == Var: x.adr += f.offset
    elif type(x) == Ref:
        if x.lev > 0: asm.append('local.get $' + x.name)
        asm.append('i32.const ' + str(f.offset))
        asm.append('i32.add')
        x.lev = -1
    x.tp = f.tp
    return x

def genIndex(x, y):
    # x[y], assuming x.tp is Array and x is global Var, local Ref, stack Ref
    # and y is Const, local Var, global Var, stack Var
    if type(x) == Var: # at x.adr
        if type(y) == Const: 
            x.adr += (y.val - x.tp.lower) * x.tp.base.size
            x.tp = x.tp.base
        else: # y is global Var, local Var, stack Var
            loadItem(y) # y on stack
            if x.tp.lower != 0:
                asm.append('i32.const ' + str(x.tp.lower))
                asm.append('i32.sub')
            asm.append('i32.const ' + str(x.tp.base.size))
            asm.append('i32.mul')
            asm.append('i32.const ' + str(x.adr))
            asm.append('i32.add')
            x = Ref(x.tp.base); x.lev = -1
    else: # x is local Ref, stack Ref; y is Const, global Var, local Var, stack Var
        if x.lev == curlev: loadItem(x); x.lev = -1
        if type(y) == Const:
            asm.append('i32.const ' + str((y.val - x.tp.lower) * x.tp.base.size))
            asm.append('i32.add')
        else:
            loadItem(y) # y on stack
            asm.append('i32.const ' + str(x.tp.lower))
            asm.append('i32.sub')
            asm.append('i32.const ' + str(x.tp.base.size))
            asm.append('i32.mul')
            asm.append('i32.add')
        x.tp = x.tp.base
    return x

def genAssign(x, y):
    if type(x) == Var:
        if x.lev == -2: asm.append('i32.const ' + str(x.adr))
        loadItem(y)
        if x.lev == 0: asm.append('global.set $' + x.name)
        elif x.lev == curlev: asm.append('local.set $' + x.name)
        elif x.lev == -2: asm.append('i32.store')
        else: mark('WASM: level!')
    elif type(x) == Ref:
        if x.lev == curlev: asm.append('local.get $' + x.name)
        loadItem(y)
        asm.append('i32.store')
    else: assert False

def genProgEntry(ident):
    asm.append('(func $program')

def genProgExit(x):
    asm.append(')\n(memory ' + str(memsize // 2** 16 + 1) + ')\n(start $program)\n)')
    return '\n'.join(l for l in asm)

def genProcStart(ident, fp):
    global curlev
    if curlev > 0: mark('WASM: no nested procedures')
    curlev = curlev + 1
    asm.append('(func $' + ident + ' ' + ' '.join('(param $' + e.name + ' i32)' for e in fp) + '')
    for p in fp:
        if p.tp in (Int, Bool) and type(p) == Ref:
            mark('WASM: only array and record reference parameters')
        elif type(p.tp) in (Array, Record) and type(p) == Var:
            mark('WASM: no structured value parameters')

def genProcEntry(ident, parsize, localsize):
    pass

def genProcExit(x, parsize, localsize):
    global curlev
    curlev = curlev - 1
    asm.append(')')

def genActualPara(ap, fp, n):
    if type(fp) == Ref:  #  reference parameter, assume ap is Var
        if ap.lev == -2: asm.append('i32.const ' + str(ap.adr))
        # else ap.lev == -1, on stack already
    elif type(ap) in (Var, Ref, Const): loadItem(ap)
    else: mark('unsupported parameter type')

def genCall(pr, ap):
    asm.append('call $' + pr.name)

def genRead(x):
    global curlev
    asm.append('call $read')
    if x.lev == 0: asm.append('global.set $' + x.name)
    elif x.lev == curlev: asm.append('local.set $' + x.name)
    elif x.lev == -2: asm.append('i32.store')
    else: mark('WASM: level!')
    y = Var(Int); y.lev = -1 # WTF is this supposed to do

def genWrite(x):
    loadItem(x)
    asm.append('call $write')

def genWriteln():
    asm.append('call $writeln')

def genPrint(x):
    loadItem(x)
    asm.append('call $print')

def genSeq(x, y):
    pass

def genThen(x):
    loadItem(x)
    asm.append('if')
    return x

def genIfThen(x, y):
    asm.append('end')

def genElse(x, y):
    asm.append('else')

def genIfElse(x, y, z):
    asm.append('end')

def genWhile():
    asm.append('loop')

def genDo(x):
    loadItem(x)
    asm.append('if')
    return x

def genWhileDo(t, x, y):
    asm.append('br 1')
    asm.append('end')
    asm.append('end')
