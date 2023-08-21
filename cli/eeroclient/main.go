package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ab623/goeero/eero"
	"github.com/urfave/cli/v2"
)

const sessionFilename = "eero_session.txt"

func main() {
	app := &cli.App{
		Usage:  "Unofficial CLI application to interact with Eero devices.",
		Before: setupApp,
		// HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "session-file",
				Aliases: []string{"s"},
				Value:   sessionFilename,
				Usage:   "Load/Save session from `FILE`.",
				EnvVars: []string{"SESSION_FILE"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "auth",
				Usage:  "Authenticate with Eero platform",
				Action: authenticate,
			},
			{
				Name:   "account",
				Usage:  "Show account information",
				Action: account,
			},
			{
				Name:   "devices",
				Usage:  "Get a list of devices",
				Action: devices,
			},
			{
				Name:   "networks",
				Usage:  "Get a list of networks",
				Action: networks,
			},
			{
				Name:   "session",
				Usage:  "Print session token if exits",
				Action: token,
			},
		},
		// Action: func(cCtx *cli.Context) error {
		// 	return nil
		// },
	}

	app.Run(os.Args)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

func setupApp(cCtx *cli.Context) error {
	// Return nil if no args were passed
	if cCtx.NArg() == 0 {
		return nil
	}

	// 1. Get the value from the given flag - defaults to ""
	var flagValue string = cCtx.Value("session-file").(string)
	// 2. Resolve the session file path to the current directory or the given directory
	sessionFilePath := resolveConfigFilePath(flagValue)
	// 3. Save the path in the metadata as the application may need it for subsequent writing
	cCtx.App.Metadata["sessionFile"] = sessionFilePath

	// If the application is authenticating, don't need the file data
	// Otherwise we should load the value of the file into the metadata to ensure it can be used in API calls
	// 1. Setup a base token
	var userToken string = ""
	// 2. Check if this is not auth, then load the token from the file
	if cCtx.Args().First() != "auth" {
		t, err := loadToken(sessionFilePath)
		if err != nil {
			return err
		}
		userToken = t
	}
	// 3. Store token in the metadata
	cCtx.App.Metadata["userToken"] = userToken

	// Setup the Client here
	client := eero.New(userToken)
	cCtx.App.Metadata["client"] = client

	return nil
}

func resolveConfigFilePath(sessionFilePath string) string {
	// If set then store filepath in var
	if sessionFilePath == sessionFilename {
		// Get the executable path
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		// Get the executable directory
		exPath := filepath.Dir(ex)

		// Create the config filepath
		configFilepath := fmt.Sprintf("%s/%s", exPath, sessionFilename)

		sessionFilePath = configFilepath
	}

	return sessionFilePath
}

func loadToken(sessionFilePath string) (string, error) {
	// Check for a file
	if _, err := os.Stat(sessionFilePath); err != nil {
		fmt.Printf("WARNING:   %s does not exit. Cannot load eero session token.\n", sessionFilePath)
		fmt.Printf("WARNING:   Use ./eeroclient auth to generate new file.\n\n")

		return "", err
	}

	// Read file
	content, err := os.ReadFile(sessionFilePath)
	if err != nil {
		return "", err
	}

	// Set the user token
	return string(content), nil
}

func authenticate(cCtx *cli.Context) error {
	var client eero.Eero = cCtx.App.Metadata["client"].(eero.Eero)

	var userId string
	var verificationCode string

	// 1. Request user mobile number
	fmt.Println("Step 1: Enter Eero login ID (phone or email address)")
	fmt.Scanln(&userId)
	userId = strings.ToLower(strings.TrimSpace(userId))

	// 2. Generate OTP
	loginResponse, err := client.Login(userId)
	if err != nil {
		log.Fatal("Could not send OTP: ", err)
	}
	sessionToken := loginResponse.UserToken

	// 3. Send the verification
	fmt.Println("Step 2: Enter OTP from Email or SMS")
	fmt.Scanln(&verificationCode)
	verificationCode = strings.TrimSpace(verificationCode)

	// 4.Check verification was successful
	err = client.LoginVerify(sessionToken, verificationCode)
	if err != nil {
		log.Fatal("Could not verify: ", err)
	}

	// 5. If successful then save
	filePath := cCtx.Value("session-file").(string)
	file, err := os.Create(filePath)

	if err != nil {
		log.Fatal("Could not create file: ", err)
	}
	defer file.Close()

	file.WriteString(sessionToken)

	// 5. Confirmation
	fmt.Println("Authentication successful.")

	return nil
}

func devices(cCtx *cli.Context) error {
	var client eero.Eero = cCtx.App.Metadata["client"].(eero.Eero)

	resp, err := client.Devices()
	if err != nil {
		log.Fatal("Could not get device data: ", err)
	}

	fmt.Println(prettyPrint(resp))
	return nil
}

func account(cCtx *cli.Context) error {
	var client eero.Eero = cCtx.App.Metadata["client"].(eero.Eero)

	resp, err := client.Account()
	if err != nil {
		log.Fatal("Could not get account data: ", err)
	}

	fmt.Println(prettyPrint(resp))
	return nil
}

func networks(cCtx *cli.Context) error {
	var client eero.Eero = cCtx.App.Metadata["client"].(eero.Eero)

	resp, err := client.Networks()
	if err != nil {
		log.Fatal("Could not get network data: ", err)
	}

	fmt.Println(prettyPrint(resp))
	return nil
}

func token(cCtx *cli.Context) error {
	var sessionToken string = cCtx.App.Metadata["userToken"].(string)

	fmt.Printf("User token: %s\n", sessionToken)
	return nil
}
