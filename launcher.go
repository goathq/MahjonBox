package main

import (
	"bufio"
	"context"
	"io"
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/runtime" // 正确导入路径
)

type Launcher struct {
	ctx       context.Context // Wails v2 必须的上下文
	skinCmd   *exec.Cmd
	helperCmd *exec.Cmd
}

// NewLauncher 构造函数
func NewLauncher() *Launcher {
	return &Launcher{}
}

// Wails v2 生命周期钩子 (替代原来的WailsInit)
func (l *Launcher) OnStartup(ctx context.Context) {
	l.ctx = ctx
}

// 启动子进程并实时捕获输出
func runCommandWithOutput(name string, arg ...string) (chan string, *exec.Cmd, error) {
	cmd := exec.Command(name, arg...)
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)
		scanner := bufio.NewScanner(io.MultiReader(stdoutPipe, stderrPipe))
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
	}()

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	return outputChan, cmd, nil
}

// 暴露给前端的方法保持不变
func (l *Launcher) WatchOutput(processType string) chan string {
	switch processType {
	case "skin":
		ch, cmd, _ := runCommandWithOutput("python", "majsoul_max.py")
		l.skinCmd = cmd
		return ch
	case "helper":
		ch, cmd, _ := runCommandWithOutput("./third_party/mahjong-helper")
		l.helperCmd = cmd
		return ch
	}
	return nil
}

func (l *Launcher) StopProcess(processType string) bool {
	switch processType {
	case "skin":
		if l.skinCmd != nil {
			if err := l.skinCmd.Process.Kill(); err != nil {
				runtime.LogError(l.ctx, "终止换肤进程失败: "+err.Error())
				return false
			}
			l.skinCmd = nil
			return true
		}
	case "helper":
		if l.helperCmd != nil {
			if err := l.helperCmd.Process.Kill(); err != nil {
				runtime.LogError(l.ctx, "终止牌助进程失败: "+err.Error())
				return false
			}
			l.helperCmd = nil
			return true
		}
	}
	return false
}

// ============= 新增功能 =============

// StartServices 同时启动多个服务（对应之前的3.2）
func (l *Launcher) StartServices(enableSkin, enableHelper bool) error {
	if enableSkin {
		ch, cmd, err := runCommandWithOutput("python", "majsoul_max.py")
		if err != nil {
			return err
		}
		l.skinCmd = cmd
		go l.forwardOutput(ch, "skin")
	}

	if enableHelper {
		ch, cmd, err := runCommandWithOutput("./third_party/mahjong-helper")
		if err != nil {
			return err
		}
		l.helperCmd = cmd
		go l.forwardOutput(ch, "helper")
	}
	return nil
}

// 私有方法：转发输出到前端事件
func (l *Launcher) forwardOutput(ch <-chan string, serviceType string) {
	for line := range ch {
		runtime.EventsEmit(l.ctx, serviceType+"_output", line)
	}
}

// SetSkin 设置皮肤（对应4.1）
func (l *Launcher) SetSkin(skinName string) error {
	// 实现你的皮肤配置逻辑
	return nil
}

// GetHelperStats 获取助手状态（对应4.2）
func (l *Launcher) GetHelperStats() (string, error) {
	// 实现获取助手状态的逻辑
	return "运行中", nil
}
