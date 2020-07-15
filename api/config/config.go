package config

var (
	DbStoreFilePath      = "./api/data_access/db_store/gateway.db"
	ExpirationDateLayout = "01-2006"
	UUIDCodeLayout       = "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"
	CvvFormatLayout      = "^[0-9]{3,4}$"
	CurrencyCodeLayout   = "^[A-Z]{3}$"
)
