package ezconfig

import "github.com/BurntSushi/toml"

// ReadConfig reads a file and will unmarshal the results into the given structure
func ReadConfig(path string, v interface{}) error {
	_, err := toml.DecodeFile(path, v)
	if err != nil {
		return err
	}
	return nil
}
