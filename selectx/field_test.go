package selectx

import (
    "fmt"
    "testing"
)

func TestGetInt(t *testing.T) {
    type User struct {
        Name string
        Age  int
    }

    var allU []*User
    allU = append(allU, &User{
        Name: "zzh1",
        Age:  1,
    }, &User{
        Name: "zzh2",
        Age:  2,
    })

    r, err := GetInt(allU, "Age")
    t.Error(fmt.Sprintf("%v", r), "; err:", err)
}

func TestGetString(t *testing.T) {
    type User struct {
        Name string
        Age  int
    }

    var allU []*User
    allU = append(allU, &User{
        Name: "zzh1",
        Age:  1,
    }, &User{
        Name: "zzh2",
        Age:  2,
    })

    r, err := GetString(allU, "Name")
    t.Error(fmt.Sprintf("%v", r), "; err:", err)
}
