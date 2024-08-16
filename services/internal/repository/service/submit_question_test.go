package service_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	cfgGlobal "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/domain/service"
	mockPkg "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/mocks"
	"github.com/stretchr/testify/mock"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/entity"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/config"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/service"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/test/testdata"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_chatRepository_SubmitQuestion(t *testing.T) {
	var (
		ctx     = context.Background()
		reqData = &entity.ChatBotReq{
			UserId:         "533be47a-9091-4524-afc9-38a9c32447e1",
			UserName:       "Yor",
			Query:          "Test?",
			ConversationId: "533be47a-9091-4524-afc9-38a9c32447e5",
		}
		cfg = config.ChatConfig{
			Port:          ":8080",
			Timeout:       2,
			Url:           "localhost:9090",
			Authorization: "12121213",
		}
		response = `{
			"event": "message",
			"task_id": "56fd9d12-f1f3-4c0b-9894-f65da887b74b",
			"id": "bc1bf28a-3e36-4553-a845-a0a1d2fba314",
			"message_id": "bc1bf28a-3e36-4553-a845-a0a1d2fba314",
			"conversation_id": "b9745476-2e80-44c9-bc4f-6643ee1957f2",
			"mode": "advanced-chat",
			"answer": "Ini contoh jawaban",
			"metadata": {
				"retriever_resources": [
					{
						"position": 1,
						"dataset_id": "af85792b-8418-4232-ab49-5e2fd71ee589",
						"dataset_name": "tSurvey-HRBOT",
						"document_id": "3e41bfc2-328b-4f3f-8695-086135ab15f9",
						"document_name": "tsurvey_id_legacy_prompt_completion.xlsx",
						"data_source_type": "upload_file",
						"segment_id": "6becaf03-93cb-44d1-93c7-51c1db5b9007",
						"retriever_from": "workflow",
						"score": null,
						"hit_count": 17,
						"word_count": 2366,
						"segment_position": 12,
						"index_node_hash": "deb598c6fc4abca2adc1e7f55ffa3b186ed4930d7e0f672bfeb90443a7afdf29",
						"content": "com, WWE | 373 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Tech? | Technologynewssites, UpdatesForSamsung | 374 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Travel? | airbnb, traveloka, hotelbvlgari | 375 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest TV? | hulu, viu, mncnow | 376 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Video? | youtube, tiktok, bigo | 377 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Wedding? | bridestory, weddingku, weddingsites | 378 | Data Knowledge |\n| Apa pengertian kolom target audience interest? | Dikatakan interest jika jumlah volume dan frekuensi akses melebihi jumlah rata-rata dibandingkan populasi | 379 | Data Knowledge |\n| Apa pengertian kolom target audience Top Categories? | Kategori terbanyak yang diakses dalam 1 bulan terakhir | 380 | Data Knowledge |\n| Apa pengertian kolom target audience Top Apps? | Aplikasi terbanyak yang diakses dalam 1 bulan terakhir | 381 | Data Knowledge |\n| Apa pengertian kolom target audience Age? | Prediction of Age based on Telco Features Modelling | 382 | Data Knowledge |\n| Apa pengertian kolom target audience Gender? | Prediction of Gender based on Telco Features Modelling | 383 | Data Knowledge |\n| Apa pengertian kolom target audience SES? | Prediction of Socio-economic and expenditure level based on telco features modelling | 384 | Data Knowledge |\n| Apa pengertian kolom target audience Mobility Segment? | Audience segmentation based on mobility behaviour, in reference of BPS dictionary (example: Commuter, Circular, etc) | 385 | Data Knowledge |\n| Apa pengertian kolom target audience Visit city? | kabupaten/kota yang audience lewati (staypoint) minimal selama 15 menit dalam sebulan terakhir | 386 | Data Knowledge |\n| Apa pengertian kolom target audience Visit kecamatan? | kecamatan yang audience lewati (staypoint) minimal selama 15 menit dalam sebulan terakhir | 387 | Data Knowledge |\n| Apa pengertian kolom target audience Visit kelurahan? | kelurahan yang audience lewati (staypoint) minimal selama 15 menit dalam sebulan terakhir | 388 | Data Knowledge |\n| Apa pengertian kolom target audience Visit province? | provinsi yang audience lewati (staypoint) minimal selama 15 menit dalam sebulan terakhir | 389 | Data Knowledge |\n| Apa pengertian kolom"
					},
					{
						"position": 1,
						"dataset_id": "af85792b-8418-4232-ab49-5e2fd71ee589",
						"dataset_name": "tSurvey-HRBOT",
						"document_id": "3e41bfc2-328b-4f3f-8695-086135ab15f9",
						"document_name": "tsurvey_id_legacy_prompt_completion.xlsx",
						"data_source_type": "upload_file",
						"segment_id": "be5f9ac2-2094-4331-8194-619e6d849c1e",
						"retriever_from": "workflow",
						"score": null,
						"hit_count": 13,
						"word_count": 2247,
						"segment_position": 9,
						"index_node_hash": "2b9c6293509b77466ddcc3a67c22c1cb5ff3d52bfeb2cb74569d491312d322da",
						"content": "pada interest Fitness? | NikeRun, Guovapass, 30DaysFitnessChallenge | 348 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Foodie? | chope, cookpad, graved | 349 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Gadget? | jakartanotebook, tokopda, okeshop | 350 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Game? | ClashofClans, MobileLegends, AngryBirds | 351 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Gardening? | Kebunbibit, GardeningSites, WitandWhistle | 352 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Homedecor? | Ikea, houseofmarble | 353 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Insurance? | chubb, axa, generali | 354 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Islam? | Alquranbahasaindonesia, jadwalsholat, kitabmuslim | 355 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Job Seek? | jobsdb, glassdoor, jobsid | 356 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Language? | duolingo, kamusku, oxfordadvanceddictionary | 357 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Medicine? | bpjskesehatan, tanyadok, webrnd | 358 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Movie? | cinema21-official, cinemaxx, robinsonsmovieworld | 359 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Music? | applemusic, indonesianjazzfestival, spotify | 360 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest News? | CNBC, Detik, Jakartaglobe | 361 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Online Banking? | bankmandiri, bankmega, bankbtpn | 362 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Online Transport? | go-jek, grab, mybluebird | 363 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Pets? | tokohewonpeliharaan, KucingMonia, PetSites | 364 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Photo? | AsosiasiDesainerGrafisIndonesia, PhotographySites | 365 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Politics? | AdvocacyGroupSites | 366 | Data Knowledge"
					},
					{
						"position": 1,
						"dataset_id": "af85792b-8418-4232-ab49-5e2fd71ee589",
						"dataset_name": "tSurvey-HRBOT",
						"document_id": "3e41bfc2-328b-4f3f-8695-086135ab15f9",
						"document_name": "tsurvey_id_legacy_prompt_completion.xlsx",
						"data_source_type": "upload_file",
						"segment_id": "e8382e16-52c5-4bef-835a-fd7b7f252772",
						"retriever_from": "workflow",
						"score": null,
						"hit_count": 14,
						"word_count": 2222,
						"segment_position": 8,
						"index_node_hash": "b1a5f4ecdc03c7695f2091aea227fe0150f2945928db5f564d019ee72043d74f",
						"content": "com, jualbelimobill23 | 332 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Autorace? | hondarnotogp, motogp, formulal | 333 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Beauty? | Femoledaily, sociolla, Wardahkosmetikcciline | 334 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Book? | cribs, Kobobooks, kindle | 335 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Cashless? | Sokuku, linkaja, paytren | 336 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Chatting? | WeChat, Whatsapp, Telegram | 337 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Christianity? | eBible, AlkitobAnakKomik, CotholicDailyReflections | 338 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest College? | InstilITB, LlniversitasNegeriSuraboya, UniversityofTexosAustin | 339 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Comic Animation? | Lovers MangaRock, Animeheaven, Crunchyroll | 340 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Cooking? | BukuResepJajanan, NasokApa, ResepSayvrSeharihori | 341 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Dating? | ZooskDatingApp, DatingSites, KotokataCinta | 342 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Drama? | Kdromaindo, dromotplus, daebokdrarna | 343 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Ecommerce? | Shopper Bukalapak, JDID, Tokopedia | 344 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Family-Oriented? | motherandbaby, parentingsites, nurseryrhymesvideo | 345 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Fashion? | FIBM, pesonohijob, rnenstylefoshion | 346 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Finance? | Enthusiasts forbes, sahobatfinansialkeluargo, osiobankersclub | 347 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Fitness? | NikeRun, Guovapass, 30DaysFitnessChallenge | 348 | Data Knowledge |\n| Apa contoh aplikasi yg digunakan pada interest Foodie? | chope, cookpad, graved | 349 | Data Knowledge |\n| Apa contoh aplikasi yg"
					}
				],
				"usage": {
					"prompt_tokens": 3456,
					"prompt_unit_price": "0.003",
					"prompt_price_unit": "0.001",
					"prompt_price": "0.0103680",
					"completion_tokens": 223,
					"completion_unit_price": "0.004",
					"completion_price_unit": "0.001",
					"completion_price": "0.0008920",
					"total_tokens": 3679,
					"total_price": "0.0112600",
					"currency": "USD",
					"latency": 4.971822835505009
				}
			},
			"created_at": 1713952679
		}`
		// responseWorkFlowStarted = `data: {\"event\": \"workflow_started\", \"conversation_id\": \"ea925695-80e0-43e2-b31e-df89fbd394a0\", \"message_id\": \"75253d6b-26e7-4fa4-807e-421e6940bdf1\", \"created_at\": 1714376166, \"task_id\": \"0f55f583-0a58-4d9a-ac08-d68ae362d027\", \"workflow_run_id\": \"08fbeb52-1aca-4da4-bd3c-c2596022de3c\", \"data\": {\"id\": \"08fbeb52-1aca-4da4-bd3c-c2596022de3c\", \"workflow_id\": \"5846d192-0f16-4f30-906b-7cee63bd46f8\", \"sequence_number\": 215, \"inputs\": {\"name\": \"Yor\", \"sys.query\": \"Display logic\", \"sys.files\": []}, \"created_at\": 1714376166}}\n`
	)

	fakeClientErrInternal := &http.Client{Transport: testdata.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(response)),
		}
	})}

	fakeClientOK := &http.Client{Transport: testdata.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(response)),
		}
	})}

	// fakeClientOKWorkFlowStarted := &http.Client{Transport: testdata.RoundTripFunc(func(req *http.Request) *http.Response {
	// 	return &http.Response{
	// 		StatusCode: http.StatusOK,
	// 		Body:       io.NopCloser(strings.NewReader(responseWorkFlowStarted)),
	// 	}
	// })}

	redis := new(mockPkg.RedisClient)
	redis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	type fields struct {
		client      *http.Client
		redisClient cfgGlobal.RedisClient
	}
	type args struct {
		ctx     context.Context
		reqData *entity.ChatBotReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error - process hit",
			fields: fields{
				client:      fakeClientErrInternal,
				redisClient: redis,
			},
			args: args{
				ctx:     ctx,
				reqData: reqData,
			},
			wantErr: true,
		},
		// {
		// 	name: "error - workflow started - error redis set",
		// 	fields: fields{
		// 		client: fakeClientOKWorkFlowStarted,
		// 	},
		// 	args: args{
		// 		ctx:     ctx,
		// 		reqData: reqData,
		// 	},
		// 	wantErr: true,
		// },
		{
			name: "success - ready 200",
			fields: fields{
				client:      fakeClientOK,
				redisClient: redis,
			},
			args: args{
				ctx:     ctx,
				reqData: reqData,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatRepoFactory := service.ChatRepositoryFactory{
				Cfg:         cfg,
				Logger:      logger.NewFakeApiLogger(),
				Client:      tt.fields.client,
				RedisClient: tt.fields.redisClient,
			}

			repo, _ := chatRepoFactory.Create()
			err := repo.SubmitQuestion(tt.args.ctx, tt.args.reqData)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
