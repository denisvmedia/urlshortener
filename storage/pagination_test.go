package storage_test

import (
	"github.com/denisvmedia/urlshortener/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pagination", func() {
	Context("SlicePaginate", func() {
		It("will paginate", func() {
			start, end := storage.SlicePaginate(0, 0, 10)
			Expect(start).To(Equal(0))
			Expect(end).To(Equal(0))

			start, end = storage.SlicePaginate(0, 2, 10)
			Expect(start).To(Equal(0))
			Expect(end).To(Equal(2))

			start, end = storage.SlicePaginate(2, 2, 10)
			Expect(start).To(Equal(4))
			Expect(end).To(Equal(6))

			start, end = storage.SlicePaginate(4, 2, 10)
			Expect(start).To(Equal(8))
			Expect(end).To(Equal(10))

			start, end = storage.SlicePaginate(5, 2, 10)
			Expect(start).To(Equal(10))
			Expect(end).To(Equal(10))

			start, end = storage.SlicePaginate(6, 2, 10)
			Expect(start).To(Equal(10))
			Expect(end).To(Equal(10))
		})
	})
})
