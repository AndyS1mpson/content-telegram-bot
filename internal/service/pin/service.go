package pin

// Service сервис, содержащий логику по работе с пинами
type Service struct {
	parser     parser
	repository repository
}

// NewService конструктор для Service
func NewService(parser parser, repository repository) *Service {
	return &Service{
		parser:     parser,
		repository: repository,
	}
}
