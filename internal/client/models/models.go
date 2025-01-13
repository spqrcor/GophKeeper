package models

// ItemData обобщенный тип для записи
type ItemData struct {
	Id        string `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Type      string `json:"type,omitempty"`
	Login     string `json:"login,omitempty"`
	Password  string `json:"password,omitempty"`
	Text      string `json:"text,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	CardNum   string `json:"card_num,omitempty"`
	CardPayer string `json:"card_payer,omitempty"`
	CardValid string `json:"card_valid,omitempty"`
	CardPin   string `json:"card_pin,omitempty"`
}

// InputDataUser тип входящих данных пользователя
type InputDataUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Pin      string `json:"pin"`
}
