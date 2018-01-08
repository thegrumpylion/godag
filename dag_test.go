package dag_test

import (
	"fmt"
	"testing"

	"."
)

func TestDag(t *testing.T) {

	nodes := []struct {
		id   string
		data map[string]string
	}{
		{"a.c", map[string]string{"licence": "gpl"}},
		{"b.c", map[string]string{"licence": "mit"}},
		{"c.c", map[string]string{"licence": "apache2"}},
		{"a.o", map[string]string{}},
		{"b.o", map[string]string{}},
		{"c.o", map[string]string{}},
		{"libcustom.so", map[string]string{}},
		{"libsysA.so", map[string]string{"licence": "apache2"}},
		{"libsysB.so", map[string]string{"licence": "custom"}},
		{"binA", map[string]string{}},
		{"binB", map[string]string{}},
	}

	edges := []struct {
		from string
		to   string
	}{
		{"binA", "libcustom.so"},
		{"binA", "libsysA.so"},
		{"binB", "libcustom.so"},
		{"binB", "libsysB.so"},
		{"libcustom.so", "a.o"},
		{"libcustom.so", "b.o"},
		{"libcustom.so", "c.o"},
		{"a.o", "a.c"},
		{"b.o", "b.c"},
		{"c.o", "c.c"},
	}

	d := dag.New()

	for _, node := range nodes {
		d.AddNode(node.id, node.data)
	}

	for _, edge := range edges {
		d.AddEdge(edge.from, edge.to)
	}

	err := d.AddEdge("a.o", "a.c")
	fmt.Println("Error:", err)

	roots := d.Roots()
	fmt.Printf("Roots found: %d\n\n", len(roots))

	fmt.Println("Walking depth first\n")
	for _, root := range roots {
		fmt.Printf("Root: %s\n\n", root.ID)
		root.Walk(dag.WalkDepthFirst, func(node *dag.Node) error {
			fmt.Println("Visiting node", node.ID)
			return nil
		})
		fmt.Println("")
	}

	fmt.Println("Walking breadth first\n")
	for _, root := range roots {
		fmt.Printf("Root: %s\n\n", root.ID)
		root.Walk(dag.WalkBreadthFirst, func(node *dag.Node) error {
			fmt.Println("Visiting node", node.ID)
			return nil
		})
		fmt.Println("")
	}
}

func TestLicences(t *testing.T) {

	nodes := []struct {
		id   string
		data map[string]string
	}{
		{"a.c", map[string]string{"licence": "gpl"}},
		{"b.c", map[string]string{"licence": "mit"}},
		{"c.c", map[string]string{"licence": "apache2"}},
		{"a.o", map[string]string{}},
		{"b.o", map[string]string{}},
		{"c.o", map[string]string{}},
		{"libcustom.so", map[string]string{}},
		{"libsysA.so", map[string]string{"licence": "apache2"}},
		{"libsysB.so", map[string]string{"licence": "custom"}},
		{"binA", map[string]string{}},
		{"binB", map[string]string{}},
	}

	edges := []struct {
		from string
		to   string
	}{
		{"binA", "libcustom.so"},
		{"binA", "libsysA.so"},
		{"binB", "libcustom.so"},
		{"binB", "libsysB.so"},
		{"libcustom.so", "a.o"},
		{"libcustom.so", "b.o"},
		{"libcustom.so", "c.o"},
		{"a.o", "a.c"},
		{"b.o", "b.c"},
		{"c.o", "c.c"},
	}

	d := dag.New()

	for _, node := range nodes {
		d.AddNode(node.id, node.data)
	}

	for _, edge := range edges {
		d.AddEdge(edge.from, edge.to)
	}

	roots := d.Roots()

	for _, root := range roots {

		licenses := map[string]struct{}{}

		root.Walk(dag.WalkDepthFirst, func(node *dag.Node) error {

			data, ok := node.Value.(map[string]string)
			if !ok {
				t.Errorf("Cannot get value as map[string]string")
			}

			if _, inMap := data["licence"]; !inMap {
				return nil
			}

			if _, inMap := licenses[data["licence"]]; !inMap {
				licenses[data["licence"]] = struct{}{}
			}
			return nil
		})

		fmt.Printf("Root: %s\n", root.ID)
		for l, _ := range licenses {
			fmt.Println(l)
		}
		fmt.Println()
	}

}
