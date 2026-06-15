package functionevaluation

import (
	"errors"
	"strings"
)

func lower(fnParams map[string]any) (string, error) {

	strParamAny, hasStrParam := fnParams["str"]

	if !hasStrParam {
		return "", errors.New("lower function requires string parameter 'str'")
	}

	strParam, isStrString := strParamAny.(string)
	if !isStrString {
		return "", errors.New("lower function requires parameter 'str' to be of type string")
	}
	return strings.ToLower(strParam), nil
}

func upper(fnParams map[string]any) (string, error) {

	strParamAny, hasStrParam := fnParams["str"]

	if !hasStrParam {
		return "", errors.New("upper function requires string parameter 'str'")
	}

	strParam, isStrString := strParamAny.(string)
	if !isStrString {
		return "", errors.New("upper function requires parameter 'str' to be of type string")
	}
	return strings.ToUpper(strParam), nil

}

func substr(fnParams map[string]any) (string, error) {
	strParamAsAny, hasStrParam := fnParams["str"]
	startIdxAsAny, hasStartIdxParam := fnParams["startIdx"]
	endIdxAsAny, hasEndIdxParam := fnParams["endIdx"]

	if !hasStrParam {
		return "", errors.New("upper function requires string parameter 'str'")
	}

	if !hasStartIdxParam {
		return "", errors.New("upper function requires int parameter 'startIdx'")
	}

	if !hasEndIdxParam {
		return "", errors.New("upper function requires int parameter 'endIdx'")
	}

	strParam, isStrString := strParamAsAny.(string)

	if !isStrString {
		return "", errors.New("substr function requires parameter 'str' to be of type string")
	}

	startIdx, isStartIdxInt := startIdxAsAny.(uint64)

	if !isStartIdxInt {
		return "", errors.New("substr function requires parameter 'startIdx' to be an unsigned integer")
	}

	endIdx, isEndIdxInt := endIdxAsAny.(uint64)

	if !isEndIdxInt {
		return "", errors.New("substr function requires parameter 'endIdx' to be an unsigned integer")
	}

	return strParam[startIdx:endIdx], nil
}
