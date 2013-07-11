package gltools

import (
	gl "github.com/chsc/gogl/gl33"
	"github.com/threeguys/math3d"
	"unsafe"
)

type Parameter struct {
	id gl.Int
	name string
}

func (p *Parameter) Matrix(m *math3d.Matrix) {
	gl.UniformMatrix4fv(p.id, 1, gl.FALSE, (*gl.Float)(m.Pointer()) )
}

type Program struct {
	id gl.Uint
}

func (p *Program) Use() {
	gl.UseProgram(p.id)
}

func (p *Program) GetParameter(name string) *Parameter {
	param := new(Parameter)
	param.name = name
	param.id = gl.GetUniformLocation(p.id, gl.GLString(name))
	return param
}

func (p *Program) Delete() {
	gl.DeleteProgram(p.id)
}

func Load(vertexShaderSource, fragmentShaderSource string) (*Program, error) {
	glVtxSrc := gl.GLString(vertexShaderSource)
	defer gl.GLStringFree(glVtxSrc)

	vShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vShader, 1, &glVtxSrc, nil)
	gl.CompileShader(vShader)
	defer gl.DeleteShader(vShader)
	
	var compiled gl.Int
	gl.GetShaderiv(vShader, gl.COMPILE_STATUS, &compiled)
	
	if compiled == gl.FALSE {
		var length gl.Int
		gl.GetShaderiv(vShader, gl.INFO_LOG_LENGTH, &length)
		msg := make([]byte, length + 1)
		gl.GetShaderInfoLog(vShader, gl.Sizei(length), nil, (*gl.Char)(unsafe.Pointer(&msg[0])))
		
		return nil, GLToolsError{ string(msg) }
	}

	glFrgSrc := gl.GLString(fragmentShaderSource)
	defer gl.GLStringFree(glFrgSrc)
	fShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fShader, 1, &glFrgSrc, nil)
	gl.CompileShader(fShader)
	defer gl.DeleteShader(fShader)	

	if compiled == gl.FALSE {
		var length gl.Int
		gl.GetShaderiv(fShader, gl.INFO_LOG_LENGTH, &length)
		msg := make([]byte, length + 1)
		gl.GetShaderInfoLog(fShader, gl.Sizei(length), nil, (*gl.Char)(unsafe.Pointer(&msg[0])))
		
		return nil, GLToolsError{ string(msg) }
	}
	
	program := gl.CreateProgram()
	gl.AttachShader(program, vShader)
	gl.AttachShader(program, fShader)

	gl.LinkProgram(program)
	
	p := new(Program)
	p.id = program
	return p, nil
}
