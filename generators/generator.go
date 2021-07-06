package generators

type CodeGenerator interface {
    Init(interface{}) error
    GetType() string
    GetName() string
    GenerateFile() ([]byte, string, error)
}