package eternus

import (
	"golang.org/x/crypto/ssh"
	"io"
	"time"
)

type Executor interface {
	Execute(command string) (string, error)
}

type sshExecutor struct {
	session *ssh.Session
	client  *ssh.Client
	in      io.Writer
	out     io.Reader
	wait    time.Duration
}

func NewExecutor(host, user, password, cipher string, wait int64) (*sshExecutor, error) {
	sshClintConfig := ssh.ClientConfig{}
	sshClintConfig.SetDefaults()
	sshClintConfig.User = user
	sshClintConfig.Auth = []ssh.AuthMethod{ssh.Password(password)}
	sshClintConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	if cipher != "" {
		sshClintConfig.Ciphers = append(sshClintConfig.Ciphers, cipher)
	}

	client, err := ssh.Dial("tcp", host, &sshClintConfig)
	if err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, err
	}
	out, err := session.StdoutPipe()
	if err != nil {
		client.Close()
		return nil, err
	}
	in, err := session.StdinPipe()
	if err != nil {
		client.Close()
		return nil, err
	}
	err = session.Shell()
	if err != nil {
		session.Close()
		client.Close()
		return nil, err
	}
	executor := sshExecutor{
		wait:    time.Duration(wait) * time.Millisecond,
		client:  client,
		session: session,
		in:      in,
		out:     out,
	}
	return &executor, nil
}

func (e *sshExecutor) Disconnect() {
	e.session.Close()
	e.client.Close()
}

func (e *sshExecutor) Execute(command string) (string, error) {
	buf := make([]byte, 100000)
	_, err := e.in.Write([]byte(command + "\n"))
	if err != nil {
		return "", err
	}
	time.Sleep(e.wait)
	n, err := e.out.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}
