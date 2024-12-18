package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type ReactionModel struct {
	ID           int64  `db:"id"`
	EmojiName    string `db:"emoji_name"`
	UserID       int64  `db:"user_id"`
	LivestreamID int64  `db:"livestream_id"`
	CreatedAt    int64  `db:"created_at"`
}

type Reaction struct {
	ID         int64      `json:"id"`
	EmojiName  string     `json:"emoji_name"`
	User       User       `json:"user"`
	Livestream Livestream `json:"livestream"`
	CreatedAt  int64      `json:"created_at"`
}

type PostReactionRequest struct {
	EmojiName string `json:"emoji_name"`
}

func getReactionsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	if err := verifyUserSession(c); err != nil {
		// echo.NewHTTPErrorが返っているのでそのまま出力
		return err
	}

	livestreamID, err := strconv.Atoi(c.Param("livestream_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "livestream_id in path must be integer")
	}

	tx, err := dbConn.BeginTxx(ctx, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to begin transaction: "+err.Error())
	}
	defer tx.Rollback()

	query := "SELECT * FROM reactions WHERE livestream_id = ? ORDER BY created_at DESC"
	if c.QueryParam("limit") != "" {
		limit, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "limit query parameter must be integer")
		}
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	reactionModels := []ReactionModel{}
	if err := tx.SelectContext(ctx, &reactionModels, query, livestreamID); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "failed to get reactions")
	}

	reactions := make([]Reaction, len(reactionModels))
	reactionMap, err := fillReactionResponse(ctx, tx, reactionModels)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fill reaction: "+err.Error())
	}
	for i := range reactionModels {
		reaction, ok := reactionMap[reactionModels[i].ID]
		if !ok {
			log.Printf("not found reaction, reactionID: %v", reactionModels[i].ID)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get reaction")
		}
		reactions[i] = reaction
	}

	// for i := range reactionModels {
	// 	reaction, err := fillReactionResponse(ctx, tx, reactionModels[i])

	// 	reactions[i] = reaction
	// }

	if err := tx.Rollback(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to commit: "+err.Error())
	}

	// return c.JSON(http.StatusOK, reactions)
	return c.JSONBlob(http.StatusOK, jsonEncode(reactions))
}

func postReactionHandler(c echo.Context) error {
	ctx := c.Request().Context()
	livestreamID, err := strconv.Atoi(c.Param("livestream_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "livestream_id in path must be integer")
	}

	if err := verifyUserSession(c); err != nil {
		// echo.NewHTTPErrorが返っているのでそのまま出力
		return err
	}

	// error already checked
	sess, _ := session.Get(defaultSessionIDKey, c)
	// existence already checked
	userID := sess.Values[defaultUserIDKey].(int64)

	var req *PostReactionRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to decode the request body as json")
	}

	tx, err := dbConn.BeginTxx(ctx, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to begin transaction: "+err.Error())
	}
	defer tx.Rollback()

	reactionModel := ReactionModel{
		UserID:       int64(userID),
		LivestreamID: int64(livestreamID),
		EmojiName:    req.EmojiName,
		CreatedAt:    time.Now().Unix(),
	}

	result, err := tx.NamedExecContext(ctx, "INSERT INTO reactions (user_id, livestream_id, emoji_name, created_at) VALUES (:user_id, :livestream_id, :emoji_name, :created_at)", reactionModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to insert reaction: "+err.Error())
	}

	reactionID, err := result.LastInsertId()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get last inserted reaction id: "+err.Error())
	}
	reactionModel.ID = reactionID

	reactionMap, err := fillReactionResponse(ctx, tx, []ReactionModel{reactionModel})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fill reaction: "+err.Error())
	}
	reaction, ok := reactionMap[reactionID]
	if !ok {
		log.Printf("not found reaction, reactionID: %v", reactionID)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get reaction")
	}

	if err := tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to commit: "+err.Error())
	}

	return c.JSON(http.StatusCreated, reaction)
}

func fillReactionResponse(ctx context.Context, tx *sqlx.Tx, reactionModels []ReactionModel) (map[int64]Reaction, error) {
	userIDs := make([]int64, len(reactionModels))
	livestreamIDs := make([]int64, len(reactionModels))
	for i := range reactionModels {
		userIDs[i] = reactionModels[i].UserID
		livestreamIDs[i] = reactionModels[i].LivestreamID
	}
	userMap, err := getUsers(ctx, tx, userIDs)
	if err != nil {
		return nil, err
	}
	// user := userMap[reactionModel.UserID]

	// var livestreamModels []LivestreamModel
	// query, params, err := sqlx.In("SELECT * FROM livestreams WHERE id IN (?)", livestreamIDs)
	// if err != nil {
	// 	return nil, err
	// }
	// if err := tx.SelectContext(ctx, &livestreamModels, query, params...); err != nil {
	// 	return nil, err
	// }

	// var livestreamMap = make(map[int64]Livestream, len(livestreamModels))
	livestreamMap, err := getLivestreams(ctx, tx, livestreamIDs)
	if err != nil {
		return nil, err
	}
	// for _, livestreamModel := range livestreamModels {
	// 	livestream, err := fillLivestreamResponse(ctx, tx, livestreamModel)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	livestreamMap[livestreamModel.ID] = livestream
	// }
	// livestreamModel := LivestreamModel{}
	// if err := tx.GetContext(ctx, &livestreamModel, "SELECT * FROM livestreams WHERE id = ?", reactionModel.LivestreamID); err != nil {
	// 	return Reaction{}, err
	// }
	var reactionMap = make(map[int64]Reaction, len(reactionModels))
	for _, reactionModel := range reactionModels {
		user, ok := userMap[reactionModel.UserID]
		if !ok {
			log.Printf("not found user, userID: %v", reactionModel.UserID)
			continue
		}
		livestream, ok := livestreamMap[reactionModel.LivestreamID]
		if !ok {
			log.Printf("not found livestream, livestreamID: %v", reactionModel.LivestreamID)
			continue
		}

		reaction := Reaction{
			ID:         reactionModel.ID,
			EmojiName:  reactionModel.EmojiName,
			User:       *user,
			Livestream: *livestream,
			CreatedAt:  reactionModel.CreatedAt,
		}
		reactionMap[reactionModel.ID] = reaction
	}

	// reaction := Reaction{
	// 	ID:         reactionModel.ID,
	// 	EmojiName:  reactionModel.EmojiName,
	// 	User:       *user,
	// 	Livestream: livestream,
	// 	CreatedAt:  reactionModel.CreatedAt,
	// }

	return reactionMap, nil
}
