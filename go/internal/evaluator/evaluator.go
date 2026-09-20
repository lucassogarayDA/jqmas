package evaluator

import (
"fmt"
"strconv"
"strings"

"github.com/lucassogarayDA/jqmas/go/internal/parser"
)

// EvaluatorError es un error durante la evaluación.
type EvaluatorError struct {
Mensaje     string
Sugerencias []string
Tips        []string
}

func (e *EvaluatorError) Error() string {
return e.Mensaje
}

// Environment guarda las variables.
type Environment struct {
Variables map[string]interface{}
}

// NewEnvironment crea un entorno vacío.
func NewEnvironment() *Environment {
return &Environment{Variables: map[string]interface{}{}}
}

// Evaluator ejecuta el AST.
type Evaluator struct {
Env    *Environment
Output strings.Builder
}

// NewEvaluator crea un evaluator con un entorno vacío.
func NewEvaluator() *Evaluator {
return &Evaluator{
Env:    NewEnvironment(),
Output: strings.Builder{},
}
}

// Evaluate ejecuta el nodo raíz.
func (e *Evaluator) Evaluate(node parser.Node) error {
switch n := node.(type) {
case *parser.Program:
for _, stmt := range n.Statements {
if err := e.Evaluate(stmt); err != nil {
return err
}
}
return nil
case *parser.PresCall:
return e.evalPres(n)
case *parser.Assignment:
return e.evalAssignment(n)
}
return nil
}

// EvalExpr evalúa una expresión y devuelve su valor.
func (e *Evaluator) EvalExpr(node parser.Node) (interface{}, error) {
switch n := node.(type) {
case *parser.NumberLiteral:
return n.Value, nil
case *parser.StringLiteral:
return n.Value, nil
case *parser.Variable:
return e.envGet(n.Name)
case *parser.BinaryOp:
return e.evalBinary(n)
}
return nil, &EvaluatorError{Mensaje: fmt.Sprintf("Nodo desconocido: %T", node)}
}

func (e *Evaluator) envGet(name string) (interface{}, error) {
if v, ok := e.Env.Variables[name]; ok {
return v, nil
}
sugerencias := buscarVariableSimilar(name, e.Env.Variables)
tips := []string{"declarala con haz[nombre = valor]"}
if len(sugerencias) > 0 {
return nil, &EvaluatorError{
Mensaje:     fmt.Sprintf("Variable no definida: '%s'", name),
Sugerencias: sugerencias,
Tips:        tips,
}
}
return nil, &EvaluatorError{
Mensaje: fmt.Sprintf("Variable no definida: '%s'", name),
Tips:    tips,
}
}

func (e *Evaluator) evalAssignment(node *parser.Assignment) error {
value, err := e.EvalExpr(node.Value)
if err != nil {
return err
}
e.Env.Variables[node.Name] = value
return nil
}

func (e *Evaluator) evalBinary(node *parser.BinaryOp) (interface{}, error) {
left, err := e.EvalExpr(node.Left)
if err != nil {
return nil, err
}
right, err := e.EvalExpr(node.Right)
if err != nil {
return nil, err
}

switch node.Operator {
case "+":
// String + String
lStr, lOk := left.(string)
rStr, rOk := right.(string)
if lOk && rOk {
return lStr + rStr, nil
}
// Mixto: convertir a string
if lOk || rOk {
return formatValue(left) + formatValue(right), nil
}
// Números
return toFloat(left) + toFloat(right), nil

case "-", "*", "/":
if _, ok := left.(string); ok {
return nil, operacionIncompatible(node.Operator, left, right)
}
if _, ok := right.(string); ok {
return nil, operacionIncompatible(node.Operator, left, right)
}

lf := toFloat(left)
rf := toFloat(right)

switch node.Operator {
case "-":
return lf - rf, nil
case "*":
return lf * rf, nil
case "/":
if rf == 0 {
return nil, &EvaluatorError{
Mensaje: "División por cero",
Tips: []string{
fmt.Sprintf("estás dividiendo %s entre 0", formatValue(left)),
"revisá que el divisor no sea 0",
},
}
}
return lf / rf, nil
}
}

return nil, &EvaluatorError{Mensaje: fmt.Sprintf("Operador desconocido: %s", node.Operator)}
}

func operacionIncompatible(op string, left, right interface{}) error {
return &EvaluatorError{
Mensaje: fmt.Sprintf("No se puede usar '%s' entre %s y %s",
op, typeName(left), typeName(right)),
Tips: []string{
fmt.Sprintf("izquierda: %s (%s)", formatValue(left), typeName(left)),
fmt.Sprintf("derecha: %s (%s)", formatValue(right), typeName(right)),
"los operadores aritméticos solo funcionan con números",
},
}
}

func (e *Evaluator) evalPres(node *parser.PresCall) error {
linea := []string{}

volcarLinea := func() {
if len(linea) > 0 {
e.Output.WriteString(strings.Join(linea, " "))
linea = []string{}
}
}

for _, arg := range node.Arguments {
if arg.NotCount > 0 {
volcarLinea()
for i := 0; i < arg.NotCount; i++ {
e.Output.WriteString("\n")
}
}

value, err := e.EvalExpr(arg.Value)
if err != nil {
return err
}
linea = append(linea, formatValue(value))
}

volcarLinea()
e.Output.WriteString("\n")
return nil
}

// ----- helpers -----

func toFloat(v interface{}) float64 {
switch n := v.(type) {
case int:
return float64(n)
case float64:
return n
}
return 0
}

func formatValue(v interface{}) string {
switch n := v.(type) {
case float64:
// Si es entero, mostrar sin decimales
if n == float64(int64(n)) {
return strconv.FormatInt(int64(n), 10)
}
return strconv.FormatFloat(n, 'f', -1, 64)
case int:
return strconv.Itoa(n)
case string:
return n
}
return fmt.Sprintf("%v", v)
}

func typeName(v interface{}) string {
switch v.(type) {
case string:
return "str"
case int, float64:
return "número"
}
return fmt.Sprintf("%T", v)
}

// buscarVariableSimilar busca variables parecidas a `nombre`.
func buscarVariableSimilar(nombre string, variables map[string]interface{}) []string {
nombreLower := strings.ToLower(nombre)
resultado := []string{}
for key := range variables {
if distanciaLevenshtein(nombreLower, strings.ToLower(key)) <= 2 {
resultado = append(resultado, key)
}
}
return resultado
}

func distanciaLevenshtein(a, b string) int {
la, lb := len(a), len(b)
if la == 0 {
return lb
}
if lb == 0 {
return la
}
prev := make([]int, lb+1)
for j := 0; j <= lb; j++ {
prev[j] = j
}
for i := 1; i <= la; i++ {
curr := make([]int, lb+1)
curr[0] = i
for j := 1; j <= lb; j++ {
costo := 1
if a[i-1] == b[j-1] {
costo = 0
}
curr[j] = min3(curr[j-1]+1, prev[j]+1, prev[j-1]+costo)
}
prev = curr
}
return prev[lb]
}

func min3(a, b, c int) int {
if a < b {
if a < c {
return a
}
return c
}
if b < c {
return b
}
return c
}
