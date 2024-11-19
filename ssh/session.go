package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
)

type SessionOpts struct {
	Username      string
	Password      string
	ServerAddress string
}

type PasswordRequiredError struct {
}

type UnknownConnectionError struct {
	Err error
}

type ErrorCreatingSession struct {
	Err error
}

func (e ErrorCreatingSession) Error() string {
	return fmt.Sprintf("Error: Error creating session: %v", e.Err)
}

func (e UnknownConnectionError) Error() string {
	return fmt.Sprintf("Error: Unknown connection error: %v", e.Err)
}

func (e PasswordRequiredError) Error() string {
	return "Error: Password authentication required."
}

type Client struct {
	client *ssh.Client
}

func (s *Client) Disconnect() {
	s.client.Close()
}

func (s *Client) RunWithOutput(command string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", ErrorCreatingSession{Err: err}
	}
	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (s *Client) RunWithOutputStream(command string, stdOut, stdErr io.Writer) error {
	session, err := s.client.NewSession()

	if err != nil {
		return ErrorCreatingSession{Err: err}
	}

	outPipe, err := session.StdoutPipe()
	if err == nil {
		go streamOutput(outPipe, stdOut)
	}
	errPipe, err := session.StderrPipe()
	if err == nil {
		go streamOutput(errPipe, stdErr)
	}

	if err := session.Run(command); err != nil {
		return err
	}

	return err
}

func OpenClient(opts SessionOpts) (*Client, error) {

	if !strings.Contains(opts.ServerAddress, ":") {
		opts.ServerAddress = fmt.Sprintf("%s:22", opts.ServerAddress)
	}

	// Configure the SSH client
	config := &ssh.ClientConfig{
		User:            opts.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if opts.Password != "" {
		config.Auth = []ssh.AuthMethod{
			ssh.Password(opts.Password),
		}
	}

	// Attempt to establish an SSH connection
	client, err := ssh.Dial("tcp", opts.ServerAddress, config)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ssh: handshake failed: ssh: unable to authenticate") {
			return nil, PasswordRequiredError{}
		} else {
			return nil, UnknownConnectionError{Err: err}
		}
	}

	return &Client{
		client: client,
	}, nil
}
