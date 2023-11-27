package redirect_test

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"urlC/internal/http-server/handlers/url/redirect"
	"urlC/internal/http-server/handlers/url/redirect/mocks"
	"urlC/internal/lib/api"
	"urlC/internal/lib/logger/handlers/slogdiscard"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/assert.v1"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()
			fmt.Println(ts.URL)
			fmt.Println(tc.alias)
			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			//TODO:why status ok error, no returning value(redirectedToURL)
			require.NoError(t, err)

			// Check the final URL after redirection.

			assert.Equal(t, tc.url, redirectedToURL)
		})
	}
}
