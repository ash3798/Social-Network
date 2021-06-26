package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppPort            int    `split_words:"true" default:"9999"`
	Hostname           string `split_words:"true" default:"localhost"`
	DatabaseUsername   string `split_words:"true" default:"postgres"`
	DatabasePassword   string `split_words:"true" default:"mysecretpassword"`
	DatabasePort       int    `split_words:"true" default:"6789"`
	DatabaseName       string `split_words:"true" default:"postgres"`
	AccessSecret       string `split_words:"true" default:"dummyAccessSecret"`
	TokenExpireTimeSec int    `split_words:"true" default:"900"`
	AuthEnabled        bool   //just for unit testing
}

const (
	like    = "like"
	dislike = "dislike"
	plusOne = "+1"
)

var (
	Manager *Config
	//list of allowed reactions
	AllowedReactions = []string{like, dislike, plusOne}
	//map of reactions to improve search performance of checking valid reaction
	ReactionMap map[string]int
)

func InitEnv() {
	Manager = &Config{}

	err := envconfig.Process("", Manager)
	if err != nil {
		log.Println("Error while loading environment variables")
	}

	Manager.AuthEnabled = true
}

//InitReactions will fill the reactionMap using the entries of AllowedReactions
func InitReactions() {
	log.Printf("Loading allowed reactions : %v \n", AllowedReactions)
	ReactionMap = make(map[string]int)
	for idx, reaction := range AllowedReactions {
		ReactionMap[reaction] = idx
	}
}
