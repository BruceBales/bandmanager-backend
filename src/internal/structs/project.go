package structs

type File struct {
	Name string
	Type string
	Path string //identifies the physical storage path to the resource. Might end up being the web URL.
}

type Folder struct {
	Name  string
	Files []File
}

type project struct {
	BandID      int //Correlates to ID from band struct
	Name        string
	Description string
	Type        string
	Folders     []Folder
}
