package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	var selected_method string = "GET"

	//api result
	result := tview.NewTextArea().
		SetPlaceholder("Result of the API will show here")
	result.SetBorder(true)
	result.SetTitle("Response result")

	var form *tview.Form

	form = tview.NewForm().
		AddDropDown("Method", []string{"GET", "POST", "PUT"}, 0,
			func(option string, option_index int) {
				selected_method = option
				if option != "GET" {
					result.SetText("Different methods than 'GET' will not work atm.", false)
				}
			}).
		AddInputField("URL", "", 40, nil, nil).
		AddButton("Send", func() {
			url_field := form.GetFormItemByLabel("URL").(*tview.InputField)
			url := url_field.GetText()

			if url == "" {
				result.SetText("Error: URL is required C:", false)
				return
			}

			//client http
			client := &http.Client{Timeout: 30 * time.Second}
			var req *http.Request
			var err error

			if selected_method == "GET" {
				req, err = http.NewRequest(selected_method, url, nil)
			} else {
				// POST / PUT todo: add body
				req, err = http.NewRequest(selected_method, url, nil)
			}

			if err != nil {
				result.SetText(fmt.Sprintf("Error creating request: %v", err), false)
				return
			}

			req.Header.Set("User-Agent", "RestMe/1.0")
			req.Header.Set("Accept", "*/*")

			// Show request info
			request_info := fmt.Sprintf("REQUEST:\nMethod: %s\nURL: %s\nHeaders:\n", selected_method, url)
			for key, values := range req.Header {
				request_info += fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", "))
			}
			result.SetText(request_info+"\nSending request...\n", false)

			app.ForceDraw() // Force UI update

			// Send the request
			resp, err := client.Do(req)

			var response_text string
			if err != nil {
				response_text = fmt.Sprintf("ERROR:\n%v", err)
			} else {
				defer resp.Body.Close()

				// Read response body
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					response_text = fmt.Sprintf("ERROR reading response: %v", err)
				} else {
					// Format response
					response_text = fmt.Sprintf(
						"RESPONSE:\nStatus: %s %d\nHeaders:\n",
						resp.Status,
						resp.StatusCode,
					)

					for key, values := range resp.Header {
						response_text += fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", "))
					}

					response_text += fmt.Sprintf("\nBody:\n%s", string(body))
				}
			}

			// Update result
			current_text := result.GetText()
			result.SetText(current_text+"\n"+response_text, false)

		}).
		AddButton("Clear", func() {
			result.SetText("", false)
		}).
		SetButtonsAlign(tview.AlignLeft)

	form.SetHorizontal(true)
	form.SetBackgroundColor(tcell.ColorBlack)
	form.SetBorder(false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true).
		AddItem(result, 0, 1, true)

	flex.SetBorder(true).
		SetTitle("RestMe, API Rest Client").
		SetTitleAlign(tview.AlignCenter)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
