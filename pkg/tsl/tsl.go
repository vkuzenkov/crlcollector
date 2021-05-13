package tsl

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Tsl struct {
	mu                 sync.Mutex
	url                *url.URL
	filename           string
	additionalCaConfig string
	Data               *QualifiedCa
	AdditionalData     *AdditionalCaConfig
	logger             *log.Logger
}

func NewTSL(tslUrl, filename, additionalCaConfig string, logger *log.Logger) (tsl *Tsl, err error) {
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
		tsl.logger.Printf("✅ Parsed %d qualified CA from file version %d", len(tsl.Data.Cas), tsl.Data.Version)
		err = tsl.parseAdditional()
		tsl.logger.Printf("✅ Parsed %d CA from additional config", len(tsl.AdditionalData.Ca))
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
	t.logger.Printf("✅ Download TSL file complete with status: %d", resp.StatusCode)

	out, err := os.Create(t.filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func (t *Tsl) Update(interval time.Duration) error {
	timer := time.NewTicker(interval)
	for {
		select {
		case <-timer.C:
			oldVersion := t.Data.Version
			t.logger.Printf("Starting update TSL current version: %d", oldVersion)
			err := t.Download()
			if err != nil {
				return err
			}
			oldData := t.Data
			err = t.parse()
			if err != nil {
				t.Data = oldData
				return err
			}
			if t.Data.Version == oldVersion {
				t.logger.Println("❕ File up to date")
			} else {
				t.logger.Printf("Update complete, new version: %d", t.Data.Version)
			}
		}
	}
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
	if t.AdditionalData == nil {
		return m
	}
	for _, aca := range t.AdditionalData.Ca {
		m[strings.ToLower(aca.KeyId)] = aca.Cdp
	}
	return m
}

func (t *Tsl) GetRootMap() map[string][]RootCert {
	m := map[string][]RootCert{}
	if t.Data == nil {
		return m
	}
	for _, ca := range t.Data.Cas {
		if ca.Status == Valid {
			for _, pak := range ca.Paks {
				for _, key := range pak.Keys {
					for _, rootcert := range key.RootCerts {
						validcerts := []RootCert{}
						if time.Until(*rootcert.NotAfter) > 0 {
							validcerts = append(validcerts, rootcert)
						}
						m[strings.ToLower(key.KeyId)] = validcerts
					}
				}
			}
		}
	}
	if t.AdditionalData == nil {
		return m
	}
	for _, aca := range t.AdditionalData.Ca {
		certs := []RootCert{}
		rootCert := RootCert{
			Base64Str: aca.Base64Str,
		}
		certs = append(certs, rootCert)
		m[strings.ToLower(aca.KeyId)] = certs
	}
	return m
}

func (t *Tsl) parse() error {
	t.mu.Lock()
	b, err := ioutil.ReadFile(t.filename)
	if err != nil {
		return err
	}
	t.Data = &QualifiedCa{}
	err = xml.Unmarshal(b, t.Data)
	defer t.mu.Unlock()
	return err
}

func (t *Tsl) parseAdditional() error {
	t.mu.Lock()
	b, err := ioutil.ReadFile(t.additionalCaConfig)
	if err != nil {
		return err
	}
	t.AdditionalData = &AdditionalCaConfig{}
	err = json.Unmarshal(b, t.AdditionalData)
	defer t.mu.Unlock()
	return err
}
