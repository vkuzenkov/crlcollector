# CRL Proxy

### Флаги
    -filename string
        TSL filename (default "tsl.xml")
    -listen string
        Address:port for API (default ":8080")
    -tsllink string
        TSL url (default "https://e-trust.gosuslugi.ru/app/scc/portal/api/v1/portal/ca/getxml")
    -update duration
        TSL file update interval (default 12h0m0s)

### Описание API

#### GET /debug
Получить список всех Authority key id и CDP

#### GET /crl/:keyid
Получить CRL для УЦ с корневым сертификатов со SKID keyid. keyid корневого сертификата УЦ можно определить в расширении [Authority Key Identifier](https://www.alvestrand.no/objectid/2.5.29.35.html) (OID 2.5.29.35) пользовательского сертификата.

