package main

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/clok/kemba"
	"github.com/pbthorste/avtool"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	version string
	k       = kemba.New("gwvault")
	kenc    = k.Extend("encrypt")
	kdec    = k.Extend("decrypt")
	kedit   = k.Extend("edit")
	ke      = k.Extend("editor")
	kencf   = k.Extend("encryptFile")
	kdecf   = k.Extend("decryptFile")
	ktf     = k.Extend("createTempFile")
	kclean  = k.Extend("cleanupFile")
)

// OG: [create|decrypt|edit|encrypt|encrypt_string|rekey|view] [options] [vaultfile.yml]
// DONE: [create|decrypt|edit|encrypt|encrypt_string|rekey|view]

func main() {
	k.Println("executing")
	app := cli.NewApp()
	app.Name = "gwvault"
	app.Version = version
	app.Usage = "encryption/decryption utility for Ansible data files"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "vault-password-file",
			Usage: "vault password file `VAULT_PASSWORD_FILE`",
		},
		cli.StringFlag{
			Name:  "new-vault-password-file",
			Usage: "new vault password file for rekey `NEW_VAULT_PASSWORD_FILE`",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "encrypt",
			Usage:     "encrypt file",
			UsageText: "[options] [vaultfile.yml]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "vault-password-file",
					Usage: "vault password file `VAULT_PASSWORD_FILE`",
				},
			},
			Action: func(c *cli.Context) error {
				vaultFileNames, err := validateAndGetVaultFile(c)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Get Vault Password
				vaultPassword := c.String("vault-password-file")
				var pw string
				pw, err = retrieveVaultPassword(vaultPassword, "Vault password:")
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Encrypt
				kenc.Printf("processing %d files. %# v", len(vaultFileNames), vaultFileNames)
				for _, file := range vaultFileNames {
					kenc.Printf("processing: %s", file)
					err := encryptFile(file, pw)
					if err != nil {
						return cli.NewExitError(err, 2)
					}
				}

				println("Encryption successful")

				return nil
			},
		},
		{
			Name:      "decrypt",
			Usage:     "decrypt file",
			UsageText: "[options] [vaultfile.yml]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "vault-password-file",
					Usage: "vault password file `VAULT_PASSWORD_FILE`",
				},
			},
			Action: func(c *cli.Context) error {
				vaultFileNames, err := validateAndGetVaultFile(c)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Get Vault Password
				vaultPassword := c.String("vault-password-file")
				var pw string
				pw, err = retrieveVaultPassword(vaultPassword, "Vault password:")
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Decrypt
				kdec.Printf("processing %d files. %# v", len(vaultFileNames), vaultFileNames)
				for _, file := range vaultFileNames {
					kdec.Printf("processing: %s", file)
					result, err := decryptFile(file, pw)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Create a new temp file
					var tempFile *os.File
					tempFile, err = createTempFile()
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Write decrypted contents to temp file
					err = ioutil.WriteFile(tempFile.Name(), []byte(result), 0644)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Move temp file to old file
					kdec.Printf("overwriting contents %s -> %s", tempFile.Name(), file)
					err = os.Rename(tempFile.Name(), file)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					err = cleanupFile(tempFile)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				}

				println("Decryption successful")

				return nil
			},
		},
		{
			Name:      "edit",
			Usage:     "edit file and re-encrypt",
			UsageText: "[options] [vaultfile.yml]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "vault-password-file",
					Usage: "vault password file `VAULT_PASSWORD_FILE`",
				},
			},
			Action: func(c *cli.Context) error {
				kedit.Println("start")
				vaultFileNames, err := validateAndGetVaultFile(c)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Get Vault Password
				vaultPassword := c.String("vault-password-file")
				var pw string
				pw, err = retrieveVaultPassword(vaultPassword, "Vault password:")
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				kedit.Printf("processing %d files. %# v", len(vaultFileNames), vaultFileNames)
				for _, file := range vaultFileNames {
					kedit.Printf("processing: %s", file)
					result, err := decryptFile(file, pw)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Create a new temp file
					var tempFile *os.File
					tempFile, err = createTempFile()
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Write decrypted contents to temp file
					err = ioutil.WriteFile(tempFile.Name(), []byte(result), 0644)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Open editor for modifications
					cmd := exec.Command(getEditor(), tempFile.Name())
					cmd.Stdout = os.Stdout
					cmd.Stdin = os.Stdin
					cmd.Stderr = os.Stderr
					err = cmd.Run()
					if err != nil {
						return cli.NewExitError(err, 2)
					}

					// Encrypt temp file
					err = encryptFile(tempFile.Name(), pw)
					if err != nil {
						return cli.NewExitError(err, 2)
					}

					// Move temp file to old file
					kedit.Printf("overwriting contents %s -> %s", tempFile.Name(), file)
					err = os.Rename(tempFile.Name(), file)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					err = cleanupFile(tempFile)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				}

				println("Vault file edited")

				return nil
			},
		},
		{
			Name:      "rekey",
			Usage:     "alter encryption password and re-encrypt",
			UsageText: "[options] [vaultfile.yml] [newvaultfile.yml]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "vault-password-file",
					Usage: "vault password file `VAULT_PASSWORD_FILE`",
				},
				cli.StringFlag{
					Name:  "new-vault-password-file",
					Usage: "new vault password file for rekey `NEW_VAULT_PASSWORD_FILE`",
				},
			},
			Action: func(c *cli.Context) error {
				vaultFileNames, err := validateAndGetVaultFile(c)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Get Vault Password
				vaultPassword := c.String("vault-password-file")
				var pw string
				pw, err = retrieveVaultPassword(vaultPassword, "Vault password:")
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Get New Vault password
				newVaultPassword := c.String("new-vault-password-file")
				var newPw string
				newPw, err = retrieveVaultPassword(newVaultPassword, "New Vault password:")
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				if newVaultPassword == "" {
					var confirmPw string
					confirmPw, err = retrieveVaultPassword("", "Confirm New Vault password:")
					if err != nil {
						return cli.NewExitError(err, 2)
					}

					if newPw != confirmPw {
						return cli.NewExitError(errors.New("ERROR! Passwords do not match"), 2)
					}
				}

				// Decrypt
				for _, file := range vaultFileNames {
					result, err := decryptFile(file, pw)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Create a new temp file
					var tempFile *os.File
					tempFile, err = createTempFile()
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Write decrypted contents to temp file
					err = ioutil.WriteFile(tempFile.Name(), []byte(result), 0644)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Encrypt temp file with new pw
					result, err = avtool.EncryptFile(tempFile.Name(), newPw)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					err = ioutil.WriteFile(tempFile.Name(), []byte(result), 0644)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Move temp file to old file
					err = os.Rename(tempFile.Name(), file)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					err = cleanupFile(tempFile)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				}

				println("Rekey successful")

				return nil
			},
		},
		{
			Name:      "create",
			Usage:     "create a new encrypted file",
			UsageText: "[options] [vaultfile.yml]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "vault-password-file",
					Usage: "vault password file `VAULT_PASSWORD_FILE`",
				},
			},
			Action: func(c *cli.Context) error {
				vaultFileName, err := validateAndGetVaultFileToCreate(c)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Get Vault Password
				vaultPassword := c.String("vault-password-file")
				var pw string
				pw, err = retrieveVaultPassword(vaultPassword, "Vault password:")
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Create a new temp file
				var tempFile *os.File
				tempFile, err = createTempFile()
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				// Open temp file for edit
				cmd := exec.Command(getEditor(), tempFile.Name())
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Encrypt temp file
				result, err := avtool.EncryptFile(tempFile.Name(), pw)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Write encrypted content to new file location
				err = ioutil.WriteFile(vaultFileName, []byte(result), 0644)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				err = cleanupFile(tempFile)
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				println("Vault file created")

				return nil
			},
		},
		{
			Name:      "view",
			Usage:     "view contents of encrypted file",
			UsageText: "[options] [vaultfile.yml]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "vault-password-file",
					Usage: "vault password file `VAULT_PASSWORD_FILE`",
				},
			},
			Action: func(c *cli.Context) error {
				vaultFileNames, err := validateAndGetVaultFile(c)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Get Vault Password
				vaultPassword := c.String("vault-password-file")
				var pw string
				pw, err = retrieveVaultPassword(vaultPassword, "Vault password:")
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Decrypt
				for _, file := range vaultFileNames {
					result, err := decryptFile(file, pw)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Create a new temp file
					var tempFile *os.File
					tempFile, err = createTempFile()
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					err = ioutil.WriteFile(tempFile.Name(), []byte(result), 0644)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					var command string

					// Check for TTY
					if terminal.IsTerminal(int(os.Stdin.Fd())) { // We have TTY!
						command = "less" // pick less to allow for search and other niceties
					} else {
						command = "cat"
					}

					cmd := exec.Command(command, tempFile.Name())
					cmd.Stdout = os.Stdout
					cmd.Stdin = os.Stdin
					cmd.Stderr = os.Stderr
					err = cmd.Run()
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					err = cleanupFile(tempFile)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				}

				return nil
			},
		},
		{
			Name:      "encrypt_string",
			Aliases:   []string{"av_encrypt_string"},
			Usage:     "encrypt provided string, output in ansible-vault format",
			UsageText: "[options] string_to_encrypt",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "vault-password-file",
					Usage: "vault password file `VAULT_PASSWORD_FILE`",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "variable name to encrypt",
				},
			},
			Action: func(c *cli.Context) error {
				var pw string
				var err error
				vaultPassword := c.String("vault-password-file")
				pw, err = retrieveVaultPassword(vaultPassword, "Vault password:")
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				var strToEncrypt string
				strToEncrypt, err = validateAndGetStringToEncrypt(c)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Encrypt
				var result string
				result, err = avtool.Encrypt(strToEncrypt, pw)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				variableName := c.String("name")
				if variableName == "" {
					fmt.Println("!vault |")
				} else {
					fmt.Println(variableName + ": !vault |")
				}

				r := strings.Split(result, "\n")
				for _, stringLine := range r {
					fmt.Println("          " + stringLine)
				}

				println("Encryption successful")

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func createTempFile() (*os.File, error) {
	t, err := ioutil.TempFile("", "vault")
	if err != nil {
		return nil, err
	}
	ktf.Printf("created temp file: %s", t.Name())
	return t, nil
}

func cleanupFile(t *os.File) error {
	kclean.Printf("deleting file: %s", t.Name())
	// Close temp file
	err := t.Close()
	if err != nil {
		return err
	}

	_, err = os.Stat(t.Name())
	if os.IsNotExist(err) {
		kclean.Printf("skipping - file no longer exists: %s", t.Name())
		return nil
	}

	// Delete the temp file
	err = os.Remove(t.Name())
	if err != nil {
		return err
	}
	kclean.Printf("delete complete: %s", t.Name())
	return nil
}

func decryptFile(file string, pw string) (string, error) {
	kdecf.Printf("attempting decryption: %s", file)
	result, err := avtool.DecryptFile(file, pw)
	if err != nil {
		if strings.Compare(err.Error(), "ERROR: runtime error: index out of range") == 0 {
			return "", fmt.Errorf("input is not a vault encrypted %s is not a vault encrypted file for %s", file, file)
		}
		return "", err
	}
	kdecf.Printf("decryption successful: %s", file)
	return result, nil
}

func encryptFile(file string, pw string) error {
	kencf.Printf("attempting encryption: %s", file)
	result, err := avtool.EncryptFile(file, pw)
	if err != nil {
		return err
	}
	kencf.Printf("encryption successful: %s", file)

	kencf.Printf("writing out encrypted contents: %s", file)
	err = ioutil.WriteFile(file, []byte(result), 0644)
	if err != nil {
		return err
	}
	return nil
}

func validateCommandArgs(c *cli.Context) error {
	if !c.Args().Present() {
		_ = cli.ShowSubcommandHelp(c)
		return errors.New("ERROR: no valid files provided")
	}
	return nil
}

func validateAndGetVaultFile(c *cli.Context) (files []string, err error) {
	// Validate CLI args
	err = validateCommandArgs(c)
	if err != nil {
		return nil, err
	}

	var warnings []string
	if c.NArg() <= 0 {
		_ = cli.ShowSubcommandHelp(c)
		return files, errors.New("ERROR: no valid files provided")
	}

	for i := 0; i < c.NArg(); i++ {
		filename := c.Args().Get(i)

		if fileInfo, err := os.Stat(filename); os.IsNotExist(err) {
			warnings = append(warnings, "WARN: skipping file "+filename+" because it does not exist")
			continue
		} else {
			if fileInfo.IsDir() {
				warnings = append(warnings, "WARN: skipping file "+filename+" because it is a directory")
				continue
			}
		}

		files = append(files, filename)
	}

	if len(warnings) > 0 {
		for i := 0; i < len(warnings); i++ {
			println(warnings[i])
		}
	}

	if len(files) <= 0 {
		_ = cli.ShowSubcommandHelp(c)
		return files, errors.New("ERROR: No supported files found")
	}

	return files, nil
}

func validateAndGetStringToEncrypt(c *cli.Context) (strToEncrypt string, err error) {
	if c.NArg() <= 0 {
		prompt := &survey.Editor{
			Message: "Open editor to input string to encrypt",
		}

		err = survey.AskOne(prompt, &strToEncrypt)
		if err != nil {
			_ = cli.ShowSubcommandHelp(c)
			return "", err
		}

		return strToEncrypt, nil
	}

	strToEncrypt = strings.TrimSpace(c.Args().First())

	return strToEncrypt, nil
}

func validateAndGetVaultFileToCreate(c *cli.Context) (filename string, err error) {
	// Validate CLI args
	err = validateCommandArgs(c)
	if err != nil {
		return "", err
	}

	if c.NArg() > 1 {
		_ = cli.ShowSubcommandHelp(c)
		return filename, errors.New("ERROR: can only create one vault file at a time")
	}

	filename = strings.TrimSpace(c.Args().First())
	if filename == "" {
		_ = cli.ShowSubcommandHelp(c)
		return filename, errors.New("ERROR: filename not specified")
	} else {
		if fileInfo, err := os.Stat(filename); os.IsNotExist(err) {
			// File does not exist, good to go
			return filename, nil
		} else {
			if fileInfo.IsDir() {
				_ = cli.ShowSubcommandHelp(c)
				return filename, errors.New("ERROR: file " + filename + " is a directory")
			}
			return filename, errors.New("ERROR: file " + filename + " already exists")
		}
	}
}

func retrieveVaultPassword(vaultPasswordFile string, msg string) (string, error) {
	if vaultPasswordFile == "" {
		// Not specified via CLI, check ANSIBLE_VAULT_PASSWORD_FILE environment variable
		if os.Getenv("ANSIBLE_VAULT_PASSWORD_FILE") != "" {
			vaultPasswordFile = os.Getenv("ANSIBLE_VAULT_PASSWORD_FILE")
		}
	}
	if vaultPasswordFile != "" {
		if _, err := os.Stat(vaultPasswordFile); os.IsNotExist(err) {
			return "", errors.New("ERROR: vault-password-file, could not find: " + vaultPasswordFile)
		}
		pw, err := ioutil.ReadFile(vaultPasswordFile)
		if err != nil {
			return "", errors.New("ERROR: vault-password-file, " + err.Error())
		}
		return strings.TrimSpace(string(pw)), nil
	}

	return readVaultPassword(msg)
}

func readVaultPassword(msg string) (password string, err error) {
	prompt := &survey.Password{
		Message: msg,
	}
	err = survey.AskOne(prompt, &password)
	return password, nil
}

func getEditor() string {
	var editorEnv = os.Getenv("EDITOR")
	ke.Printf("found editor: %s", editorEnv)
	if editorEnv == "" {
		ke.Printf("using default of vim")
		return "vim"
	}
	return editorEnv
}
