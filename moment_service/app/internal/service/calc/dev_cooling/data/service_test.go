package data

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/Alexander272/sealur/moment_service/internal/config"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur/moment_service/pkg/database/postgres"
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

// Думаю это очень корявые тесты, но пока так пойдет
func Test_DataService(t *testing.T) {
	if err := gotenv.Load("../../../../../../../.env"); err != nil {
		t.Fatalf("error loading env variables")
	}
	conf, err := config.Init("../../../../../configs")
	if err != nil {
		t.Fatalf("error initializing configs")
	}
	logger.Init(os.Stdout, conf.Environment)

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     conf.Postgres.Host,
		Port:     conf.Postgres.Port,
		Username: conf.Postgres.Username,
		Password: conf.Postgres.Password,
		DBName:   conf.Postgres.DbName,
		SSLMode:  conf.Postgres.SSLMode,
	})
	if err != nil {
		t.Log(err)
		t.Fatalf("failed to initialize db")
	}

	//* Services, Repos & API Handlers

	repos := repository.NewRepo(db)

	flange := flange.NewFlangeService(repos.Flange)
	materials := materials.NewMaterialsService(repos.Materials)
	gasket := gasket.NewGasketService(repos.Gasket)
	service := NewDataService(flange, materials, gasket, graphic.NewGraphicService())

	t.Run("gasket", func(t *testing.T) {
		type testCase struct {
			name string
			arg  *dev_cooling_model.GasketData
			want struct {
				str *dev_cooling_model.GasketResult
				t   dev_cooling_model.GasketData_Type
			}
		}

		cases := []testCase{
			{
				name: "standart data",
				arg: &dev_cooling_model.GasketData{
					EnvId:     "1",
					GasketId:  "1",
					Thickness: 2.3,
					Width:     2,
					SizeLong:  14,
					SizeTrans: 14,
				},
				want: struct {
					str *dev_cooling_model.GasketResult
					t   dev_cooling_model.GasketData_Type
				}{
					str: &dev_cooling_model.GasketResult{
						Gasket:          "СНП",
						Env:             "Жидкость",
						Thickness:       2.3,
						Width:           2,
						SizeLong:        14,
						SizeTrans:       14,
						M:               3,
						Pres:            69,
						Compression:     1,
						Epsilon:         2000,
						PermissiblePres: 400,
						Type:            "Metal",
					},
					t: dev_cooling_model.GasketData_Metal,
				},
			},
			{
				name: "not standart data",
				arg: &dev_cooling_model.GasketData{
					GasketId:  "another",
					Thickness: 2.3,
					Width:     2,
					SizeLong:  14,
					SizeTrans: 14,
					Data: &dev_cooling_model.GasketData_Data{
						Title:           "СНП",
						Type:            dev_cooling_model.GasketData_Metal,
						Qo:              69,
						M:               3,
						Compression:     1,
						PermissiblePres: 400,
						Epsilon:         2000,
					},
				},
				want: struct {
					str *dev_cooling_model.GasketResult
					t   dev_cooling_model.GasketData_Type
				}{
					str: &dev_cooling_model.GasketResult{
						Gasket:          "СНП",
						Thickness:       2.3,
						Width:           2,
						SizeLong:        14,
						SizeTrans:       14,
						M:               3,
						Pres:            69,
						Compression:     1,
						Epsilon:         2000,
						PermissiblePres: 400,
						Type:            "Metal",
					},
					t: dev_cooling_model.GasketData_Metal,
				},
			},
		}

		for _, tc := range cases {
			gotSt, gotType, err := service.getGasketData(context.TODO(), tc.arg)
			if err != nil {
				t.Fatalf(err.Error())
			}

			if gotType != tc.want.t {
				t.Log(gotType, " != ", tc.want.t)
				t.Fatalf("expected type doesn't match the result")
			}

			if !reflect.DeepEqual(gotSt, tc.want.str) {
				t.Fatalf("expected value doesn't match the result")
			}
		}
	})
}
