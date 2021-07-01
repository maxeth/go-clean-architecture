package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/maxeth/go-account-api/library"
	"github.com/maxeth/go-account-api/model"
	"github.com/maxeth/go-account-api/model/mocks"
	"github.com/stretchr/testify/require"
)

func TestSignup(t *testing.T) {
	//user := randomUser(t)
	pw := library.RandomString(15)
	email := "somemail@gmail.com"

	randomAT := "at"
	randomRT := "rt"

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(us *mocks.MockUserService, ts *mocks.MockTokenService)
		checkResponse func(resRec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    email,
				"password": pw,
			},
			buildStubs: func(us *mocks.MockUserService, ts *mocks.MockTokenService) {
				u := &model.User{
					Email:    email,
					Password: pw,
				}
				us.EXPECT().
					Signup(gomock.Any(), email, pw).
					Times(1).Return(u, nil)

				tp := &model.TokenPair{
					AccessToken:  randomAT,
					RefreshToken: randomRT,
				}
				ts.EXPECT().NewPairFromUser(gomock.Any(), gomock.Eq(u), "").Times(1).Return(tp, nil)
			},
			checkResponse: func(resRec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, resRec.Code)
				requireResponseBodyJWTMatch(t, resRec.Body, model.TokenPair{
					AccessToken:  randomAT,
					RefreshToken: randomRT,
				})
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			us := mocks.NewMockUserService(ctrl)
			ts := mocks.NewMockTokenService(ctrl)
			tc.buildStubs(us, ts)

			router := gin.Default()
			hc := Config{
				R:               router,
				UserService:     us,
				TokenService:    ts,
				TimeOutDuration: time.Duration(5 * time.Second),
			}
			NewHandler(&hc)

			recorder := httptest.NewRecorder()

			url := "/account/signup"

			body, err := json.Marshal(tc.body)
			//fmt.Println("passing body: ", string(body))
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(recorder, req)

			// individual response checker function
			tc.checkResponse(recorder)
		})
	}
}

func requireResponseBodyJWTMatch(t *testing.T, body *bytes.Buffer, tp model.TokenPair) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	fmt.Println("data is", string(data))

	type resStruct struct {
		Tokens model.TokenPair `json:"tokens"`
	}
	var gotRes resStruct
	//var test model.TokenPair
	err = json.Unmarshal(data, &gotRes)

	fmt.Println("got unmarshalled: ", gotRes)
	require.NoError(t, err)

	//err = json.Unmarshal(data, &test)
	//require.NoError(t, err)
	//fmt.Println("got data: ", test)

	require.Equal(t, tp.AccessToken, gotRes.Tokens.AccessToken)
	require.Equal(t, tp.RefreshToken, gotRes.Tokens.RefreshToken)
}

// require the response body of a server response matches a given user struct
func requireResponseBodyUserMatch(t *testing.T, body *bytes.Buffer, user model.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser model.User

	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.Email, gotUser.Email)

	// the hashed pw shouldnt be returned from the server
	require.Empty(t, gotUser.Password)
}
