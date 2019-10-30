package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/dbaumgarten/yodk/nolol"
	"github.com/dbaumgarten/yodk/parser"

	"github.com/abiosoft/ishell"
	"github.com/dbaumgarten/yodk/vm"
	"github.com/spf13/cobra"
)

var yvm *vm.YololVM
var debugShell *ishell.Shell
var inputProg string
var inputFileName string

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug [file]",
	Short: "Debug a yolol/nolol program",
	Long:  `Execute program interactively in debugger`,
	Run: func(cmd *cobra.Command, args []string) {
		inputFileName = args[0]
		inputProg = loadInputFile(args[0])

		if !strings.HasSuffix(inputFileName, ".yolol") && !strings.HasSuffix(inputFileName, ".nolol") {
			fmt.Println("Unknown file-extension for file: ", inputFileName)
			os.Exit(1)
		}

		debugShell.Println("Loaded and paused programm. Enter 'run' to execute")
		debugShell.Run()
	},
	Args: cobra.MinimumNArgs(1),
}

func run() {
	if strings.HasSuffix(inputFileName, ".yolol") {
		debugShell.Println("--Started--")
		go yvm.RunSource(inputProg)
		return
	}
	if strings.HasSuffix(inputFileName, ".nolol") {
		debugShell.Println("--Started--")
		converter := nolol.NewConverter()
		yololcode, err := converter.ConvertFromSource(inputProg)
		if err != nil {
			exitOnError(err, "pasrsing nolol code")
		}
		go yvm.Run(yololcode, inputProg)
		return
	}
}

func init() {
	rootCmd.AddCommand(debugCmd)

	debugShell = ishell.New()

	yvm = vm.NewYololVM()
	yvm.Pause()

	yvm.SetBreakpointHandler(func(x *vm.YololVM) bool {
		debugShell.Println("--Hit Breakpoint at line: ", x.CurrentSourceLine(), "--")
		return false
	})

	yvm.SetErrorHandler(func(x *vm.YololVM, err error) bool {
		debugShell.Println("--A runtime error occured--")
		debugShell.Println(err)
		debugShell.Println("--Execution paused--")
		return false
	})

	yvm.SetFinishHandler(func(x *vm.YololVM) {
		debugShell.Println("--Program finished--")
		debugShell.Println("--Enter r to restart--")
	})

	debugShell.AddCmd(&ishell.Cmd{
		Name:    "run",
		Aliases: []string{"r"},
		Help:    "run programm from start",
		Func: func(c *ishell.Context) {
			run()
		},
	})
	debugShell.AddCmd(&ishell.Cmd{
		Name:    "pause",
		Aliases: []string{"p"},
		Help:    "pause execution",
		Func: func(c *ishell.Context) {
			yvm.Pause()
			debugShell.Println("--Paused--")
		},
	})
	debugShell.AddCmd(&ishell.Cmd{
		Name:    "continue",
		Aliases: []string{"c"},
		Help:    "continue paused execution",
		Func: func(c *ishell.Context) {
			err := yvm.Resume()
			if err == nil {
				debugShell.Println("--Resumed--")
			} else {
				debugShell.Println(err)
			}
		},
	})
	debugShell.AddCmd(&ishell.Cmd{
		Name:    "step",
		Aliases: []string{"s"},
		Help:    "execute the next line and pause again",
		Func: func(c *ishell.Context) {
			if yvm.Step() == nil {
				debugShell.Println("--Line executed. Paused again--")
			}
		},
	})
	debugShell.AddCmd(&ishell.Cmd{
		Name:    "break",
		Aliases: []string{"b"},
		Help:    "add breakpoint at line",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				debugShell.Println("You must enter a line number for the breakpoint.")
				return
			}
			line, err := strconv.Atoi(c.Args[0])
			if err != nil {
				debugShell.Println("Error parsing line-number: ", err)
				return
			}
			yvm.AddBreakpoint(line)
			debugShell.Println("--Breakpoint added--")
		},
	})
	debugShell.AddCmd(&ishell.Cmd{
		Name:    "delete",
		Aliases: []string{"d"},
		Help:    "delete breakpoint at line",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				debugShell.Println("You must enter a line number for the breakpoint.")
				return
			}
			line, err := strconv.Atoi(c.Args[0])
			if err != nil {
				debugShell.Println("Error parsing line-number: ", err)
				return
			}
			yvm.RemoveBreakpoint(line)
			debugShell.Println("--Breakpoint removed--")
		},
	})
	debugShell.AddCmd(&ishell.Cmd{
		Name:    "vars",
		Aliases: []string{"v"},
		Help:    "print all current variables",
		Func: func(c *ishell.Context) {
			debugShell.Println("--Variables--")
			vars := sortVariables(yvm.GetVariables())
			for _, variable := range vars {
				if variable.val.IsString() {
					debugShell.Println(variable.name, "'"+variable.val.String()+"'")
				}
				if variable.val.IsNumber() {
					debugShell.Println(variable.name, variable.val.Itoa())
				}
			}
		},
	})
	debugShell.AddCmd(&ishell.Cmd{
		Name:    "info",
		Aliases: []string{"i"},
		Help:    "show vm-state",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)
			debugShell.Printf("--State: %d\n", yvm.State())
		},
	})
	debugShell.AddCmd(&ishell.Cmd{
		Name:    "list",
		Aliases: []string{"l"},
		Help:    "show programm source code",
		Func: func(c *ishell.Context) {
			current := yvm.CurrentSourceLine()
			bps := yvm.ListBreakpoints()
			progLines := strings.Split(inputProg, "\n")
			debugShell.Println("--Programm--")
			pfx := ""
			for i, line := range progLines {
				if i+1 == current {
					pfx = ">"
				} else {
					pfx = " "
				}
				if contains(bps, i+1) {
					pfx += "x"
				} else {
					pfx += " "
				}
				pfx += fmt.Sprintf("%3d ", i+1)
				debugShell.Println(pfx + line)
			}
		},
	})
	debugShell.AddCmd(&ishell.Cmd{
		Name:    "disas",
		Aliases: []string{"d"},
		Help:    "show yolol code for nolol source",
		Func: func(c *ishell.Context) {
			if !strings.HasSuffix(inputFileName, ".nolol") {
				debugShell.Print("Disas is only available when debugging nolol code")
			}
			current := yvm.CurrentAstLine()
			conv := nolol.NewConverter()
			ast, _ := conv.ConvertFromSource(inputProg)
			yolol, _ := (&parser.Printer{}).Print(ast)
			progLines := strings.Split(yolol, "\n")
			debugShell.Println("--Programm--")
			pfx := ""
			for i, line := range progLines {
				if i+1 == current {
					pfx = ">"
				} else {
					pfx = " "
				}
				pfx += fmt.Sprintf("%3d ", i+1)
				debugShell.Println(pfx + line)
			}
		},
	})
}

type namedVariable struct {
	name string
	val  vm.Variable
}

func sortVariables(vars map[string]vm.Variable) []namedVariable {
	sorted := make([]namedVariable, 0, len(vars))
	for k, v := range vars {
		sorted = append(sorted, namedVariable{
			k,
			v,
		})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].name < sorted[j].name
	})
	return sorted
}

func contains(arr []int, val int) bool {
	for _, e := range arr {
		if e == val {
			return true
		}
	}
	return false
}
