package repository

// PinterestPin описывает пин из pinterest
type PinterestPin struct {
	ID        int64  `db:"id"`
	ImageURL  string `db:"image_url"`
	Status    int64  `db:"status"`
	TgChannel string `db:"tg_channel"`
}
