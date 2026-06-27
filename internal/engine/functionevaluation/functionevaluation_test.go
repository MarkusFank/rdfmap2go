package functionevaluation

import (
	"testing"
)

func Test_Lower_Func(t *testing.T) {
	/*
		Tests whether the lower function is executed correctly and delivers the correct result
	*/

	fe := FunctionEvaluator{}

	fe.Init(map[string]string{"fn": "rdfmap2go.strings"})

	res, err := fe.EvaluateFunction("fn::lower", map[string]any{"str": "I AM A TEXT"})

	if err != nil {
		t.Fatal(err)
	}

	expected := "i am a text"
	if res != expected {
		t.Fatalf("Expected result to be %q but was %q instead!", expected, res)
	}
}

func Test_Uninitialized(t *testing.T) {
	/*
		Tests whether the function evaluator returns an error when it wasn't initialized before usage
	*/

	fe := FunctionEvaluator{}

	_, err := fe.EvaluateFunction("fn::lower", map[string]any{"str": "ASDF"})

	if err == nil {
		t.Fatal("Expected an error as function evaluator was not initialized")
	}
}

func Test_WrongNamespace(t *testing.T) {
	/*
		Tests whether the function evaluator returns an error when it is unable to resolve the namespace
	*/

	fe := FunctionEvaluator{}

	fe.Init(map[string]string{"fn": "rdfmap2go.strings"})

	_, err := fe.EvaluateFunction("wrong_ns::lower", map[string]any{"str": "I AM A TEXT"})

	if err == nil {
		t.Fatal("Expected an error as function call uses namespace that is not defined")
	}
}

func Test_UnknownNamespace(t *testing.T) {
	/*
		Tests whether the function evaluator returns an error when the mapping specifies an unknown namespace
	*/

	fe := FunctionEvaluator{}

	fe.Init(map[string]string{"fn": "some.wrong.namespace"})

	_, err := fe.EvaluateFunction("fn::lower", map[string]any{})

	if err == nil {
		t.Fatal("Expected an error as function call uses an unknown namespace")
	}
}

func Test_UnknownFunction(t *testing.T) {
	/*
		Tests whether the function evaluator returns an error when the function does not exist within the specified namespace
	*/

	fe := FunctionEvaluator{}

	fe.Init(map[string]string{"fn": "rdfmap2go.strings"})

	_, err := fe.EvaluateFunction("fn::wrooong", map[string]any{"str": "I AM A TEXT"})

	if err == nil {
		t.Fatal("Expected an error as function call uses namespace that is not defined")
	}
}
