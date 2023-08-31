package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"segment/pkg/config"
	"segment/pkg/delivery/router"
	"segment/pkg/dtos"
	"segment/pkg/models"
	"segment/pkg/repo"
	"segment/pkg/storage"
	"segment/pkg/storage/postgres"
	"sort"
	"testing"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

func init() {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal("No .env file found")
	}
}

func InitializeStorage(cfg *config.Config) (storage.IStorage, error) {
	db := postgres.NewStorage(cfg.Database)
	if err := db.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func TestSegmentHandler(t *testing.T) {
	suite.Run(t, &TestSuiteHandler{})
}

type TestSuiteHandler struct {
	suite.Suite

	storage storage.IStorage
	server  *httptest.Server
}

func (s *TestSuiteHandler) SetupTest() {
	ctx := context.Background()
	var cfg config.Config
	err := confita.NewLoader(
		env.NewBackend(),
	).Load(ctx, &cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse config: %v", err))
	}
	cfg.Database.Host = "localhost"
	log.Println(cfg)
	db, err := InitializeStorage(&cfg)
	for err != nil {
		db, err = InitializeStorage(&cfg)
	}
	if err != nil {
		log.Fatal(fmt.Errorf("failed to init database: %v", err))
	}
	s.storage = db
	err = db.MakeMigrations()
	if err != nil {
		log.Fatalf("failed to make migrations: %v", err)
	}
	repo := repo.NewRepository(db)
	s.server = httptest.NewServer(router.SetupRouter(repo))
}

func (s *TestSuiteHandler) TearDownTest() {
	if s.storage == nil {
		return
	}
	db := s.storage.Init()
	db.Model(&models.User{}).Association("Segments").Unscoped().Clear()
	db.Exec("DROP TABLE IF EXISTS segment_users")
	err := db.Migrator().DropTable(&models.User{}, &models.Segment{})
	if err != nil {
		log.Printf("Can't drop database: %s\n", err.Error())
	}
}

func (s *TestSuiteHandler) Test_Ping() {
	resp, err := http.Get(fmt.Sprintf("%s/ping", s.server.URL))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *TestSuiteHandler) Test_CreateOK() {
	resp, err := http.Post(fmt.Sprintf("%s/segment/TEST", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *TestSuiteHandler) Test_CreateFail() {
	resp, err := http.Post(fmt.Sprintf("%s/segment/TEST", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Post(fmt.Sprintf("%s/segment/TEST", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *TestSuiteHandler) Test_Delete() {
	resp, err := http.Post(fmt.Sprintf("%s/segment/TEST", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Post(fmt.Sprintf("%s/segment/TEST", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/segment/TEST", s.server.URL), nil)
	s.Require().NoError(err)
	resp, err = client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Post(fmt.Sprintf("%s/segment/TEST", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *TestSuiteHandler) Test_CreateSegmentToUserOK() {
	resp, err := http.Post(fmt.Sprintf("%s/segment/TEST", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	data := []byte(`
	{
		"add":["TEST"],
		"TTL": 2,
		"TTLUnit": "DAYS"
	}
	`)
	r := bytes.NewReader(data)
	resp, err = http.Post(fmt.Sprintf("%s/user/1", s.server.URL), "application/json", r)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	var UserDTO dtos.UserDTO
	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	json.Unmarshal(body, &UserDTO)
	log.Println(UserDTO)
}

func (s *TestSuiteHandler) Test_CreateSegmentToUserFail() {
	data := []byte(`
	{
		"add":["TEST"],
		"TTL": 2,
		"TTLUnit": "DAYS"
	}
	`)
	r := bytes.NewReader(data)
	resp, err := http.Post(fmt.Sprintf("%s/user/1", s.server.URL), "application/json", r)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *TestSuiteHandler) Test_CreateSegmentsUserFail() {
	data := []byte(`
	{
		"add":["TEST"],
		"TTL": 2,
		"TTLUnit": "DAYS"
	}
	`)
	r := bytes.NewReader(data)
	resp, err := http.Post(fmt.Sprintf("%s/user/1", s.server.URL), "application/json", r)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *TestSuiteHandler) Test_GetUserSegments() {
	resp, err := http.Post(fmt.Sprintf("%s/segment/TEST1", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Post(fmt.Sprintf("%s/segment/TEST2", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Post(fmt.Sprintf("%s/segment/TEST3", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Post(fmt.Sprintf("%s/segment/TEST4", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	data := []byte(`
	{
		"add":["TEST1", "TEST2", "TEST3"],
		"TTL": 2,
		"TTLUnit": "DAYS"
	}
	`)
	r := bytes.NewReader(data)
	resp, err = http.Post(fmt.Sprintf("%s/user/1", s.server.URL), "application/json", r)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(fmt.Sprintf("%s/user/1", s.server.URL))
	s.Require().NoError(err)
	var UserDTO dtos.UserDTO
	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	json.Unmarshal(body, &UserDTO)
	log.Println(UserDTO)
	expect := []string{"TEST1", "TEST2", "TEST3"}
	s.Require().Equal(len(expect), len(UserDTO.Segments))
	sort.Slice(UserDTO.Segments, func(i, j int) bool {
		return UserDTO.Segments[i] < UserDTO.Segments[j]
	})
	for i := 0; i < len(UserDTO.Segments); i++ {
		s.Require().Equal(expect[i], UserDTO.Segments[i])
	}
	data = []byte(`
	{
		"add":["TEST4"],
		"delete":["TEST1", "TEST2"],
		"TTL": 2,
		"TTLUnit": "DAYS"
	}
	`)
	r = bytes.NewReader(data)
	resp, err = http.Post(fmt.Sprintf("%s/user/1", s.server.URL), "application/json", r)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	resp, err = http.Get(fmt.Sprintf("%s/user/1", s.server.URL))
	s.Require().NoError(err)
	body, err = io.ReadAll(resp.Body)
	s.Require().NoError(err)
	json.Unmarshal(body, &UserDTO)
	log.Println(UserDTO)
	expect = []string{"TEST3", "TEST4"}
	s.Require().Equal(len(expect), len(UserDTO.Segments))
	sort.Slice(UserDTO.Segments, func(i, j int) bool {
		return UserDTO.Segments[i] < UserDTO.Segments[j]
	})
	for i := 0; i < len(UserDTO.Segments); i++ {
		s.Require().Equal(expect[i], UserDTO.Segments[i])
	}
}

func (s *TestSuiteHandler) Test_TTLUserSegments() {
	resp, err := http.Post(fmt.Sprintf("%s/segment/TEST1", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Post(fmt.Sprintf("%s/segment/TEST2", s.server.URL), "", nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	data := []byte(`
	{
		"add":["TEST1"],
		"TTL": 1,
		"TTLUnit": "SECONDS"
	}
	`)
	r := bytes.NewReader(data)
	resp, err = http.Post(fmt.Sprintf("%s/user/1", s.server.URL), "application/json", r)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	data = []byte(`
	{
		"add":["TEST2"],
		"TTL": 1,
		"TTLUnit": "DAYS"
	}
	`)
	r = bytes.NewReader(data)
	resp, err = http.Post(fmt.Sprintf("%s/user/1", s.server.URL), "application/json", r)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	time.Sleep(1 * time.Second)
	resp, err = http.Get(fmt.Sprintf("%s/user/1", s.server.URL))
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().NoError(err)
	var UserDTO dtos.UserDTO
	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	json.Unmarshal(body, &UserDTO)
	expect := []string{"TEST2"}
	for i := 0; i < len(UserDTO.Segments); i++ {
		s.Require().Equal(expect[i], UserDTO.Segments[i])
	}
}
