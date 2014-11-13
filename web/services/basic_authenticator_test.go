package services_test

import (
	"encoding/base64"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BasicAuthenticator", func() {
	var auth services.BasicAuthenticator
	var header http.Header

	BeforeEach(func() {
		auth = services.NewBasicAuthenticator("user", "password")
		header = make(map[string][]string)
	})

	It("Returns true when the basic auth credentials are correct", func() {
		encoded := base64.StdEncoding.EncodeToString([]byte("user:password"))
		header.Set("Authorization", "Basic "+encoded)
		result := auth.Verify(header)
		Expect(result).To(Equal(true))
	})

	It("Returns false when the credentials do not match the correct username and password", func() {
		encoded := base64.StdEncoding.EncodeToString([]byte("username:Incorrectpassword"))
		header.Set("Authorization", "Basic "+encoded)
		result := auth.Verify(header)
		Expect(result).To(Equal(false))
	})

	Context("The basic auth is improperly formatted", func() {
		It("Returns false when it cannot be base64 decoded", func() {
			header.Set("Authorization", "Basic This is invalid")
			result := auth.Verify(header)
			Expect(result).To(Equal(false))
		})

		It("Returns false when the basic auth is not formatted as username:password", func() {
			encoded := base64.StdEncoding.EncodeToString([]byte("thereisnocolon"))
			header.Set("Authorization", "Basic "+encoded)
			result := auth.Verify(header)
			Expect(result).To(Equal(false))
		})
	})
})
