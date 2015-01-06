package application_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/application"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	PrivatePEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA1eUDU+9MaLnurRstr0z3iMeiUeYawSVZwSWv3LHRirLiJrsx
uSFSZUoqe2q9XQdW9PZZXe/sJHUx5fHSbTamXytWJF4e9hPTV4PHa1FI7DhXTikM
U49zK3++KPuquHNINyO/dFo4oQRd/LwUHNI1/D44wpGnLX/Fj/7mpJbBt76FM7AT
TV2o386NVzVsRZfh7dg3Nlqz7gaccOv4NQoOe6bAeb9Ev1rGVAE1J8R6RFg7a0QG
tUTVaHOeeQbJYWd/vHPqsA6+gFSp6RGklrCrUMhmxTFz8nDFY+8+Gw4Rpsj55Dg9
bRJe1ceAvHb3bbRDXtKdAnZntdad1WX7wiqc1QIDAQABAoIBAQCM80xdHF1ayePZ
mQZi9bJNJVj25U6OJwHdgOtB4L/3yWe0JDlhCO9WJiiWicFds39/D3PWrksLv1rx
b4i+RXwfTNyIPKnkeS4VBpb8RmVqnLoTnIMt8WLwsE5sjWN0Byv4ggshLm8Q1GLn
MD2EvJWssogO15K2LfPPryboIRxo6uQ2GzK4SU5Vkd3owY7+XV0ELKBwYl1RioU/
6xxnuA2YKFg4MJTTmF+kfDiSZQXzNVAxUXPupdlJQgXGzrK0szBSPRSk/cHoiQ8b
uNaxVospZBLDUm6Xe8VKqxB5BzC+HcZnfqw5Gb/ZZBJrcbtl9l825KSPrnPl7NER
iwGhF1JFAoGBAO/a7skDFc8hXbSL1J4mBGIRHu6/P/ll5AoTF2L0y88OW72L6GzA
ox6EwSKKTBwUYY2i/mrQ+BSRUm7BIeid7lbdtrYT16tuPtnLfXi7qa29WiWvgbUX
QvcC/dcBd5pkT7ZoohgSPNxiyYjE2Ux8XGdelwQ2l92WENUGLJP4qLg/AoGBAORK
vVBQ1ijKaFyIuunsvO6QnrkprYYyU5PN6lLtQ6z5xi/gLtcqABc4wiZo9wXzt/2t
nHXUUewx6qa5WvqHahLRA5kGcvyIITtZXzXwBBO6Z2suU6vB8Ukb+9u+sDYVKkYW
38+tUv5ZbqfPSGkR5kiNsUhtMqBo+b8rO9QuP8XrAoGBAIsvEk2W+rrdc9DnK5Qy
H20A3mBQnsEMfU3TUkjcIMCgZmARpeglyQJWqvRuKEhLE2jrYpN9e9gDlEAs1o5z
xvCla/cwgIA8U0BzMvYyf/4P2RXxSGVbgEJye/aeJVd0SkVhZl5tht+ke6pgAHC/
4aciXqPVQj8Va+MR2CBttQplAoGBANyrliVtjiWtyYUwsaRurw3Xg2WucMpYGUu0
7n0sVY99fOJITF61fZL0zU79hVIejMpMqAGJs4qhkZWJc/TZMmJv4Y9omXubRqws
rojfscE0HMWQ6VYMSWSHBUQbJg+RE+TeNYd0ndW4surIxdCyeavGMwi0bQx7jHYK
n3FxJznzAoGASTYzSg3QDM7rq6eg035r8zGT+PhQNOn3kBs12KXXEBdcNrhN0mUN
yeYSS/n4TwYAcD/YoYvgVqYGpZS2mBwPujoZpi2SIOdqjR4G5G0UUqSDZNYfLJ+i
o0u1iFbMKbAu97LkpGYLYCvprAa6JmEAVPFoy5wEgrTQJM2CLyMsrR4=
-----END RSA PRIVATE KEY-----`
	PublicPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1eUDU+9MaLnurRstr0z3
iMeiUeYawSVZwSWv3LHRirLiJrsxuSFSZUoqe2q9XQdW9PZZXe/sJHUx5fHSbTam
XytWJF4e9hPTV4PHa1FI7DhXTikMU49zK3++KPuquHNINyO/dFo4oQRd/LwUHNI1
/D44wpGnLX/Fj/7mpJbBt76FM7ATTV2o386NVzVsRZfh7dg3Nlqz7gaccOv4NQoO
e6bAeb9Ev1rGVAE1J8R6RFg7a0QGtUTVaHOeeQbJYWd/vHPqsA6+gFSp6RGklrCr
UMhmxTFz8nDFY+8+Gw4Rpsj55Dg9bRJe1ceAvHb3bbRDXtKdAnZntdad1WX7wiqc
1QIDAQAB
-----END PUBLIC KEY-----`
)

type NotificationsHandler struct {
	SourceName string
}

func (handler *NotificationsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Authorization")
	token, err := jwt.Parse(strings.TrimPrefix(header, "Bearer "), jwt.Keyfunc(func(token *jwt.Token) ([]byte, error) {
		return []byte(PublicPEM), nil
	}))
	if err != nil {
		panic(err)
	}

	scopes := token.Claims["scope"].([]interface{})
	if !handler.authenticated(scopes) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var body struct {
		SourceName string `json:"source_name"`
	}

	err = json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		panic(err)
	}

	handler.SourceName = body.SourceName

	w.WriteHeader(http.StatusNoContent)
}

func (handler NotificationsHandler) authenticated(scopes []interface{}) bool {
	for _, scope := range scopes {
		if scope.(string) == "notifications.write" {
			return true
		}
	}

	return false
}

func generateTokenWithScopes(scopes ...string) string {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token.Claims["client_id"] = "notifications-sendgrid-receiver"
	token.Claims["scope"] = scopes

	tokenString, err := token.SignedString([]byte(PrivatePEM))
	if err != nil {
		panic(err)
	}

	return tokenString
}

var _ = Describe("Registrar", func() {
	var handler *NotificationsHandler
	var server *httptest.Server

	BeforeEach(func() {
		handler = &NotificationsHandler{}
		router := mux.NewRouter()
		router.Handle("/notifications", handler).Methods("PUT")
		server = httptest.NewServer(router)
	})

	AfterEach(func() {
		server.Close()
	})

	It("registers the service with notifications", func() {
		verifySSL := true
		registrar := application.NewRegistrar(server.URL, verifySSL)

		token := generateTokenWithScopes("notifications.write")
		err := registrar.Register(token)
		Expect(err).NotTo(HaveOccurred())
		Expect(handler.SourceName).To(Equal("3rd Party Services"))
	})
})
