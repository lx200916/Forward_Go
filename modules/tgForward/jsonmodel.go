package tgForward

type AppModel struct {
	Desc   string `json:"desc"`
	View   string `json:"view"`
	Prompt string `json:"prompt"`
	Meta   struct {
		Music struct {
			Action         string `json:"action"`
			AndroidPkgName string `json:"android_pkg_name"`
			AppType        int    `json:"app_type"`
			Appid          int    `json:"appid"`
			Desc           string `json:"desc"`
			JumpURL        string `json:"jumpUrl"`
			MusicURL       string `json:"musicUrl"`
			Preview        string `json:"preview"`
			SourceMsgID    string `json:"sourceMsgId"`
			SourceIcon     string `json:"source_icon"`
			SourceURL      string `json:"source_url"`
			Tag            string `json:"tag"`
			Title          string `json:"title"`
		} `json:"music"`
		News struct {
			Action         string `json:"action"`
			AndroidPkgName string `json:"android_pkg_name"`
			AppType        int    `json:"app_type"`
			Appid          int    `json:"appid"`
			Desc           string `json:"desc"`
			JumpURL        string `json:"jumpUrl"`
			Preview        string `json:"preview"`
			SourceIcon     string `json:"source_icon"`
			SourceURL      string `json:"source_url"`
			Tag            string `json:"tag"`
			Title          string `json:"title"`
		} `json:"news"`
		Detail1 struct {
			Appid         string `json:"appid"`
			Desc          string `json:"desc"`
			GamePoints    string `json:"gamePoints"`
			GamePointsURL string `json:"gamePointsUrl"`
			Host          struct {
				Nick string `json:"nick"`
				Uin  int64  `json:"uin"`
			} `json:"host"`
			Icon              string `json:"icon"`
			Preview           string `json:"preview"`
			Qqdocurl          string `json:"qqdocurl"`
			Scene             int    `json:"scene"`
			ShareTemplateData struct {
			} `json:"shareTemplateData"`
			ShareTemplateID string `json:"shareTemplateId"`
			ShowLittleTail  string `json:"showLittleTail"`
			Title           string `json:"title"`
			URL             string `json:"url"`
		} `json:"detail_1"`
	} `json:"meta"`
}
