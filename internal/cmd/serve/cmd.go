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

package serve

import (
	"fmt"
	"github.com/happystraw/text-player/internal/cmd/play"
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"
)

var serveFlags = struct {
	Port int
}{
	Port: 8888,
}

var Cmd = &cobra.Command{
	Use:   "serve",
	Short: "Run on server",
	Long:  "Run on server, play the value of the `msg` field of the POST request form",
	RunE:  run,
}

func init() {
	Cmd.Flags().IntVarP(&serveFlags.Port, "port", "p", 8888, "port to listen on")
}

func run(*cobra.Command, []string) error {
	if err := startServer(); err != nil {
		return fmt.Errorf("start server error: %s", err)
	}
	return nil
}

func startServer() error {
	slog.Info(fmt.Sprintf("Starting server on: http://localhost:%d", serveFlags.Port))
	http.HandleFunc("/", playInServer)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", serveFlags.Port), nil); err != nil {

		return err
	}
	return nil
}

func playInServer(w http.ResponseWriter, r *http.Request) {
	msg := r.PostFormValue("msg")
	if len(msg) > 1 {
		slog.Info(fmt.Sprintf("Got message: %s", msg))
		if err := play.Play(msg); err != nil {
			w.Write([]byte("Fail"))
			slog.Info(fmt.Sprintf("Fail: %s", err))
		} else {
			w.Write([]byte("OK"))
			slog.Info("Success")
		}
	} else {
		w.Write([]byte("No message"))
	}

}
