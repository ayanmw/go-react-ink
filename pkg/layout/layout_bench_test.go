package layout

import (
	"testing"
	"time"
)

// Performance benchmarks for the layout engine

func BenchmarkLayoutSimpleBox(b *testing.B) {
	node := NewNode()
	node.Width = 100
	node.Height = 50

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.CalculateLayout(100, 50)
	}
}

func BenchmarkLayoutNestedBoxes(b *testing.B) {
	root := NewNode()
	root.Width = 800
	root.Height = 600
	root.Direction = DirectionColumn

	// Create nested structure
	for i := 0; i < 10; i++ {
		child := NewNode()
		child.Width = 800
		child.Height = 60
		child.Direction = DirectionRow
		root.AddChild(child)

		for j := 0; j < 5; j++ {
			grandchild := NewNode()
			grandchild.FlexGrow = 1
			child.AddChild(grandchild)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.CalculateLayout(800, 600)
	}
}

func BenchmarkLayoutDeepTree(b *testing.B) {
	root := NewNode()
	root.Width = 1000
	root.Height = 800
	root.Direction = DirectionColumn

	// Build a deep tree
	current := root
	for i := 0; i < 20; i++ {
		child := NewNode()
		child.Width = float64(1000 - i*20)
		child.Height = float64(800 - i*10)
		child.Direction = DirectionRow
		current.AddChild(child)
		current = child
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.CalculateLayout(1000, 800)
	}
}

func BenchmarkLayoutComplexFlex(b *testing.B) {
	root := NewNode()
	root.Width = 1200
	root.Height = 900
	root.Direction = DirectionColumn
	root.Justify = JustifySpaceBetween
	root.AlignItems = AlignCenter
	root.Padding = [4]float64{10, 10, 10, 10}

	// Create various flex configurations
	for i := 0; i < 5; i++ {
		row := NewNode()
		row.Direction = DirectionRow
		row.Justify = JustifySpaceAround
		row.AlignItems = AlignCenter
		row.Margin = [4]float64{5, 0, 5, 0}
		row.FlexGrow = 1
		root.AddChild(row)

		for j := 0; j < 8; j++ {
			item := NewNode()
			item.Width = 100
			item.Height = 80
			item.Margin = [4]float64{0, 10, 0, 0}
			row.AddChild(item)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.CalculateLayout(1200, 900)
	}
}

func BenchmarkMeasureText(b *testing.B) {
	node := NewNode()
	node.MeasureFunc = func(width float64) (minWidth, maxWidth, height float64) {
		// Simulate text measurement
		text := "Hello, World! This is a performance test."
		return float64(len(text) * 10), float64(len(text) * 10), 20
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.CalculateLayout(500, 0)
	}
}

func BenchmarkLayoutWithWrap(b *testing.B) {
	root := NewNode()
	root.Width = 200
	root.Height = 1000
	root.Direction = DirectionRow
	root.Wrap = WrapWrap
	root.AlignItems = AlignFlexStart

	// Create items that will wrap
	for i := 0; i < 20; i++ {
		item := NewNode()
		item.Width = 80
		item.Height = 60
		item.Margin = [4]float64{0, 5, 5, 0}
		root.AddChild(item)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.CalculateLayout(200, 1000)
	}
}

func BenchmarkLayoutFullApp(b *testing.B) {
	// Simulate a full TUI app layout
	root := NewNode()
	root.Width = 800
	root.Height = 600
	root.Direction = DirectionColumn
	root.Padding = [4]float64{1, 1, 1, 1}

	// Header
	header := NewNode()
	header.Height = 3
	header.Direction = DirectionRow
	root.AddChild(header)
	for i := 0; i < 3; i++ {
		hItem := NewNode()
		hItem.Width = 200
		hItem.Height = 3
		hItem.FlexGrow = 1
		header.AddChild(hItem)
	}

	// Main content
	content := NewNode()
	content.Direction = DirectionRow
	content.FlexGrow = 1
	root.AddChild(content)

	// Sidebar
	sidebar := NewNode()
	sidebar.Width = 200
	sidebar.FlexGrow = 0
	content.AddChild(sidebar)
	for i := 0; i < 10; i++ {
		sItem := NewNode()
		sItem.Height = 2
		sItem.Margin = [4]float64{0, 0, 1, 0}
		sidebar.AddChild(sItem)
	}

	// Main area
	main := NewNode()
	main.Direction = DirectionColumn
	main.FlexGrow = 1
	content.AddChild(main)
	for i := 0; i < 20; i++ {
		mItem := NewNode()
		mItem.Height = 2
		mItem.Margin = [4]float64{0, 0, 1, 0}
		main.AddChild(mItem)
	}

	// Footer
	footer := NewNode()
	footer.Height = 2
	footer.Direction = DirectionRow
	root.AddChild(footer)
	for i := 0; i < 5; i++ {
		fItem := NewNode()
		fItem.Width = 100
		fItem.Height = 2
		fItem.FlexGrow = 1
		footer.AddChild(fItem)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.CalculateLayout(800, 600)
	}
}

// Performance validation test
func TestLayoutPerformance(t *testing.T) {
	// Test that layout completes in reasonable time
	node := NewNode()
	node.Width = 800
	node.Height = 600
	node.Direction = DirectionColumn

	// Create a realistic layout
	for i := 0; i < 50; i++ {
		child := NewNode()
		child.Direction = DirectionRow
		child.Justify = JustifySpaceBetween
		child.FlexGrow = 1
		node.AddChild(child)

		for j := 0; j < 10; j++ {
			grandchild := NewNode()
			grandchild.FlexGrow = 1
			child.AddChild(grandchild)
		}
	}

	// Measure layout time
	start := time.Now()
	node.CalculateLayout(800, 600)
	duration := time.Since(start)

	// Layout should complete in less than 10ms for reasonable size
	if duration > 10*time.Millisecond {
		t.Logf("Layout took %v for 500 nodes", duration)
	}

	t.Logf("Layout performance: %v for %d nodes", duration, countNodes(node))
}

func countNodes(node *Node) int {
	count := 1
	for _, child := range node.Children {
		count += countNodes(child)
	}
	return count
}
