package main

import (
	"errors"
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
// DONE: [create|decrypt|edit|encrypt|view]
// TODO: [encrypt_string|rekey]

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
				vaultFileName, err := validateAndGetVaultFile(c)
				if err != nil {
					return err
				}
				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Encrypt
				result, err := avtool.EncryptFile(vaultFileName, pw)
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				err = ioutil.WriteFile(vaultFileName, []byte(result), 0644)
				if err != nil {
					return cli.NewExitError(err, 2)
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
				vaultFileName, err := validateAndGetVaultFile(c)
				if err != nil {
					return err
				}

				// Get Vault password
				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Decrypt
				result, err := avtool.DecryptFile(vaultFileName, pw)
				if err != nil {
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
				vaultFileName, err := validateAndGetVaultFile(c)
				if err != nil {
					return err
				}

				// Get Vault password
				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}

				// Decrypt
				result, err := avtool.DecryptFile(vaultFileName, pw)
				if err != nil {
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
				cmd := exec.Command("vim", tempFile.Name())
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

				println("Vault file edited")

				return nil
			},
		},
		{
			Name:      "create",
			Usage:     "create a new encypted file",
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
				cmd := exec.Command("vim", tempFile.Name())
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
				vaultFileName, err := validateAndGetVaultFile(c)
				if err != nil {
					return err
				}

				// Get Vault password
				pw, err := retrieveVaultPassword(vaultPassword)
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				result, err := avtool.DecryptFile(vaultFileName, pw)
				if err != nil {
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

				// Open 'more' stream of contents
				cmd := exec.Command("more", tempFile.Name())
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

				return nil
			},
		},
	}
	app.Run(os.Args)
}

func validateCommandArgs(c *cli.Context) (err error) {
	if !c.Args().Present() {
		cli.ShowSubcommandHelp(c)
		return cli.NewExitError(errors.New("ERROR: Empty or Invalid inputs! Please ref. to usage instructions!"), 2)
	}
	return nil
}

func validateAndGetVaultFile(c *cli.Context) (filename string, err error) {
	filename = strings.TrimSpace(c.Args().First())
	if filename == "" {
		cli.ShowSubcommandHelp(c)
		return filename, cli.NewExitError(errors.New("ERROR: Filename not specified!  Please ref. to usage instructions!"), 2)
	} else {
		if fileInfo, err := os.Stat(filename); os.IsNotExist(err) {
			cli.ShowSubcommandHelp(c)
			return filename, cli.NewExitError(errors.New("ERROR: file "+filename+" "+"doesn't exist!"), 2)
		} else {
			if fileInfo.IsDir() {
				cli.ShowSubcommandHelp(c)
				return filename, cli.NewExitError(errors.New("ERROR: file "+filename+" is a "+"directory!"), 2)
			}
		}
	}
	return filename, nil
}

func validateAndGetVaultFileToCreate(c *cli.Context) (filename string, err error) {
	filename = strings.TrimSpace(c.Args().First())
	if filename == "" {
		cli.ShowSubcommandHelp(c)
		return filename, cli.NewExitError(errors.New("ERROR: Filename not specified!  Please ref. to usage instructions!"), 2)
	} else {
		if fileInfo, err := os.Stat(filename); os.IsNotExist(err) {
			// File does not exist, good to go
			return filename, nil
		} else {
			if fileInfo.IsDir() {
				cli.ShowSubcommandHelp(c)
				return filename, cli.NewExitError(errors.New("ERROR: file "+filename+" is a directory!"), 2)
			}
			return filename, cli.NewExitError(errors.New("ERROR: file "+filename+" already exists!"), 2)
		}
	}
	// return filename, error on error; nil if no error;
	return filename, nil
}

func retrieveVaultPassword(vaultPasswordFile string) (string, error) {
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
