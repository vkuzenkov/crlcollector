package tsl

type QualifiedCa struct {
	Version int `xml:"Версия"`
	Cas []Ca `xml:"УдостоверяющийЦентр"`
}

type Ca struct {
	Name string `xml:"Название"`
	Paks []Pak   `xml:"ПрограммноАппаратныеКомплексы>ПрограммноАппаратныйКомплекс"`
}

type Pak struct {
	Pseudo string `xml:"Псевдоним"`
	Keys   []Key   `xml:"КлючиУполномоченныхЛиц>Ключ"`
}

type Key struct {
	KeyId string `xml:"ИдентификаторКлюча"`
	Cdp  []string   `xml:"АдресаСписковОтзыва>Адрес"`
}
