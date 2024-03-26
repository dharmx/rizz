package main

type item struct {
	title    string
	filetype bool
}

func (file item) FilterValue() string {
	return file.title
}
