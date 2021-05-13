package tsl

import (
	"encoding/base64"
	"time"
)

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
	KeyId     string     `xml:"ИдентификаторКлюча"`
	Cdp       []string   `xml:"АдресаСписковОтзыва>Адрес"`
	RootCerts []RootCert `xml:"Сертификаты>ДанныеСертификата"`
}

type RootCert struct {
	Thumbprint   string     `xml:"Отпечаток"`
	SerialNumber string     `xml:"СерийныйНомер"`
	NotAfter     *time.Time `xml:"ПериодДействияДо"`
	Base64Str    string     `xml:"Данные"`
}

type AdditionalCaConfig struct {
	Ca []AdditionalCa
}

type AdditionalCa struct {
	KeyId     string   `json:"keyId"`
	Cdp       []string `json:"cdp"`
	Base64Str string   `json:"base64Str"`
}

func (r *RootCert) ToDER() []byte {
	b, err := base64.StdEncoding.DecodeString(r.Base64Str)
	if err != nil {
		return nil
	}
	return b
}
