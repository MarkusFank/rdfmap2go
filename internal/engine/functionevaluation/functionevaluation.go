package functionevaluation

import (
	"errors"
	"fmt"
	"strings"
)

type functions = map[string]func(fnParams map[string]any) (string, error)

type FunctionEvaluator struct {
	functionMapping         map[string]string
	namespacesWithFunctions map[string]functions
	isInitialized           bool
}

func (evaluator *FunctionEvaluator) Init(functionMapping map[string]string) {
	evaluator.functionMapping = functionMapping

	evaluator.namespacesWithFunctions = map[string]functions{}

	evaluator.namespacesWithFunctions["rdfmap2go.strings"] = functions{}
	evaluator.namespacesWithFunctions["rdfmap2go.strings"]["lower"] = lower
	evaluator.namespacesWithFunctions["rdfmap2go.strings"]["upper"] = upper
	evaluator.namespacesWithFunctions["rdfmap2go.strings"]["substr"] = substr

	evaluator.isInitialized = true
}

func (evaluator *FunctionEvaluator) EvaluateFunction(functionName string, params map[string]any) (string, error) {
	if !evaluator.isInitialized {
		return "", errors.New("Function evaluator must be initialized before usage")
	}

	namespace, fnNameWithoutNamespace, err := extractNamespaceAndFunctionName(functionName, evaluator.functionMapping)

	if err != nil {
		return "", err
	}

	functionsOfNamespace, hasNamespace := evaluator.namespacesWithFunctions[namespace]

	if !hasNamespace {
		return "", fmt.Errorf("Unknown namespace '%s' for function '%s'", namespace, fnNameWithoutNamespace)
	}

	function, hasFunction := functionsOfNamespace[fnNameWithoutNamespace]

	if !hasFunction {
		return "", fmt.Errorf("Unable to find function '%s' in namespace '%s'", fnNameWithoutNamespace, namespace)
	}

	evalResult, fnErr := function(params)

	return evalResult, fnErr
}

func extractNamespaceAndFunctionName(functionName string, functionMapping map[string]string) (string, string, error) {

	for fnPrefix, fnNamespace := range functionMapping {
		if fnWithoutPrefix, replaced := strings.CutPrefix(functionName, fnPrefix+"::"); replaced {

			return fnNamespace, fnWithoutPrefix, nil
		}
	}

	return "", "", fmt.Errorf("Unable to resolve function name of '%s'. Make sure the prefix is declared!", functionName)
}
