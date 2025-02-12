package httputils

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetIdFromRequestPath(r *http.Request) (int, error) {
	var id int

	paramId := r.PathValue("id")
	if paramId == "" {
		return 0, fmt.Errorf("missing {id} param in url")
	}

	id, err := strconv.Atoi(paramId)

	if err != nil {
		return 0, err
	}
	return id, nil
}
