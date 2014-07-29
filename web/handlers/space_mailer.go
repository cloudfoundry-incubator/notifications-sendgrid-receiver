package handlers

type SpaceMailerApiInterface interface {
    PostToSpace(map[string]string) error
}

type SpaceMailerApi struct{}

func (api *SpaceMailerApi) PostToSpace(params map[string]string) error {
    return nil
}
