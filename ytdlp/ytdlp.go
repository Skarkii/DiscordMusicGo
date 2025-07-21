package ytdlp

import (
	"context"

	"github.com/lrstanley/go-ytdlp"
)

func GetYTDLPCommand(user_search_query string) *ytdlp.Result {
	ytdlp.MustInstall(context.TODO(), nil)

	f := ytdlp.New().
		DefaultSearch("ytsearch:").
		Print("urls").
		Format("bestaudio")

	fisk, err := f.Run(context.TODO(), user_search_query)

	if err != nil {
		panic(err)
	}

	return fisk
}
