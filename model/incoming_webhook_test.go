// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIncomingWebhookJson(t *testing.T) {
	o := IncomingWebhook{ID: NewID()}
	json := o.ToJson()
	ro := IncomingWebhookFromJson(strings.NewReader(json))

	require.Equal(t, o.ID, ro.ID)
}

func TestIncomingWebhookIsValid(t *testing.T) {
	o := IncomingWebhook{}

	require.NotNil(t, o.IsValid())

	o.ID = NewID()
	require.NotNil(t, o.IsValid())

	o.CreateAt = GetMillis()
	require.NotNil(t, o.IsValid())

	o.UpdateAt = GetMillis()
	require.NotNil(t, o.IsValid())

	o.UserID = "123"
	require.NotNil(t, o.IsValid())

	o.UserID = NewID()
	require.NotNil(t, o.IsValid())

	o.ChannelID = "123"
	require.NotNil(t, o.IsValid())

	o.ChannelID = NewID()
	require.NotNil(t, o.IsValid())

	o.TeamID = "123"
	require.NotNil(t, o.IsValid())

	o.TeamID = NewID()
	require.Nil(t, o.IsValid())

	o.DisplayName = strings.Repeat("1", 65)
	require.NotNil(t, o.IsValid())

	o.DisplayName = strings.Repeat("1", 64)
	require.Nil(t, o.IsValid())

	o.Description = strings.Repeat("1", 501)
	require.NotNil(t, o.IsValid())

	o.Description = strings.Repeat("1", 500)
	require.Nil(t, o.IsValid())

	o.Username = strings.Repeat("1", 65)
	require.NotNil(t, o.IsValid())

	o.Username = strings.Repeat("1", 64)
	require.Nil(t, o.IsValid())

	o.IconURL = strings.Repeat("1", 1025)
	require.NotNil(t, o.IsValid())

	o.IconURL = strings.Repeat("1", 1024)
	require.Nil(t, o.IsValid())
}

func TestIncomingWebhookPreSave(t *testing.T) {
	o := IncomingWebhook{}
	o.PreSave()
}

func TestIncomingWebhookPreUpdate(t *testing.T) {
	o := IncomingWebhook{}
	o.PreUpdate()
}

func TestIncomingWebhookRequestFromJson(t *testing.T) {
	texts := []string{
		`this is a test`,
		`this is a test
			that contains a newline and tabs`,
		`this is a test \"foo
			that contains a newline and tabs`,
		`this is a test \"foo\"
			that contains a newline and tabs`,
		`this is a test \"foo\"
		\"			that contains a newline and tabs`,
		`this is a test \"foo\"

		\"			that contains a newline and tabs
		`,
	}

	for _, text := range texts {
		// build a sample payload with the text
		payload := `{
        "text": "` + text + `",
        "attachments": [
            {
                "fallback": "` + text + `",

                "color": "#36a64f",

                "pretext": "` + text + `",

                "author_name": "` + text + `",
                "author_link": "http://flickr.com/bobby/",
                "author_icon": "http://flickr.com/icons/bobby.jpg",

                "title": "` + text + `",
                "title_link": "https://api.slack.com/",

                "text": "` + text + `",

                "fields": [
                    {
                        "title": "` + text + `",
                        "value": "` + text + `",
                        "short": false
                    }
                ],

                "image_url": "http://my-website.com/path/to/image.jpg",
                "thumb_url": "http://example.com/path/to/thumb.png"
            }
        ]
    }`

		// try to create an IncomingWebhookRequest from the payload
		data := strings.NewReader(payload)
		iwr, _ := IncomingWebhookRequestFromJson(data)

		// After it has been decoded, the JSON string won't contain the escape char anymore
		expected := strings.Replace(text, `\"`, `"`, -1)
		require.NotNil(t, iwr)
		require.Equal(t, expected, iwr.Text)

		attachment := iwr.Attachments[0]
		require.Equal(t, expected, attachment.Text)
	}
}

func TestIncomingWebhookNullArrayItems(t *testing.T) {
	payload := `{"attachments":[{"fields":[{"title":"foo","value":"bar","short":true}, null]}, null]}`
	iwr, _ := IncomingWebhookRequestFromJson(strings.NewReader(payload))
	require.NotNil(t, iwr)
	require.Len(t, iwr.Attachments, 1)
	require.Len(t, iwr.Attachments[0].Fields, 1)
}
