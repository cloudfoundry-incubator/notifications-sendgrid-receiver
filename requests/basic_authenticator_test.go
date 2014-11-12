package requests_test

import (
	"encoding/base64"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/requests"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BasicAuthenticator", func() {
	var BasicAuthenticator requests.BasicAuthenticator
	var Header http.Header

	BeforeEach(func() {
		BasicAuthenticator = requests.NewBasicAuthenticator()
		Header = make(map[string][]string)
	})

	It("Returns true when the basic auth credentials are correct", func() {
		auth := base64.StdEncoding.EncodeToString([]byte("test_user:password"))
		Header.Set("Authorization", "Basic "+auth)
		result := BasicAuthenticator.Verify(Header)
		Expect(result).To(Equal(true))
	})

	It("Returns false when the credentials do not match the correct username and password", func() {
		auth := base64.StdEncoding.EncodeToString([]byte("username:Incorrectpassword"))
		Header.Set("Authorization", "Basic "+auth)
		result := BasicAuthenticator.Verify(Header)
		Expect(result).To(Equal(false))
	})

	Context("The basic auth is improperly formatted", func() {
		It("Returns false when it cannot be base64 decoded", func() {
			Header.Set("Authorization", "Basic This is invalid")
			result := BasicAuthenticator.Verify(Header)
			Expect(result).To(Equal(false))
		})

		It("Returns false when the basic auth is not formatted as username:paswword", func() {
			auth := base64.StdEncoding.EncodeToString([]byte("thereisnocolon"))
			Header.Set("Authorization", "Basic "+auth)
			result := BasicAuthenticator.Verify(Header)
			Expect(result).To(Equal(false))
		})
	})
})
