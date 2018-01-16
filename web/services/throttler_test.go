package services_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Throttler", func() {
	It("returns true when the less than maxRequests are active", func() {
		throttler := services.NewThrottler(1)

		Expect(throttler.Throttle()).To(BeFalse())
		Expect(throttler.Throttle()).To(BeTrue())

		throttler.Finish()
		Expect(throttler.Throttle()).To(BeFalse())
	})

	It("is thread safe", func(done Done) {
		throttler := services.NewThrottler(2)
		Expect(throttler.Throttle()).To(BeFalse())
		go func() {
			time.Sleep(time.Second)
			Expect(throttler.Throttle()).To(BeTrue())
			close(done)
		}()
		Expect(throttler.Throttle()).To(BeFalse())
	}, 2)
})
