package services

type Service struct {
}

func NewService(repos ...interface{}) *Service {
	service := &Service{}
	for _, repo := range repos {
		switch repo.(type) {

		default:
			panic("Unknown repository type provided")
		}
	}
	return service
}
