package gopybuf

//#cgo LDFLAGS: -Llibpython.3.13 -lpython3.13
//#cgo LDFLAGS: -lpython3.13 -ldl

/*
#include <Python.h>
*/
import "C"
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"io"
	"unsafe"
)

// gcc -c -o example_wrapper.o example_wrapper.c -I/path/to/python/include
// gcc -o libexample_wrapper.so -shared example_wrapper.o -L/path/to/python/lib -lpython3.13
// gcc -c -o example.o example.c -I/usr/include/python3.13
// gcc -o libexample.so -shared example.o -Llibpython.3.13 -lpython3.13
// CGO_CFLAGS=-I/usr/include/python3.13 go build

var (
	ErrPythonModuleNotInitialized   = errors.New("python module not initialized")
	ErrPythonFunctionNotInitialized = errors.New("python function not initialized")
	ErrPythonCantBuildArgs          = errors.New("failed to build arguments")
	ErrPythonCallFailed             = errors.New("failed to call Python function")
	ErrPythonCantUnmarshalResultErr = errors.New("failed to unmarshal result")
)

type pythonErr struct {
	Err       string `json:"error,omitempty"`
	Traceback string `json:"traceback,omitempty"`
}

func (c pythonErr) Error() string {
	return fmt.Sprintf("%s\n%s", c.Err, c.Traceback)
}

type Client interface {
	io.Closer
	grpc.ClientConnInterface
}

type clientConn struct {
	//needToDeref      []unsafe.Pointer
	moduleName, funcName, libsPath string
	targetPyFunction               *C.PyObject
	initialized                    bool
}

func (c *clientConn) Invoke(_ context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	if marshal, err := proto.Marshal(args.(proto.Message)); err != nil {
		return err
	} else {
		// Create Python bytes objects from the byte slices
		pyArg1 := c.strToPy(method)
		pyArg2 := c.bytesToPy(marshal)
		defer func() {
			// Release the references to the arguments and result
			c.pyObjectDecRef(pyArg1)
			c.pyObjectDecRef(pyArg2)
		}()

		// Create a tuple of arguments
		pyArgs := C.PyTuple_New(2)
		if pyArgs == nil {
			return ErrPythonCantBuildArgs
		}
		C.PyTuple_SetItem(pyArgs, 0, pyArg1)
		C.PyTuple_SetItem(pyArgs, 1, pyArg2)

		// Call the Python function with the arguments
		result := C.PyObject_CallObject(c.targetPyFunction, pyArgs)
		defer c.pyObjectDecRef(result)
		if result == nil {
			C.PyErr_Print()
			return ErrPythonCallFailed
		}
		// Convert the Python bytes object to Go byte slice
		pyResult := c.pyToBytes(C.PyTuple_GetItem(result, 0))
		pyErr := c.pyToBytes(C.PyTuple_GetItem(result, 1))
		if len(pyErr) != 0 {
			pyErrM := &pythonErr{}
			pyErrJsonErr := json.Unmarshal(pyErr, pyErrM)
			if pyErrJsonErr != nil {
				return ErrPythonCantUnmarshalResultErr
			}
			return pyErrM
		}
		return proto.Unmarshal(pyResult, reply.(proto.Message))
	}
}

func (c *clientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	//TODO implement me
	panic("implement me")
}

// Close finalizes the Python interpreter
func (c *clientConn) Close() error {
	if !c.initialized {
		return nil
	}
	c.pyObjectDecRef(c.targetPyFunction)
	// Cleanup and finalize the Python interpreter
	C.Py_Finalize()
	return nil
}

// strToPy converts a Go string to a Python string object
func (c *clientConn) strToPy(str string) *C.PyObject {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))
	return C.PyUnicode_DecodeUTF8(cStr, C.Py_ssize_t(len(str)), nil)
}

// bytesToPy converts a Go byte slice to a Python bytes object
func (c *clientConn) bytesToPy(data []byte) *C.PyObject {
	cData := C.CBytes(data)
	defer C.free(cData)
	return C.PyBytes_FromStringAndSize((*C.char)(cData), C.Py_ssize_t(len(data)))
}

// pyToBytes converts a Python bytes object to a Go byte slice
func (c *clientConn) pyToBytes(obj *C.PyObject) []byte {
	cStr := C.PyBytes_AsString(obj)
	size := C.PyBytes_Size(obj)
	return C.GoBytes(unsafe.Pointer(cStr), C.int(size))
}

// pyObjectDecRef releases the reference to a Python object
func (c *clientConn) pyObjectDecRef(obj *C.PyObject) {
	C.Py_DecRef(obj)
}

func (c *clientConn) init() error {
	C.Py_Initialize()
	moduleName := C.CString(c.moduleName)
	defer C.free(unsafe.Pointer(moduleName))
	pyModuleName := C.PyUnicode_DecodeFSDefault(moduleName)
	module := C.PyImport_Import(pyModuleName)
	if module == nil {
		C.PyErr_Print()
		return ErrPythonModuleNotInitialized
	}
	funcName := C.CString(c.funcName)
	defer C.free(unsafe.Pointer(funcName))
	c.targetPyFunction = C.PyObject_GetAttrString(module, funcName)
	if c.targetPyFunction == nil {
		C.PyErr_Print()
		return ErrPythonFunctionNotInitialized
	}
	c.initialized = true
	return nil
}

func NewClient(moduleName, funcName, libsPath string) (Client, error) {
	c := &clientConn{
		moduleName: moduleName,
		funcName:   funcName,
		libsPath:   libsPath,
	}
	err := c.init()
	if err != nil {
		_ = c.Close()
		return nil, err
	}
	return c, nil
}

//func main() {
//	client, createClientErr := NewClient("example", "call_go_py", "./")
//	if createClientErr != nil {
//		panic(createClientErr)
//		return
//	}
//	stub := gopybuf.NewTestServiceClient(client)
//	test, err := stub.Test(context.Background(), &gopybuf.TestMessage{
//		Name:   "test",
//		Age:    1,
//		Scores: []int32{1, 2, 3},
//	})
//	if err != nil {
//		println(err.Error())
//		return
//	}
//	fmt.Printf("test: %v\n", test)
//}

//C.initializePython()
//println("Hello World from Golang")
//C.callPythonHelloWorld()
//C.finalizePython()
//println("END FROM GO")
