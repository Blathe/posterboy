package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
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
	route  string
	json   bool
	form   bool
	method string
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
	flag.StringVar(&f.route, "r", "/", "add a route to the end of the root domain - ex. /posts")
	flag.BoolVar(&f.json, "j", false, "specifies the POST request to use Content-Type application/json")
	flag.BoolVar(&f.form, "f", false, "specifies the POST request to use Content-Type multipart/form-data")
	flag.StringVar(&f.method, "m", "GET", "specify the method of the request")
}

func main() {

	app_flags := flags{}
	app_flags.loadFlags()
	loadConfig()
	flag.Parse()

	fmt.Printf("Config loaded: %s, %s", viper.Get("domain"), viper.Get("port"))

	switch app_flags.method {
	case "GET":
		{
			url := fmt.Sprintf("%s:%s%s", viper.Get("domain"), viper.Get("port"), app_flags.route)
			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			val, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(val))
		}
	}
}
