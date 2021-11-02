package greeting

import "context"

type (
	ConnectMock struct {
		OutError error
	}

	DisconnectMock struct {
		InContext context.Context
		OutError  error
	}

	MockClient struct {
		StringOut string

		ConnectIndex int
		ConnectMocks []ConnectMock

		DisconnectIndex int
		DisconnectMocks []DisconnectMock
	}
)

func (m *MockClient) String() string {
	return m.StringOut
}

func (m *MockClient) Connect() error {
	i := m.ConnectIndex
	m.ConnectIndex++
	return m.ConnectMocks[i].OutError
}

func (m *MockClient) Disconnect(ctx context.Context) error {
	i := m.DisconnectIndex
	m.DisconnectIndex++
	m.DisconnectMocks[i].InContext = ctx
	return m.DisconnectMocks[i].OutError
}

type (
	HealthCheckMock struct {
		InContext context.Context
		OutError  error
	}

	MockChecker struct {
		StringOut string

		HealthCheckIndex int
		HealthCheckMocks []HealthCheckMock
	}
)

func (m *MockChecker) String() string {
	return m.StringOut
}

func (m *MockChecker) HealthCheck(ctx context.Context) error {
	i := m.HealthCheckIndex
	m.HealthCheckIndex++
	m.HealthCheckMocks[i].InContext = ctx
	return m.HealthCheckMocks[i].OutError
}

type (
	TranslateMock struct {
		InContext  context.Context
		InLanguage string
		InText     string
		OutString  string
		OutError   error
	}

	MockTranslateGateway struct {
		MockClient
		MockChecker

		StringOut string

		TranslateIndex int
		TranslateMocks []TranslateMock
	}
)

func (m *MockTranslateGateway) String() string {
	return m.StringOut
}

func (m *MockTranslateGateway) Translate(ctx context.Context, lang, text string) (string, error) {
	i := m.TranslateIndex
	m.TranslateIndex++
	m.TranslateMocks[i].InContext = ctx
	m.TranslateMocks[i].InLanguage = lang
	m.TranslateMocks[i].InText = text
	return m.TranslateMocks[i].OutString, m.TranslateMocks[i].OutError
}

type (
	StoreMock struct {
		InContext  context.Context
		InLanguage string
		InGreeting string
		OutError   error
	}

	LookupMock struct {
		InContext   context.Context
		InLanguage  string
		OutGreeting string
		OutError    error
	}

	MockGreetingCacheRepository struct {
		MockClient
		MockChecker

		StringOut string

		StoreIndex int
		StoreMocks []StoreMock

		LookupIndex int
		LookupMocks []LookupMock
	}
)

func (m *MockGreetingCacheRepository) String() string {
	return m.StringOut
}

func (m *MockGreetingCacheRepository) Store(ctx context.Context, lang, greeting string) error {
	i := m.StoreIndex
	m.StoreIndex++
	m.StoreMocks[i].InContext = ctx
	m.StoreMocks[i].InLanguage = lang
	m.StoreMocks[i].InGreeting = greeting
	return m.StoreMocks[i].OutError
}

func (m *MockGreetingCacheRepository) Lookup(ctx context.Context, lang string) (string, error) {
	i := m.LookupIndex
	m.LookupIndex++
	m.LookupMocks[i].InContext = ctx
	m.LookupMocks[i].InLanguage = lang
	return m.LookupMocks[i].OutGreeting, m.LookupMocks[i].OutError
}
