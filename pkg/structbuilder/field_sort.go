package structbuilder

import "sort"

// SortFieldsByName sorts the fields by their name.
func SortFieldsByName(fields []*Field) {
	sorter := byNameFieldSorter(fields)
	sort.Sort(sorter)
}

type byNameFieldSorter []*Field

func (s byNameFieldSorter) Len() int {
	return len(s)
}

func (s byNameFieldSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byNameFieldSorter) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
