// +build darwin linux

package app

// GuessWowDir return wow directory name
func GuessWowDir() (string, error) {
	return "random", nil
}
