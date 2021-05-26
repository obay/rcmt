package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var webclientCmd = &cobra.Command{
	Use:   "webclient",
	Short: "Run a GUI in your web browser",
	Long:  `Run a GUI in your web browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		webclientCmdMain(args)
	},
}

func init() {
	rootCmd.AddCommand(webclientCmd)
}

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func webclientCmdMain(args []string) {
	portNumber := "8090"
	// start a web server
	http.HandleFunc("/", hello)
	http.HandleFunc("/headers", headers)

	// launch the default browser pointing to the local port you just listed to
	openbrowser("http://localhost:" + portNumber)

	http.ListenAndServe(":"+portNumber, nil)
}
