package stringcase_test

import (
	"fmt"

	"github.com/alextanhongpin/pkg/stringcase"
)

func ExampleCamelCase() {
	for _, s := range []string{
		"user_service",
		"party pooper",
		"THE AMAZING SPIDER-MAN",
		"A!@#$SAFS ridiculous",
		"property changer",
		"userID",
		"USER ID.",
		"created+at",
		"!created+at!",
		"this--is--a--slug",
		"address.home",
	} {
		fmt.Println(s + " -> " + stringcase.CamelCase(s))
	}
	// Output:
	//user_service -> userService
	//party pooper -> partyPooper
	//THE AMAZING SPIDER-MAN -> theAmazingSpiderMan
	//A!@#$SAFS ridiculous -> aSafsRidiculous
	//property changer -> propertyChanger
	//userID -> userId
	//USER ID. -> userId
	//created+at -> createdAt
	//!created+at! -> createdAt
	//this--is--a--slug -> thisIsASlug
	//address.home -> addressHome
}

func ExampleSnakeCase() {
	for _, s := range []string{
		"user_service",
		"party pooper",
		"THE AMAZING SPIDER-MAN",
		"A!@#$SAFS ridiculous",
		"property changer",
		"userID",
		"USER ID.",
		"created+at",
		"!created+at!",
		"this--is--a--slug",
		"address.home",
	} {
		fmt.Println(s + " -> " + stringcase.SnakeCase(s))
	}
	// Output:
	//user_service -> user_service
	//party pooper -> party_pooper
	//THE AMAZING SPIDER-MAN -> the_amazing_spider_man
	//A!@#$SAFS ridiculous -> a_safs_ridiculous
	//property changer -> property_changer
	//userID -> user_id
	//USER ID. -> user_id
	//created+at -> created_at
	//!created+at! -> created_at
	//this--is--a--slug -> this_is_aslug
	//address.home -> address_home
}

func ExampleKebabCase() {
	for _, s := range []string{
		"user_service",
		"party pooper",
		"THE AMAZING SPIDER-MAN",
		"A!@#$SAFS ridiculous",
		"property changer",
		"userID",
		"USER ID.",
		"created+at",
		"!created+at!",
		"this--is--a--slug",
		"address.home",
	} {
		fmt.Println(s + " -> " + stringcase.KebabCase(s))
	}
	// Output:
	//user_service -> user-service
	//party pooper -> party-pooper
	//THE AMAZING SPIDER-MAN -> the-amazing-spider-man
	//A!@#$SAFS ridiculous -> a-safs-ridiculous
	//property changer -> property-changer
	//userID -> user-id
	//USER ID. -> user-id
	//created+at -> created-at
	//!created+at! -> created-at
	//this--is--a--slug -> this-is-aslug
	//address.home -> address-home
}

func ExamplePascalCase() {
	for _, s := range []string{
		"user_service",
		"party pooper",
		"THE AMAZING SPIDER-MAN",
		"A!@#$SAFS ridiculous",
		"property changer",
		"userID",
		"USER ID.",
		"created+at",
		"!created+at!",
		"this--is--a--slug",
		"address.home",
	} {
		fmt.Println(s + " -> " + stringcase.PascalCase(s))
	}
	// Output:
	//user_service -> UserService
	//party pooper -> PartyPooper
	//THE AMAZING SPIDER-MAN -> TheAmazingSpiderMan
	//A!@#$SAFS ridiculous -> ASafsRidiculous
	//property changer -> PropertyChanger
	//userID -> UserId
	//USER ID. -> UserId
	//created+at -> CreatedAt
	//!created+at! -> CreatedAt
	//this--is--a--slug -> ThisIsASlug
	//address.home -> AddressHome
}
