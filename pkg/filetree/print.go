package filetree

import (
	"fmt"
	"io"
	"strings"
)

func PrintTree(writer io.Writer, entry Entry) {
	var printTreeRecursive func(writer io.Writer, entry Entry, level int)

	printTreeRecursive = func(writer io.Writer, entry Entry, level int) {
		tab := strings.Repeat("\t", level)

		write := func(str string) {
			_, _ = writer.Write([]byte(str)) //nolint:errcheck
		}

		if entry.IsDir() {
			childrenCount := len(entry.Children())

			if childrenCount > 0 {
				write(fmt.Sprintf("%s- [%s] (%d children):\n", tab, entry.Name(), childrenCount))

				for _, child := range entry.Children() {
					printTreeRecursive(writer, child, level+1)
				}
			} else {
				write(fmt.Sprintf("%s- [%s]\n", tab, entry.Name()))
			}
		} else {
			write(fmt.Sprintf("%s- %s\n", tab, entry.Name()))
		}
	}

	printTreeRecursive(writer, entry, 0)
}
