package protocol

import (
	"crypto/x509"
	"database/sql"
	"errors"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var (
	ErrInvalidCert = errors.New("gemini: Invalid TOFU certificate")
)

func InitCertsDB() (*CertStore, error) {
	db, err := sql.Open("sqlite3", "certs.db")
	if err != nil {
		return nil, err
	}

	db.Exec(`
		CREATE TABLE IF NOT EXISTS certs (
			hostname TEXT PRIMARY KEY,
			cert BLOB
		)
	`)

	return (*CertStore)(db), nil
}

type CertStore sql.DB

func (cs *CertStore) CheckCertificate(hostname string, cert *x509.Certificate) (bool, error) {

	db := (*sql.DB)(cs)

	stmt, err := db.Prepare(`
		SELECT cert
		FROM certs
		WHERE hostname = ?
	`)
	if err != nil {
		return false, err
	}

	res := stmt.QueryRow(hostname)
	var stored_raw []byte
	err = res.Scan(&stored_raw)
	if err == sql.ErrNoRows {
		stmt, err = db.Prepare(`
			INSERT INTO certs (hostname, cert)
			VALUES (?, ?)
		`)
		if err != nil {
			return false, err
		}
		_, err := stmt.Exec(hostname, cert.Raw)
		if err != nil {
			return false, err
		}
		return true, nil
	} else if err != nil {
		return false, err
	}

	stored, err := x509.ParseCertificate(stored_raw)
	if err != nil {
		return false, err
	}

	return stored.Equal(cert), nil
}
