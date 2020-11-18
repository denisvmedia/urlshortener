package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/denisvmedia/urlshortener/cmd"
	"github.com/denisvmedia/urlshortener/shortener"
	"github.com/denisvmedia/urlshortener/storage/linkstorage"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"

	"github.com/denisvmedia/urlshortener/model"
	"github.com/denisvmedia/urlshortener/resource"
	"github.com/go-extras/api2go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Functional Tests", func() {
	var api *api2go.API
	var linkStorage linkstorage.Storage
	var dbData cmd.Mysql // a little bit ugly borrowing this structure from `cmd`, but it works...

	BeforeEach(func() {
		log.SetOutput(ioutil.Discard)
		api = api2go.NewAPIWithBaseURL("api", "")
		if v, ok := os.LookupEnv("TEST_STORAGE"); ok && v == "mysql" {
			dbData = cmd.Mysql{}
			dbData.Host, ok = os.LookupEnv("MYSQL_HOST")
			Expect(ok).To(BeTrue(), "MYSQL_HOST must be set for this test")
			dbData.Name, ok = os.LookupEnv("MYSQL_DBNAME")
			Expect(ok).To(BeTrue(), "MYSQL_DBNAME must be set for this test")
			dbData.User, ok = os.LookupEnv("MYSQL_USER")
			Expect(ok).To(BeTrue(), "MYSQL_USER must be set for this test")
			dbData.Password, ok = os.LookupEnv("MYSQL_PASSWORD")
			Expect(ok).To(BeTrue(), "MYSQL_PASSWORD must be set for this test")
			err := dbData.Validate()
			Expect(err).ToNot(HaveOccurred(), "all MYSQL_* env vars must be set in order to run tests using 'mysql' storage")
			err = linkstorage.MysqlInitStorage(dbData.User, dbData.Password, dbData.Host, dbData.Name, true)
			Expect(err).ToNot(HaveOccurred())
			dbh, err := linkstorage.MysqlConnect(dbData.User, dbData.Password, dbData.Host, dbData.Name)
			Expect(err).ToNot(HaveOccurred())
			linkStorage = linkstorage.NewMysqlStorage(dbh)
		} else {
			linkStorage = linkstorage.NewInMemoryStorage()
		}
		api.AddResource(model.Link{}, resource.NewLinkResource(linkStorage))
	})

	AfterEach(func() {
		if v, ok := os.LookupEnv("TEST_STORAGE"); ok && v == "mysql" {
			err := linkstorage.MysqlDropDB(dbData.User, dbData.Password, dbData.Host, dbData.Name)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	var jsonMustMarshal = func(v interface{}) []byte {
		result, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return result
	}

	var newLinkRequest = func(shortName, originalUri, comment string) *http.Request {
		// "https://example.com/my-super-puper/url?withArgs=val%20with%20space#and-hash"
		data := jsonMustMarshal(map[string]interface{}{
			"data": map[string]interface{}{
				"type": "links",
				"attributes": map[string]interface{}{
					"shortName":   shortName,
					"originalUrl": originalUri,
					"comment":     comment,
				},
			},
		})
		req, err := http.NewRequest("POST", "/api/links", bytes.NewReader(data))
		Expect(err).ToNot(HaveOccurred())
		return req
	}

	var updateLinkRequest = func(id, shortName, originalUri, comment string) *http.Request {
		// "https://example.com/my-super-puper/url?withArgs=val%20with%20space#and-hash"
		data := jsonMustMarshal(map[string]interface{}{
			"data": map[string]interface{}{
				"type": "links",
				"id":   id,
				"attributes": map[string]interface{}{
					"shortName":   shortName,
					"originalUrl": originalUri,
					"comment":     comment,
				},
			},
		})
		req, err := http.NewRequest("PATCH", "/api/links/"+id, bytes.NewReader(data))
		Expect(err).ToNot(HaveOccurred())
		return req
	}

	When("Using API", func() {
		It("API Creates a new link", func() {
			By("Creating the first link", func() {
				rec := httptest.NewRecorder()
				req := newLinkRequest(
					"my-cool-shortName",
					"https://example.com/my-super-puper/url?withArgs=val%20with%20space#and-hash",
					"And some cool comment",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusCreated))
				Expect(rec.Body.String()).To(MatchJSON(`
			{
				"data": {
					"id": "1",
					"type": "links",
					"attributes": {
					  "comment": "And some cool comment",
					  "originalUrl": "https://example.com/my-super-puper/url?withArgs=val%20with%20space#and-hash",
					  "shortName": "my-cool-shortName"
					}
				}
			}
			`))
			})

			By("Creating another link", func() {
				rec := httptest.NewRecorder()
				req := newLinkRequest(
					"another-link",
					"https://example.com/another-link",
					"", // empty comment, why not
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusCreated))
				Expect(rec.Body.String()).To(MatchJSON(`
			{
				"data": {
					"id": "2",
					"type": "links",
					"attributes": {
					  "comment": "",
					  "originalUrl": "https://example.com/another-link",
					  "shortName": "another-link"
					}
				}
			}
			`))
			})

			By("Should create a link with an empty shortname and generate its name for the user", func() {
				rec := httptest.NewRecorder()
				req := newLinkRequest(
					"",
					"https://example.com/and-another-link",
					"",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusCreated))
				m := make(map[string]interface{})
				err := json.Unmarshal(rec.Body.Bytes(), &m)
				Expect(err).ToNot(HaveOccurred())
				Expect(m).To(HaveKey("data"))
				Expect(m["data"]).To(HaveKey("attributes"))
				Expect(m["data"].(map[string]interface{})["attributes"]).To(HaveKey("shortName"))
				Expect(m["data"].(map[string]interface{})["attributes"].(map[string]interface{})["shortName"]).ToNot(BeEmpty())
			})

			By("Should fail when giving an existing short link", func() {
				rec := httptest.NewRecorder()
				req := newLinkRequest(
					"my-cool-shortName",
					"not a valid url",
					"",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})

			By("Should fail when giving an invalid url", func() {
				rec := httptest.NewRecorder()
				req := newLinkRequest(
					"invalid-url-link",
					"not a valid url",
					"",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})

			By("Should fail when giving an invalid short link", func() {
				rec := httptest.NewRecorder()
				req := newLinkRequest(
					"an invalid short link (with spaces and other signs)",
					"https://example.com/valid-url",
					"",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})

			By("Should fail when passing an invalid json", func() {
				rec := httptest.NewRecorder()
				data := []byte(`invalid json{}`)
				req, err := http.NewRequest("POST", "/api/links", bytes.NewReader(data))
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusNotAcceptable))
			})

			By("Should fail when passing an invalid json structure", func() {
				rec := httptest.NewRecorder()
				data := []byte(`{
				"wrong_json_structure": {
					"id": "2",
					"type": "links"
				}
			}`)
				req, err := http.NewRequest("POST", "/api/links", bytes.NewReader(data))
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusNotAcceptable))
			})
		})

		It("API Updates a link", func() {
			By("Creating a link", func() {
				rec := httptest.NewRecorder()
				req := newLinkRequest(
					"my-cool-link",
					"https://example.com/my-cool-link",
					"",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusCreated))
			})

			By("Creating another link", func() {
				rec := httptest.NewRecorder()
				req := newLinkRequest(
					"my-another-link",
					"https://example.com/my-another-link",
					"",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusCreated))
			})

			By("Updating a link", func() {
				rec := httptest.NewRecorder()
				req := updateLinkRequest(
					"1",
					"my-updated-cool-link",
					"https://example.com/my-updated-cool-link",
					"add a comment",
				)

				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(rec.Body.String()).To(MatchJSON(`
			{
				"data": {
					"id": "1",
					"type": "links",
					"attributes": {
					  "comment": "add a comment",
					  "originalUrl": "https://example.com/my-updated-cool-link",
					  "shortName": "my-updated-cool-link"
					}
				}
			}
			`))
			})

			By("Should fail when giving an invalid url", func() {
				rec := httptest.NewRecorder()
				req := updateLinkRequest(
					"1",
					"my-updated-cool-link",
					"invalid url",
					"add a comment",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})

			By("Should fail when giving an invalid short name", func() {
				rec := httptest.NewRecorder()
				req := updateLinkRequest(
					"1",
					"inva lid$% short name",
					"https://example.com/my-updated-cool-link",
					"add a comment",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})

			By("Should fail when using another link's short name", func() {
				rec := httptest.NewRecorder()
				req := updateLinkRequest(
					"1",
					"my-another-link",
					"https://example.com/my-updated-cool-link",
					"add a comment",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})

			By("Should fail when object id is not the same as url id", func() {
				rec := httptest.NewRecorder()
				req := updateLinkRequest(
					"2",
					"my-another-link",
					"https://example.com/my-updated-cool-link",
					"add a comment",
				)
				req.URL.Path = "/api/links/1"
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusConflict))
			})

			By("Should fail when passing an invalid json", func() {
				rec := httptest.NewRecorder()
				data := []byte(`invalid json{}`)
				req, err := http.NewRequest("PATCH", "/api/links/1", bytes.NewReader(data))
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusNotAcceptable))
			})

			By("Should fail when passing an invalid json structure", func() {
				rec := httptest.NewRecorder()
				data := []byte(`{
				"wrong_json_structure": {
					"id": "2",
					"type": "links"
				}
			}`)
				req, err := http.NewRequest("PATCH", "/api/links/1", bytes.NewReader(data))
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusNotAcceptable))
			})
		})

		It("API Gets a link", func() {
			By("Creating a link", func() {
				rec := httptest.NewRecorder()
				req := newLinkRequest(
					"my-cool-link",
					"https://example.com/my-cool-link",
					"",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusCreated))
			})

			By("Should get a link by id", func() {
				rec := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/api/links/1", nil)
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(rec.Body.String()).To(MatchJSON(`
			{
				"data": {
					"id": "1",
					"type": "links",
					"attributes": {
					  "comment": "",
					  "originalUrl": "https://example.com/my-cool-link",
					  "shortName": "my-cool-link"
					}
				}
			}
			`))
			})

			By("Should fail when getting a missing link", func() {
				rec := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/api/links/100500", nil)
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusNotFound))

				rec = httptest.NewRecorder()
				req, err = http.NewRequest("GET", "/api/links/missing-link", nil)
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusNotFound))
			})
		})

		It("API Gets links", func() {
			By("Should get an empty link list initially", func() {
				rec := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/api/links", nil)
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(rec.Body.String()).To(MatchJSON(`
			{
			  "data": [],
			  "links": {
				"first": "/api/links?page[number]=1&page[size]=10"
			  },
			  "meta": {
				"links": 0
			  }
			}
			`))
			})

			By("Creating a link", func() {
				rec := httptest.NewRecorder()
				req := newLinkRequest(
					"my-cool-link",
					"https://example.com/my-cool-link",
					"",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusCreated))
			})

			By("Should get a link list", func() {
				rec := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/api/links", nil)
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(rec.Body.String()).To(MatchJSON(`
			{
			  "links": {
				"first": "/api/links?page[number]=1&page[size]=10"
			  },
			  "data": [
				{
				  "type": "links",
				  "id": "1",
				  "attributes": {
					"shortName": "my-cool-link",
					"originalUrl": "https://example.com/my-cool-link",
					"comment": ""
				  }
				}
			  ],
			  "meta": {
				"links": 1
			  }
			}
			`))
			})

			By("Creating more links", func() {
				var wg sync.WaitGroup
				wg.Add(10)
				for i := 1; i <= 10; i++ {
					go func(i int) {
						rec := httptest.NewRecorder()
						req := newLinkRequest(
							fmt.Sprintf("my-cool-link-%d", i),
							fmt.Sprintf("https://example.com/my-cool-link-%d", i),
							fmt.Sprintf("and now with a comment #%d", i),
						)
						api.Handler().ServeHTTP(rec, req)
						Expect(rec.Code).To(Equal(http.StatusCreated))
						wg.Done()
					}(i)
				}
				wg.Wait()
			})

			By("Should get a link list paginated", func() {
				rec := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/api/links?page[number]=2", nil)
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusOK))

				m := make(map[string]interface{})
				err = json.Unmarshal(rec.Body.Bytes(), &m)
				Expect(err).ToNot(HaveOccurred())
				Expect(m).To(HaveKey("data"))
				Expect(m["data"]).To(HaveLen(1))
				Expect(m).To(HaveKey("meta"))
				Expect(m["meta"].(map[string]interface{})).To(HaveKey("links"))
				Expect(m["meta"].(map[string]interface{})["links"]).To(Equal(float64(11)))
			})
		})

		It("API Deletes links", func() {
			By("Creating a link", func() {
				rec := httptest.NewRecorder()
				req := newLinkRequest(
					"my-cool-link",
					"https://example.com/my-cool-link",
					"",
				)
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusCreated))
			})

			By("Should delete a link", func() {
				rec := httptest.NewRecorder()
				req, err := http.NewRequest("DELETE", "/api/links/1", nil)
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusNoContent))
			})

			By("Should fail to get a deleted link", func() {
				rec := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/api/links/1", nil)
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusNotFound))
			})

			By("Should fail to delete a deleted link", func() {
				rec := httptest.NewRecorder()
				req, err := http.NewRequest("DELETE", "/api/links/1", nil)
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusNotFound))
			})

			By("Creating more links", func() {
				var wg sync.WaitGroup
				wg.Add(10)
				for i := 1; i <= 10; i++ {
					go func(i int) {
						rec := httptest.NewRecorder()
						req := newLinkRequest(
							fmt.Sprintf("my-cool-link-%d", i),
							fmt.Sprintf("https://example.com/my-cool-link-%d", i),
							fmt.Sprintf("and now with a comment #%d", i),
						)
						api.Handler().ServeHTTP(rec, req)
						Expect(rec.Code).To(Equal(http.StatusCreated))
						wg.Done()
					}(i)
				}
				wg.Wait()
			})

			By("Deleting more links in an async way", func() {
				var wg sync.WaitGroup
				wg.Add(10)

				items, _, err := linkStorage.PaginatedGetAll(1, 10)
				Expect(err).ToNot(HaveOccurred())
				ids := make([]string, 0, len(items))
				for _, item := range items {
					ids = append(ids, item.ID)
				}
				for _, id := range ids {
					go func(id string) {
						rec := httptest.NewRecorder()
						req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/links/%s", id), nil)
						Expect(err).ToNot(HaveOccurred())
						api.Handler().ServeHTTP(rec, req)
						Expect(rec.Code).To(Equal(http.StatusNoContent))
						wg.Done()
					}(id)
				}
				wg.Wait()
			})

			By("Should have empty link list", func() {
				rec := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/api/links", nil)
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusOK))

				m := make(map[string]interface{})
				err = json.Unmarshal(rec.Body.Bytes(), &m)
				Expect(err).ToNot(HaveOccurred())
				Expect(m).To(HaveKey("data"))
				Expect(m["data"]).To(HaveLen(0))
				Expect(m).To(HaveKey("meta"))
				Expect(m["meta"].(map[string]interface{})).To(HaveKey("links"))
				Expect(m["meta"].(map[string]interface{})["links"]).To(Equal(float64(0)))
			})

			By("Should fail to delete a non-existent link", func() {
				rec := httptest.NewRecorder()
				req, err := http.NewRequest("DELETE", "/api/links/100500", nil)
				Expect(err).ToNot(HaveOccurred())
				api.Handler().ServeHTTP(rec, req)
				Expect(rec.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	When("Using redirector service", func() {
		var handler echo.HandlerFunc
		var router *echo.Echo

		BeforeEach(func() {
			_, err := linkStorage.Insert(model.Link{
				ShortName:   "my-cool-link",
				OriginalUrl: "https://example.com/my-cool-link",
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
})
