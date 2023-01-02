package command

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

// TestCommand go调用系统支持
func TestCommand(t *testing.T) {
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)

	// 创建command
	//cmd = exec.Command("/bin/bash", "-c", "sleep 5; ls -l")
	cmd = exec.Command("/bin/bash", "-c", "sleep 5; ls -l; echo hello; /bin/bash -c ls")

	// 执行命令，捕获子进程输出（pipe）
	if output, err = cmd.CombinedOutput(); err != nil {
		t.Error(err)
		return
	}

	// 打印子进程输出
	t.Logf("%s", output)
}

type result struct {
	output []byte
	err    error
}

func TestKillCommand(t *testing.T) {
	// 执行一个cmd，它在协程中执行，让它执行2秒，sleep 2; echo hello;
	// 过 1 秒，杀掉上面执行的这个cmd
	var (
		ctx        context.Context
		cancelFunc context.CancelFunc
		cmd        *exec.Cmd
		resultChan chan *result
	)

	resultChan = make(chan *result, 10)

	// context: chan struct{}
	// cancelFunc: close(chan struct{})

	ctx, cancelFunc = context.WithCancel(context.TODO())

	go func() {
		var (
			output []byte
			err    error
		)
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 2; echo hello")

		// select {
		//	case <- ctx.Done():
		//}
		// kill pid, 进程ID, 杀死子进程
		/*if output, err = cmd.CombinedOutput(); err != nil {
			t.Log(err)
			return
		}*/
		//t.Logf("%s", output)
		output, err = cmd.CombinedOutput()
		resultChan <- &result{
			output: output,
			err:    err,
		}
	}()

	// 继续执行
	time.Sleep(1 * time.Second)
	// 取消上下文
	cancelFunc()

	// main协程等着子协程里面的执行结果，并等待结果
	res := <-resultChan
	if res.err != nil {
		t.Log(res.err)
		return
	}

	t.Logf("%s", res.output)
}
