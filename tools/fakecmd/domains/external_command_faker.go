package domains

type ExternalCommandFaker struct {
	filePathFakeCMD string
}

func (t *ExternalCommandFaker) Add(
	dirPathCommand DirPathFakeCommand,
	behaviors Behaviors,
) (*FakeCommand, error) {
	return &FakeCommand{
		filePathFakeCMD: t.filePathFakeCMD,
		dirPath:         dirPathCommand,
		behaviors:       behaviors,
	}, nil
}

func NewExternalCommandFaker(
	filePathFakeCMD string,
) *ExternalCommandFaker {
	return &ExternalCommandFaker{
		filePathFakeCMD: filePathFakeCMD,
	}
}
