package handlers

type FakeSpaceMailerAPI struct {
    Header    map[string]string
    SpaceGuid string
    Params    string
}

func (fake *FakeSpaceMailerAPI) PostToSpace(uaaAccessToken string, params map[string]string) error {
    fake.SpaceGuid, _ = parseSpaceGuid(params["to"])
    fake.Params = "blank"

    return nil
}

func FakeOAuth() (token string) {
    return "fakeTokenThatNeedsToBeCreatedOrSomething"
}
