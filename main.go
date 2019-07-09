package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go-graphql-cloud-api/ciphers"
	"go-graphql-cloud-api/gql"

	"go-graphql-cloud-api/postgres"
	"go-graphql-cloud-api/server"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	config()

	// Initialize our api and return a pointer to our router for http.ListenAndServe
	// and a pointer to our db to defer its closing when main() is finished
	router, db := initializeAPI()
	//initPem()
	defer db.Close()

	// Listen on port and if there's an error log it and exit
	fmt.Println("Listening on port:", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func config() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		fmt.Println("Loaded config file")
	}
}

func initializeAPI() (*chi.Mux, *postgres.Db) {
	// Create a new router
	router := chi.NewRouter()

	// Create a new connection to our pg database

	// converting the str1 variable into an int using Atoi method
	dbPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Fatal("Error loading POSTGRES_PORT")
	}
	db, err := postgres.New(
		postgres.ConnString(os.Getenv("POSTGRES_HOST"), dbPort, os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_NAME"), os.Getenv("POSTGRES_SSLMODE")),
	)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("database has been set up")
	}

	// Create our root query for graphql
	rootQuery := gql.NewRoot(db)
	// Create a new graphql schema, passing in the the root query
	sc, err := graphql.NewSchema(
		graphql.SchemaConfig{Query: rootQuery.Query},
	)
	if err != nil {
		fmt.Println("Error creating schema: ", err)
	} else {
		fmt.Println("GraphQL Schema has been set up")
	}

	// Create a server struct that holds a pointer to our database as well
	// as the address of our graphql schema
	s := server.Server{
		GqlSchema: &sc,
		Context:   rootQuery.Context,
	}

	// Add some middleware to our router
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		middleware.Logger,          // log api request calls
		middleware.DefaultCompress, // compress results, mostly gzipping assets and json
		middleware.StripSlashes,    // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer,       // recover from panics without crashing server
	)

	// Create the graphql route with a Server method to handle it
	router.Post(os.Getenv("GRAPHQL_LINK"), s.GraphQL())

	return router, db
}

func initPem() {
	// err := ciphers.GenerateKeyPair(1024, "private.pem", "public.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err2 := ciphers.GenerateKeyPair(8192, "srprivate.pem", "srpublic.pem")
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }

	fmt.Println("------------------------------------")
	// Load Keys
	privateKey := ciphers.LoadRSAPrivatePemKey("srprivate.pem")
	publicKey := ciphers.LoadRSAPublicPemKey("srpublic.pem")

	fmt.Println(string(ciphers.PublicKeyToBytes(publicKey)))

	// msg := string(ciphers.PublicKeyToBytes(ciphers.LoadRSAPublicPemKey("public.pem")))
	// msg := string(ciphers.PublicKeyToBytes(&ciphers.LoadRSAPrivatePemKey("testprivate.pem").PublicKey))
	// msg := string(ciphers.PrivateKeyToBytes(ciphers.LoadRSAPrivatePemKey("private.pem")))
	//msg := "{users(phoneDeviceID:\"1234\"){id, staffID, stripeID, phoneDeviceID, name, remark}}"
	msg := "S0000232"
	key := "LKHlhb899Y09olUiLKHlhb899Y09olUi"
	fmt.Println("MSG:")
	fmt.Println(msg)
	fmt.Println("KEY:")
	fmt.Println(key)

	fmt.Println("------------------------------------")
	fmt.Println("ENCRYPTING USING AES.......")
	aesEncrypted, _ := ciphers.EncryptWithAes(msg, key)

	fmt.Println(aesEncrypted)

	fmt.Println("------------------------------------")
	fmt.Println("ENCRYPTING USING PUBLIC KEYS.......")
	pubEncryptedMsg, _ := ciphers.EncryptWithPublicKey(key, publicKey)

	fmt.Println(pubEncryptedMsg)

	fmt.Println("------------------------------------")
	fmt.Println("DECRYPTING USING PRIVATE KEYS.......")
	priDecryptedMsg, _ := ciphers.DecryptWithPrivateKey(pubEncryptedMsg, privateKey)

	fmt.Println(priDecryptedMsg)

	fmt.Println("------------------------------------")
	fmt.Println("DECRYPTING USING AES.......")
	aesDecrypted, _ := ciphers.DecryptWithAes(aesEncrypted, priDecryptedMsg)

	fmt.Println(aesDecrypted)

	// REVERSE
	// 	fmt.Println("------------------------------------")

	// 	fmt.Println("ENCRYPTING USING PRIVATE KEYS.......")
	// 	priEncryptedMsg, err := ciphers.EncryptWithPrivateKey(msg, privateKey)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Println(priEncryptedMsg)

	// 	fmt.Println("------------------------------------")
	// 	fmt.Println("DECRYPTING USING PUBLIC KEYS.......")
	// 	pubDecryptedMsg := ciphers.DecryptWithPublicKey(priEncryptedMsg, publicKey)

	// 	fmt.Println(pubDecryptedMsg)

	// 	fmt.Println("------------------------------------")
	// 	fmt.Println("DECRYPTING USING PRVATE KEYS AS STRING.......")
	// 	var privateKeyString = `
	// -----BEGIN RSA PRIVATE KEY-----
	// MIISKQIBAAKCBAEApfvnykm6Nfm1EOODT4DN1UQ0w4g1C3KivTz4rWfqnGoSh33z
	// Vmb3u5GmEGdhH5EDDssY6qo1Pq9QDV8RyhSF/M9xeFh6Y6Mx7hgsW56PnNOeE6f8
	// wLpQ/1N56MIB7OqckI55eqyLTuwt2Pj+UbLb3LYQPfK0jFgbHXhX3IgtHsEzpi+6
	// +YjAg2x4QAO1mqbQoSaUZDe6Xsl1jY/uvgfo29uhwFyN1cHVQOEZtbyDwOfVHmju
	// OXy/QS6Fht9SxAWis24D9esi/+LnLGOcYTWqj40jsU3OK5c7hY0mIY4HEiZqGqFn
	// 6bUEId3J3gy9kbLJ29gB27HGzGIl2LbdXde7q1cfhJXNMQXiEXDttU2emK7yzMkH
	// MYsh4G4wSU6bPQLB+MB1MLYEBjz9EHcIHxrQJA4u3XGpQbf+WQph6/8YS4eLrQPe
	// uYl48vQ6GzPhOstojmMFr/RqY3Xs+qhpCHJebz5vnSib3KjpuLExTsQz13SvFo7U
	// bGHSXnmzQGa7WXA+LM+I3H9CiPYRmDWZkcx0xhBl+x364fqp/YTcOMOSb8u5rxJ1
	// MwWLkz/yIxDAmShk+RGeTn4ozGR3sJnLuVIr2IV+S38WfIKa00I/HNiPjmqT+M2U
	// xcZsb87YI/WDae2ic9YEtHEIKPHkElVeJ88YKvjQkWNRgZLu5+v+XWe5FFUqNidi
	// o8wUczPmWl5Vz6eKpsLODcCLpvhyL485Ws2dpF0iyTTa74uahsy/HgiRUb6oX7y5
	// mvzXGmmiWj/81oX8qCP8/66vIQqIMta0dWdBpsReYdt2SIBOOgpnzGjoALJ5h8sT
	// leogr123rWhXptBtvL7dmljt/Ml2Jmhewrf3ieKrHpvD5K2/7sjhlQ8/ufVhd1gy
	// foVy+1vZZvQwW8GbQtm7+KVZkOgAhKOzZzwLUzW5bXXIgwSR1skWViaTe8WILw5g
	// qOKyKIi8C78SCOKNZIBXCcxo1RlJI0JHisV19xvJdGHRoOlbD5Drj8//Khpr6QRg
	// aRnnR/XIq2XQa+rP5PbMFEhO2FIwmeJ8F9sFaT0BBNvonFV9S0f1ZJK6hbUYVOkC
	// 47AY7i5g3FYfc/mGjbWA9oBPfo91PKzzCfOFbwMKPndXe4RUhKg5hR7lGQ3U1Fyl
	// odEMQSGDha4ukWP1B6bnXkoyNUEl7SnFV7Mz09HWCBIBNDPevcUp/1D1apMTAEUR
	// LcGSxDVQbgUgIZAnYlhPqvopyHCMqLjwsLctu2E5gbT7XHmc2eFOcu/RoMXr/p/z
	// XG7NxiR54I31202CGKZ8BUVGuUiMcxbMdM/7Er6i5DZAgBR+A3/qRGupndXhXnWY
	// U/p+vDjRmMEomXrvd8e0sqB0G1f6S704jR1E2wIDAQABAoIEAQCLazQe7h7DB6PX
	// M8MwHHjIEipfsyEbJIcdfQFEZmziRRabSGgEwyiDHKvoS22BHkT7QA0Fz0I6DPeo
	// w5olq2eAGp6s/2SOHJ3r59Iagu0ORZIZ4JilKdqvHGbh41PKtKYdpg7Tw+qfuDIq
	// dbro1NtoYGT0I9ETsU3VSyhBzQBZPN8tUBo1T3Ck849t9a+s30GOmuE57LVWuNFM
	// vwFCDVLCmZEcCqH4Un/Vpz35MQdaGWEh0fGHyAi9r13pE9xpWnJ3VIlvHrqfRO/4
	// kRLhdDp3qY25eCkc5ELSwuSseqflzu6Av+KQo8N8ztao9PqmTw7O4nltF75sHmS7
	// y53VHUg9e0i2DT31V4O0Mq74AuKcTRcUfCbWro8d/HACo0wxywcy8jV9XOW0d7X5
	// BoKD3fnGz0L+LmKFGQ9DqzzlBbhaeNgR05ZZMeet6waDougkgoaoUvSL/zj9oigT
	// wzG2LpqzxjXqRBnnFtp6kBtildsVtqT6xa4X+3LAeGjgR2pGjckryoai7yHr+QR+
	// XUdI9/t+AngmneAFwI7Yq6hu5Zcy4GRvcYmP+bBrVuMEISYlf5CQxQO72qxaQmlm
	// pUhBuYuqqzaSI4W3v9xcZuCsUhsC/PKdodC87xZWM9/JFIf1ie1Sbwo0ply8R7HS
	// 1HHp7kHWlnhaJppnxAt4MvnNazpCoIwr1ExCLdSlfLb2fRJ+fa2awwz9z8HjWlxj
	// h8odAeyQMPLlVpcGsYuzVYRnDilN66u45VG4d8jDOt0CWgeJdObnejg2ngP23JlU
	// KsAqEikJYVaGu8GBzHRuwPjWRyJmmDZlxUjIUQxrHyAVQEoxPRpraTgxy+Y9vcnb
	// 9N2BPpt+e3b5dxyKgzEbcFUJy5rjlLZYN5Ow+gY6op5w6F+NYoN+I8EhnBszLsTr
	// yfI5dApiDjMEUSAyFfLML1zt+xNmRMIb5t4N+qaRnjZzsHnbfC7hAv0NIeZIWGzD
	// iTyyLeNbcHFHWIzoaymJsryZarL6HFdFrOJ6aYREPTp/8k2pp02x+89k6rYAjuCZ
	// XqFJ0ine1OZMQhMbooDiE0O1tWY71VR857dBpLcW7Nb8JSRSVd2CLTgBo4xSH0PH
	// DmbmP2+6fg9d35bHPkX084w3mX8ML+gC1kwhL9jZQsdLNCuCniY7+1szWhXfhOJQ
	// WU9PbBoNdoiA6kScwBhEdhplLCuNV3QBnnQCm0O3ylEYodwemRhgNgw5LovKZQRy
	// ewnf8SroAG39bwQK+TvzzIQ6470yKN2lORtvfIJgnpXj4nn2MLju/kvBOLtj6v3d
	// B0sOMT1lSkySwzA3aQ2Z/RXn6gzSZdwuOGj8ZFAYy4s5he1HQxJMn/Tk3cfIO3/B
	// Dl/NONgBAoICAQDQYgA5hoXpLrYpin2U/c4lkwlGFv5x7KlS3mMHnMYxLaObKCSK
	// +Q8qhBqlJRj7j+wHKKicQsTqW/zqf00ZsugQLKpvHbdlRbMtN/OjeGCcJ9WdMNd1
	// tTMWBsM7AwrwRM7maP/BYsQWhM/HEOdiptDrR0Q4RmRzKAIhXqh031cnhk/l69nN
	// k7mCTMzKpY5SMqoPIKdyyU1xsLMODc+2EjskWJA2N7LuaRxh8C/9y9Zat9uHbQPR
	// 6TfwoAuG8qxXtDHVJ1GpxQGZnV2lLD5c3IQckkOQzSUwQLR4uMeADzu7//vBcVxk
	// fTM0mtY9s5Ht64WlD4MWtJCRhCNiwbObVHz0fbIZSRa1kBfp+75ktOZtaB9OP4m8
	// wupwvA8KHcvOX2/uoBm6NXXE6JKK8sBrZ3wdW/+bcdCAlh96Qg7o+UbWlzajnkIa
	// xmhOBG8gqBFk1ztil3Z1MfWLZqu0tU3Qb68LepPFqnJ75tUZf0lRDePT1IVUVa2F
	// b4TOeQjMGxPykovT+qd8dhkBLiD7Lphda+568iegvBwXa/Rpfv83Somk+/4fayjj
	// ZZ0c41Cri4e3BM8MQAQBbcMSA+DKwepTnGyyl3tDdD7sGx75Zp+RBbeuaPMGXKnY
	// CKif9sV4UsD1WWrUmPTps7Hshc+xIdN1y9atfyVA1aktS99ozLmwhXZzgQKCAgEA
	// y+mnlonAPY5yOImLN4VR6pmLSa45yW1jjGig/QCH4aVtW+wr5bidBFtiVmhP+oYD
	// Um7qscxmM2FJVg8mevO2/1yNm8sE97AbkV5Z6aQ1nttpFNzZjS4nWqxjhrY+kIDy
	// IqmMMVwjhmhode4EB2in6iTN2qoQH0PGvN5gJN17F/rV0tREx8RAC6MYGhaHC3G3
	// wxSpxgesWGbIE+qnBiWmrdM5aqTx/KnG0dXABhm2Z4wYlM23jQBOAI63tsM9gqWS
	// ZTva87uyH5evDR3mEAG+CzF4fuEZa8yKXcW2fYrkJLJ3XTWvNIK63bMjMVtVj05b
	// 2S6Ip1BA3jO5Zydbm5tny7wMjSPXQQ526a5nFU0QQZy+4QKFCC5wrtvX08JTdd9h
	// UyQarUYJ2R7aFmHYhPN6jARPIRZJSKZiQwz8ZdzbdSr/I2rTQfVpWR2gLwW6K6qO
	// Tm9bA6VWkHlNz6RPmD/yG9dAfKwK3J5vUGBa201wryjtYr7q1pZq4fc14UfNeHA3
	// KJU/k+TdVOR6J+txIOhaeFjCUOKdUFatZpQrEHFAH3SEcijdkDKypoUqU1PbLNlS
	// lR9YV66kp/FJH20/02OMeTFevE7H8Qx6vrnZYCp4kQN+0F5A1Uaam/+oQLy50voc
	// CRG/0iIswksdxS4KbOkFpnXaNtOKeIYp/mhwkm4lNlsCggIAJnA9bF8VKotUTNtg
	// CPC8aw/xYP1y2B74qPGewNO5d27cxA0mmIASvRv7MkVik1zcBAILADAM1NI9GCB9
	// X2UyC3HBypdnsgWmZIvypE/bOlLaw5Ez3WGcV87NioAoL2Px6myp6MlINIMxNJ44
	// oGj3Fr9hpSPO+bMCmZc1acGc+Qx8hoP3/ETorlFuYV8inqT726iBWtopL/SRQuFn
	// VEWOpUi06Vti5Tl8Y9CDecR0+Oz7UMLaNolXqUpMdgsjoVX67z++f5t4jRaLJKHB
	// qy1+LSsJzZcZ6gm1gKNNIaEtuqMglHFNwreZj3WUWXNbUtDkKStwWera1AzU3i2u
	// 7NMr7YmWJbXmu30l/CJU34zDCKFxTOTJiolF55+5Anu2kY3nPtRHiqK746OnDHx0
	// 5BsFuQC0aL8Y3+9RRSMUTwNUftcL3CigNwIsUV/eA98cvkY5jRk1X67khV9eAwqC
	// nRVM68gjpA2p6pMM8hQZuRrf0L5eDc10Ia4NiGypI8VG67L+mg6dgIbkW+RfbxYM
	// QwdayVEL4ElywRtHzYIAn/le5JEbMj55H2EqFx8TTWGCmk8HxlG0q3LMzfDrMVuN
	// 1vh+0H7C3RF/Xdwu3YCBpKWyWTjAw4aqSdBnEoLZsDVLpLZmz2qbmtRr2J5iTmFX
	// +V2tbvmKLKh8+X8KafGBuz8TN4ECggIBAKzYYjY4jbugAzHpZRiq6cTyYh8Sj+kt
	// 0DQ0fEH5MTUUBJ7mkk8nLfGDhNe5lBHHhwlY+5e6qubikEdikDuoWnYa/e4RKSTL
	// IpMWZyClEM6XLXuBuObzOGijt7l4wY+7vYOkNDGi6HoGBoXwO3oMPEk86UFw5jma
	// Odveo6CcHLs25AjR5Xtv3tn+ZW+3qMEKo7m4yjASUQSeGvrG7M1JHQA7C3BaYMTL
	// JLy5cwu4v+UF29BGkE/08imG0DMbTAhyUKV/FxmDAHnHGZsEvP7Ue0MBrCHgrKuz
	// tZXHAelo0fpJTOS1nM3mAn/qWOp9gQj4YuK1jyBD5rhzH7VQp1Y/ZLxNqlX5C30Z
	// xbJPIxHrQNxdAW+5swo8U0flEg1tq1E/CEBxWiuGMtLYVASk6+PGfxowjOJ/AZ+T
	// abdsSoYNr/iheaOL6Yb/f00Kz70MHcLlz7qsqXltrPLJ4CsiHFXx3ke7Le5/5rUb
	// lz+KbXl3Kw1Oo+oP6v/6gJ8J7aIIOyLb++VqnXm/hd1vzhwIdtxuGEHq+suVS8PN
	// tS4+akKFvNTs71fVS/nu/4AnYLlBW4eor1e4OioY33l8DX+WG4QT3f7c2nDEVrb8
	// 3Q7OenQJjClfzR7SxiICtwuEqho13032s8z8PpzgWxBWm+WIIU4wN5niTBYtOjLv
	// hqkf3kxBEUoBAoICAEbJI8zfQxr7XmP9H2y41tXv4PmBEE/96RuKw2A+wk39NdmG
	// NFjQ0UPjVPitRwRbrj8zI3VPvwjfSmmuoD2bebaALtXvbn+pFfMj6ZewszMujPvc
	// vhuuOFQwhukXQpRDz2XKaVg4ECgvcXxNozcUjTH7niXdm6+3ABSDofRsm791NeH9
	// QpYRXRlJcV+VkanURYEyAJ8p0HX7jmcKBlREMbV3J7PPtmvrUMJ9Q8n2+cqW0B0T
	// le5m7R/+uZo5pyZw5DWbUNx376XKceEO8j1JIiBmzluVpTGzorrMgPIVdeGuXIhX
	// HnOj8kum1QSloijmzFHlTedE8Omd5vNJKr5ln1JS5ir2BwTGRQwCtSdFVJ0Eulp5
	// R4NfvPke2lb8MbsD1MOrkDDY9yO61B5fQQWxUFGjnKfIAaod0YK+hM1GCUcaGsEt
	// Pia/uM6WPPYssVXeRR0HvmR4qIo5QiozokFPv1Qf2rXhs4PBzhpRIBLCyMaDQcoq
	// I5Sk1f/Dviua1I2giLgFzZLPLLe/ahHZwoOBml5rac4h/WIjbE+7sEFpvjpcqidO
	// x9Nt3IoS/XYFwPWVpZZVnSJxARsav/oL269iDVtvOHY0wqLsmpjEYSVjLZXCdV3b
	// j+Mw5PJUY+7Y2YOPaXhVW1fLxwsIBd+x5jQyeT47Sr3F1trlzJRdAHWhsoKP
	// -----END RSA PRIVATE KEY-----
	// `
	// 	// newPublicKey := ciphers.BytesToPublicKey([]byte(publicKeyString))
	// 	// newPubDecryptedMsg := ciphers.DecryptWithPublicKey(priEncryptedMsg, newPublicKey)

	// 	newPrivateKey := ciphers.BytesToPrivateKey([]byte(privateKeyString))
	// 	newPriDecryptedMsg, _ := ciphers.DecryptWithPrivateKey(pubEncryptedMsg, newPrivateKey)

	// 	fmt.Println(newPriDecryptedMsg)

}
