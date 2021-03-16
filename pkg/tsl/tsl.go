package tsl

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Tsl struct {
	url      *url.URL
	filename string
	Data     *QualifiedCa
	logger   *log.Logger
}

func NewTSL(tslUrl, filename string, logger *log.Logger) (tsl *Tsl, err error) {
	tsl = &Tsl{
		filename: filename,
		logger:   logger,
	}
	if u, err := url.Parse(tslUrl); err == nil {
		tsl.url = u
	} else {
		return nil, err
	}
	go func() {
		err = tsl.Download()
		if err != nil {
			return
		}
		err = tsl.parse()
	}()

	if err != nil {
		return nil, err
	}
	return
}

func (t *Tsl) Download() error {
	resp, err := http.Get(t.url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	t.logger.Printf("Download TSL file complete with status: %d", resp.StatusCode)

	out, err := os.Create(t.filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func (t *Tsl) parse() error {
	b, err := ioutil.ReadFile(t.filename)
	if err != nil {
		return err
	}
	t.Data = &QualifiedCa{}
	err = xml.Unmarshal(b, t.Data)
	t.logger.Printf("Parsed %d qualified CA from file version %d", len(t.Data.Cas), t.Data.Version)
	return err
}

func (t *Tsl) GetCDPMap() map[string][]string {
	m := map[string][]string{}
	if t.Data == nil {
		return m
	}
	for _, ca := range t.Data.Cas {
		if ca.Status == Valid {
			for _, pak := range ca.Paks {
				for _, key := range pak.Keys {
					m[strings.ToLower(key.KeyId)] = key.Cdp
				}
			}
		}
	}
	return m
}
