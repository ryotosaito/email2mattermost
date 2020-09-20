package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/mattermost/mattermost-server/model"
)

func main() {
	mattermostURL := flag.String("mattermostURL", "", "Base URL of MatterMost")
	channelID := flag.String("channelID", "", "Channel ID of MatterMost")
	bearerToken := flag.String("bearerToken", "", "Secret Token of MatterMost")
	myAddress := flag.String("myAddress", "", "Your Address")

	flag.Parse()

	if *mattermostURL == "" || *channelID == "" || *bearerToken == "" {
		flag.PrintDefaults()
		log.Fatalln("Args mattermostURL, channelID and bearerToken are all required.")
	}

	email, err := parsemail.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	Client := model.NewAPIv4Client(*mattermostURL)
	Client.MockSession(*bearerToken)

	fileIds := make([]string, 0)
	for _, attachment := range email.Attachments {
		data, err := ioutil.ReadAll(attachment.Data)
		if err != nil {
			log.Fatal(err)
		}

		fileUploadResponse, _ := Client.UploadFile(data, *channelID, attachment.Filename)
		fileIds = append(fileIds, fileUploadResponse.FileInfos[0].Id)
	}

	var emailText string
	if strings.ToLower(email.Header.Get("Content-Transfer-Encoding")) == "base64" {
		bytes, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(email.TextBody)))
		if err != nil {
			log.Fatalln(err)
		}
		emailText = string(bytes)
	} else {
		emailText = email.TextBody
	}

	post := model.Post{
		ChannelId: *channelID,
		Message:   "New email received.",
		Props: model.StringInterface{
			"attachments": []map[string]string{{
				"author_name": fmt.Sprint(email.From[0]),
				"title":       email.Subject,
				"text":        emailText,
				"footer":      fmt.Sprintf("To: %s", *myAddress),
			}},
		},
		FileIds: fileIds,
	}
	Client.CreatePost(&post)
}
