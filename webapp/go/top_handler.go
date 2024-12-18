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

var tagsMap = map[int]string{
	43:  "DIY",
	66:  "DIY電子工作",
	31:  "DJセット",
	10:  "FPS",
	19:  "Q&Aセッション",
	9:   "RPG",
	94:  "UFO",
	38:  "お料理配信",
	56:  "アウトドア",
	11:  "アクションゲーム",
	25:  "アコースティック",
	4:   "アドバイス",
	45:  "アニメトーク",
	41:  "アート配信",
	17:  "イベント生放送",
	93:  "オカルト",
	24:  "オリジナル楽曲",
	23:  "カバーソング",
	64:  "ガジェット紹介",
	62:  "ガーデニング",
	89:  "キャリア",
	57:  "キャンプ",
	28:  "ギター",
	2:   "ゲーム実況",
	15:  "ゲーム解説",
	75:  "コメディ",
	98:  "コラボ配信",
	96:  "コンサート",
	77:  "サッカー",
	102: "サプライズ",
	14:  "シングルプレイ",
	90:  "スピリチュアル",
	76:  "スポーツ",
	54:  "ダンス",
	20:  "チャット交流",
	63:  "テクノロジー",
	32:  "トーク配信",
	67:  "ニュース解説",
	79:  "バスケットボール",
	30:  "バンドセッション",
	83:  "ビジネス",
	50:  "ビューティー",
	29:  "ピアノ",
	48:  "ファッション",
	97:  "ファンミーティング",
	65:  "プログラミング",
	6:   "プロゲーマー",
	58:  "ペットと一緒",
	16:  "ホラーゲーム",
	74:  "マジック",
	13:  "マルチプレイ",
	49:  "メイク",
	53:  "ヨガ",
	80:  "ライフハック",
	1:   "ライブ配信",
	40:  "レシピ紹介",
	8:   "レトロゲーム",
	52:  "ワークアウト",
	88:  "不動産",
	86:  "仮想通貨",
	51:  "健康",
	5:   "初心者歓迎",
	91:  "占い",
	101: "周年記念",
	34:  "夜ふかし",
	82:  "子育て",
	72:  "宇宙",
	12:  "対戦ゲーム",
	71:  "心理学",
	39:  "手料理",
	92:  "手相",
	44:  "手芸",
	85:  "投資",
	81:  "教育",
	69:  "文化",
	7:   "新作ゲーム",
	18:  "新情報発表",
	55:  "旅行記",
	35:  "日常話",
	46:  "映画レビュー",
	33:  "朝活",
	87:  "株式投資",
	103: "椅子",
	27:  "楽器演奏",
	26:  "歌配信",
	68:  "歴史",
	60:  "犬",
	59:  "猫",
	3:   "生放送",
	100: "生誕祭",
	70:  "社会問題",
	73:  "科学",
	42:  "絵描き",
	21:  "視聴者参加",
	99:  "記念配信",
	37:  "語学学習",
	47:  "読書感想",
	84:  "起業",
	36:  "趣味の話",
	95:  "都市伝説",
	78:  "野球",
	61:  "釣り",
	22:  "音楽ライブ",
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

	tags := make([]*Tag, 0, len(tagsMap))
	for id, name := range tagsMap {
		tags = append(tags, &Tag{
			ID:   int64(id),
			Name: name,
		})
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

	if err := tx.Rollback(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to commit: "+err.Error())
	}

	theme := Theme{
		ID:       themeModel.ID,
		DarkMode: themeModel.DarkMode,
	}

	return c.JSON(http.StatusOK, theme)
}
