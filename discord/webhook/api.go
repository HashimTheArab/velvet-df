package webhook

func Send(url string, message Message) {
	webhook <- request{url, message}
}
