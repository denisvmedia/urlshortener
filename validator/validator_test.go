package validator_test

import (
	"fmt"
	"github.com/go-playground/validator/v10"

	. "github.com/denisvmedia/urlshortener/validator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Functional Tests", func() {
	Context("ValidateURLShortName", func() {
		var validate *validator.Validate

		BeforeEach(func() {
			validate = validator.New()
			err := validate.RegisterValidation("shortname", ValidateURLShortName)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should successfully validate valid values", func() {
			validValues := map[string]string{
				"lowercase latin letters":                        "abcd",
				"uppercase latin letters":                        "ABCD",
				"mixedcase latin letters":                        "aBcD",
				"mixedcase latin letters with dashes":            "a-B-c-D",
				"mixedcase latin letters with digits":            "123fdafDFA",
				"only digits":                                    "123123123",
				"digits with dashes":                             "123-123-123",
				"mixedcase latin letters with digits and dashes": "123fd-afD-FA",
			}
			for typ, value := range validValues {
				By(fmt.Sprintf("should accept %s", typ))
				err := validate.Var(value, "shortname")
				Expect(err).NotTo(HaveOccurred(), "should have accepted %s", typ)
			}
		})

		It("Should fail to validate invalid values", func() {
			validValues := []string{
				/* 1 */ "a$",
				/* 2 */ "===",
				/* 3 */ "   ",
				/* 4 */ "abcd.asd",
				/* 5 */ "abcd/efdh",
				/* 6 */ "abcd%20%40",
				/* 7 */ "<script>",
				/* 8 */ "(asfdas)",
			}
			for id, value := range validValues {
				By(fmt.Sprintf("should BOT accept invalid value #%d", id))
				err := validate.Var(value, "shortname")
				Expect(err).To(HaveOccurred(), "should NOT have accepted invalid value #%s", id)
			}
		})

		It("Should fail to validate blacklisted values", func() {
			validValues := []string{
				"api",
				"swagger",
				"metrics",
			}
			for _, value := range validValues {
				By(fmt.Sprintf("should BOT accept blacklisted value %s", value))
				err := validate.Var(value, "shortname")
				Expect(err).To(HaveOccurred(), "should NOT have accepted invalid value %s", value)
			}
		})
	})

	Context("ValidateURLScheme", func() {
		var validate *validator.Validate

		BeforeEach(func() {
			validate = validator.New()
			err := validate.RegisterValidation("urlscheme", ValidateURLScheme)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should successfully validate allowed url schemes", func() {
			validValues := []string{
				"http://example.com",
				"https://example.com",
			}
			for _, value := range validValues {
				By(fmt.Sprintf("should accept %s", value))
				err := validate.Var(value, "urlscheme")
				Expect(err).NotTo(HaveOccurred(), "should have accepted %s", value)
			}
		})

		It("Should fail to validate invalid values", func() {
			validValues := []string{
				/* 1 */ "ftp://example.com/",
				/* 2 */ "ssh://example.com/",
				/* 3 */ "telnet://example.com/",
				/* 4 */ "news://example.com/",
				/* 5 */ "magnet:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a",
				/* 6 */ "wz0r-in\\/a1iD://aha!",
				/* 7 */ "(another%20invalid) ",
				/* 8 */ "//example.com",
				/* 9 */ "example.com",
			}
			for id, value := range validValues {
				By(fmt.Sprintf("should NOT accept invalid value #%d", id))
				err := validate.Var(value, "urlscheme")
				Expect(err).To(HaveOccurred(), "should NOT have accepted invalid value #%s", id)
			}
		})
	})
})
