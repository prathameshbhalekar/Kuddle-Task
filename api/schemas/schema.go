package schemas

type Class struct {
	ClassUuid string `json:"class_uuid"`
	Type      string `json:"type"`
	Capacity  int    `json:"capacity"`
	Members   int    `json:"members"`
}

type Registration struct {
	UserUuid    *string `json:"user_uuid"`
	ClassUuid   *string `json:"class_uuid"`
	IsWaiting   bool    `json:"is_waiting"`
	IsCancelled bool    `json:"is_cancelled"`
	BookedAt    int     `json:"booked_at"`
}
