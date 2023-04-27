package config_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Prastiwar/Go-flow/config"
)

func Example() {
	// Provide creates new Source instance with provided configs.
	cfg := config.Provide(
	// // { "queryTimeout": "10s" }
	// config.NewFileProvider("config.json", decoders.NewJson()),
	// // --dbName="my-collection" --errorDetails=true
	// config.NewFlagProvider(
	// 	config.StringFlag("dbName", "name for database"),
	// 	config.BoolFlag("errorDetails", "should show error details"),
	// ),
	// // CONNECTION_STRING="mongodb://localhost:8089"; ERROR_DETAILS="false"
	// config.NewEnvProvider(),
	)

	// ShareOptions provides options to set default options for all provider.
	cfg.ShareOptions(
		// KeyInterceptor allows to intercept field name before it's used to find it in provider.
		// It's useful when you want to use different field names than they're defined in struct.
		// For example you can use tag to define field name and intercept it there.
		config.WithInterceptor(func(providerName string, field reflect.StructField) string {
			if providerName == config.EnvProviderName {
				return strings.ToUpper(field.Name)
			}
			return field.Name
		}),
	)

	// Use default values for options in case they are not included in providers.
	err := cfg.SetDefault(
		config.Opt("connectionString", "mongodb://localhost:27017"),
		config.Opt("dbName", "go-flow"),
		config.Opt("errorDetails", true),
		config.Opt("queryTimeout", time.Second*15),
		config.Opt("access-key", "ABC123EFGH456IJK789"),
	)
	if err != nil {
		// some default value couldn't be marshaled as json
		panic(err)
	}

	type DbOptions struct {
		DbName           string
		ConnectionString string
		ErrorDetails     bool
		QueryTimeout     time.Duration
		AccessKey        string `json:"access-key"`
	}

	var dbOptions DbOptions

	// dbOptions will be loaded starting from the first passed provider up to the last one.
	// All values will be also overridden by each provider in this order.
	// The default value is not overridden by provider if it doesn't exist in it.
	err = cfg.Load(context.Background(), &dbOptions)
	if err != nil {
		// One of the providers failed to load config values
		panic(err)
	}

	fmt.Println(dbOptions.DbName)
	fmt.Println(dbOptions.ConnectionString)
	fmt.Println(dbOptions.ErrorDetails)
	fmt.Println(dbOptions.QueryTimeout)
	fmt.Println(dbOptions.AccessKey)

	// Output:
	// go-flow
	// mongodb://localhost:27017
	// true
	// 15s
	// ABC123EFGH456IJK789
}

func ExampleBind() {
	type DbOptions struct {
		AccessKey string `json:"access-key"`
	}

	type AccessOptions struct {
		AccessKey string
	}

	dbOptions := DbOptions{
		AccessKey: "ABC123EFGH456IJK789",
	}
	var aOptions AccessOptions

	fmt.Println(dbOptions.AccessKey)
	fmt.Println(aOptions.AccessKey)

	// Bind will try to copy corresponding field from dbOptions to aOptions
	err := config.Bind(&dbOptions, &aOptions)
	if err != nil {
		// Probably field type mismatch
		panic(err)
	}

	fmt.Println(aOptions.AccessKey)

	// Output:
	// ABC123EFGH456IJK789
	//
	// ABC123EFGH456IJK789
}
