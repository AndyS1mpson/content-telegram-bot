package telegram

import "errors"

var (
	CommandStart   Command = "/start"   // Старт бота
	CommandCollect Command = "/collect" // /collect <query> — собрать пины с Pinterest по теме
	CommandView    Command = "/view"    // /view <query> — показать следующий пин по теме
	CommandPublish Command = "/publish" // /publish <query> — опубликовать отобранные пины
)

// Command команда для телеграм-бота
type Command string

const (
	callbackLike    = "like"
	callbackDislike = "dislike"
	callbackSkip    = "skip"
)

var (
	ErrAccessDenied    = errors.New("you do not have access")
	ErrIncorrectAction = errors.New("incorrect action")
	ErrQueryRequired   = errors.New("query is required, use: /collect <query>")
)
