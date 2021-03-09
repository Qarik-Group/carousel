package action

type Action interface {
	Name() string
	Description() string
	//	Do() func() error
}
