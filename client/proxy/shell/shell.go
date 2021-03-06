package shell

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gookit/color"
	commonnet "github.com/hellgate75/go-tcp-common/net"
	"github.com/hellgate75/go-tcp-client/common"
	"github.com/hellgate75/go-tcp-common/log"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type shell struct{
	logger log.Logger
}

func existsFile(file string) bool {
	_, err1 := os.Stat(file)
	if err1 != nil {
		return false
	}
	return true
}

func loadFile(path string) ([]byte, error) {
	file, err1 := os.Open(path)
	if err1 != nil {
		return nil, err1
	}
	return ioutil.ReadAll(file)
}

var serverCommand string = "shell"

func (shell *shell) SetLogger(logger log.Logger) {
	shell.logger = logger
}


func (shell *shell) SendMessage(conn *tls.Conn, params ...interface{}) error {
	var paramsLen int = len(params)
	var interactive string = "true"
	if paramsLen > 0 {
		if "true" != fmt.Sprintf("%v", params[0]) {
			interactive = "false"
		}

	}

	var shellCommandOrScript string = ""
	var isScriptFile bool = false
	if paramsLen > 1 && params[1] != "" {
		if "" != fmt.Sprintf("%v", params[1]) {
			shellCommandOrScript = fmt.Sprintf("%v", params[1])
			isScriptFile = len(shellCommandOrScript) > 5 && strings.Index(shellCommandOrScript, ".") >= len(shellCommandOrScript)-5
			interactive = "false"
		}
	}
	var stdin io.Reader
	if paramsLen > 2 && params[2] != nil {
		stdin = params[2].(io.Reader)
	} else {
		stdin = os.Stdin
	}
	var stdout io.Writer
	if paramsLen > 3 && params[3] != nil {
		stdout = params[3].(io.Writer)
	}
	var stderr io.Writer
	if paramsLen > 4 && params[4] != nil {
		stderr = params[4].(io.Writer)
	}

	//	fmt.Printf("Shell Script: %s, Is Script: %v\n", shellCommandOrScript, isScriptFile)
	n0, err3b := common.WriteString(serverCommand, conn)
	if err3b != nil {
		return err3b
	}
	if n0 == 0 {
		return errors.New(fmt.Sprintf("Unable to send command: %s", serverCommand))
	}
	time.Sleep(3 * time.Second)
	n1, err4 := common.WriteString(interactive, conn)
	if err4 != nil {
		return err4
	}
	if n1 == 0 {
		return errors.New(fmt.Sprintf("Unable to send interactive: %s", interactive))
	}
	time.Sleep(3 * time.Second)
	if "" != shellCommandOrScript {
		var script string = ""
		if isScriptFile {
			if !existsFile(shellCommandOrScript) {
				commonnet.WriteString("exit", conn)
				return errors.New(fmt.Sprintf("Script File %s doesn't exists!!", shellCommandOrScript))
			}
			n2, err5 := common.WriteString("script", conn)
			if err5 != nil {
				commonnet.WriteString("exit", conn)
				return err5
			}
			if n2 == 0 {
				commonnet.WriteString("exit", conn)
				return errors.New(fmt.Sprintf("Unable to send script file type: %v", isScriptFile))
			}
			fileName := shellCommandOrScript
			if strings.Contains(shellCommandOrScript, "/") {
				listX := strings.Split(shellCommandOrScript, "/")
				fileName = listX[len(listX)-1]
			} else if strings.Contains(shellCommandOrScript, "\\") {
				listX := strings.Split(shellCommandOrScript, "\\")
				fileName = listX[len(listX)-1]
			}
			n2, err5 = common.WriteString(fileName, conn)
			if err5 != nil {
				commonnet.WriteString("exit", conn)
				return err5
			}
			if n2 == 0 {
				commonnet.WriteString("exit", conn)
				return errors.New(fmt.Sprintf("Unable to send script file type: %v", isScriptFile))
			}
			content, errReadScript := loadFile(shellCommandOrScript)
			if errReadScript != nil {
				commonnet.WriteString("exit", conn)
				return errors.New(fmt.Sprintf("Cannot read script File %s -> Details: %s", shellCommandOrScript, errReadScript.Error()))
			}
			script = string(content)
		} else {
			n2, err5 := common.WriteString("command", conn)
			if err5 != nil {
				commonnet.WriteString("exit", conn)
				return err5
			}
			if n2 == 0 {
				commonnet.WriteString("exit", conn)
				return errors.New(fmt.Sprintf("Unable to send COMMAND -> script file type: %v", isScriptFile))
			}
			script = shellCommandOrScript
			n3, err6 := common.Write([]byte(script), conn)
			if err6 != nil {
				common.WriteString("exit", conn)
				return err6
			}
			if n3 == 0 {
				commonnet.WriteString("exit", conn)
				return errors.New(fmt.Sprintf("Unable to send data -> shell command: %v", script))
			}
		}
		state, errContinue := common.ReadString(conn)
		if errContinue  != nil {
			return errors.New(fmt.Sprintf("Receive pre-conditions error -> shell command: %v, Details: %s", script, errContinue.Error()))
		}
		if shell.logger != nil {
			shell.logger.Debugf("Pre-conditions message: <%s>", state)
		}
		if len(state) > 2 && "ko" == state[:2] {
			return errors.New(fmt.Sprintf("Pre-conditions failed -> shell command: %v, Details: %s", script, state))
		}
		content, errAnswer := common.Read(conn)
		if errAnswer != nil {
			return errors.New(fmt.Sprintf("Receive data -> shell command: %v, Details: %s", script, errAnswer.Error()))
		}
		if nil != shell.logger {
			shell.logger.Debugf("Response: %s", string(content))
		} else {
			color.LightWhite.Printf("Response: %s\n", string(content))
		}
		if stdout != nil {
			_, err := stdout.Write(content)
			if err != nil {
				return err
			}
		}
	} else {
		n2, err5 := common.WriteString("shell", conn)
		if err5 != nil {
			commonnet.WriteString("exit", conn)
			if stderr != nil {
				if nil != shell.logger {
					shell.logger.Error("Error: exit shell: " + err5.Error() + "!!")
				} else {
					stderr.Write([]byte("Error: exit shell: " + err5.Error() + "!!\n"))
				}
			} else {
				if nil != shell.logger {
					shell.logger.Error("Error: exit shell: " + err5.Error() + "!!")
				} else {
					color.Red.Println("Error: exit shell: " + err5.Error() + "!!")
				}
			}
			return err5
		}
		if n2 == 0 {
			commonnet.WriteString("exit", conn)
			if stderr != nil {
				if nil != shell.logger {
					shell.logger.Error("Error: exit shell!!")
				} else {
					stderr.Write([]byte("Error: exit shell!!\n"))
				}
			} else {
				if nil != shell.logger {
					shell.logger.Error("Error: exit shell!!")
				} else {
					color.Red.Println("Error: exit shell!!")
				}
			}
			return errors.New("Unable to send shell command")
		}
		if stdout != nil {
			if nil != shell.logger {
				shell.logger.Info("Shell mode : type exit command to exit the interactive mode\n")
			} else {
				stdout.Write([]byte("Shell mode : type exit command to exit the interactive mode\n"))
			}
		} else {
			if nil != shell.logger {
				shell.logger.Info("Shell mode : type exit command to exit the interactive mode\n")
			} else {
				color.LightYellow.Printf("Shell mode : type exit command to exit the interactive mode\n")
			}
		}
		time.Sleep(3 * time.Second)
		color.Green.Printf("shell> ")
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {
			var currentCommand string = scanner.Text()
			if "" == currentCommand {
				color.Yellow.Println("Empty command, try again...")
				color.Green.Printf("shell> ")
				continue
			}
			color.Yellow.Printf("Sending request to the server...\n")
			if "exit" == strings.ToLower(currentCommand) {
				if stdout != nil {
					if nil != shell.logger {
						shell.logger.Debug("Request: exit shell!!")
					} else {
						stdout.Write([]byte("Request: exit shell!!"))
					}
				} else {
					if nil != shell.logger {
						shell.logger.Debug("Request: exit shell!!")
					} else {
						color.Yellow.Println("Request: exit shell!!")
					}
				}
				commonnet.WriteString("exit", conn)
				break
			}
			n3, err6 := common.WriteString(currentCommand, conn)
			if err6 != nil {
				commonnet.WriteString("exit", conn)
				if stderr != nil {
					if nil != shell.logger {
						shell.logger.Error("Error: exit shell: " + err6.Error() + "!!")
					} else {
						stderr.Write([]byte("Error: exit shell: " + err6.Error() + "!!\n"))
					}
				} else {
					if nil != shell.logger {
						shell.logger.Error("Error: exit shell: " + err6.Error() + "!!")
					} else {
						color.Red.Println("Error: exit shell: " + err6.Error() + "!!")
					}
				}
				return err6
			}
			if n3 == 0 {
				commonnet.WriteString("exit", conn)
				if stderr != nil {
					if nil != shell.logger {
						shell.logger.Error("Error: exit shell!!")
					} else {
						stderr.Write([]byte("Error: exit shell!!\n"))
					}
				} else {
					if nil != shell.logger {
						shell.logger.Error("Error: exit shell!!")
					} else {
						color.Red.Println("Error: exit shell!!")
					}
				}
				return errors.New(fmt.Sprintf("Unable to send command ->  %v", currentCommand))
			}
			//time.Sleep(3 * time.Second)
			content, errAnswer := common.Read(conn)
			if errAnswer != nil {
				commonnet.WriteString("exit", conn)
				if stderr != nil {
					if nil != shell.logger {
						shell.logger.Error("Error: exit shell: " + errAnswer.Error() + "!!")
					} else {
						stderr.Write([]byte("Error: exit shell: " + errAnswer.Error() + "!!\n"))
					}
				} else {
					if nil != shell.logger {
						shell.logger.Error("Error: exit shell: " + errAnswer.Error() + "!!")
					} else {
						color.Red.Println("Error: exit shell: " + errAnswer.Error() + "!!")
					}
				}
				return errAnswer
			}
			//color.LightYellow.Printf("Answer: %s\n", content)
			if stdout != nil {
				if nil != shell.logger {
					shell.logger.Debug("Response: ", string(content))
				} else {
					stdout.Write([]byte(fmt.Sprintf("Response: ", string(content)) + "\n"))
				}
			} else {
				if nil != shell.logger {
					shell.logger.Debug("Response: ", string(content))
				} else {
					color.LightWhite.Println("Response: ", string(content))
				}
			}
			color.Green.Printf("shell> ")
		}

		if err := scanner.Err(); err != nil {
			if stderr != nil {
				if nil != shell.logger {
					shell.logger.Error("Error: exit shell: " + err.Error() + "!!")
				} else {
					stderr.Write([]byte("Error: exit shell: " + err.Error() + "!!\n"))
				}
			} else {
				if shell.logger != nil {
					shell.logger.Error("Error: exit shell: " + err.Error() + "!!")
				} else {
					color.Red.Println("Error: exit shell: " + err.Error() + "!!")
				}
			}
		}

	}
	return nil
}
func (shell *shell) Helper() string {
	return "shell [interactive] [script file|command]\n  Parameters:\n    [interactive]      (optional) interactive shell[true/false] (default: true)\n    [script file]      (optional) full path of local script file\n    [command]          (optional) shell command\n"
}

func New() common.Sender {
	return &shell{
		logger: nil,
	}
}
