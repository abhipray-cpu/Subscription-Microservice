package test

import (
	"context"
	"log"
	"logger-service/data"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogSuite struct {
	client *mongo.Client
}

func (suite *LogSuite) SetupSuite() {
	const mongoURL = "mongodb://mongo:27017"
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	suite.client = c
	log.Println("Connected to MongoDB")
}

func (suite *LogSuite) Teardown() {
	err := suite.client.Disconnect(context.Background())
	if err != nil {
		log.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
	log.Println("Disconnected from MongoDB")
}

func (suite *LogSuite) TestInsertLog(t *testing.T) {

	testCases := []struct {
		name    string
		entry   data.LogEntry
		wantErr bool
	}{
		{
			name: "Valid Log Entry",
			entry: data.LogEntry{
				Service: "TestService",
				Message: "Test message",
			},
			wantErr: false,
		},
		{
			name: "Empty Service Name",
			entry: data.LogEntry{
				Service: "",
				Message: "Test message",
			},
			wantErr: true,
		},
		{
			name: "Empty Message",
			entry: data.LogEntry{
				Service: "TestService",
				Message: "",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log := data.LogEntry{}

			err := log.Insert(suite.client, tc.entry)

			if tc.wantErr && err == nil {
				t.Errorf("Expected error, but got nil")
			}

		})
	}
}

func (suite *LogSuite) TestGetAllLogs(t *testing.T) {}

func (suite *LogSuite) TestGetOneLogs(t *testing.T) {}

func (suite *LogSuite) TestUpdateLog(t *testing.T) {}

func TestLogSuite(t *testing.T) {
	suite := LogSuite{}
	suite.SetupSuite()
	defer suite.Teardown()

	t.Run("TestInserLog", suite.TestInsertLog)
	t.Run("TestGetAllLogs", suite.TestGetAllLogs)
	t.Run("TestGetOneLogs", suite.TestGetOneLogs)
	t.Run("TestUpdateLog", suite.TestUpdateLog)
}
