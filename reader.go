package ezconfig

import "github.com/BurntSushi/toml"

func ReadConfig(path string, v interface{}) error {
	_, err := toml.DecodeFile(path, v)
	if err != nil {
		return err
	}
	return nil
}
