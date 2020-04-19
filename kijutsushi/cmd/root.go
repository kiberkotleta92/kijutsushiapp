/*
Copyright © 2020 Kirill Denisov <kirill.denisov700@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kijutsushi [video file]",
	Short: "Video to text cl app",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("invalid args")
		}
		if err := SendVideo(args[0]); err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func SendVideo(fname string) error {
	c := &http.Client{
		Timeout: 10 * time.Minute,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   1 * time.Minute,
				KeepAlive: 1 * time.Minute,
			}).Dial,
			TLSHandshakeTimeout:   1 * time.Minute,
			ResponseHeaderTimeout: 1 * time.Minute,
			ExpectContinueTimeout: 1 * time.Minute,
			IdleConnTimeout:       1 * time.Minute,
		},
	}

	path := strings.Split(fname, "/")
	name := path[len(path)-1]

	file, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("error opening %w", err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading %w", err)
	}
	reqBody := bytes.NewBuffer(content)
	req, _ := http.NewRequest("POST", "http://130.193.36.65:8080/video", reqBody)
	//req, _ := http.NewRequest("POST", "http://localhost:8080/video", reqBody)
	//resp, err := c.Post("http://130.193.36.65:8080/video", "", reqBody)
	req.Header.Add("name", name)
	resp, err := c.Do(req)
	i := time.Now()
	if err != nil {
		fmt.Println(time.Since(i))
		return fmt.Errorf("error posting %w", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading answer %w", err)
	}
	fmt.Println(string(respBody))
	return nil
}
