package utils

import (
	"static-api/models"
	"strconv"
	"strings"
)

func ParseContentRange(contentRange, limitStr string) models.Meta {
	meta := models.Meta{
		TotalCount: 0,
		PageSize:   0,
		Page:       1,
	}

	if contentRange == "" {
		return meta
	}

	// Example: "0-19/100"
	parts := strings.Split(contentRange, "/")
	if len(parts) != 2 {
		return meta
	}

	rangePart := parts[0]
	totalStr := parts[1]

	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return meta
	}

	rangeParts := strings.Split(rangePart, "-")
	if len(rangeParts) != 2 {
		return meta
	}

	start, err1 := strconv.Atoi(rangeParts[0])
	end, err2 := strconv.Atoi(rangeParts[1])

	if err1 != nil || err2 != nil {
		return meta
	}

	pageSize := end - start + 1
	page := 1

	if limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			page = (start / limit) + 1
		}
	}

	meta.TotalCount = total
	meta.PageSize = pageSize
	meta.Page = page

	return meta
}
