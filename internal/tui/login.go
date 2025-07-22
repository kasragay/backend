package tui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kasragay/backend/internal/ports"
	"github.com/rivo/tview"
)

func WelcomePage() {
	app := tview.NewApplication()
	pages := tview.NewPages()
	var password string
	var username string

	// --- Username Page ---
	usernameForm := tview.NewForm().
		AddInputField("Username", "", 20, nil, func(text string) {
			username = text
		}).
		AddButton("Next", func() {
			userAuth := ports.AuthCheckPostRequest{
				Username: username,
				UserType: ports.ClientUserType,
			}
			result, err := checkAuth(userAuth)
			if err != nil {
				modal := tview.NewModal().
					SetText(fmt.Sprintf("Error: %v", err)).
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						pages.RemovePage("error")
					})
				pages.AddPage("error", modal, true, true)
				return
			}

			if result.Deleted {
				pages.SwitchToPage("sicktir")
			} else if result.Exists && result.HasPassword {
				pages.SwitchToPage("password")
			} else if result.Exists == false {
				pages.SwitchToPage("signup")
			}
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	usernameForm.SetBorder(true).SetTitle(" Welcome to kasragay platform ").SetTitleAlign(tview.AlignCenter)

	// --- Password Page ---
	passwordForm := tview.NewForm().
		AddPasswordField("Password", "", 20, '*', func(text string) {
			password = text
		}).
		AddButton("Login", func() {
			pages.SwitchToPage("confirmation")
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("username")
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	passwordForm.SetBorder(true).SetTitle("Login Page - Enter Password ").SetTitleAlign(tview.AlignCenter)

	// --- Singup Page ---
	signupForm := tview.NewForm().
		AddInputField("Username", "", 20, nil, func(text string) {
			username = text
		}).
		AddPasswordField("passwrord", "", 20, '*', func(text string) {
			password = text
		}).
		AddButton("Login", func() {
			// #TODO: it has to be the Signup Funtion
			pages.SwitchToPage("confirmation")
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	signupForm.SetBorder(true).SetTitle(" Signup Page ").SetTitleAlign(tview.AlignCenter)

	// --- Confirmation Page ---
	confirmationText := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	confirmationText.SetText(fmt.Sprintf("[green]Welcome, %s!\n\nLogin successful.", username))

	confirmationPage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(confirmationText, 0, 1, false).
		AddItem(tview.NewButton("Exit").SetSelectedFunc(func() {
			app.Stop()
		}), 3, 1, true)

	// --- "Sicktir" Page ---
	sicktirText := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	sicktirText.SetText(fmt.Sprintf("[red]Go away %s!\n\nAccess denied. You've been login before.", username))

	sicktirPage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(sicktirText, 0, 1, false).
		AddItem(tview.NewButton("Exit").SetSelectedFunc(func() {
			app.Stop()
		}), 3, 1, true)

	// --- Pages ---
	pages.AddPage("username", usernameForm, true, true)
	pages.AddPage("password", passwordForm, true, false)
	pages.AddPage("signup", signupForm, true, false)
	pages.AddPage("confirmation", confirmationPage, true, false)
	pages.AddPage("sicktir", sicktirPage, true, false)

	// Run the app
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
	println(password)
}

func checkAuth(auth ports.AuthCheckPostRequest) (ports.AuthCheckPostResponse, error) {
	baseURL := "https://api.kasragay.com/v1/auth/check"

	params := url.Values{}
	params.Set("username", auth.Username)
	params.Set("user_type", string(auth.UserType))

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return ports.AuthCheckPostResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ports.AuthCheckPostResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ports.AuthCheckPostResponse{}, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var result ports.AuthCheckPostResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ports.AuthCheckPostResponse{}, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return result, nil
}
