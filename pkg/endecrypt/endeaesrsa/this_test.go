package endeaesrsa

import (
	"fmt"
	"odin.ren/endecrypt"
	"testing"
)

var (
	//指定密钥
	text = []byte("hello world")
	err  error
	str  string
)

// 公钥加密
func BenchmarkPubEncrypt(b *testing.B) {
	str, err = PubEncrypt(text, endecrypt.PubkeyCode, endecrypt.AesKeyServer1)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(str)
	}
}

// 私钥解密
func BenchmarkPirDecrypt(b *testing.B) {
	str, err = PubEncrypt(text, endecrypt.PubkeyCode, endecrypt.AesKeyServer1)
	str = `0g/UNJ1NKHRgAJ6O+gEbd1xNOXdv+8E9cqa1nfy/HFqzyYwstwxrwsZYe91ifOJlRytSgH3JvQ6cQ89S9IX8iVCW6MwtC5X/NZ3sST3+JQyevaStJfGP96HyynEhj+Sgqgm2zobkZ5ZBpXzxZqSvRQ04qfvJKo+BQ0GO3TMziVAYvqtAm3IQXpDlPWRed1aK6EYBs0yGiSnjK/Z4bXgHwOamfXkhehCtsRsE+VfkZmyMcZ9P/s3XJlS+w24VMqABPrtjy9HcJa7PGXo/VP9qzLKDWrRc99ch7D0luT/6a0ZR8IPW4Smq1ijGMhm93q87LLIYQYILznLh715ujIYNbJnLq8SNXPJir/paR0tJMdswWvrO40lhsN4MQDzLp7PtkEjUrMH2Dw3BE2wzStt3kqeKGRODx+Z1Aue+JD1veiffV4cPh/Dr194NlwKB5z4Lm1WAa4WM4Daf3J0DdIi3z4LU1jcLfrVjxvs7QRjCUXdKKz22GWskr7Veas6gAEfqCTnQ2exRV9OvcPKXkw36MFddBHohMF49HH4Et/ZEQk4wOtOL7SN60J1miDzOKxoRlwho+KTUyPFwyGsWqq86vhmpd80lllwpXK+SGDbMJDwTtnnyupWCY4EDC3CUHS0EloyXvyVjDr1dDyS2y6Gm25USRLgpA9g3hH0Co8WT9GTtlDcgBqmUKQYGKfbjDGNcPK78ySl/KEtZpy/M+lATSSQesyms7Y1mVYcC8Es8JnY8G/p3dzd9rZog8yYRNbgZGkrPwqOY+A+ZIdTbe6Y9Ij1i9YzWxlHqucod6TZ6GSG5buwtf35aqcm9hMDpc5R2x+OSHoi6zMh6QNqRWreRyOQM0nSm4GDDbia/wChdzgRvc7S73lPE8eJRX5+oRYMtPcobutJ+dG4Xk4+cuiBGx9PyJLsvP4Pzxj0F4cfXCnFwKY2dWoVwEhP3Wf12ztvAjhxzOeLgg7e6ZkLs091v6cieRHm+XXYWq9D82LvjwKbT3wjZzLGke7OwmcWtRL+VDeW4MAYIR4A3QtYGOekDqw==`
	str, err = PirDecrypt(str, endecrypt.PirvatekeyServer, endecrypt.AesKeyServer1)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(str)
	}
}

// 私钥加密
func BenchmarkPriEncrypt(b *testing.B) {
	str, err = PriEncrypt(text, endecrypt.PirvatekeyServer, endecrypt.AesKeyServer1)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(str)
	}
}

// 公钥解密
func BenchmarkPubDecrypt(b *testing.B) {
	str, err = PriEncrypt(text, endecrypt.PirvatekeyServer, endecrypt.AesKeyServer1)
	str, err = PubDecrypt(str, endecrypt.PubkeyCode, endecrypt.AesKeyServer1)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(str)
	}
}

// 私钥加密2
func BenchmarkPriEncrypt2(b *testing.B) {
	// ZVjJ+eA/qXBpmSQuMEl2d6rLpCt8X9JvMo8G9fNK3nBX4pXFM6vYUI69yfk6fsqj2VT6Mg1GJhY1Bch+aFGG/4u+/iKo+re69LRJyYSZyaOap0mDlmE3B6upSeBt63Zi/R+OnHBd4uCncT+wVZ//6Bu0WHnjm7iCgQzdNoXrk145/9HEwdlXM/8xEZevS5N/Zxht2P0OjOE86rljintl6cKCNL+KGt4o3LpKKBZzUte8GxHBPl68i2uK6hVzeSVNYk/PmOINt+LgMk+nnJKgIk+hAjXwBwzOl9yCf7ZWgedY3rCjGEoyJ7UzlC2RUFPAU327V9iI/91+CVaYOq+HNs2MzosE6wMECMY1SCsUj4PWsWp+xnB7iHZL2kJGFT8aMDMQDFGuGSpJ8uuo56vqO6lUEaqVfs9+juEvF//2GsTPewF5nhrEVcz42sKTTA8cYoNVhziIuSW30ea9ORk2PLdcdCe2c1Rm+4ZT5ErACmc9CpyYG59GDa7Kf/Ac8MSt+cttZScIEkLUzovrXIMb0kVAhzKRBh19qscCbKDaU2OncvxMhTS4aXipJ/sx6gtHZX5VbCTJUK2Z2bYvWoGyA8WPhtTDiu7niYDVhE5lc1p11lMscC+I9U/SseDsz/trDCj7gpgxwK01/sLbJ2cCKmUHInlApAWtgS0FQ5QrqpnHm5+udFKqEyLGG33C12b6
	text2 := []byte("4444")
	str, err = PriEncrypt(text2, endecrypt.PirvatekeyClient1024, endecrypt.AesKeyAuth)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(str)
	}
}

// 公钥解密2
func BenchmarkPubDecrypt2(b *testing.B) {
	// ladB3Hp1X0f6WMTrrKmFAST48W4ETNWPlw91AnsjQHROPEYUQljTT2R0NmUmAqfZuaN91DaRJek+04bYrdL56mD2nytYi15QA3HUZ6F0LsSAPPppR1dYc4WO4PwYgFmi7a3O51U3vlcthQcEZ720EHgrdi+OyPrrvG9mtcphkRraIngXjXJBmW2VxZq4TDjI9pzgG5PvA27ihBby/TmyGgTHuqmKcPuGM192a0wPAj1F71EWEjSDWGUtgUiBLzRpNX3u9LBmeoqW07AV2y+qo5ThyKJAaXotUbNudLvjK8hQ4fsq/XYF3wyYQt+QiUfmjKGrLtexiOESRIsVq+242xYvo5WEoUoFCcwf/3XrZPQ=
	text2 := []byte("4444")
	str, err = PriEncrypt(text2, endecrypt.PirvatekeyClient1024, endecrypt.AesKeyAuth)
	str, err = PubDecrypt(str, endecrypt.PubkeyClient1024, endecrypt.AesKeyAuth)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(str)
	}
}
