package webhook

type Message struct {
	Username  string  `json:"username,omitempty"`
	AvatarUrl string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string    `json:"title,omitempty"`
	Url         string    `json:"url,omitempty"`
	Description string    `json:"description,omitempty"`
	Color       int       `json:"color,omitempty"`
	Author      Author    `json:"author,omitempty"`
	Fields      []Field   `json:"fields,omitempty"`
	Thumbnail   Thumbnail `json:"thumbnail,omitempty"`
	Image       Image     `json:"image,omitempty"`
	Footer      Footer    `json:"footer,omitempty"`
}

// Author is like a smaller version of a title and supports an icon. It is seen on embeds.
type Author struct {
	Name    string `json:"name,omitempty"`
	Url     string `json:"url,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}

type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// Thumbnail is a small image that appears in an embed.
type Thumbnail struct {
	Url string `json:"url,omitempty"`
}

// Image is a large image that appears in an embed.
type Image struct {
	Url string `json:"url,omitempty"`
}

// Footer is a small message that appears on the bottom of an embed.
type Footer struct {
	Text    string `json:"text,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}
