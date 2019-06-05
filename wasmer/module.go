package wasmer

import (
	"io/ioutil"
	"unsafe"
)

// ReadBytes reads a `.wasm` file and returns its content as an array of bytes.
func ReadBytes(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// Validate validates a sequence of bytes that is supposed to represent a valid
// WebAssembly module.
func Validate(bytes []byte) bool {
	return true == cWasmerValidate((*cUchar)(unsafe.Pointer(&bytes[0])), cUint(len(bytes)))
}

// ModuleError represents any kind of errors related to a WebAssembly
// module.
type ModuleError struct {
	// Error message.
	message string
}

// NewModuleError constructs a new `ModuleError`.
func NewModuleError(message string) *ModuleError {
	return &ModuleError{message}
}

// `ModuleError` is an actual error. The `Error` function returns the
// error message.
func (error *ModuleError) Error() string {
	return error.message
}

// Module represents a WebAssembly module.
type Module struct {
	module *cWasmerModuleT
}

// Compile compiles a WebAssembly module from bytes.
func Compile(bytes []byte) (Module, error) {
	var module *cWasmerModuleT

	var compileResult = cWasmerCompile(
		&module,
		(*cUchar)(unsafe.Pointer(&bytes[0])),
		cUint(len(bytes)),
	)

	var emptyModule = Module{module: nil}

	if compileResult != cWasmerOk {
		return emptyModule, NewModuleError("Failed to compile the module.")
	}

	return Module{module}, nil
}

// Instantiate creates a new instance of the WebAssembly module.
func (module *Module) Instantiate() (Instance, error) {
	return module.InstantiateWithImports(NewImports())
}

// InstantiateWithImports creates a new instance with imports of the WebAssembly module.
func (module *Module) InstantiateWithImports(imports *Imports) (Instance, error) {
	return newInstanceWithImports(
		imports,
		func(wasmImportsCPointer *cWasmerImportT, numberOfImports int) (*cWasmerInstanceT, error) {
			var instance *cWasmerInstanceT

			var instantiateResult = cWasmerModuleInstantiate(
				module.module,
				&instance,
				wasmImportsCPointer,
				cInt(numberOfImports),
			)

			if instantiateResult != cWasmerOk {
				return nil, NewModuleError("Failed to instantiate the module.")
			}

			return instance, nil
		},
	)
}

// Close closes/frees a `Module`.
func (module *Module) Close() {
	if module.module != nil {
		cWasmerModuleDestroy(module.module)
	}
}
