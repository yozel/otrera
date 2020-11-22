package gatherer

type RawObjectInterface interface {
	Name() string
	Content() interface{}
	Copy() RawObject
}

type RawObject struct {
	IName    string
	IContent interface{}
}

func (r *RawObject) Name() string {
	return r.IName
}

func (r *RawObject) Content() interface{} {
	return r.IContent
}

func (r *RawObject) Copy() RawObject {
	return RawObject{r.Name(), r.Content()}
}
