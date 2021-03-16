package tsl

type QualifiedCa struct {
	Version int  `xml:"Версия"`
	Cas     []Ca `xml:"УдостоверяющийЦентр"`
}

// Тип для обработки статусов аккредитации УЦ
type Status string

const (
	Canceled   Status = "Аннулирована"
	Terminated Status = "Прекращена"
	Valid      Status = "Действует"
)

type Ca struct {
	Name   string `xml:"Название"`
	Status Status `xml:"СтатусАккредитации>Статус"`
	Paks   []Pak  `xml:"ПрограммноАппаратныеКомплексы>ПрограммноАппаратныйКомплекс"`
}

type Pak struct {
	Pseudo string `xml:"Псевдоним"`
	Keys   []Key  `xml:"КлючиУполномоченныхЛиц>Ключ"`
}

type Key struct {
	KeyId string   `xml:"ИдентификаторКлюча"`
	Cdp   []string `xml:"АдресаСписковОтзыва>Адрес"`
}
