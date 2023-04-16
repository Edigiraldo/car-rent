package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type citiesDependencies struct {
	citiesService *mocks.MockCitiesService
}

func NewCitiesDependencies(citiesSrv *mocks.MockCitiesService) *citiesDependencies {
	return &citiesDependencies{
		citiesService: citiesSrv,
	}
}

func TestCitiesDelete(t *testing.T) {
	CitiesName := []string{"CityName A", "CityName B", "CityName C"}

	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		wants    wants
		setMocks func(*citiesDependencies)
	}{
		{
			name: "returns status code 200 when cities name list was successfully retrieved",
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *citiesDependencies) {
				d.citiesService.EXPECT().ListNames(gomock.Any()).Return(CitiesName, nil)
			},
		},
		{
			name: "returns 500 status code when there was a server error",
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *citiesDependencies) {
				d.citiesService.EXPECT().ListNames(gomock.Any()).Return(nil, errors.New("error getting cities name list"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			citiesSrv := mocks.NewMockCitiesService(mockCtlr)
			d := NewCitiesDependencies(citiesSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			urlObj, _ := url.Parse(baseURL + "cities/names")
			URL := urlObj.String()

			req, err := http.NewRequest(http.MethodDelete, URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			citiesHandler := NewCities(citiesSrv)
			citiesHandler.ListNames(rr, req)

			if rr.Result().StatusCode == http.StatusOK {
				body := dtos.ListCitiesNameResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), &body)
				if err != nil {
					t.Fatal(err)
				}
				if len(body.CitiesName) != 0 {
					assert.Equal(t, CitiesName[0], body.CitiesName[0])
					assert.Equal(t, len(CitiesName), len(body.CitiesName))
				}
			}

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}
