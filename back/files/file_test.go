package files

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/crypto/ring"
	"github.com/qwertyqq2/filebc/user"
)

func TestCreateDb(t *testing.T) {
	l, err := NewLevelDB()
	if err != nil {
		t.Fatal(err)
	}
	f1 := NewFile("Здарова епты бля! \n\n\n я знаю")
	err = l.insertFile(f1)
	if err != nil {
		t.Fatal(err)
	}
	f2 := NewFile("Че нада")
	err = l.insertFile(f2)
	if err != nil {
		t.Fatal(err)
	}
	f3 := NewFile("Ниче")
	err = l.insertFile(f3)
	if err != nil {
		t.Fatal(err)
	}
	pk1 := ring.GeneratePrivate()
	u1 := user.NewUser(pk1)
	pk2 := ring.GeneratePrivate()
	u2 := user.NewUser(pk2)
	err = l.newUser(u1.Addr.String())
	if err != nil {
		t.Fatal(err)
	}
	err = l.newUser(u2.Addr.String())
	if err != nil {
		t.Fatal(err)
	}
	addrs, err := l.getUsers()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(addrs))
	fmt.Println("TEST ADD BALANCE")
	pk3 := ring.GeneratePrivate()
	u3 := user.NewUser(pk3)
	err = l.newUser(u3.Addr.String())
	if err != nil {
		t.Fatal(err)
	}
	addrs, err = l.getUsers()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(addrs))
	err = l.addBalance(addrs[0].Addr, 10)
	if err != nil {
		t.Fatal(err)
	}
	addrs, err = l.getUsers()
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range addrs {
		fmt.Println("Address: ", a.Addr)
		fmt.Println("Balance: ", a.Bal)
	}
	fmt.Println("New user add")
	pk4 := ring.GeneratePrivate()
	u4 := user.NewUser(pk4)
	err = l.addBalance(u4.Addr.String(), 10)
	if err != nil {
		t.Fatal(err)
	}
	addrs, err = l.getUsers()
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range addrs {
		fmt.Println("Address: ", a.Addr)
		fmt.Println("Balance: ", a.Bal)
	}
}

func TestGetBalance(t *testing.T) {
	l, err := NewLevelDB()
	if err != nil {
		t.Fatal(err)
	}
	pk1 := ring.GeneratePrivate()
	u1 := user.NewUser(pk1)
	err = l.addBalance(u1.Addr.String(), 10)
	if err != nil {
		t.Fatal(err)
	}
	bal, _, err := l.getBalance(u1.Addr.String())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Balance ", bal)
}

func TestGetFiles(t *testing.T) {
	l, err := LoadLevel()
	if err != nil {
		t.Fatal(err)
	}
	files, err := l.allFiles()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		fmt.Println(crypto.Base64EncodeString(f.Id), string(f.Data))
	}
}
