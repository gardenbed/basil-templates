package translate

import "net/http"

type (
	DoMock struct {
		InRequest   *http.Request
		OutResponse *http.Response
		OutError    error
	}

	MockHTTPClient struct {
		DoIndex int
		DoMocks []DoMock
	}
)

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	i := m.DoIndex
	m.DoIndex++
	m.DoMocks[i].InRequest = req
	return m.DoMocks[i].OutResponse, m.DoMocks[i].OutError
}
