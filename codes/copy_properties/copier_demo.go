package main

import (
	"fmt"
	"github.com/jinzhu/copier"
)

type RoleA struct {
	RoleId   int
	RoleName string
}

type RoleB struct {
	RoleId   int
	RoleName string
}

type User struct {
	Name         string
	Role         RoleA
	Age          int32
	EmployeeCode int64
	Salary int
}

type Employee struct {
	Name         string
	Role         RoleB
	Age          int32
	EmployeeCode int64
	Salary int
}

func CopyProperties() {
	user := User{Name: "KevinYan", Age: 18, Salary: 200000}
	user.Role = RoleA{
		RoleId:   2,
		RoleName: "Admin",
	}
	employee := new(Employee)


	copier.Copy(employee, &user)

	fmt.Printf( "%v\n", employee)
	fmt.Printf("%v\n", user)
}

func main() {
	CopyProperties()
}
