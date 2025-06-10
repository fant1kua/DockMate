package app

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TerminalSession struct {
	stdinWriter io.WriteCloser
	done        chan struct{}
	lock        sync.Mutex
}

type DockerContainersTerminal struct {
	DockerBaseService
	lock sync.RWMutex
}

var sessions = map[string]*TerminalSession{}

func NewDockerTerminalService() *DockerContainersTerminal {
	return &DockerContainersTerminal{}
}

func StartupDockerTerminalService(s *DockerContainersTerminal, ctx context.Context, cli *client.Client) {
	s.ctx = ctx
	s.cli = cli
}

func (s *DockerContainersTerminal) StartInteractiveTerminal(id string, sessionID string) error {
	if s.cli == nil || s.ctx == nil {
		return fmt.Errorf("Docker client not initialized")
	}

	execConfig := container.ExecOptions{
		Cmd:          []string{"bash"}, // or "sh" if bash is unavailable
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

	session := &TerminalSession{
		stdinWriter: resp.Conn,
		done:        make(chan struct{}),
	}
	s.lock.Lock()
	sessions[sessionID] = session
	s.lock.Unlock()

	// Read output and emit to frontend
	go func() {
		defer resp.Close()
		reader := bufio.NewReader(resp.Reader)
		buf := make([]byte, 4096)
		for {
			select {
			case <-session.done:
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
func (s *DockerContainersTerminal) SendToTerminal(sessionID string, input string) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	session, ok := sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found")
	}

	_, err := session.stdinWriter.Write([]byte(input))
	return err
}

// Close session cleanly
func (s *DockerContainersTerminal) CloseTerminal(sessionID string) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if session, ok := sessions[sessionID]; ok {
		close(session.done)
		session.stdinWriter.Close()
		delete(sessions, sessionID)
	}
}
