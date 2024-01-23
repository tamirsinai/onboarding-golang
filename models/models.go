package models

type Input struct {
	CloneUrl string `json:"clone_url"`
	Size     int `json:"size"`
}

type File struct {
	Name string
	Size int
}

type Scan struct {
	Total int
	Files []File
}
