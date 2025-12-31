package parser

// AST represents a full Kyra module.
type AST struct {
	ModuleName string
	Imports    []*UseStmt
	TopLevel   []Stmt
}

// ---------------------------
// Statements
// ---------------------------

type Stmt interface {
	stmtNode()
}

// UseStmt: use kyra/add
type UseStmt struct {
	Path      []string
	IsStdlib  bool
}
func (*UseStmt) stmtNode() {}

// LetStmt: let x = expr
type LetStmt struct {
	Name string
	Expr Expr
}
func (*LetStmt) stmtNode() {}

// ReturnStmt: return expr
type ReturnStmt struct {
	Value Expr
}
func (*ReturnStmt) stmtNode() {}

// ExitStmt: exit
type ExitStmt struct{}
func (*ExitStmt) stmtNode() {}

// PassStmt: pass
type PassStmt struct{}
func (*PassStmt) stmtNode() {}

// ExprStmt: expression as statement
type ExprStmt struct {
	Expr Expr
}
func (*ExprStmt) stmtNode() {}

// IfStmt: if cond: block  OR  if cond { block }
type IfStmt struct {
	Cond Expr
	Then []Stmt
	Else []Stmt
}
func (*IfStmt) stmtNode() {}

// WhileStmt: while cond: block
type WhileStmt struct {
	Cond Expr
	Body []Stmt
}
func (*WhileStmt) stmtNode() {}

// ForStmt: for i 10   or   for i range
type ForStmt struct {
	VarName string
	Limit   Expr
	Body    []Stmt
}
func (*ForStmt) stmtNode() {}

// FuncDef: def name(args): block
type FuncDef struct {
	Name       string
	Args       []FuncArg
	ReturnType *TypeNode
	Body       []Stmt
}
func (*FuncDef) stmtNode() {}

// FuncExprDef: func name(args) { body }
type FuncExprDef struct {
	Name       string
	Args       []FuncArg
	ReturnType *TypeNode
	Body       []Stmt
}
func (*FuncExprDef) stmtNode() {}

// FuncOneLiner: func name(args) = expr
type FuncOneLiner struct {
	Name       string
	Args       []FuncArg
	ReturnType *TypeNode
	Expr       Expr
}
func (*FuncOneLiner) stmtNode() {}

// ---------------------------
// Expressions
// ---------------------------

type Expr interface {
	exprNode()
}

// Identifier: x, foo, bar
type IdentExpr struct {
	Name string
}
func (*IdentExpr) exprNode() {}

// Literal numbers
type NumberExpr struct {
	Value string
}
func (*NumberExpr) exprNode() {}

// Literal strings
type StringExpr struct {
	Value string
}
func (*StringExpr) exprNode() {}

// Literal booleans
type BoolExpr struct {
	Value bool
}
func (*BoolExpr) exprNode() {}

// Unary: -x, !x
type UnaryExpr struct {
	Op   string
	Expr Expr
}
func (*UnaryExpr) exprNode() {}

// Binary: x + y, x == y, x && y
type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
}
func (*BinaryExpr) exprNode() {}

// Assignment: x = expr
type AssignExpr struct {
	Name string
	Expr Expr
}
func (*AssignExpr) exprNode() {}

// Call: foo(x, y)
type CallExpr struct {
	Callee Expr
	Args   []Expr
}
func (*CallExpr) exprNode() {}

// Member access: a.b.c
type MemberExpr struct {
	Object Expr
	Name   string
}
func (*MemberExpr) exprNode() {}

// Parenthesized: (expr)
type ParenExpr struct {
	Expr Expr
}
func (*ParenExpr) exprNode() {}

// ---------------------------
// Function arguments
// ---------------------------

type FuncArg struct {
	Name string
	Type *TypeNode
}

// ---------------------------
// Types
// ---------------------------

type TypeNode struct {
	Name string
}
