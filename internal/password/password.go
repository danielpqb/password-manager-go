package password_manager

import (
	"crypto/aes"
	"database/sql"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var pool *sql.DB

// Salva uma chave com valor no SQLite
func Add(k string, v string) {

}

// Recupera o valor de uma chave no SQLite
func Get(k string) {

}

// Registra um usuário
func Register(u string, p string) {
	hash, _ := HashPassword(p)

	match := CheckPasswordHash(p, hash)
	fmt.Println("Match: ", match)
}

func ValidatePassword(u string, p string) {
	if u != "" || p != "" {
		fmt.Println("É necessário passar as flags de usuário (-u) e senha (-p)")
		return
	}

	var userPassword string
	//Find User Password Hash
	err := pool.QueryRow("SELECT password FROM users WHERE name=? LIMIT 1", u).Scan(&userPassword)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("Usuário não encontrado.")
		return
	case err != nil:
		fmt.Println("Erro ao consultar o banco de dados:", err)
		return
	}

	valid := CheckPasswordHash(p, userPassword)
	if !valid {
		fmt.Println("Usuário e/ou senha inválidos.")
		return
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func EncryptMessage(key string, message string) string {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	msgByte := make([]byte, len(message))
	c.Encrypt(msgByte, []byte(message))
	return hex.EncodeToString(msgByte)
}

func DecryptMessage(key string, encryptedMessage string) string {
	txt, _ := hex.DecodeString(encryptedMessage)
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	msgByte := make([]byte, len(txt))
	c.Decrypt(msgByte, []byte(txt))

	msg := string(msgByte[:])
	return msg
}
