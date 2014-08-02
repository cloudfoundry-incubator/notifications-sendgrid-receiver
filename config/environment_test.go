package config_test

import (
    "os"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
    variables := map[string]string{
        "UAA_HOST":           os.Getenv("UAA_HOST"),
        "UAA_CLIENT_ID":      os.Getenv("UAA_CLIENT_ID"),
        "UAA_CLIENT_SECRET":  os.Getenv("UAA_CLIENT_SECRET"),
        "CC_HOST":            os.Getenv("CC_HOST"),
        "VERIFY_SSL":         os.Getenv("VERIFY_SSL"),
        "ROOT_PATH":          os.Getenv("ROOT_PATH"),
        "NOTIFICATIONS_HOST": os.Getenv("NOTIFICATIONS_HOST"),
    }

    AfterEach(func() {
        for key, value := range variables {
            os.Setenv(key, value)
        }
    })

    Describe("SpaceMailer configuration", func() {
        It("loads the values when they are set", func() {
            os.Setenv("NOTIFICATIONS_HOST", "https://notifications.example.com")

            env := config.NewEnvironment()

            Expect(env.NotificationsHost).To(Equal("https://notifications.example.com"))
        })

        It("panics when the values are missing", func() {
            os.Setenv("NOTIFICATIONS_HOST", "")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())
        })
    })

    Describe("UAA configuration", func() {
        It("loads the values when they are set", func() {
            os.Setenv("UAA_HOST", "https://uaa.example.com")
            os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
            os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

            env := config.NewEnvironment()

            Expect(env.UAAHost).To(Equal("https://uaa.example.com"))
            Expect(env.UAAClientID).To(Equal("uaa-client-id"))
            Expect(env.UAAClientSecret).To(Equal("uaa-client-secret"))
        })

        It("panics when the values are missing", func() {
            os.Setenv("UAA_HOST", "")
            os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
            os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())

            os.Setenv("UAA_HOST", "https://uaa.example.com")
            os.Setenv("UAA_CLIENT_ID", "")
            os.Setenv("UAA_CLIENT_SECRET", "uaa-client-secret")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())

            os.Setenv("UAA_HOST", "https://uaa.example.com")
            os.Setenv("UAA_CLIENT_ID", "uaa-client-id")
            os.Setenv("UAA_CLIENT_SECRET", "")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())
        })
    })

    Describe("CloudController configuration", func() {
        It("loads the values when they are present", func() {
            os.Setenv("CC_HOST", "https://api.example.com")

            env := config.NewEnvironment()

            Expect(env.CCHost).To(Equal("https://api.example.com"))
        })

        It("panics when any of the values are missing", func() {
            os.Setenv("CC_HOST", "")

            Expect(func() {
                config.NewEnvironment()
            }).To(Panic())
        })
    })

    Describe("SSL verification configuration", func() {
        It("set the value to true by default", func() {
            os.Setenv("VERIFY_SSL", "")

            env := config.NewEnvironment()

            Expect(env.VerifySSL).To(BeTrue())
        })

        It("can be set to false", func() {
            os.Setenv("VERIFY_SSL", "false")

            env := config.NewEnvironment()

            Expect(env.VerifySSL).To(BeFalse())
        })

        It("can be set to true", func() {
            os.Setenv("VERIFY_SSL", "TRUE")

            env := config.NewEnvironment()

            Expect(env.VerifySSL).To(BeTrue())
        })

        It("sets the value to true if the value is non-boolean", func() {
            os.Setenv("VERIFY_SSL", "banana")

            env := config.NewEnvironment()

            Expect(env.VerifySSL).To(BeTrue())
        })
    })

    Describe("RootPath config", func() {
        It("loads the config value", func() {
            os.Setenv("ROOT_PATH", "bananaDAMAGE")
            env := config.NewEnvironment()

            Expect(env.RootPath).To(Equal("bananaDAMAGE"))
        })
    })
})
