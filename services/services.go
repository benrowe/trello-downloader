package services

// Service baseline service
type Service struct {
	Name      string
	baseURL   string
	endpoints struct {
		search   string
		add      string
		callback string
	}
}

func (v *Service) DoSomething(value string) {
	v.Name = value
}

func (v Service) String() string {
	return v.Name
}
