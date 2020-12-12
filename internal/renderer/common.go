package renderer

import (
	"fmt"
	"io/ioutil"
)

type RenderableType string

const (
	Text       RenderableType = "Text"
	GoTemplate                = "GoTemplate"
)

type Renderable struct {
	Type    RenderableType
	Content *string
	Path    string
}

func NewRenderableWithContent(rtype RenderableType, content string) *Renderable {
	return &Renderable{Type: rtype, Content: &content}
}

func NewRenderableWithPath(rtype RenderableType, path string) *Renderable {
	return &Renderable{Type: rtype, Path: path}
}

func (r *Renderable) Render() error {
	if r.Content == nil {
		if r.Path == "" {
			return fmt.Errorf("Invalid path \"%s\". A valid path is required when Content is undefined.", r.Path)
		}
		c, err := ioutil.ReadFile(r.Path)
		if err != nil {
			return err
		}
		s := string(c)
		r.Content = &s
	}

	switch r.Type {
	case Text:
		fmt.Print(*r.Content)
	case GoTemplate:
		result, err := renderGoTemplate(*r.Content)
		if err != nil {
			return err
		}
		fmt.Print(result)
	}
	return nil
}
