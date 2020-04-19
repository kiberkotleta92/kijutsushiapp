package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Client struct {
	internalClient *http.Client
	Token          string
	FolderID       string
}

func SplitAudio(fname string) string {
	outname := strings.Split(fname, ".")[0]

	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(ffmpeg, "-y", "-loglevel", "quiet", "-i", fname, "-vn", "-f", "segment", "-segment_time", "29", "-c:a", "libvorbis", "-ar", "8000", "tmp/"+outname+"%03d.oga")
	if err = cmd.Run(); err != nil {
		fmt.Printf("ffmpeg %s", err.Error())
		log.Fatal(err)
	}

	return outname
}

func (c *Client) AudioToText(fname string) string {
	file, err := os.Open(fname)
	if err != nil {
		fmt.Printf("opening %w", err)
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("reading %w", err)
		log.Fatal(err)
	}
	reqBody := bytes.NewBuffer(content)

	u, _ := url.Parse("https://stt.api.cloud.yandex.net/speech/v1/stt:recognize")
	q := url.Values{}
	q.Add("lang", "ru-RU")
	q.Add("format", "oggopus")
	q.Add("sampleRateHertz", "8000")
	q.Add("folderId", c.FolderID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("POST", u.String(), reqBody)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.internalClient.Do(req)
	if err != nil {
		fmt.Printf("posting %w", err)
		log.Fatal(err)
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	var res map[string]string

	json.Unmarshal(respBody, &res)
	fmt.Println(res)

	return res["result"]
}

func (c *Client) FileToText(fname string) string {
	var name string
	var res []string
	targetName := SplitAudio(fname)

	d, _ := ioutil.ReadDir("./tmp")
	for _, f := range d {
		name = "./tmp/" + f.Name()
		if !strings.Contains(name, targetName) {
			continue
		}
		t := c.AudioToText(name)
		os.Remove(name)
		res = append(res, t)
	}

	return strings.Join(res, " ")
}

func (c *Client) FileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got " + r.Header.Get("name"))

	defer r.Body.Close()

	filename := "./videotmp/" + time.Now().String() + r.Header.Get("name")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "eror", http.StatusBadRequest)
	}
	fmt.Println("read")
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		fmt.Println("eroror")
		http.Error(w, "eror", http.StatusBadRequest)
	}
	fmt.Println("written")
	text := c.FileToText(filename)
	os.Remove(filename)
	fmt.Println("filetotexted " + text)
	fmt.Fprintln(w, text)
}

func main() {
	c := Client{
		internalClient: &http.Client{},
		Token:          TOKEN,
		FolderID:       FOLDERID,
	}

	srv := &http.Server{
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		Addr:         ":8080",
		IdleTimeout:  1 * time.Minute,
	}

	http.HandleFunc("/video", c.FileHandler)
	fmt.Println("start listening")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
