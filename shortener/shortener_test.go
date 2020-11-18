package shortener_test

import (
	"github.com/denisvmedia/urlshortener/model"
	"github.com/denisvmedia/urlshortener/shortener"
	"github.com/denisvmedia/urlshortener/storage/linkstorage"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Functional Tests", func() {
	var linkStorage linkstorage.Storage
	var handler echo.HandlerFunc
	var router *echo.Echo

	BeforeEach(func() {
		linkStorage = linkstorage.NewInMemoryStorage()
		_, err := linkStorage.Insert(model.Link{
			ShortName:   "my-cool-link",
			OriginalURL: "https://example.com/my-cool-link",
		})
		Expect(err).ToNot(HaveOccurred())
		handler = shortener.Handler(linkStorage)
		router = echo.New()
		router.GET("/*", handler)
	})

	When("Link with given shortname exists", func() {
		It("Should be able to redirect", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/my-cool-link", nil)
			Expect(err).ToNot(HaveOccurred())
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusMovedPermanently))
			Expect(rec.Header().Get("location")).To(Equal("https://example.com/my-cool-link"))
		})
	})

	When("Link with given shortname does not exist", func() {
		It("Should return Not Found", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/nonexistent", nil)
			Expect(err).ToNot(HaveOccurred())
			router.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})
	})
})
