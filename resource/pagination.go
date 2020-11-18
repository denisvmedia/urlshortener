package resource

import (
	"github.com/go-extras/api2go"
	"strconv"
)

const (
	pageNumberDefault = 1
	pageSizeDefault   = 10
	pageSizeMax       = 1000
)

// Page holds page number and page size used in pagination
type Page struct {
	Number int
	Size   int
}

func parsePageArgs(params map[string][]string) (result Page) {
	v := params["page[number]"]
	if len(v) == 0 {
		v = append(v, "")
	}
	number, err := strconv.ParseInt(v[0], 10, 32)
	if err != nil || number < 1 {
		number = pageNumberDefault
	}
	v = params["page[size]"]
	if len(v) == 0 {
		v = append(v, "")
	}
	size, err := strconv.ParseInt(v[0], 10, 32)
	// size can be used as 0 for meta reasons (e.g. getting total size, without getting the results)
	if err != nil || size < 0 || size > pageSizeMax {
		size = pageSizeDefault
	}

	return Page{
		Number: int(number),
		Size:   int(size),
	}
}

func getPagination(number, size, total int) (result api2go.Pagination) {
	var totalPages int
	if size > 0 {
		totalPages = total / size
		if reminder := total % size; reminder > 0 {
			totalPages++
		}
	} else {
		totalPages = 0
	}
	sizeStr := strconv.Itoa(size)
	if total > size*number {
		result.Last = map[string]string{
			"number": strconv.Itoa(totalPages),
			"size":   sizeStr,
		}
	}
	if number < totalPages {
		result.Next = map[string]string{
			"number": strconv.Itoa(number + 1),
			"size":   sizeStr,
		}
	}
	if number > 1 {
		result.Prev = map[string]string{
			"number": strconv.Itoa(number - 1),
			"size":   sizeStr,
		}
	}
	result.First = map[string]string{
		"number": "1",
		"size":   sizeStr,
	}

	return
}
