package app

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TerminalSession struct {
	stdin io.WriteCloser
	done        chan struct{}
}

type DockerContainersTerminal struct {
	DockerBaseService
	lock sync.RWMutex
	session *TerminalSession
}

func NewDockerTerminalService() *DockerContainersTerminal {
	return &DockerContainersTerminal{}
}

func StartupDockerTerminalService(s *DockerContainersTerminal, ctx context.Context, cli *client.Client) {
	s.ctx = ctx
	s.cli = cli
}

func (s *DockerContainersTerminal) StartInteractiveTerminal(id string) error {
	if s.cli == nil || s.ctx == nil {
		return fmt.Errorf("Docker client not initialized")
	}

	shell := "sh"
	checkCmd := []string{"which", "bash"}
	execCheck, err := s.cli.ContainerExecCreate(s.ctx, id, container.ExecOptions{
		Cmd:          checkCmd,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err == nil {
		checkResp, err := s.cli.ContainerExecAttach(s.ctx, execCheck.ID, container.ExecAttachOptions{})
		if err == nil {
			defer checkResp.Close()
			output, _ := io.ReadAll(checkResp.Reader)
			if strings.Contains(string(output), "bash") {
				shell = "bash"
			}
		}
	}

	execConfig := container.ExecOptions{
		Cmd:          []string{shell},
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	execResp, err := s.cli.ContainerExecCreate(s.ctx, id, execConfig)
	if err != nil {
		return err
	}

	resp, err := s.cli.ContainerExecAttach(s.ctx, execResp.ID, container.ExecAttachOptions{Tty: true})
	if err != nil {
		return err
	}

	s.lock.Lock()
	s.session = &TerminalSession{
		stdin: resp.Conn,
		done:        make(chan struct{}),
	}
	s.lock.Unlock()

	// Read output and emit to frontend
	go func() {
		defer resp.Close()
		reader := bufio.NewReader(resp.Reader)
		buf := make([]byte, 4096)
		for {
			select {
			case <-s.session.done:
				return
			default:
				n, err := reader.Read(buf)
				if err != nil {
					if err != io.EOF {
						runtime.EventsEmit(s.ctx, "docker:output", fmt.Sprintf("Error: %v", err))
					}
					return
				}
				runtime.EventsEmit(s.ctx, "docker:output", string(buf[:n]))
			}
		}
	}()

	return nil
}

// Send user input from frontend to container
func (s *DockerContainersTerminal) SendToTerminal(input string) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.session == nil {
		return nil
	}

	_, err := s.session.stdin.Write([]byte(input))
	return err
}

// Close session cleanly
func (s *DockerContainersTerminal) CloseTerminal() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.session != nil {
		close(s.session.done)
		s.session.stdin.Close()
		s.session = nil
	}
}
