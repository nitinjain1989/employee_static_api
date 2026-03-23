package utils

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func BuildPath(c *gin.Context) string {

	limit := c.Query("limit")
	offset := c.Query("offset")
	search := c.Query("search")
	designation := c.Query("designation")
	department := c.Query("department")
	status := c.Query("status")

	path := "/employees?select=*,mobiles!fk_employee(*)"

	if search != "" {
		path += "&name=ilike.*" + url.QueryEscape(search) + "*"
		path += "&order=name.asc,created_at.desc,id.desc"
	} else {
		path += "&order=created_at.desc,id.desc"
	}

	path += BuildInFilter("designation", designation)
	path += BuildInFilter("department", department)

	if status == "active" {
		path += "&is_active=eq.true"
	} else if status == "inactive" {
		path += "&is_active=eq.false"
	}

	if limit != "" {
		path += "&limit=" + limit
	}
	if offset != "" {
		path += "&offset=" + offset
	}

	return path
}

func BuildInFilter(key, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	// Split comma-separated values
	items := strings.Split(value, ",")

	var cleaned []string
	for _, item := range items {
		v := strings.TrimSpace(item)
		if v != "" {
			// URL encode each value
			cleaned = append(cleaned, url.QueryEscape(v))
		}
	}

	if len(cleaned) == 0 {
		return ""
	}

	return "&" + key + "=in.(" + strings.Join(cleaned, ",") + ")"
}
