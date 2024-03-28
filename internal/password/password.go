package password_manager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	db *sql.DB
}

func NewPassword(db *sql.DB) *Password {
	return &Password{db: db}
}

// Salva uma chave com valor no SQLite
func (pw *Password) Add(k string, v string, p string) {
	encryptedValue, err := EncryptMessage([]byte(p), v)
	if err != nil {
		fmt.Println("crypted")
		panic(err.Error())
	}
	fmt.Println("enc..", encryptedValue)

	_, err = pw.db.Exec("INSERT INTO key_values (key, value) VALUES (?, ?);", k, encryptedValue)
	if err != nil {
		panic(errors.New("Add" + err.Error()))
	}
}

// Recupera o valor de uma chave no SQLite
func (pw *Password) Get(k string, p string) string {
	var encryptedValue string
	err := pw.db.QueryRow("SELECT value FROM key_values WHERE key=?;", k).Scan(&encryptedValue)
	if err != nil {
		panic(errors.New("Add" + err.Error()))
	}

	value, err := DecryptMessage([]byte(p), encryptedValue)
	if err != nil {
		fmt.Println("decrypted")
		panic(err.Error())
	}
	fmt.Println(k, "=", value)

	return value
}

// Registra um usuário
func (pw *Password) Register(u string, p string) {
	hash, _ := HashPassword(p)

	_, err := pw.db.Exec("INSERT INTO users (name, password) VALUES (?, ?);", u, hash)
	if err != nil {
		panic(errors.New("Register" + err.Error()))
	}

	fmt.Println("Success! You are registered.")
}

func (pw *Password) ValidatePassword(u string, p string) {
	if u != "" || p != "" {
		var userPassword string
		//Find User Password Hash
		err := pw.db.QueryRow("SELECT password FROM users WHERE name=? LIMIT 1", u).Scan(&userPassword)
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("Usuário não encontrado.")
			panic(err)
		case err != nil:
			fmt.Println("Erro ao consultar o banco de dados:", err)
			panic(err)
		}

		valid := CheckPasswordHash(p, userPassword)
		if !valid {
			fmt.Println("Usuário e/ou senha inválidos.")
			panic(errors.New("authError"))
		}
	} else {
		fmt.Println("É necessário passar as flags de usuário (-u) e senha (-p)")
		panic(errors.New("flagError"))
	}

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func EncryptMessage(password []byte, message string) (string, error) {
	byteMsg := []byte(message)
	block, err := aes.NewCipher(password)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptMessage(password []byte, message string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	block, err := aes.NewCipher(password)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
