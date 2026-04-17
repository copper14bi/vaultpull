// Package prompt provides lightweight interactive confirmation helpers
// for CLI commands in vaultpull.
//
// Usage:
//
//	c := prompt.New()
//	ok, err := c.Ask("Overwrite existing .env file?", false)
//	if err != nil {
//		// handle unrecognised input
//	}
//	if !ok {
//		fmt.Println("Aborted.")
//		return
//	}
//
// The Confirmer reads from stdin by default but accepts any io.Reader
// and io.Writer via NewWithIO, making it straightforward to test without
// mocking os.Stdin.
package prompt
