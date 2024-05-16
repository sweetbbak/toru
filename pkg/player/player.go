package player

import (
	"bytes"
	"text/template"
)

// GenericPlayer represents most players. The stream URL will be appended to the arguments.
type GenericPlayer struct {
	Name string
	Args []string

	GetHelperScriptPath func() string
}

// Player opens a stream URL in a video player.
type Player interface {
	Open(url string) error
}

type MediaEntry struct {
	Title string
	URL   string
}

func (p *GenericPlayer) expandArgs(data MediaEntry) ([]string) {
	var res []string
	var buffer bytes.Buffer

	for _, u := range p.Args {
		temp := template.New("")
		temp.Funcs(template.FuncMap{"GetHelperScriptPath": p.GetHelperScriptPath})
		temp, _ = temp.Parse(u)

		temp.Execute(&buffer, data)

		res = append(res, buffer.String())
		buffer.Reset()
	}

	return res
}
