package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Loader struct {
	filenames []string
}

type Option func(*Loader)

func WithFilenames(filenames ...string) Option {
	return func(l *Loader) {
		l.filenames = filenames
	}
}

func NewLoader(opts ...Option) *Loader {
	loader := &Loader{
		filenames: []string{
			".env.local",
			".env.prod",
			".env.test",
			".env",
		},
	}

	for _, opt := range opts {
		opt(loader)
	}

	return loader
}

func (l *Loader) LoadEnv() error {
	var loaded bool
	for _, file := range l.filenames {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}
		err := godotenv.Load(file)
		if err != nil {
			return fmt.Errorf("error loading %s: %w", file, err)
		}
		log.Printf("config: loaded %s", file)
		loaded = true
	}
	if !loaded {
		log.Printf("config: no .env file found, relying on system environment variables")
	}
	return nil
}
