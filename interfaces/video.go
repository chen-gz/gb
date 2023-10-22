package interfaces

type VideoItem struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id"`
	Title   string `json:"title"`
	Tag     string `json:"tag"`
	Md5     string `json:"md5"`
	Sha256  string `json:"sha256"`
	Deleted bool   `json:"deleted"`
	Ext     string `json:"ext"` // Extension of video file. e.g. mp4, mov, avi, etc.
}
