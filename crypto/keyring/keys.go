package keyring

const (
	BIP44Prefix = "44'/118'/"
	PartialPath = "0'/0/0"
	//FullPath    = BIP44Prefix + PartialPath
	FullPath = "m/" + BIP44Prefix + PartialPath
)
