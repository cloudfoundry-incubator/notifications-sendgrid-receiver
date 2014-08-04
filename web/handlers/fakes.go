package handlers

type FakeNotificationsReceiver struct {
    Header    map[string]string
    SpaceGuid string
    Params    string
}

func (fake *FakeNotificationsReceiver) PostToSpace(uaaAccessToken string, params map[string]string) error {
    fake.SpaceGuid, _ = params["to"]
    fake.Params = "blank"

    return nil
}

func FakeOAuth() (token string) {
    return "fakeTokenThatNeedsToBeCreatedOrSomething"
}
