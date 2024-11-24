package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TagModel struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type TagsResponse struct {
	Tags []*Tag `json:"tags"`
}

func getTagHandler(c echo.Context) error {
	// ctx := c.Request().Context()

	// tx, err := dbConn.BeginTxx(ctx, nil)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "failed to begin new transaction: : "+err.Error()+err.Error())
	// }
	// defer tx.Rollback()

	// var tagModels []*TagModel
	// if err := tx.SelectContext(ctx, &tagModels, "SELECT * FROM tags"); err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "failed to get tags: "+err.Error())
	// }

	// if err := tx.Commit(); err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "failed to commit: "+err.Error())
	// }

	// tags := make([]*Tag, len(tagModels))
	// for i := range tagModels {
	// 	tags[i] = &Tag{
	// 		ID:   tagModels[i].ID,
	// 		Name: tagModels[i].Name,
	// 	}
	// }
	tags := []*Tag{
		{ID: 43, Name: "DIY"},
		{ID: 66, Name: "DIY電子工作"},
		{ID: 31, Name: "DJセット"},
		{ID: 10, Name: "FPS"},
		{ID: 19, Name: "Q&Aセッション"},
		{ID: 9, Name: "RPG"},
		{ID: 94, Name: "UFO"},
		{ID: 38, Name: "お料理配信"},
		{ID: 56, Name: "アウトドア"},
		{ID: 11, Name: "アクションゲーム"},
		{ID: 25, Name: "アコースティック"},
		{ID: 4, Name: "アドバイス"},
		{ID: 45, Name: "アニメトーク"},
		{ID: 41, Name: "アート配信"},
		{ID: 17, Name: "イベント生放送"},
		{ID: 93, Name: "オカルト"},
		{ID: 24, Name: "オリジナル楽曲"},
		{ID: 23, Name: "カバーソング"},
		{ID: 64, Name: "ガジェット紹介"},
		{ID: 62, Name: "ガーデニング"},
		{ID: 89, Name: "キャリア"},
		{ID: 57, Name: "キャンプ"},
		{ID: 28, Name: "ギター"},
		{ID: 2, Name: "ゲーム実況"},
		{ID: 15, Name: "ゲーム解説"},
		{ID: 75, Name: "コメディ"},
		{ID: 98, Name: "コラボ配信"},
		{ID: 96, Name: "コンサート"},
		{ID: 77, Name: "サッカー"},
		{ID: 102, Name: "サプライズ"},
		{ID: 14, Name: "シングルプレイ"},
		{ID: 90, Name: "スピリチュアル"},
		{ID: 76, Name: "スポーツ"},
		{ID: 54, Name: "ダンス"},
		{ID: 20, Name: "チャット交流"},
		{ID: 63, Name: "テクノロジー"},
		{ID: 32, Name: "トーク配信"},
		{ID: 67, Name: "ニュース解説"},
		{ID: 79, Name: "バスケットボール"},
		{ID: 30, Name: "バンドセッション"},
		{ID: 83, Name: "ビジネス"},
		{ID: 50, Name: "ビューティー"},
		{ID: 29, Name: "ピアノ"},
		{ID: 48, Name: "ファッション"},
		{ID: 97, Name: "ファンミーティング"},
		{ID: 65, Name: "プログラミング"},
		{ID: 6, Name: "プロゲーマー"},
		{ID: 58, Name: "ペットと一緒"},
		{ID: 16, Name: "ホラーゲーム"},
		{ID: 74, Name: "マジック"},
		{ID: 13, Name: "マルチプレイ"},
		{ID: 49, Name: "メイク"},
		{ID: 53, Name: "ヨガ"},
		{ID: 80, Name: "ライフハック"},
		{ID: 1, Name: "ライブ配信"},
		{ID: 40, Name: "レシピ紹介"},
		{ID: 8, Name: "レトロゲーム"},
		{ID: 52, Name: "ワークアウト"},
		{ID: 88, Name: "不動産"},
		{ID: 86, Name: "仮想通貨"},
		{ID: 51, Name: "健康"},
		{ID: 5, Name: "初心者歓迎"},
		{ID: 91, Name: "占い"},
		{ID: 101, Name: "周年記念"},
		{ID: 34, Name: "夜ふかし"},
		{ID: 82, Name: "子育て"},
		{ID: 72, Name: "宇宙"},
		{ID: 12, Name: "対戦ゲーム"},
		{ID: 71, Name: "心理学"},
		{ID: 39, Name: "手料理"},
		{ID: 92, Name: "手相"},
		{ID: 44, Name: "手芸"},
		{ID: 85, Name: "投資"},
		{ID: 81, Name: "教育"},
		{ID: 69, Name: "文化"},
		{ID: 7, Name: "新作ゲーム"},
		{ID: 18, Name: "新情報発表"},
		{ID: 55, Name: "旅行記"},
		{ID: 35, Name: "日常話"},
		{ID: 46, Name: "映画レビュー"},
		{ID: 33, Name: "朝活"},
		{ID: 87, Name: "株式投資"},
		{ID: 103, Name: "椅子"},
		{ID: 27, Name: "楽器演奏"},
		{ID: 26, Name: "歌配信"},
		{ID: 68, Name: "歴史"},
		{ID: 60, Name: "犬"},
		{ID: 59, Name: "猫"},
		{ID: 3, Name: "生放送"},
		{ID: 100, Name: "生誕祭"},
		{ID: 70, Name: "社会問題"},
		{ID: 73, Name: "科学"},
		{ID: 42, Name: "絵描き"},
		{ID: 21, Name: "視聴者参加"},
		{ID: 99, Name: "記念配信"},
		{ID: 37, Name: "語学学習"},
		{ID: 47, Name: "読書感想"},
		{ID: 84, Name: "起業"},
		{ID: 36, Name: "趣味の話"},
		{ID: 95, Name: "都市伝説"},
		{ID: 78, Name: "野球"},
		{ID: 61, Name: "釣り"},
		{ID: 22, Name: "音楽ライブ"},
	}

	return c.JSON(http.StatusOK, &TagsResponse{
		Tags: tags,
	})
}

// 配信者のテーマ取得API
// GET /api/user/:username/theme
func getStreamerThemeHandler(c echo.Context) error {
	ctx := c.Request().Context()

	if err := verifyUserSession(c); err != nil {
		// echo.NewHTTPErrorが返っているのでそのまま出力
		c.Logger().Printf("verifyUserSession: %+v\n", err)
		return err
	}

	username := c.Param("username")

	tx, err := dbConn.BeginTxx(ctx, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to begin transaction: "+err.Error())
	}
	defer tx.Rollback()

	userModel := UserModel{}
	err = tx.GetContext(ctx, &userModel, "SELECT id FROM users WHERE name = ?", username)
	if errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "not found user that has the given username")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user: "+err.Error())
	}

	themeModel := ThemeModel{}
	if err := tx.GetContext(ctx, &themeModel, "SELECT * FROM themes WHERE user_id = ?", userModel.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user theme: "+err.Error())
	}

	if err := tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to commit: "+err.Error())
	}

	theme := Theme{
		ID:       themeModel.ID,
		DarkMode: themeModel.DarkMode,
	}

	return c.JSON(http.StatusOK, theme)
}
