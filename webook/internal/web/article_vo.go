package web

type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type Req struct {
	Id int64
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleVo struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	// 摘要
	Abstract string `json:"abstract"`
	// 内容
	Content string `json:"content"`
	Status  uint8  `json:"status"`
	Author  string `json:"author"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`

	// 点赞之类的信息
	LikeCnt    int64 `json:"like_cnt"`
	CollectCnt int64 `json:"collect_cnt"`
	ReadCnt    int64 `json:"read_cnt"`

	// 个人是否点赞的信息
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

type LikeReq struct {

	//点赞和取消点赞
	Id   int64 `json:"id"`
	Like bool  `json:"like"`
}
