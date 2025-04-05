package api

// func TestWebsocketConnection(t *testing.T) {

// 	testCases := []struct {
// 		name          string
// 		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
// 		buildRequest  func() (*http.Request, error)
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
// 			},
// 			buildRequest: func() (*http.Request, error) {

// 				req, err := http.NewRequest(http.MethodGet, "/ws", nil)
// 				require.NoError(t, err)
// 				req.Header.Set("Content-Type", "application/json")
// 				return req, nil

// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {

// 			},
// 		},
// 	}

// }
