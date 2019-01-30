package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"os"
	"strings"
)

type User struct {
	username string
	password string
	token string
}

const (
	MAXUSER = 100
)

var UserNameOrPasswordError = errors.New("Username or Password error")
var UserAlreadyExists = errors.New("Username Already Exists")
var UuidGenerateError = errors.New("Uuid Generate Error")

var UserList = make([]User,0)
var message  = make(chan User)
type response struct {
	status bool
	token string
	description string
}
var resp = make(chan response)

func init(){
	go func() {

        var (
        	msg User
        	token string
        	err error
		)
        for{
			msg = <- message
			fmt.Println("New User Register Request Arrive:", msg.username)
			if token, err = NewUser(msg.username, msg.password);err != nil{
				resp<- response{false, "", err.Error()}
			}else{
				resp<- response{true, token, ""}
			}
		}
	}()
}

func NewUser(username string, password string) (string, error){
	var (
		value User
		uuidNew uuid.UUID
		err error
	)
	if username == "" || password == ""{
		return "", UserNameOrPasswordError
	}
	for _, value = range UserList{
		if value.username == username{
			return "", UserAlreadyExists
		}
	}
	if uuidNew, err = uuid.NewV4();err != nil{
		return "", UuidGenerateError
	}
	UserList = append(UserList, User{username, password, uuidNew.String()})
	return uuidNew.String(), nil
}

func main(){
	var (
		count int
		userInfo []string
	)

	input := bufio.NewScanner(os.Stdin)
START:
	fmt.Println("input register info here: ")
	input.Scan()
	count = len(input.Text())
	if count != 0{
		userInfo = strings.Split(input.Text(), ";")
		if len(userInfo) == 1{
			if userInfo[0] == "show"{
				fmt.Println(UserList)
				goto START
			}else{
				fmt.Println("cmd not exist")
				goto START
			}
		}else if len(userInfo) == 2{
			message <- User{userInfo[0], userInfo[1], ""}
		}else{
			fmt.Println("Input Format Error")
			goto START
		}
		var responseTmp response = <- resp
		if responseTmp.status == true{
			fmt.Println("Create user successfully, token:", responseTmp.token)
		}else{
			fmt.Println("Create user failed, reason:", responseTmp.description)
		}
		goto START
	}else{
		goto START
	}
}
