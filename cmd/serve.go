// Copyright © 2019 Yutao Fang <fangyutao1993@hotmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"net/http"

	"github.com/spf13/cobra"
)

var serveFlags = struct {
	Port string
}{
	Port: "8888",
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run on server",
	Long:  "Run on server, play the value of the `msg` field of the POST request form",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&serveFlags.Port, "port", "p", "8888", "port to listen to on the web server")
}

func startServer() {
	runningEnv = runningInServer
	recordf("Ready on http://localhost:%s \n", serveFlags.Port)
	http.HandleFunc("/", playInServer)
	err := http.ListenAndServe(":"+serveFlags.Port, nil)
	if err != nil {
		panic(err)
	}
}

func playInServer(w http.ResponseWriter, r *http.Request) {
	msg := r.PostFormValue("msg")
	if len(msg) > 1 {
		recordf("Got message: [%s] .\n", msg)
		if err := Play(msg); err != nil {
			record("Fail: ", err)
			w.Write([]byte("Fail"))
		} else {
			w.Write([]byte("OK"))
			record("Success")
		}
	} else {
		w.Write([]byte("No message"))
	}

}
