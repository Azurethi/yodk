package testing

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dbaumgarten/yodk/pkg/nolol"

	yaml "gopkg.in/yaml.v2"

	"github.com/dbaumgarten/yodk/pkg/vm"
)

// Test defines a test-run
type Test struct {
	// The path where the test-file was located. Used to retrieve the script files.
	Path string
	// Scripts to use in this test
	Scripts []Script
	// Cases for this test
	Cases []Case
}

// Script contains run-options for a script in the test
type Script struct {
	// The absolut path where the test-file was located. Used to retrieve the script files.
	TestPath string
	// Name of the script to run
	Name string
	// Maximum number of iterations for the script (0=infinite)
	Iterations int
	// Maximum number of lines to run from the script (0=infinite)
	MaxLines int
	// the content of the script. If empty, it is loaded from disk at run-time
	Content string
}

// Case defines inputs and expected outputs for a run
type Case struct {
	// Name of the testcase
	Name string
	// Values of gloal variables before run
	Inputs map[string]interface{}
	// Expected values of global vars after run
	Outputs map[string]interface{}
}

func prefixVarname(inp string) string {
	if !strings.HasPrefix(inp, ":") {
		return ":" + inp
	}
	return inp
}

// Parse parses a yaml file into a Test
// path is the path from where the test was loaded. This is needed as the scripts are located relatice to the test-file
func Parse(file []byte, path string) (Test, error) {
	var test Test
	err := yaml.Unmarshal(file, &test)
	if err != nil {
		return test, fmt.Errorf("The provided test-file is invalid: %s", err.Error())
	}
	test.Path = path
	for i, script := range test.Scripts {
		if script.Iterations == 0 {
			test.Scripts[i].Iterations = 1
		}
		test.Scripts[i].TestPath = path
	}
	return test, nil
}

// InitializeVariables adds the variables required for the testcase
// to the variables of the given Coordinator
func (c Case) InitializeVariables(coord *vm.Coordinator) error {
	for key, value := range c.Inputs {
		//key = strings.ToLower(key)
		variable, err := vm.VariableFromType(value)
		if err != nil {
			return err
		}
		coord.SetVariable(prefixVarname(key), variable)
	}
	return nil
}

// GetCode returns the code for script either from the script struct itself or from the referenced file
func (script Script) GetCode() (string, error) {
	file := filepath.Join(filepath.Dir(script.TestPath), script.Name)
	if script.Content == "" {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		return string(f), nil
	}
	return script.Content, nil
}

// CreateVMs creates and sets up the required vms for this test
// coord is the coordinato to use with the VMs
// Run() has been called on the returned VMs, but they are paused until coord.Run() is called
// The error handler of the VMs is set to errF
// Also returns variable-name translation-tables for nolol scripts
func (t Test) CreateVMs(coord *vm.Coordinator, errF vm.ErrorHandlerFunc) ([]*vm.VM, []map[string]string, error) {
	vms := make([]*vm.VM, len(t.Scripts))
	translationTables := make([]map[string]string, len(t.Scripts))
	for i, script := range t.Scripts {
		var v *vm.VM

		if strings.HasSuffix(script.Name, ".nolol") {
			conv := nolol.NewConverter()
			file := filepath.Join(filepath.Dir(script.TestPath), script.Name)
			prog, err := conv.ConvertFile(file)
			translationTables[i] = conv.GetVariableTranslations()
			if err != nil {
				return nil, nil, err
			}
			v = vm.Create(prog)
		} else {
			scriptContent, err := script.GetCode()
			if err != nil {
				return nil, nil, err
			}
			v, err = vm.CreateFromSource(string(scriptContent))
			if err != nil {
				return nil, nil, err
			}
		}
		v.SetIterations(script.Iterations)
		v.SetMaxExecutedLines(script.MaxLines)
		v.SetErrorHandler(errF)
		v.SetCoordinator(coord)
		vms[i] = v
		v.Resume()
	}
	return vms, translationTables, nil
}

// CheckResults compares the global variables of coord with the expected results for c
// and returns found errors
func (c Case) CheckResults(coord *vm.Coordinator) []error {
	fails := make([]error, 0)
	for key, value := range c.Outputs {
		//key = strings.ToLower(key)
		key = prefixVarname(key)
		var fail error
		expected, err := vm.VariableFromType(value)
		if err != nil {
			fail = fmt.Errorf("Invalid type for expected var: %T", value)
			fails = append(fails, fail)
			continue
		}
		actual, exists := coord.GetVariable(key)

		if !exists {
			fail = fmt.Errorf("Expected output variable %s does not exist", key)
		} else {
			if !actual.SameType(expected) {
				fail = fmt.Errorf("Case '%s': Output '%s' has type '%s' but should be '%s' ", c.Name, key, actual.TypeName(), expected.TypeName())

			} else if !actual.Equals(expected) {
				fail = fmt.Errorf("Case '%s': Output '%s' has value %s but should be %s ", c.Name, key, actual.Repr(), expected.Repr())
			}
		}
		if fail != nil {
			fails = append(fails, fail)
		}
	}
	return fails
}

// Run runs a the given test and return found errors
// caseCallback is called before executing a case. Can be used for logging.
// Main method of the test class
func (t Test) Run(caseCallback func(c Case)) []error {

	fails := make([]error, 0)
	flock := &sync.Mutex{}

	for _, c := range t.Cases {
		if caseCallback != nil {
			caseCallback(c)
		}
		coord := vm.NewCoordinator()
		c.InitializeVariables(coord)

		errHandler := func(vm *vm.VM, err error) bool {
			flock.Lock()
			defer flock.Unlock()
			fails = append(fails, err)
			go coord.Terminate()
			return true
		}

		_, _, err := t.CreateVMs(coord, errHandler)
		if err != nil {
			return []error{err}
		}

		coord.Run()
		coord.WaitForTermination()

		caseFails := c.CheckResults(coord)
		fails = append(fails, caseFails...)
	}

	return fails
}
