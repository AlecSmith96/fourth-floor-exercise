package adapters

import (
	"fmt"
	"log"
	"os"

	"github.com/AlecSmith96/fourth-floor-exercise/entities"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	flag "github.com/spf13/pflag"
)

// NewConfig reads in config file parsed through command line and returns the config struct
func NewConfig() (*entities.Config, error) {
	// parse flags into FlagSet
	k := koanf.New(".")
	f := flag.NewFlagSet("configSet", flag.ContinueOnError)
	f.String("config", "", "config file name")
	if err := f.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("parsing flags: %v", err)
	}

	// read config file name from FlagSet.
	fileName, err := f.GetString("config")
	if err != nil {
		log.Fatalf("%v", err)
	}
	if err := k.Load(file.Provider(fileName), yaml.Parser()); err != nil {
		log.Fatalf("error loading file: %v", err)
	}

	// // load into koanf instance
	// if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
	// 	log.Fatalf("error loading config: %v", err)
	// }

	// unmarshal into config struct
	config := &entities.Config{}
	err = k.Unmarshal("", config)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling config: %v", err)
	}

	return config, nil
}
