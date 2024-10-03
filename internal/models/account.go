package models

// Account настройки аккаунта
type Account struct {
	Channel  Channel // тг канал куда будет поститься контент
	Login    string  // логин от аккаунта на сайте откуда будет парситься контент
	Password string  // пароль от аккаунта
}
