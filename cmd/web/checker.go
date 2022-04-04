package main

import (
	"fmt"
	"net/mail"
	"regexp"

	"git.01.alem.school/quazar/forum/pkg/models"
)

var err error = fmt.Errorf("error:incorrect form")

func userSignUpForm(user models.User) error {
	usernameConvention := "^[a-zA-Z0-9]*[-]?[a-zA-Z0-9]*$"

	emailRegexp := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if (models.User{} == user) {
		fmt.Println(err)
		return err
	}

	if user.Password != user.Confirm {
		fmt.Println(err)
		return err
	}

	if user.Username == "" || user.Email == "" || user.Password == "" || user.Confirm == "" {
		fmt.Println(err)
		return err
	}

	if re, _ := regexp.Compile(usernameConvention); !re.MatchString(user.Username) {
		fmt.Println(err)
		return err
	}

	_, err := mail.ParseAddress(user.Email)
	if !emailRegexp.MatchString(user.Email) || err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
