package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/heroiclabs/nakama-common/runtime"
	"log"
	"os"
	"testing"
)

const (
	TEST_TMP_FOLDER = "nakama_tmp"
)

type MockLogger struct {
	runtime.Logger
}

func (m *MockLogger) Error(format string, v ...interface{}) {
}

func (m *MockLogger) Info(format string, v ...interface{}) {
}

type MockNakamaModule struct {
	runtime.NakamaModule
}

func TestProcessPayloadRpc(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		payloadType string
		want        string
		wantErr     bool
	}{
		{
			name:        "Valid `score` payload with the correct hash. Returns a response with content",
			payload:     `{"type": "score","version": "1.0.0","hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}`,
			payloadType: "score",
			want:        `{"type":"score","version":"1.0.0","hash":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","content":"{\n  \"data\": [\n    {\n      \"user_id\": 1,\n      \"game_id\": 2,\n      \"score\": 999\n    },\n    {\n      \"user_id\": 2,\n      \"game_id\": 2,\n      \"score\": 111\n    }\n  ]\n}\n"}`,
		},
		{
			name:        "Valid `score` payload with incorrect hash. Returns a response with empty content and a correct hash",
			payload:     `{"type": "score","version": "1.0.0","hash": "incorrect_hash"}`,
			payloadType: "score",
			want:        `{"type":"score","version":"1.0.0","hash":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","content":null}`,
		},
		{
			name:        "Valid `core` payload with incorrect hash. Returns a response with empty content and a correct hash",
			payload:     `{"type": "core","version": "1.0.0","hash": "incorrect_hash"}`,
			payloadType: "core",
			want:        `{"type":"core","version":"1.0.0","hash":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","content":null}`,
		},
		{
			name:        "Valid empty payload. Uses defaults and returns `core` type",
			payload:     `{}`,
			payloadType: "core",
			want:        `{"type":"core","version":"1.0.0","hash":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","content":null}`,
		},
		{
			name:        "Invalid payload with a non-existent file. Returns error",
			payload:     `{"type": "score","version": "2.0.0","hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}`,
			payloadType: "",
			wantErr:     true,
		},
	}

	if _, err := os.Stat(TEST_TMP_FOLDER); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(TEST_TMP_FOLDER, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	defer os.Remove(TEST_TMP_FOLDER)
	logger := new(MockLogger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createContext()
			db := mockDb(t, tt.payloadType)
			defer db.Close()

			got, err := ProcessPayloadRpc(ctx, logger, db, &MockNakamaModule{}, tt.payload)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessPayloadRpc() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("ProcessPayloadRpc() error = %v, wantErr %v", err.Error(), "File does not exist")
			}
			if got != tt.want {
				t.Errorf("ProcessPayloadRpc() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockDb(t *testing.T, payloadType string) *sql.DB {
	db, mock, dbErr := sqlmock.New()
	if dbErr != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", dbErr)
	}

	if payloadType == "score" {
		mock.ExpectPrepare("INSERT INTO game_score")
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO game_score").WithArgs(1, 2, 999).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO game_score").WithArgs(2, 2, 111).WillReturnResult(sqlmock.NewResult(2, 1))
		mock.ExpectCommit()
	} else {
		mock.ExpectPrepare("INSERT INTO core")
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO core").WithArgs(1, "roqueta").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	return db
}

func createContext() context.Context {
	ctx := context.WithValue(
		context.Background(),
		runtime.RUNTIME_CTX_USER_ID, "1",
	)

	return context.WithValue(ctx, FILE_PATH_PREFIX_KEY, "./user_files")
}
