package get_categories

// @Description get categories response
type response struct {
	Categories []category `json:"categories"`
}

// @Description data of category
type category struct {
	Id     int    `json:"id" example:"1"`
	Name   string `json:"name" example:"Programing"`
	MinAge int    `json:"min_age" example:"12"`
}
