package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/pbthorste/avtool"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/urfave/cli.v1"
)

var (
	version string
)

// OG: [create|decrypt|edit|encrypt|encrypt_string|rekey|view] [options] [vaultfile.yml]
// DONE: [create|decrypt|edit|encrypt|encrypt_string|view]
// TODO: [rekey]

func main() {
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
				vaultPassword := c.String("vault-password-file")
				// Validate CLI args
				err := validateCommandArgs(c)
				if err != nil {
					return err
				}
				vaultFileNames, err := validateAndGetVaultFile(c)
				if err != nil {
					return err
				}
				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Encrypt
				for i := 0; i < len(vaultFileNames); i++ {
					vaultFileName := vaultFileNames[i]
					result, err := avtool.EncryptFile(vaultFileName, pw)
					if err != nil {
						return cli.NewExitError(err, 2)
					}
					err = ioutil.WriteFile(vaultFileName, []byte(result), 0644)
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
				vaultPassword := c.String("vault-password-file")
				// Validate CLI args
				err := validateCommandArgs(c)
				if err != nil {
					return err
				}
				vaultFileNames, err := validateAndGetVaultFile(c)
				if err != nil {
					return err
				}

				// Get Vault password
				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Decrypt
				for i := 0; i < len(vaultFileNames); i++ {
					vaultFileName := vaultFileNames[i]
					result, err := avtool.DecryptFile(vaultFileName, pw)
					if err != nil {
						if strings.Compare(err.Error(), "ERROR: runtime error: index out of range") == 0 {
							return cli.NewExitError("input is not a vault encrypted "+vaultFileName+" is not a vault encrypted file for "+vaultFileName, 2)
						}
						return cli.NewExitError(err, 1)
					}

					// Create a temp file
					tempFile, err := ioutil.TempFile("", "vault")
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Write decrypted contents to temp file
					err = ioutil.WriteFile(tempFile.Name(), []byte(result), 0644)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Move temp file to old file
					err = os.Rename(tempFile.Name(), vaultFileName)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Close file
					err = tempFile.Close()
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
				vaultPassword := c.String("vault-password-file")
				// Validate CLI args
				err := validateCommandArgs(c)
				if err != nil {
					return err
				}
				vaultFileNames, err := validateAndGetVaultFile(c)
				if err != nil {
					return err
				}

				// Get Vault password
				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Decrypt
				for i := 0; i < len(vaultFileNames); i++ {
					vaultFileName := vaultFileNames[i]
					result, err := avtool.DecryptFile(vaultFileName, pw)
					if err != nil {
						if strings.Compare(err.Error(), "ERROR: runtime error: index out of range") == 0 {
							return cli.NewExitError("input is not a vault encrypted "+vaultFileName+" is not a vault encrypted file for "+vaultFileName, 2)
						}
						return cli.NewExitError(err, 1)
					}

					// Create a new temp file
					tempFile, err := ioutil.TempFile("", "vault")
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
					result, err = avtool.EncryptFile(tempFile.Name(), pw)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					err = ioutil.WriteFile(tempFile.Name(), []byte(result), 0644)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Move temp file to old file
					err = os.Rename(tempFile.Name(), vaultFileName)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Close file
					err = tempFile.Close()
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
				vaultPassword := c.String("vault-password-file")
				newVaultPassword := c.String("new-vault-password-file")
				// Validate CLI args
				err := validateCommandArgs(c)
				if err != nil {
					return err
				}
				vaultFileNames, err := validateAndGetVaultFile(c)
				if err != nil {
					return err
				}

				// Get Vault password
				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Get New Vault password
				newPw, err := retrieveVaultPassword(newVaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Decrypt
				for i := 0; i < len(vaultFileNames); i++ {
					vaultFileName := vaultFileNames[i]
					result, err := avtool.DecryptFile(vaultFileName, pw)
					if err != nil {
						if strings.Compare(err.Error(), "ERROR: runtime error: index out of range") == 0 {
							return cli.NewExitError("input is not a vault encrypted "+vaultFileName+" is not a vault encrypted file for "+vaultFileName, 2)
						}
						return cli.NewExitError(err, 1)
					}

					// Create a new temp file
					tempFile, err := ioutil.TempFile("", "vault")
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
					err = os.Rename(tempFile.Name(), vaultFileName)
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Close file
					err = tempFile.Close()
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				}

				println("Vault file edited")

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
				vaultPassword := c.String("vault-password-file")
				// Validate CLI args
				err := validateCommandArgs(c)
				if err != nil {
					return err
				}
				vaultFileName, err := validateAndGetVaultFileToCreate(c)
				if err != nil {
					return err
				}

				// Get Vault password
				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Create a new temp file
				tempFile, err := ioutil.TempFile("", "vault")
				if err != nil {
					return cli.NewExitError(err, 2)
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

				// Close temp file
				err = tempFile.Close()
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Delete the temp file
				err = os.Remove(tempFile.Name())
				if err != nil {
					return cli.NewExitError(err, 2)
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
				vaultPassword := c.String("vault-password-file")
				// Validate CLI args
				err := validateCommandArgs(c)
				if err != nil {
					return err
				}
				vaultFileNames, err := validateAndGetVaultFile(c)
				if err != nil {
					return err
				}

				// Get Vault password
				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Decrypt
				for i := 0; i < len(vaultFileNames); i++ {
					vaultFileName := vaultFileNames[i]
					result, err := avtool.DecryptFile(vaultFileName, pw)
					if err != nil {
						if strings.Compare(err.Error(), "ERROR: runtime error: index out of range") == 0 {
							return cli.NewExitError("input is not a vault encrypted "+vaultFileName+" is not a vault encrypted file for "+vaultFileName, 2)
						}
						return cli.NewExitError(err, 1)
					}

					// Create a new temp file
					tempFile, err := ioutil.TempFile("", "vault")
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
						command = "more"
					} else {
						command = "cat"
					}

					cmd := exec.Command(command, tempFile.Name())
					cmd.Stdout = os.Stdout
					cmd.Stdin = os.Stdin
					cmd.Stderr = os.Stderr
					cmd.Run()

					// Close temp file
					err = tempFile.Close()
					if err != nil {
						return cli.NewExitError(err, 1)
					}

					// Delete the temp file
					err = os.Remove(tempFile.Name())
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				}

				return nil
			},
		},
		{
			Name:      "encrypt_string",
			Usage:     "encrypt provided string",
			UsageText: "[options] string_to_encrypt",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "vault-password-file",
					Usage: "vault password file `VAULT_PASSWORD_FILE`",
				},
			},
			Action: func(c *cli.Context) error {
				vaultPassword := c.String("vault-password-file")
				// Validate CLI args
				strToEncrypt, err := validateAndGetStringToEncrypt(c)
				if err != nil {
					return err
				}
				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Encrypt
				result, err := avtool.Encrypt(strToEncrypt, pw)
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				println(result)

				println("Encryption successful")

				return nil
			},
		},
		{
			Name:      "av_encrypt_string",
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
				vaultPassword := c.String("vault-password-file")
				variableName := c.String("name")
				// Validate CLI args
				strToEncrypt, err := validateAndGetStringToEncrypt(c)
				if err != nil {
					return err
				}

				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Encrypt
				result, err := avtool.Encrypt(strToEncrypt, pw)
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				if variableName == "" {
					fmt.Println("!vault |")
				} else {
					fmt.Println(variableName + ": !vault |")

				}
				r := strings.Split(result, "\n")
				for _, stringLine := range r {
					fmt.Println("        " + stringLine)
				}

				println("Encryption successful")

				return nil
			},
		},
	}
	app.Run(os.Args)
}

func validateCommandArgs(c *cli.Context) (err error) {
	if !c.Args().Present() {
		cli.ShowSubcommandHelp(c)
		return cli.NewExitError(errors.New("ERROR: no vaild files provided"), 2)
	}
	return nil
}

func validateAndGetVaultFile(c *cli.Context) (files []string, err error) {
	var warnings []string
	if c.NArg() <= 0 {
		cli.ShowSubcommandHelp(c)
		return files, cli.NewExitError(errors.New("ERROR: no vaild files provided"), 2)
	}

	for i := 0; i < c.NArg(); i++ {
		filename := c.Args().Get(i)

		println(filename)

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
		cli.ShowSubcommandHelp(c)
		return files, cli.NewExitError(errors.New("ERROR: No supported files found"), 2)
	}

	return files, nil
}

func validateAndGetStringToEncrypt(c *cli.Context) (strToEncrypt string, err error) {
	if c.NArg() <= 0 {
		println("String to encrypt: ")
		byteInput, err2 := terminal.ReadPassword(int(syscall.Stdin))
		if err2 != nil {
			err2 = errors.New("ERROR: could not input string, " + err2.Error())
			return
		}
		strToEncrypt = string(byteInput)

		cli.ShowSubcommandHelp(c)
		return strToEncrypt, nil
	}

	strToEncrypt = strings.TrimSpace(c.Args().First())

	return strToEncrypt, nil
}

func validateAndGetVaultFileToCreate(c *cli.Context) (filename string, err error) {
	if c.NArg() > 1 {
		cli.ShowSubcommandHelp(c)
		return filename, cli.NewExitError(errors.New("ERROR: can only create one vault file at a time"), 2)
	}

	filename = strings.TrimSpace(c.Args().First())
	if filename == "" {
		cli.ShowSubcommandHelp(c)
		return filename, cli.NewExitError(errors.New("ERROR: filename not specified"), 2)
	} else {
		if fileInfo, err := os.Stat(filename); os.IsNotExist(err) {
			// File does not exist, good to go
			return filename, nil
		} else {
			if fileInfo.IsDir() {
				cli.ShowSubcommandHelp(c)
				return filename, cli.NewExitError(errors.New("ERROR: file "+filename+" is a directory"), 2)
			}
			return filename, cli.NewExitError(errors.New("ERROR: file "+filename+" already exists"), 2)
		}
	}
}

func retrieveVaultPassword(vaultPasswordFile string) (string, error) {
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

	return readVaultPassword()
}

func readVaultPassword() (password string, err error) {
	println("Vault password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		err = errors.New("ERROR: could not input password, " + err.Error())
		return
	}
	password = string(bytePassword)
	return
}

func getEditor() string {
	var editorEnv = os.Getenv("EDITOR")
	if editorEnv == "" {
		return "vim"
	}
	return editorEnv
}
