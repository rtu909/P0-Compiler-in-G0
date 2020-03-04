import scanner as SC  #  used for SC.init, SC.sym, SC.val, SC.error
from scanner import TIMES, DIV, MOD, AND, PLUS, MINUS, OR, EQ, NE, LT, GT, \
    LE, GE, PERIOD, COMMA, COLON, RPAREN, RBRAK, OF, THEN, DO, LPAREN, \
    LBRAK, NOT, BECOMES, NUMBER, IDENT, SEMICOLON, END, ELSE, IF, WHILE, \
    ARRAY, RECORD, CONST, TYPE, VAR, PROCEDURE, BEGIN, PROGRAM, EOF, \
    getSym, mark
import symboltable as ST  #  used for ST.init
from symboltable import Var, Ref, Const, Type, Proc, StdProc, Int, Bool, Enum, \
    Record, Array, newDecl, find, openScope, topScope, closeScope

FIRSTFACTOR = {IDENT, NUMBER, LPAREN, NOT}
FOLLOWFACTOR = {TIMES, DIV, MOD, AND, OR, PLUS, MINUS, EQ, NE, LT, LE, GT, GE,
                COMMA, SEMICOLON, THEN, ELSE, RPAREN, RBRAK, DO, PERIOD, END}
FIRSTEXPRESSION = {PLUS, MINUS, IDENT, NUMBER, LPAREN, NOT}
FIRSTSTATEMENT = {IDENT, IF, WHILE, BEGIN}
FOLLOWSTATEMENT = {SEMICOLON, END, ELSE}
FIRSTTYPE = {IDENT, RECORD, ARRAY, LPAREN}
FOLLOWTYPE = {SEMICOLON}
FIRSTDECL = {CONST, TYPE, VAR, PROCEDURE}
FOLLOWDECL = {BEGIN}
FOLLOWPROCCALL = {SEMICOLON, END, ELSE}
STRONGSYMS = {CONST, TYPE, VAR, PROCEDURE, WHILE, IF, BEGIN, EOF}


def selector(x):
    while SC.sym in {PERIOD, LBRAK}:
        if SC.sym == PERIOD:  #  x.f
            getSym()
            if SC.sym == IDENT:
                if type(x.tp) == Record:
                    for f in x.tp.fields:
                        if f.name == SC.val:
                            x = CG.genSelect(x, f); break
                    else: mark("not a field")
                    getSym()
                else: mark("not a record")
            else: mark("identifier expected")
        else:  #  x[y]
            getSym(); y = expression()
            if type(x.tp) == Array:
                if y.tp == Int:
                    if type(y) == Const and \
                       (y.val < x.tp.lower or y.val >= x.tp.lower + x.tp.length):
                        mark('index out of bounds')
                    else: x = CG.genIndex(x, y)
                else: mark('index not integer')
            else: mark('not an array')
            if SC.sym == RBRAK: getSym()
            else: mark("] expected")
    return x


def factor():
    if SC.sym not in FIRSTFACTOR:
        mark("expression expected")
        while SC.sym not in FIRSTFACTOR | FOLLOWFACTOR | STRONGSYMS: getSym()
    if SC.sym == IDENT:
        x = find(SC.val)
        if type(x) in {Var, Ref}: x = CG.genVar(x); getSym()
        elif type(x) == Const: x = Const(x.tp, x.val); x = CG.genConst(x); getSym()
        else: mark('expression expected')
        x = selector(x)
    elif SC.sym == NUMBER:
        x = Const(Int, SC.val); x = CG.genConst(x); getSym()
    elif SC.sym == LPAREN:
        getSym(); x = expression()
        if SC.sym == RPAREN: getSym()
        else: mark(") expected")
    elif SC.sym == NOT:
        getSym(); x = factor()
        if x.tp != Bool: mark('not boolean')
        elif type(x) == Const: x.val = 1 - x.val # constant folding
        else: x = CG.genUnaryOp(NOT, x)
    else: x = Const(None, 0)
    return x


def term():
    x = factor()
    while SC.sym in {TIMES, DIV, MOD, AND}:
        op = SC.sym; getSym();
        if op == AND and type(x) != Const: x = CG.genUnaryOp(AND, x)
        y = factor() # x op y
        if x.tp == Int == y.tp and op in {TIMES, DIV, MOD}:
            if type(x) == Const == type(y): # constant folding
                if op == TIMES: x.val = x.val * y.val
                elif op == DIV: x.val = x.val // y.val
                elif op == MOD: x.val = x.val % y.val
            else: x = CG.genBinaryOp(op, x, y)
        elif x.tp == Bool == y.tp and op == AND:
            if type(x) == Const: # constant folding
                if x.val: x = y # if x is true, take y, else x
            else: x = CG.genBinaryOp(AND, x, y)
        else: mark('bad type')
    return x


def simpleExpression():
    if SC.sym == PLUS:
        getSym(); x = term()
    elif SC.sym == MINUS:
        getSym(); x = term()
        if x.tp != Int: mark('bad type')
        elif type(x) == Const: x.val = - x.val # constant folding
        else: x = CG.genUnaryOp(MINUS, x)
    else: x = term()
    while SC.sym in {PLUS, MINUS, OR}:
        op = SC.sym; getSym()
        if op == OR and type(x) != Const: x = CG.genUnaryOp(OR, x)
        y = term() # x op y
        if x.tp == Int == y.tp and op in {PLUS, MINUS}:
            if type(x) == Const == type(y): # constant folding
                if op == PLUS: x.val = x.val + y.val
                elif op == MINUS: x.val = x.val - y.val
            else: x = CG.genBinaryOp(op, x, y)
        elif x.tp == Bool == y.tp and op == OR:
            if type(x) == Const: # constant folding
                if not x.val: x = y # if x is false, take y, else x
            else: x = CG.genBinaryOp(OR, x, y)
        else: mark('bad type')
    return x


def expression():
    x = simpleExpression()
    while SC.sym in {EQ, NE, LT, LE, GT, GE}:
        op = SC.sym; getSym(); y = simpleExpression() # x op y
        if x.tp == y.tp in (Int, Bool):
            if type(x) == Const == type(y): # constant folding
                if op == EQ: x.val = x.val == y.val
                elif op == NE: x.val = x.val != y.val
                elif op == LT: x.val = x.val < y.val
                elif op == LE: x.val = x.val <= y.val
                elif op == GT: x.val = x.val > y.val
                elif op == GE: x.val = x.val >= y.val
                x.tp = Bool
            else: x = CG.genRelation(op, x, y)
        else: mark('bad type')
    return x


def compoundStatement():
    if SC.sym == BEGIN: getSym()
    else: mark("'begin' expected")
    x = statement()
    while SC.sym == SEMICOLON or SC.sym in FIRSTSTATEMENT:
        if SC.sym == SEMICOLON: getSym()
        else: mark("; missing")
        y = statement(); x = CG.genSeq(x, y)
    if SC.sym == END: getSym()
    else: mark("'end' expected")
    return x


def statement():
    if SC.sym not in FIRSTSTATEMENT:
        mark("statement expected"); getSym()
        while SC.sym not in FIRSTSTATEMENT | FOLLOWSTATEMENT | STRONGSYMS : getSym()
    if SC.sym == IDENT:
        x = find(SC.val); getSym()
        if type(x) in {Var, Ref}:
            x = CG.genVar(x); x = selector(x)
            if SC.sym == BECOMES:
                getSym(); y = expression()
                if x.tp == y.tp in {Bool, Int}: x = CG.genAssign(x, y)
                else: mark('incompatible assignment')
            elif SC.sym == EQ:
                mark(':= expected'); getSym(); y = expression()
            else: mark(':= expected')
        elif type(x) in {Proc, StdProc}:
            fp, ap, i = x.par, [], 0   #  list of formals, list of actuals
            if SC.sym == LPAREN:
                getSym()
                if SC.sym in FIRSTEXPRESSION:
                    y = expression()
                    if i < len(fp):
                        if (type(fp[i]) == Var or type(y) == Var) and \
                           fp[i].tp == y.tp:
                            if type(x) == Proc:
                                ap.append(CG.genActualPara(y, fp[i], i))
                        else: mark('illegal parameter mode')
                    else: mark('extra parameter')
                    i = i + 1
                    while SC.sym == COMMA:
                        getSym()
                        y = expression()
                        if i < len(fp):
                            if (type(fp[i]) == Var or type(y) == Var) and \
                               fp[i].tp == y.tp:
                                if type(x) == Proc:
                                    ap.append(CG.genActualPara(y, fp[i], i))
                            else: mark('illegal parameter mode')
                        else: mark('extra parameter')
                        i = i + 1
                if SC.sym == RPAREN: getSym()
                else: mark("')' expected")
            if i < len(fp): mark('too few parameters')
            elif type(x) == StdProc:
                if x.name == 'read': x = CG.genRead(y)
                elif x.name == 'write': x = CG.genWrite(y)
                elif x.name == 'writeln': x = CG.genWriteln()
                elif x.name == 'print': x = CG.genPrint(y)
            else: x = CG.genCall(x, ap)
        else: mark("variable or procedure expected")
    elif SC.sym == BEGIN: x = compoundStatement()
    elif SC.sym == IF:
        getSym(); x = expression();
        if x.tp == Bool: x = CG.genThen(x)
        else: mark('boolean expected')
        if SC.sym == THEN: getSym()
        else: mark("'then' expected")
        y = statement()
        if SC.sym == ELSE:
            if x.tp == Bool: y = CG.genElse(x, y)
            getSym(); z = statement()
            if x.tp == Bool: x = CG.genIfElse(x, y, z)
        else:
            if x.tp == Bool: x = CG.genIfThen(x, y)
    elif SC.sym == WHILE:
        getSym(); t = CG.genWhile(); x = expression()
        if x.tp == Bool: x = CG.genDo(x)
        else: mark('boolean expected')
        if SC.sym == DO: getSym()
        else: mark("'do' expected")
        y = statement()
        if x.tp == Bool: x = CG.genWhileDo(t, x, y)
    else: x = None
    return x


def typ():
    if SC.sym not in FIRSTTYPE:
        mark("type expected")
        while SC.sym not in FIRSTTYPE | FOLLOWTYPE | STRONGSYMS: getSym()
    if SC.sym == IDENT:
        ident = SC.val; x = find(ident); getSym()
        if type(x) == Type: x = Type(x.val)
        else: mark('not a type'); x = Type(None)
    elif SC.sym == ARRAY:
        getSym()
        if SC.sym == LBRAK: getSym()
        else: mark("'[' expected")
        x = expression()
        if SC.sym == PERIOD: getSym()
        else: mark("'.' expected")
        if SC.sym == PERIOD: getSym()
        else: mark("'.' expected")
        y = expression()
        if SC.sym == RBRAK: getSym()
        else: mark("']' expected")
        if SC.sym == OF: getSym()
        else: mark("'of' expected")
        z = typ().val;
        if type(x) != Const or x.val < 0:
            mark('bad lower bound'); x = Type(None)
        elif type(y) != Const or y.val < x.val:
            mark('bad upper bound'); x = Type(None)
        else: x = Type(CG.genArray(Array(z, x.val, y.val - x.val + 1)))
    elif SC.sym == RECORD:
        getSym(); openScope(); typedIds(Var)
        while SC.sym == SEMICOLON:
            getSym(); typedIds(Var)
        if SC.sym == END: getSym()
        else: mark("'end' expected")
        r = topScope(); closeScope()
        x = Type(CG.genRec(Record(r)))
    else: x = Type(None)
    return x


def typedIds(kind):
    if SC.sym == IDENT: tid = [SC.val]; getSym()
    else: mark("identifier expected"); tid = []
    while SC.sym == COMMA:
        getSym()
        if SC.sym == IDENT: tid.append(SC.val); getSym()
        else: mark('identifier expected')
    if SC.sym == COLON:
        getSym(); tp = typ().val
        if tp != None:
            for i in tid: newDecl(i, kind(tp))
    else: mark("':' expected")


def declarations(allocVar):
    if SC.sym not in FIRSTDECL | FOLLOWDECL:
        mark("'begin' or declaration expected")
        while SC.sym not in FIRSTDECL | FOLLOWDECL | STRONGSYMS: getSym()
    while SC.sym == CONST:
        getSym()
        if SC.sym == IDENT:
            ident = SC.val; getSym()
            if SC.sym == EQ: getSym()
            else: mark("= expected")
            x = expression()
            if type(x) == Const: newDecl(ident, x)
            else: mark('expression not constant')
        else: mark("constant name expected")
        if SC.sym == SEMICOLON: getSym()
        else: mark("; expected")
    while SC.sym == TYPE:
        getSym()
        if SC.sym == IDENT:
            ident = SC.val; getSym()
            if SC.sym == EQ: getSym()
            else: mark("= expected")
            x = typ(); newDecl(ident, x)  #  x is of type ST.Type
            if SC.sym == SEMICOLON: getSym()
            else: mark("; expected")
        else: mark("type name expected")
    start = len(topScope())
    while SC.sym == VAR:
        getSym(); typedIds(Var)
        if SC.sym == SEMICOLON: getSym()
        else: mark("; expected")
    varsize = allocVar(topScope(), start)
    while SC.sym == PROCEDURE:
        getSym()
        if SC.sym == IDENT: getSym()
        else: mark("procedure name expected")
        ident = SC.val; newDecl(ident, Proc([])) #  entered without parameters
        sc = topScope()
        openScope() # new scope for parameters and body
        if SC.sym == LPAREN:
            getSym()
            if SC.sym in {VAR, IDENT}:
                if SC.sym == VAR: getSym(); typedIds(Ref)
                else: typedIds(Var)
                while SC.sym == SEMICOLON:
                    getSym()
                    if SC.sym == VAR: getSym(); typedIds(Ref)
                    else: typedIds(Var)
            else: mark("formal parameters expected")
            fp = topScope()
            sc[-1].par = fp[:] #  procedure parameters updated
            if SC.sym == RPAREN: getSym()
            else: mark(") expected")
        else: fp = []
        parsize = CG.genProcStart(ident, fp)
        if SC.sym == SEMICOLON: getSym()
        else: mark("; expected")
        localsize = declarations(CG.genLocalVars)
        CG.genProcEntry(ident, parsize, localsize)
        x = compoundStatement(); CG.genProcExit(x, parsize, localsize)
        closeScope() #  scope for parameters and body closed
        if SC.sym == SEMICOLON: getSym()
        else: mark("; expected")
    return varsize


def program():
    newDecl('boolean', Type(CG.genBool(Bool)))
    newDecl('integer', Type(CG.genInt(Int)))
    newDecl('true', Const(Bool, 1))
    newDecl('false', Const(Bool, 0))
    newDecl('read', StdProc([Ref(Int)]))
    newDecl('write', StdProc([Var(Int)]))
    newDecl('writeln', StdProc([]))
    newDecl('print', StdProc([Var(Int)]))
    CG.genProgStart()
    if SC.sym == PROGRAM: getSym()
    else: mark("'program' expected")
    ident = SC.val
    if SC.sym == IDENT: getSym()
    else: mark('program name expected')
    if SC.sym == SEMICOLON: getSym()
    else: mark('; expected')
    declarations(CG.genGlobalVars); CG.genProgEntry(ident)
    x = compoundStatement()
    return CG.genProgExit(x)


def compileString(src, dstfn = None, target = 'wat'):
    global CG
    if target == 'wat': import wasmgenerator as CG
    elif target == 'mips': import CGmips as CG
    elif target == 'ast': import CGast as CG
    else: print('unknown target'); return
    SC.init(src)
    ST.init()
    p = program()
    if p != None and not SC.error:
        if dstfn == None: print(p)
        else:
            with open(dstfn, 'w') as f: f.write(p);


def compileFile(srcfn, target = 'wat'):
    if srcfn.endswith('.p'):
        with open(srcfn, 'r') as f: src = f.read()
        dstfn = srcfn[:-2] + '.s'
        compileString(src, dstfn, target)
    else: print("'.p' file extension expected")
