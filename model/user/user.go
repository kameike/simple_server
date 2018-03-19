package user

type User struct {
	Id      int64 `db:"post_id"`
	Created int64
	Name    string `db:",size:100"`
}
