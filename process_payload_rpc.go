package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/heroiclabs/nakama-common/runtime"
	"io"
	"os"
)

const (
	FILE_PATH_PREFIX_KEY     = "file_path_prefix"
	DEFAULT_FILE_PATH_PREFIX = "/user_files"

	SQL_INSERT_GAME_SCORE_QUERY = `INSERT INTO game_score (user_id, game_id, score) VALUES ($1, $2, $3) ON CONFLICT (user_id, game_id) DO UPDATE SET score = EXCLUDED.score`
	SQL_INSERT_CORE_QUERY       = `INSERT INTO core (user_id, nickname) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET nickname = EXCLUDED.nickname`
)

type Payload struct {
	Type    string `json:"type,omitempty"`
	Version string `json:"version,omitempty"`
	Hash    string `json:"hash,omitempty"`
}

type Response struct {
	Type    string      `json:"type"`
	Version string      `json:"version"`
	Hash    string      `json:"hash"`
	Content interface{} `json:"content"`
}

type Score struct {
	UserId int `json:"user_id"`
	GameId int `json:"game_id"`
	Score  int `json:"score"`
}

type Core struct {
	UserID   int    `json:"user_id"`
	Nickname string `json:"nickname"`
}

func ProcessPayloadRpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	_, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errUserNotFound
	}

	payloadData, err := parsePayload(payload)
	if err != nil {
		logger.Error("%v", err)
		return "", err
	}

	filePath := getFilePath(ctx, payloadData.Type, payloadData.Version)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File does not exist")
		return "", errFileNotFound
	}

	jsonData, hashBytes, err := readFile(filePath)
	if err != nil {
		logger.Error("Error while opening the file: %v", err)
		return "", err
	}

	switch payloadData.Type {
	case "core":
		err := processCore(jsonData, ctx, db)
		if err != nil {
			return "", err
		}
	case "score":
		err := processScore(jsonData, ctx, db)
		if err != nil {
			return "", err
		}
	default:
		logger.Error("Unsupported type:", payloadData.Type)
		return "", errBadInput
	}

	calculatedHash := hex.EncodeToString(hashBytes)
	sendContent := calculatedHash == payloadData.Hash
	var content interface{}
	if !sendContent {
		logger.Info("Hashes are not equal. Calculated: %s, requested: %s", calculatedHash, payloadData.Hash)
		content = nil
	} else {
		content = jsonData
	}

	response := Response{
		Type:    payloadData.Type,
		Version: payloadData.Version,
		Hash:    calculatedHash,
		Content: content,
	}

	return serializeResponse(response)
}

func parsePayload(payload string) (Payload, error) {
	var data Payload
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return Payload{}, err
	}

	if data.Type == "" {
		data.Type = "core"
	}
	if data.Version == "" {
		data.Version = "1.0.0"
	}
	if data.Hash == "" {
		data.Hash = ""
	}

	return data, nil
}

func readFile(filePath string) (string, []byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	jsonData, err := io.ReadAll(f)
	if err != nil {
		return "", nil, err
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", nil, err
	}

	sum := h.Sum(nil)

	return string(jsonData), sum, nil
}

func processScore(jsonData string, ctx context.Context, db *sql.DB) error {
	scores, err := deserializeScores(jsonData)
	if err != nil {
		return err
	}

	// Prepare the upsert statement
	stmt, err := db.PrepareContext(ctx, SQL_INSERT_GAME_SCORE_QUERY)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Iterate over the data and execute the upsert statement for each item
	for _, data := range scores {
		_, err := tx.Stmt(stmt).ExecContext(ctx, data.UserId, data.GameId, data.Score)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func processCore(jsonData string, ctx context.Context, db *sql.DB) error {
	core, err := deserializeCore(jsonData)
	if err != nil {
		return err
	}

	// Prepare the upsert statement
	stmt, err := db.PrepareContext(ctx, SQL_INSERT_CORE_QUERY)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Iterate over the data and execute the upsert statement for each item
	for _, data := range core {
		_, err := tx.Stmt(stmt).ExecContext(ctx, data.UserID, data.Nickname)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func deserializeScores(jsonData string) ([]Score, error) {
	var data map[string][]Score
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, err
	}

	return data["data"], nil
}

func deserializeCore(jsonData string) ([]Core, error) {
	var data map[string][]Core
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, err
	}
	return data["data"], nil
}

func serializeResponse(response interface{}) (string, error) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	// Return the JSON response
	return string(jsonResponse), nil
}

func getFilePath(ctx context.Context, payloadType string, version string) string {
	prefix, ok := ctx.Value(FILE_PATH_PREFIX_KEY).(string)
	if !ok {
		prefix = DEFAULT_FILE_PATH_PREFIX
	}

	return prefix + "/" + payloadType + "/" + version + ".json"
}
