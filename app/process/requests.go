// func makeP4Request(posit *Posit, url string) {
// 	response, err := http.PostForm(url)
// }

package process

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/amiwx/p4_runner_golang/app/model"
)

type P4ResponseData struct {
	MsgType string      `json:"msgtype"`
	Data    model.Posit `json:"data"`
	Version string      `json:"version"`
}

func makeClient() *http.Client {
	client := http.Client{Timeout: 10 * time.Second}
	// {
	// 	//setup a mocked http client.
	// 	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		b, err := httputil.DumpRequest(r, true)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		fmt.Printf("%s", b)
	// 	}))
	// 	defer ts.Close()
	// 	client = ts.Client()
	// }

	return &client
}

func makeP4Request(positPath string, url string, client *http.Client) (P4ResponseData, error) {

	var data P4ResponseData

	//prepare the reader instances to encode
	values := map[string]io.Reader{
		"posit": mustOpen(positPath), // lets assume its this file
		// "filename": strings.NewReader(positPath),
	}
	res, err := Upload(client, url, values)
	if err != nil {
		return data, err
	}

	// decoding
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Printf("%T\n%s\n%#v\n", err, err, err)
		// switch v := err.(type) {
		// case *json.SyntaxError:
		// 	fmt.Println(string(res.Body[v.Offset-40 : v.Offset]))
		// }
	}

	return data, nil
}

func Upload(client *http.Client, url string, values map[string]io.Reader) (res *http.Response, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return nil, err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err = client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
