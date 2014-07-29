package handlers

type FakeSpaceMailerApi struct {
    SpaceGuid string
    Params    string
}

func (fake *FakeSpaceMailerApi) PostToSpace(params map[string]string) error {
    fake.SpaceGuid = parseSpaceGuid(params["to"])
    fake.Params = "blank"

    return nil
}
