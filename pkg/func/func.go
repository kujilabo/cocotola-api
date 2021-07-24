package f

import (
	"errors"
	"fmt"
	"math/rand"
)

func Abc(a int) int {
	return a * a
}

func Def(a int) int {
	return a * a
}

func Ghi(a int) int {
	b := 0
	b = b + 1
	b = b + 1
	fmt.Println(b)
	return a * a
}

func Jkl(a int) int {
	b := 0
	b = b + 1
	b = b + 1
	b=rand.Intn(32)
	validate(b)
	fmt.Println(b)
	return a * a
}

func validate(age int) error {
    if age < 20 {
        return errors.New("age should be 20 or more")
    }
    fmt.Println("ok~")
    return nil
}


func cyclo(age int) error {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
    if age < 20 {
        return errors.New("age should be 20 or more")
    }
}}}}}}}}}}}}}}}}}}
    fmt.Println("ok~")
    return nil
}

