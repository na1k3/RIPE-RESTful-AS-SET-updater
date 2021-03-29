package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func getFromRipe(adddelete, ASorASSET, password, assetname string) (err error, RipeErrorSeverityString, RipeErrorMesageString string) {
	var (
		ver3, ver4     string
		resp, resp1    *http.Response
		PutRequestBody io.Reader
	)

	client := &http.Client{}
	request := "https://rest.db.ripe.net/ripe/as-set/" + assetname
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err = client.Do(req)
	defer resp.Body.Close()

	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		ver1, _ := sjson.Delete(string(body), "terms-and-conditions")
		ver1, _ = sjson.Delete(ver1, "version")
		ver1, _ = sjson.Delete(ver1, "objects.object.0.primary-key")
		namesArr := gjson.Get(ver1, "objects.object.0.attributes.attribute.#.name")
		i := 0
		for _, name1 := range namesArr.Array() {
			if name1.String() == "last-modified" {
				ver3, _ = sjson.Delete(ver1, "objects.object.0.attributes.attribute."+strconv.Itoa(i))
			}
			i++
		}

		switch adddelete {

		case "delete":

			values := gjson.Get(ver3, "objects.object.0.attributes.attribute.#.value")
			i = 0
			for _, name := range values.Array() {
				if name.String() == ASorASSET {
					ver4, _ = sjson.Delete(ver3, "objects.object.0.attributes.attribute."+strconv.Itoa(i))
				}
				i++
			}
			if len(ver4) == 0 {
				err = errors.New("No matches to delete!")
				return err, RipeErrorSeverityString, RipeErrorMesageString
			} else {
				//fmt.Println(ver4)//for debug
				PutRequestBody = bytes.NewReader([]byte(ver4))
			}

		case "add":
			QWE := map[string]string{"href": "https://rest.db.ripe.net/ripe/aut-num/" + ASorASSET, "type": "locator"}
			value, _ := sjson.Set(ver3, "objects.object.0.attributes.attribute.-1", map[string]interface{}{"link": QWE, "name": "members", "value": ASorASSET, "referenced-type": "aut-num"})

			//fmt.Println(value)//for debug
			PutRequestBody = bytes.NewReader([]byte(value))

		}
	}

	if err == nil {
		//PutRequestBody = bytes.NewReader([]byte{1}) //TEST ERROR

		request1 := "https://rest.db.ripe.net/ripe/as-set/" + assetname + "?password=" + password
		req1, _ := http.NewRequest("PUT", request1, PutRequestBody)
		req1.Header.Add("Content-Type", "application/json")
		req1.Header.Add("Accept", "application/json")
		resp1, err = client.Do(req1)
		body, _ := ioutil.ReadAll(resp1.Body)
		// fmt.Println(string(body))  //for debug
		defer resp1.Body.Close()

		RipeErrors := gjson.Get(string(body), "errormessages")
		if RipeErrors.Exists() {
			//fmt.Println(RipeErrors.String()) //for debug
			RipeErrorSeverity := gjson.Get(RipeErrors.String(), "errormessage.0.severity")
			RipeErrorMesage := gjson.Get(RipeErrors.String(), "errormessage.0.args.0.value")
			RipeErrorSeverityString = RipeErrorSeverity.String()
			RipeErrorMesageString = RipeErrorMesage.String()
		}
	}

	return err, RipeErrorSeverityString, RipeErrorMesageString

}

func main() {

	if len(os.Args) < 4 {
		fmt.Println("USAGE:  " + os.Args[0] + " add/delete AS/AS-SET password as-set-name")
	} else {

		err, RipeError, RipeErrorMess := getFromRipe(os.Args[1], os.Args[2], os.Args[3], os.Args[4])

		if err != nil {
			fmt.Printf("-----------------\n\n %s", err.Error())
		}

		if len(RipeError) > 0 {
			fmt.Println(RipeError)
		}
		if len(RipeError) > 0 {
			fmt.Println(RipeErrorMess)
		}

	}
}
