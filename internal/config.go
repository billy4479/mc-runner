package internal

type Config struct {
	GhHookSecret string `required:"true" split_words:"true"`
	TgToken      string `required:"true" split_words:"true"`
	TgHookSecret string `required:"true" split_words:"true"`
	McSecret     string `required:"true" split_words:"true"`
	WorldDir     string `required:"true" split_words:"true"`
	DbPath       string `required:"true" split_words:"true"`

	// GhHookSecret string `split_words:"true"`
	// TgToken      string `split_words:"true"`
	// TgHookSecret string `split_words:"true"`
	// McSecret     string `split_words:"true"`
	// WorldDir     string `split_words:"true"`
	// DbPath       string `split_words:"true"`
}
