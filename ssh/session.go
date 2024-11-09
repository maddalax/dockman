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

type Session struct {
	client *ssh.Client
	s      *ssh.Session
}

func (s *Session) Disconnect() {
	s.s.Close()
	s.client.Close()
}

func (s *Session) RunWithOutput(command string) (string, error) {
	output, err := s.s.CombinedOutput(command)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (s *Session) RunWithOutputStream(command string, stdOut, stdErr io.Writer) error {
	// Set up pipes to capture the output
	stdout, err := s.s.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := s.s.StderrPipe()
	if err != nil {
		return err
	}

	if err := s.s.Start(command); err != nil {
		return err
	}

	// Stream stdout and stderr
	go streamOutput(stdout, stdOut)
	go streamOutput(stderr, stdErr)

	// Wait for the command to complete
	if err := s.s.Wait(); err != nil {
		return err
	}
	return err
}

func OpenSession(opts SessionOpts) (*Session, error) {

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

	session, err := client.NewSession()

	if err != nil {
		return nil, ErrorCreatingSession{Err: err}
	}

	return &Session{
		client: client,
		s:      session,
	}, nil
}
