package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	Domain string `toml:"domain"`
	Port   string `toml:"port"`
}

type flags struct {
	route string
	json  bool
	form  bool
}

func loadConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // look for config in the working directory
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Looks like this is your first time using PosterBoy. Let's get some stuff set up...")
			fmt.Println("Please enter a base domain (ex. 'http://localhost'):")
			reader := bufio.NewReader(os.Stdin)
			// ReadString will block until the delimiter is entered
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				return
			}
			input = strings.TrimSuffix(input, "\n")
			viper.Set("domain", input)
			fmt.Println("Please enter a port (ex. '8080'):")
			reader = bufio.NewReader(os.Stdin)
			// ReadString will block until the delimiter is entered
			input, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				return
			}
			input = strings.TrimSuffix(input, "\n")
			viper.Set("port", input)

			viper.SafeWriteConfig()

			fmt.Println("Awesome. Your config is all set!")
		} else {
			// Config file was found but another error was produced
			fmt.Println("An error has occured...")
		}
	}
}

func (f *flags) loadFlags() {
	flag.StringVar(&f.route, "r", "/", "add a route to your root domain - ex. /posts")
	flag.BoolVar(&f.json, "j", false, "specifies the POST request to use Content-Type application/json")
	flag.BoolVar(&f.form, "f", false, "specifies the POST request to use Content-Type multipart/form-data")
	flag.Parse()
}

func main() {

	app_flags := flags{}
	app_flags.loadFlags()
	loadConfig()

	switch os.Args[1] {
	case "get":
		{
			fmt.Println(app_flags.route)
			domain := viper.GetString("domain")
			port := viper.GetString("port")

			resp, err := http.Get(domain + ":" + port)
			if err != nil {
				fmt.Println(err)
				return
			}
			val, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(val)

		}
	case "post":
		{

		}
	case "del":
		{

		}
	}
}
