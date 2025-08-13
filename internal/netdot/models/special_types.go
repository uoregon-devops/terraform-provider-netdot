package models

import (
	"strconv"
	"strings"
)

type Xlink struct {
	Type string
	ID   int64
}

func parseXlink(s string) (Xlink, error) {
	xlink := Xlink{}
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return Xlink{}, nil
	}
	xlink.Type = parts[0]
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return Xlink{}, err
	}
	xlink.ID = id
	return xlink, nil
}
