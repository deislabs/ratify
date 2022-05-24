/*
Copyright The Ratify Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	exConfig "github.com/deislabs/ratify/pkg/executor/config"
	"github.com/deislabs/ratify/pkg/homedir"
	pcConfig "github.com/deislabs/ratify/pkg/policyprovider/config"
	rsConfig "github.com/deislabs/ratify/pkg/referrerstore/config"
	vfConfig "github.com/deislabs/ratify/pkg/verifier/config"
	"github.com/sirupsen/logrus"
)

const (
	ConfigFileName = "config.json"
	ConfigFileDir  = ".ratify"
	PluginsFolder  = "plugins"
)

type Config struct {
	StoresConfig    rsConfig.StoresConfig    `json:"stores,omitempty"`
	PoliciesConfig  pcConfig.PoliciesConfig  `json:"policies,omitempty"`
	VerifiersConfig vfConfig.VerifiersConfig `json:"verifiers,omitempty"`
	ExecutorConfig  exConfig.ExecutorConfig  `json:"executor,omitempty"`
	FileHash        string                   `json:"-"`
}

var (
	initConfigDir         = new(sync.Once)
	homeDir               string
	configDir             string
	defaultConfigFilePath string
	defaultPluginsPath    string
)

func InitDefaultPaths() {
	if configDir != "" {
		return
	}
	configDir = os.Getenv("RATIFY_CONFIG")
	if configDir == "" {
		configDir = filepath.Join(getHomeDir(), ConfigFileDir)

	}
	defaultPluginsPath = filepath.Join(configDir, PluginsFolder)
	defaultConfigFilePath = filepath.Join(configDir, ConfigFileName)
}

func getHomeDir() string {
	if homeDir == "" {
		homeDir = homedir.Get()
	}
	return homeDir
}

func Load(configFilePath string) (Config, error) {

	config := Config{}
	if configFilePath == "" {

		if configDir == "" {
			initConfigDir.Do(InitDefaultPaths)
		}

		configFilePath = defaultConfigFilePath
	}

	file, err := os.OpenFile(configFilePath, os.O_RDONLY, 0644)
	//s, _ := ioutil.ReadAll(file) // copy the file and get content

	if err != nil {
		if os.IsNotExist(err) {
			return config, fmt.Errorf("could not find config file at path %s", configFilePath)
		}
		return config, err
	}

	if err := json.NewDecoder(file).Decode(&config); err != nil && !errors.Is(err, io.EOF) {

		return config, err
	}

	config.FileHash, _ = GetFileHash(file) // todo: file pointer could be pointing at EOF ,need to point back at begining of file , handle error
	defer file.Close()
	return config, nil
}

func GetFileHash(file io.Reader) (fileHash string, err error) {
	hash := sha256.New()
	s, readErr := ioutil.ReadAll(file)

	if readErr != nil {
		log.Fatal(readErr)
	}
	hash.Write(s)
	logrus.Infof("hash of file %v", hex.EncodeToString(hash.Sum(nil)))
	return hex.EncodeToString(hash.Sum(nil)), nil

}

func GetDefaultPluginPath() string {
	if defaultPluginsPath == "" {
		initConfigDir.Do(InitDefaultPaths)
	}
	return defaultPluginsPath
}
