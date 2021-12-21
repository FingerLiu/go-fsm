package singletonfsm

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"log"
	"reflect"
	"runtime"
)

func (f *FSM) RenderGraphvizDot() string {
	g, graph := f.buildGraphviz()
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()
	var buf bytes.Buffer
	if err := g.Render(graph, "dot", &buf); err != nil {
		log.Fatal(err)
	}
	dot := buf.String()
	fmt.Println(dot)
	return dot
}

func (f *FSM) RenderGraphvizImage(filename string) {
	g, graph := f.buildGraphviz()
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()
	if filename == "" {
		imageName := f.name
		if f.name == "" {
			imageName = "demo"
		}
		filename = fmt.Sprintf("./%s.png", imageName)
	}

	// 1. write encoded PNG data to buffer
	var imageBuf bytes.Buffer
	if err := g.Render(graph, graphviz.PNG, &imageBuf); err != nil {
		log.Fatal(err)
	}

	// 2. get as image.Image instance
	_, err := g.RenderImage(graph)
	if err != nil {
		log.Fatal(err)
	}

	// 3. write to file directly
	if err := g.RenderFilename(graph, graphviz.PNG, filename); err != nil {
		log.Fatal(err)
	}
}

func (f *FSM) buildGraphviz() (*graphviz.Graphviz, *cgraph.Graph) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}

	for _, state := range f.states {
		graph.CreateNode(state.Name)
	}

	for _, transition := range f.transitions {
		var label string
		if transition.Condition != nil {
			label = getFunctionName(transition.Condition)
		}
		fromNode, _ := graph.Node(transition.From.Name)
		toNode, _ := graph.Node(transition.To.Name)
		e, _ := graph.CreateEdge(transition.Key, fromNode, toNode)
		e.SetLabel(label)
	}
	return g, graph
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
