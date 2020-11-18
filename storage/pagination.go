package storage

// SlicePaginate is a handy function to paginate a slice taking into account its length
func SlicePaginate(pageNum int, pageSize int, sliceLength int) (start int, end int) {
	start = pageNum * pageSize

	if start > sliceLength {
		start = sliceLength
	}

	end = start + pageSize
	if end > sliceLength {
		end = sliceLength
	}

	return start, end
}
