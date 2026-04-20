package models

// Account настройки аккаунта
type Account struct {
	Channel        Channel // тг канал куда будет поститься контент
	TelegramChatID int64   // chat_id канала в Telegram (куда слать media group)
	Login          string  // логин от аккаунта на сайте откуда будет парситься контент
	Password       string  // пароль от аккаунта
}
