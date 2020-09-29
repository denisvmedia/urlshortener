package storage

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
