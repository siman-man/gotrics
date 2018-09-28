package gotrics

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestFuncLength(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue int
	}{
		{`
package t
func add(a, b int) int {
	return a + b
}
`, 3},
	}

	for _, tt := range tests {
		var actual int
		fset := token.NewFileSet()
		f, _ := parser.ParseFile(fset, "example.go", tt.input, parser.ParseComments)

		ast.Inspect(f, func(n ast.Node) bool {
			if r, ok := n.(*ast.FuncDecl); ok {
				actual = FuncLength(fset, r)
			}
			return true
		})

		if actual != tt.expectedValue {
			t.Errorf("expected=%d, got=%d", tt.expectedValue, actual)
		}
	}
}

func TestFuncNesting(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue int
	}{
		{`
package t
func add(a, b int) int {
	return a + b
}
`, 1},
		{`
package t
func add(a, b int) int {
	if x > 1 {
		return 1
	}
	return 0
}
`, 2},
		{`
package t
func add(a, b int) int {
	if x > 1 {
		return 1
	}
	if x < 1 {
		return 2
	}
	return 0
}
`, 2},
		{`
package t
func main() {
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		fmt.Printf("%s.", os)
	}
}
`, 2},
		{`
package t
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
`, 3},
		{`
package t
func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}
`, 2},
		{`
package t
func do(i interface{}) {
	noteFrequency := map[string]float32{
		"C0": 16.35,
		"G0": 24.50
	}
}
`, 2},
		{`
package t
func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	} else {
		fmt.Printf("%g >= %g\n", v, lim)
	}
	return lim
}
`, 2},
	}

	for i, tt := range tests {
		var actual int
		fset := token.NewFileSet()
		f, _ := parser.ParseFile(fset, "example.go", tt.input, parser.ParseComments)

		ast.Inspect(f, func(n ast.Node) bool {
			if r, ok := n.(*ast.FuncDecl); ok {
				// ast.Print(fset, r)
				actual = FuncNesting(fset, r)
			}
			return true
		})

		if actual != tt.expectedValue {
			t.Errorf("case %d: expected=%d, got=%d", i, tt.expectedValue, actual)
		}
	}
}

func TestParameterList(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue int
	}{
		{`
package t
func f(a int, b int) {
}
`, 2},
		{`
package t
func f(a, b, c int) {
}
`, 3},
		{`
package t
func f() {
}
`, 0},
		{`
package t
func f(x int) {
}
`, 1},
		{`
package t
func f(a, _ int, z float32) {
}
`, 2},
		{`
package t
func f(a, b int, z float32) {
}
`, 3},
		{`
package t
func f(prefix string, values ...int) {
}
`, 2},
		{`
package t
func f(a, b int, z float64, opt ...interface{}) {
}
`, 4},
		{`
package t
func f(int, int, float64) {
}
`, 0},
		{`
package t
func f(n int) {
}
`, 1},
	}

	for _, tt := range tests {
		var actual int
		fset := token.NewFileSet()
		f, _ := parser.ParseFile(fset, "example.go", tt.input, parser.ParseComments)

		ast.Inspect(f, func(n ast.Node) bool {
			if r, ok := n.(*ast.FuncDecl); ok {
				actual = ParameterList(r)
			}
			return true
		})

		if actual != tt.expectedValue {
			t.Errorf("expected=%d, got=%d", tt.expectedValue, actual)
		}
	}
}

func TestABCSize(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue float64
	}{
		{`
package t
func add(a, b int) int {
	return a + b
}
`, 0.0},
		{`
package t
func add() int {
	var a = 1
	var b = 2
	return a + b
}
`, 2.0},
		{`
package t
func add() int {
	var a, b = 1
	return a + b
}
`, 2.0},
		{`
package t
func add() int {
	var a, b int
	return a + b
}
`, 0.0},
		{`
package t
func add() int {
	var _, b = 10, 5
	return 3 + b
}
`, 1.0},
		{`
package t
func add() int {
	a, b := 10, 5
	return a + b
}
`, 2.0},
		{`
package t
func add() int {
	_, b := 10, 5
	return b
}
`, 1.0},
		{`
package t
func add() int {
	const a, b = 10, 5
	return a + b
}
`, 0.0},
		{`
package t
func add() int {
	const (
		a = 10
		b = 5
	)
	return a + b
}
`, 0.0},
		{`
package t
func add() int {
	const _, b = 10, 5
	return b
}
`, 0.0},
		{`
package t
func add() int {
	var a, b int
	a++
	b--
	return a + b
}
`, 2.0},
		{`
package t
func add() int {
	var a, b int
	a*=2
	b%=5
	return a + b
}
`, 2.0},
		{`
package t
func f() {
	math.Atan2(x, y)
	Greeting("hello:", "World")
}
`, 2.0},
		{`
package t
func f() {
	var pt *Point
	pt.Scale(3.5) 
}
`, 1.0},
		{`
package t
func f() {
	1 + 2 + 3
}
`, 0.0},
		{`
package t
func f() {
	goto L
L:
	x := 1
	_ = x
}
`, 1.41},
		{`
package t
func f() {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln("Error")
    }
}
`, 2.45},
		{`
package t
func f() {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintln("Error")
    }
}
`, 2.45},
		{`
package t
func f() {
	switch {
	case x > 0:
	case x < 0:
	default:
	}
}
`, 3.0},
		{`
package t
func f() {
	if x > 1 {
		true
	} else {
		false
	}
}
`, 2.0},
		{`
package t
func f() {
	if x > 1 {
		true
	} else if x < 1 {
		true
	} else {
		false
	}
}
`, 3.0},
		{`
package t
func fibonacci(c, quit chan int) {
	for {
		select {
		case c <- x:
			return
		case <-quit:
			return
		default:
		}
	}
}
`, 3.0},
		{`
package t
func add() {
	for i := 1; i < 10; i++ {
	}
}
`, 2.24},
	}

	for i, tt := range tests {
		var actual float64
		fset := token.NewFileSet()
		f, _ := parser.ParseFile(fset, "example.go", tt.input, parser.ParseComments)

		ast.Inspect(f, func(n ast.Node) bool {
			if r, ok := n.(*ast.FuncDecl); ok {
				// ast.Print(fset, r)
				actual = ABCSize(r)
			}
			return true
		})

		if actual != tt.expectedValue {
			t.Errorf("case %d: expected=%f, got=%f", i, tt.expectedValue, actual)
		}
	}
}
