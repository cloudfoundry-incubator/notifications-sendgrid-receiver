package application_test

import (
	"os"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/application"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
	var variables map[string]string

	BeforeEach(func() {
		variables = map[string]string{
			"BASIC_AUTH_PASSWORD":  os.Getenv("BASIC_AUTH_PASSWORD"),
			"BASIC_AUTH_USER_NAME": os.Getenv("BASIC_AUTH_USER_NAME"),
			"CC_HOST":              os.Getenv("CC_HOST"),
			"LOG_FILE":             os.Getenv("LOG_FILE"),
			"NOTIFICATIONS_HOST":   os.Getenv("NOTIFICATIONS_HOST"),
			"PORT":                 os.Getenv("PORT"),
			"UAA_CLIENT_ID":        os.Getenv("UAA_CLIENT_ID"),
			"UAA_CLIENT_SECRET":    os.Getenv("UAA_CLIENT_SECRET"),
			"UAA_HOST":             os.Getenv("UAA_HOST"),
			"VERIFY_SSL":           os.Getenv("VERIFY_SSL"),
		}
	})

	AfterEach(func() {
		for key, value := range variables {
			os.Setenv(key, value)
		}
	})

	Describe("Notifications configuration", func() {
		It("loads the values when they are set", func() {
			os.Setenv("NOTIFICATIONS_HOST", "https://notifications.example.com")

			env := application.NewEnvironment()

			Expect(env.NotificationsHost).To(Equal("https://notifications.example.com"))
		})

		It("panics when the values are missing", func() {
			os.Setenv("NOTIFICATIONS_HOST", "")

			Expect(func() {
				application.NewEnvironment()
			}).To(Panic())
		})
	})

	Describe("UAA configuration", func() {
		It("loads the values when they are set", func() {
			os.Setenv("UAA_HOST", "https://uaa.example.com")
			os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
			os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

			env := application.NewEnvironment()

			Expect(env.UAAHost).To(Equal("https://uaa.example.com"))
			Expect(env.UAAClientID).To(Equal("uaa-client-id"))
			Expect(env.UAAClientSecret).To(Equal("uaa-client-secret"))
		})

		It("panics when the values are missing", func() {
			os.Setenv("UAA_HOST", "")
			os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
			os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

			Expect(func() {
				application.NewEnvironment()
			}).To(Panic())

			os.Setenv("UAA_HOST", "https://uaa.example.com")
			os.Setenv("UAA_CLIENT_ID", "")
			os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

			Expect(func() {
				application.NewEnvironment()
			}).To(Panic())

			os.Setenv("UAA_HOST", "https://uaa.example.com")
			os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
			os.Setenv("UAA_CLIENT_SECRET", "")

			Expect(func() {
				application.NewEnvironment()
			}).To(Panic())
		})
	})

	Describe("CloudController configuration", func() {
		It("loads the values when they are present", func() {
			os.Setenv("CC_HOST", "https://api.example.com")

			env := application.NewEnvironment()

			Expect(env.CCHost).To(Equal("https://api.example.com"))
		})

		It("panics when any of the values are missing", func() {
			os.Setenv("CC_HOST", "")

			Expect(func() {
				application.NewEnvironment()
			}).To(Panic())
		})
	})

	Describe("Port configuration", func() {
		It("sets the value when it is present", func() {
			os.Setenv("PORT", "5555")
			env := application.NewEnvironment()

			Expect(env.Port).To(Equal("5555"))
		})

		It("panics if the value is not set", func() {
			os.Setenv("PORT", "")

			Expect(func() {
				application.NewEnvironment()
			}).To(Panic())
		})
	})

	Describe("SSL verification configuration", func() {
		It("set the value to true by default", func() {
			os.Setenv("VERIFY_SSL", "")

			env := application.NewEnvironment()

			Expect(env.VerifySSL).To(BeTrue())
		})

		It("can be set to false", func() {
			os.Setenv("VERIFY_SSL", "false")

			env := application.NewEnvironment()

			Expect(env.VerifySSL).To(BeFalse())
		})

		It("can be set to true", func() {
			os.Setenv("VERIFY_SSL", "TRUE")

			env := application.NewEnvironment()

			Expect(env.VerifySSL).To(BeTrue())
		})
	})

	Describe("Basic auth configuration", func() {
		It("loads the username and password", func() {
			os.Setenv("BASIC_AUTH_USER_NAME", "happy")
			os.Setenv("BASIC_AUTH_PASSWORD", "glad")

			env := application.NewEnvironment()

			Expect(env.BasicAuthUsername).To(Equal("happy"))
			Expect(env.BasicAuthPassword).To(Equal("glad"))
		})
	})
})
