package httputils

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetIntParamFromRequestPath(r *http.Request, paramName string) (int, error) {
	var number int

	param := r.PathValue(paramName)
	if param == "" {
		return 0, fmt.Errorf("missing {%s} param in url", paramName)
	}

	number, err := strconv.Atoi(param)

	if err != nil {
		return 0, err
	}
	return number, nil
}
