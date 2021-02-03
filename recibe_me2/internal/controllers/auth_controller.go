package controllers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/mailjet/mailjet-apiv3-go"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"recibe_me/configs"
	"recibe_me/configs/constants"
	"recibe_me/internal/helpers"
	"recibe_me/internal/models"
	"recibe_me/internal/services"
	"recibe_me/pkg/crypto"
	"strings"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var userData models.UserModel

	err := decoder.Decode(&userData)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	// Codificamos la contraseña
	encPass, err := crypto.EncodePassword([]byte(userData.Password), configs.SecurityCfg.PasswordEncKey)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	userData.Password         = base64.StdEncoding.EncodeToString([]byte(encPass))
	userData.Email 	          = strings.ToLower(userData.Email)
	userData.VerificationCode = helpers.EncodeToString(5)
	userData.Verified         = false

	// Verificamos que el correo no exista
	result, err := helpers.UsersCollection.Find(bson.M{"email": userData.Email}).Count()
	if result > 0 {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_USER_EXIST, err, nil)
		return
	}

	// Insertamos el usuario en la base de datos
	err = helpers.UsersCollection.Insert(userData)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_INSERT_DATA, err, nil)
		return
	}

	// Enviamos por correo el codigo de verificacion
	resp, err := SendVerificationCode(userData)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_SEND_EMAIL, err, resp)
		return
	}

	helpers.Response(w, http.StatusOK, constants.SUCCESS, err, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var userAuthData models.UserAuthModel

	err := decoder.Decode(&userAuthData)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	userAuthData.Email = strings.ToLower(userAuthData.Email)

	user := models.UserModel{}

	// Buscamos el usuario en la base de datos
	err = helpers.UsersCollection.Find(bson.M{"email": userAuthData.Email}).One(&user)
	if err != nil {
		helpers.Response(w, http.StatusUnauthorized, constants.ERR_USER_NOT_FOUND, err, nil)
		return
	}

	// Decodificamos la contraseña
	decPass, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_DECODE, err, nil)
		return
	}

	// Chequeamos la contraseña
	isValid, err := crypto.CheckPassword([]byte(userAuthData.Password), []byte(decPass), configs.SecurityCfg.PasswordEncKey)
	if !isValid {
		helpers.Response(w, http.StatusUnauthorized, constants.ERR_PASS_NOT_VALID, err, nil)
		return
	}

	user.Password = ""

	// Generamos el token de autenticacion
	token, err := crypto.CreateTokenString(&crypto.Claims{
		Type: constants.User,
		StandardClaims: jwt.StandardClaims{
			Id: user.ID.Hex(),
		},
	}, configs.SecurityCfg.TokenSecret, configs.SecurityCfg.TokenDuration)

	helpers.Response(w, http.StatusOK, constants.SUCCESS, err, models.UserLogResponseModel{Token: token, User: user})
}

func ResendVerificationCode(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)

	if !bson.IsObjectIdHex(params["userId"]) {
		helpers.Response(w, http.StatusNotFound, constants.ERR_USER_NOT_FOUND, nil, nil)
		return
	}

	userId := bson.ObjectIdHex(params["userId"])

	userData := models.UserModel{}

	// Se busca la informacion del usuario
	err := helpers.UsersCollection.FindId(userId).One(&userData)
	if err != nil{
		helpers.Response(w, http.StatusBadRequest, constants.ERR_USER_NOT_FOUND, err, nil)
		return
	}

	// Se actualiza el codigo de verificación en la base de datos
	var verificationCode = helpers.EncodeToString(5)
	err = helpers.UsersCollection.UpdateId(userId, bson.M{"$set": bson.M{"verification_code":verificationCode}})
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_UPDATE_DATA, err, nil)
		return
	}

	// Se reenvia el codigo al correo del usuario
	userData.VerificationCode = verificationCode
	resp, err := SendVerificationCode(userData)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_SEND_EMAIL, err, resp)
		return
	}

	helpers.Response(w, http.StatusOK, constants.SUCCESS, nil, nil)
}

func VerificationAccount(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)

	if !bson.IsObjectIdHex(params["userId"]) {
		helpers.Response(w, http.StatusNotFound, constants.ERR_USER_NOT_FOUND, nil, nil)
		return
	}
	userId := bson.ObjectIdHex(params["userId"])

	code := params["verificationCode"]

	userData := models.UserModel{}

	// Se busca la informacion del usuario
	err := helpers.UsersCollection.FindId(userId).One(&userData)
	if err != nil{
		helpers.Response(w, http.StatusBadRequest, constants.ERR_USER_NOT_FOUND, err, nil)
		return
	}

	// Se valida el codigo ingresado
	if userData.VerificationCode != code{
		helpers.Response(w, http.StatusBadRequest, constants.ERR_INVALID_CODE, nil, nil)
		return
	}

	// Se actualiza el codigo de verificación en la base de datos
	err = helpers.UsersCollection.UpdateId(userId, bson.M{"$set": bson.M{"verified":true}})
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_UPDATE_DATA, err, nil)
		return
	}

	helpers.Response(w, http.StatusOK, constants.SUCCESS, nil, nil)
}

func SendVerificationCode(userData models.UserModel)(*mailjet.ResultsV31, error){
	var to = services.Recipient{Name:userData.FirstName, Email:userData.Email}

	var info = services.Info{
		ApiKeyPrivate:  configs.SecurityCfg.MjApiKeyPrivate,
		ApiKeyPublic:   configs.SecurityCfg.MjApiKeyPublic,
		FromRecipient:  services.Recipient{Name:"Cristobal Muñoz", Email:"cmunoz21x@gmail.com"},
		ToRecipient:    to,
		Code:           userData.VerificationCode,
	}

	resp, err := info.SendVerificationCode()
	if err != nil {
		return nil, err
	}

	return resp, err
}
