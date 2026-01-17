package tests

import (
	"fmt"
	"testing"

	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
)

func TestError(t *testing.T) {
	az := MyName("Binar", false)
	fmt.Println(az.Value())

	ax := MyName("Binar", true)
	if ax.IsError() {
		r1 := ax.ErrorRoot()
		aerr := Err[string](ax, "wkwkwk")
		r2 := aerr.ErrorRoot()
		fmt.Println(aerr)
		fmt.Printf("%p\n%p\n", r1, r2)
		return
	}
}

func MyName(name string, b bool) Result[string] {
	if b {
		return Err[string]("Error woyilah!")
	}

	return Ok(&name)
}
