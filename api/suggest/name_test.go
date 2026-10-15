package suggest

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type requesterMock struct {
	method string
	params interface{}
}

func (m *requesterMock) Get(_ context.Context, apiMethod string, params url.Values, _ interface{}) error {
	m.method, m.params = apiMethod, params
	return nil
}

func (m *requesterMock) Post(_ context.Context, apiMethod string, params, _ interface{}) error {
	m.method, m.params = apiMethod, params
	return nil
}

func TestApi_Name_OmitsEmptyParams(t *testing.T) {
	mock := &requesterMock{}
	api := Api{Client: mock}

	_, err := api.Name(context.Background(), &RequestParams{Query: "Ив"})
	require.NoError(t, err)

	body, err := json.Marshal(mock.params)
	require.NoError(t, err)
	assert.Equal(t, "suggest/fio", mock.method)
	assert.JSONEq(t, `{"query": "Ив"}`, string(body))
}

func TestApi_Name_WithOptions(t *testing.T) {
	mock := &requesterMock{}
	api := Api{Client: mock}

	_, err := api.Name(context.Background(), &RequestParams{Query: "Ив", Count: 5},
		WithNameParts(NamePartName), WithNameGender(NameGenderMale))
	require.NoError(t, err)

	body, err := json.Marshal(mock.params)
	require.NoError(t, err)
	assert.JSONEq(t, `{"query": "Ив", "count": 5, "parts": ["NAME"], "gender": "MALE"}`, string(body))
}
