package npapi

import "fmt"

func ExampleUserInfo() {
	user, err := UserInfo("0x39d27d66c14f7372553b1ba59833c6ba8981a76a")
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have mined %.3f ETH!\n", user.Balance)
	// Output: You have mined 20.344 ETH.
}
