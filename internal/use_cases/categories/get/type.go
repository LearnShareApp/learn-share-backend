package get

type response struct {
	Categories []category `json:"categories"`
}

type category struct {
	Name   string `json:"name"`
	MinAge int    `json:"min_age"`
}
