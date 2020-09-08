package container

type mockContainer struct {
	has     bool
	service interface{}
	error   error
}

func (m mockContainer) Get(string) (interface{}, error) {
	return m.service, m.error
}

func (m mockContainer) MustGet(i string) interface{} {
	s, err := m.Get(i)
	if err != nil {
		panic(err)
	}
	return s
}

func (m mockContainer) Has(string) bool {
	return m.has
}
