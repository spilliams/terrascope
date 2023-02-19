package graphviz

type attributeKey string

const (
	LabelAttributeKey attributeKey = "label"
)

func (ak attributeKey) String() string {
	return string(ak)
}
