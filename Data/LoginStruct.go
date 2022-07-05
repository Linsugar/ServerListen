package Data

type MagiciDemo struct {
	AlarmType string `json:"alarmType"`
	Boxes     []struct {
		Height int `json:"height"`
		Width  int `json:"width"`
		X      int `json:"x"`
		Y      int `json:"y"`
	} `json:"boxes"`
	CameraId string `json:"cameraId"`
	Extra    struct {
		Cei        interface{} `json:"cei"`
		ItemsInBox []struct {
			Confidence float64 `json:"confidence"`
			Type       string  `json:"type"`
		} `json:"itemsInBox"`
	} `json:"extra"`
	Scene string `json:"scene"`
	Ts    int64  `json:"ts"`
	Url   string `json:"url"`
}
