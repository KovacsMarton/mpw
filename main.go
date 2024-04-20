package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"

	//"reflect"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Numbers         bool `yaml:numbers"`
	Lowercase       bool `yaml:lowercase"`
	Uppercase       bool `yaml:uppercase"`
	BeginWithLetter bool `yaml:beginwithletter"`
	IncludeSymbols  bool `yaml:includesymbols"`
	NoSimilar       bool `yaml:nosimilar"`
	NoDuplicate     bool `yaml:noduplicate"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var p = fmt.Println
	args := os.Args[1:]
	var menu = `      .--------.
    / .------. \
   / /        \ \
   | |        | |
  _| |________| |_
.' |_|        |_| '.
'._____ ____ _____.'
|     .'____'.     |
'.__.'.'    '.'.__.'	Marci
'.__  | MPW |  __.'	Password
|   '.'.____.'.'   |    Generator
'.____'.____.'____.'
'.________________.'`

	var config Config
	defaultconfig := Config{
		Numbers:         false,
		Lowercase:       true,
		Uppercase:       false,
		BeginWithLetter: true,
		IncludeSymbols:  false,
		NoSimilar:       false,
		NoDuplicate:     true,
	}

	f, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("Loading config file unsuccessful, creating config.yaml with default settings.\n")
		writeToConfig(&defaultconfig)
	}

	if err := yaml.Unmarshal(f, &config); err != nil {
		fmt.Println("Error during parsing of config file, reverting to default settings.")
		config = *writeToConfig(&defaultconfig)
	}

	switch a := len(args); a {
	case 0:
		p(menu)
		fmt.Printf("\n\nCurrent Options:\n\n")
		displayConfig(&config)
		printUsage()
	case 1:
		if args[0] != "generate" && args[0] != "config" {
			printUsage()
		} else if args[0] == "generate" {
			generate(&config, 10, 5)
		} else {
			changeConfig(&config)
		}
	case 2:
		if args[0] == "generate" {
			if length, err := strconv.Atoi(args[1]); err != nil {
				printUsage()
			} else {
				if length >= 6 && length <= 50 {
					generate(&config, length, 4)
				} else {
					fmt.Println("Length must be between 6 and 50 (inclusive).")
				}
			}
		} else {
		}
	case 3:
		if args[0] == "generate" {
			if length, err := strconv.Atoi(args[1]); err != nil {
				printUsage()
			} else {
				if length >= 6 && length <= 50 {
					if quantity, err := strconv.Atoi(args[2]); err != nil {
						printUsage()
					} else if quantity <= 0 || quantity > 100 {
						fmt.Println("Number of generated passwords must be greater than 0, and 100 at maximum.")
					} else {
						generate(&config, length, quantity)
					}
				} else {
					fmt.Println("Length must be between 6 and 50 (inclusive).")
				}
			}
		}
	default:
		printUsage()
	}
}

func changeConfig(config *Config) {
	answers := []string{}
	prompt := &survey.MultiSelect{
		Message: "\nCreate a new configuration. The generated passwords:",
		Help:    "Currently active config options are marked with a checkmark, but they do not carry over when creating a new configuration. Please select all the options that you need.\nTo keep current changes, send an interrupt.",
		Options: []string{"1. Can include Numbers", "2. Can include lowercase letters", "3. Can include uppercase letters", "4. Begin with a letter", "5. Can include symbols", "6. Don't use similar characters", "7. Don't use a character more than once"},
	}

	configstruct := reflect.ValueOf(*config)
	//typeOfC := configstruct.Type()
	for i := 0; i < configstruct.NumField(); i++ {
		if configstruct.Field(i).Interface() == true {
			fmt.Printf("asd")
			prompt.Options[i] += strings.Repeat(" ", 39-len(prompt.Options[i]))
			prompt.Options[i] += "✔"
		}
		//fmt.Printf("Field: %s\tValue: %v\n", typeOfC.Field(i).Name, configstruct.Field(i).Interface())
		configstruct.Field(i).Interface()
	}
	fmt.Println(prompt.Options[0])

	err := survey.AskOne(prompt, &answers, survey.WithRemoveSelectAll(), survey.WithRemoveSelectNone()) // the strings corresponding to the ticked checkboxes are added to 'answers',
	if err != nil {
	} else {

		//so I assigned a number to each option so I can see which ones are selected.
		boolvalues := [7]bool{}
		for _, v := range answers {
			if index, err := strconv.Atoi(string(v[0])); err != nil {
				log.Fatal()
			} else {
				boolvalues[index-1] = true
			}
		}
		//We have the array of boolean values, we need to update the current config with these values.
		//This doesn't work because reflect doesn't change the original object.

		// configstruct := reflect.ValueOf(*config)
		// typeOfC := configstruct.Type()
		// for i := 0; i < configstruct.NumField(); i++ {
		// 	fmt.Printf("Field: %s\tValue: %v\n", typeOfC.Field(i).Name, configstruct.Field(i).Interface())
		// 	configstruct.Field(i).Interface()

		// }
		//This is something I assumed I should avoid
		config.Numbers = boolvalues[0]
		config.Lowercase = boolvalues[1]
		config.Uppercase = boolvalues[2]
		config.BeginWithLetter = boolvalues[3]
		config.IncludeSymbols = boolvalues[4]
		config.NoSimilar = boolvalues[5]
		config.NoDuplicate = boolvalues[6]
		writeToConfig(config)
	}
}

func writeToConfig(conf *Config) *Config {
	out, err := yaml.Marshal(&conf)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./config.yaml", out, 0644)
	check(err)

	return conf
}

func displayConfig(conf *Config) {

	configstruct := reflect.ValueOf(*conf)

	for i := 0; i < configstruct.NumField(); i++ {
		switch option := i; option {
		case 0:
			fmt.Printf("Include numbers                                   ")
			if configstruct.Field(i).Interface() == true {
				fmt.Println("[x]")
			} else {
				fmt.Println("[ ]")
			}
		case 1:
			fmt.Printf("Include lowercase letters                         ")
			if configstruct.Field(i).Interface() == true {
				fmt.Println("[x]")
			} else {
				fmt.Println("[ ]")
			}
		case 2:
			fmt.Printf("Include uppercase letters                         ")
			if configstruct.Field(i).Interface() == true {
				fmt.Println("[x]")
			} else {
				fmt.Println("[ ]")
			}
		case 3:
			fmt.Printf("Password begins with a letter                     ")
			if configstruct.Field(i).Interface() == true {
				fmt.Println("[x]")
			} else {
				fmt.Println("[ ]")
			}
		case 4:
			fmt.Printf("Include symbols                                   ")
			if configstruct.Field(i).Interface() == true {
				fmt.Println("[x]")
			} else {
				fmt.Println("[ ]")
			}
		case 5:
			fmt.Printf("Exclude similar looking characters like (o0|Il)   ")
			if configstruct.Field(i).Interface() == true {
				fmt.Println("[x]")
			} else {
				fmt.Println("[ ]")
			}
		case 6:
			fmt.Printf("No duplicate characters                           ")
			if configstruct.Field(i).Interface() == true {
				fmt.Println("[x]")
			} else {
				fmt.Println("[ ]")
			}
		default:

			fmt.Printf("\nasd")
		}
		//typeOfC := configstruct.Type()
		//fmt.Printf("Field: %s\tValue: %v\n", typeOfC.Field(i).Name, configstruct.Field(i).Interface())

	}
}

func generate(config *Config, length int, quantity int) {
	if !config.Numbers && !config.Uppercase && !config.Lowercase {
		fmt.Println("You must select at least one set of characters.")
		return
	}
	if !config.Uppercase && !config.Lowercase && config.BeginWithLetter {
		fmt.Println("Password cannot begin with a letter if no letters are included in the pool of possible characters.")
		return
	}

	var pool = ""
	var letters = ""

	if config.Numbers && !config.NoSimilar {
		numbers := "1234567890"
		pool += numbers
	} else if config.Numbers { // i, l, 1, L, o, 0, O,
		numbers := "23456789"
		pool += numbers
	}

	if config.Uppercase && !config.NoSimilar {
		uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		pool += uppercase
		letters += uppercase
	} else if config.Uppercase {
		uppercaseNoSimilar := "ABCDEFGHIJKMNPQRSTUVWXYZ"
		pool += uppercaseNoSimilar
		letters += uppercaseNoSimilar
	}

	if config.Lowercase && !config.NoSimilar {
		lowercase := "abcdefghijklmnopqrstuvwxyz"
		pool += lowercase
		letters += lowercase
	} else if config.Lowercase {
		lowercaseNoSimilar := "abcdefghjkmnpqrstuvwxyz"
		pool += lowercaseNoSimilar
		letters += lowercaseNoSimilar
	}

	if config.IncludeSymbols && !config.NoSimilar {
		symbols := `"!";#$%&'()*+,-./:;<=>?@[]^_/{|}~"`
		symbols += "`"
		pool += symbols
	} else if config.IncludeSymbols {
		symbols := `"!";#$%&'()*+,-./:;<=>?@[]^_/{}~"`
		symbols += "`"
		pool += symbols
	}

	if config.NoDuplicate && length > len(pool) {
		fmt.Println("The pool of available characters is not big enough to avoid duplication with the given length. Shorten the length, or extend the set of characters.")
		return
	}

	if config.NoDuplicate && config.BeginWithLetter {
		fmt.Printf("a")
		fmt.Printf("\n")
		for i := 0; i < quantity; i++ {
			rand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))
			pw := ""
			pw += string(letters[rand.Intn(len(letters))])
			currentLength := 1
			var duplicateCheck = make([]bool, len(pool))

			for currentLength != length {
				randomChar := rand.Intn(len(pool) - 1)
				if !duplicateCheck[randomChar] {
					pw += string(pool[randomChar])
					duplicateCheck[randomChar] = true
					currentLength++
				}
			}

			fmt.Println(pw)
			pw = ""
		}

	} else if config.NoDuplicate {
		fmt.Printf("b")
		fmt.Printf("\n")
		for i := 0; i < quantity; i++ {
			rand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))
			pw := ""
			currentLength := 0
			var duplicateCheck = make([]bool, len(pool))

			for currentLength != length {
				randomChar := rand.Intn(len(pool) - 1)
				if !duplicateCheck[randomChar] {
					pw += string(pool[randomChar])
					duplicateCheck[randomChar] = true
					currentLength++
				}
			}

			fmt.Println(pw)
			pw = ""
		}
	} else if config.BeginWithLetter {
		fmt.Printf("c")
		fmt.Printf("\n")
		for i := 0; i < quantity; i++ {
			rand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))
			pw := ""
			pw += string(letters[rand.Intn(len(letters))])
			currentLength := 1

			for currentLength != length {
				randomChar := rand.Intn(len(pool) - 1)
				pw += string(pool[randomChar])
				currentLength++
				if currentLength == length {
					fmt.Println(pw)
					pw = ""
				}
			}
		}

	} else {
		fmt.Printf("\n")
		for i := 0; i < quantity; i++ {
			rand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))
			pw := ""
			pw += string(letters[rand.Intn(len(letters))])
			currentLength := 0

			for currentLength != length {
				randomChar := rand.Intn(len(pool) - 1)
				pw += string(pool[randomChar])
				currentLength++
			}

			fmt.Println(pw)
			pw = ""
		}
	}

}
func printUsage() {
	fmt.Printf("\n\nUsage:\n'mpw config'   			     - to update password generation rules.\n'mpw generate <length> <quantity>'   - to generate passwords.\n\n<length> is an integer between 6 and 50, if not specified, the value defaults to 16\n<Quantity> specifies the number of passwords to be generated. Default value: 4, Maximum: 100\n")
}
