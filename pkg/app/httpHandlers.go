package app

type WebhookMessage struct {
	Description string `json:"description"`
	Event       string `json:"event"`
	Messagetype string `json:"messagetype"`
	Message     string `json:"message"`
	Usertgid    int64  `json:"usertgid"`
	Chattgid    int    `json:"chattgid"`
}

func (a *App) registerHttpHandlers() {

}
