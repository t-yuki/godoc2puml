package ast

import "sort"

type PackageSlice []*Package
type ClassSlice []*Class
type InterfaceSlice []*Interface
type RelationSlice []*Relation

func (p PackageSlice) Len() int           { return len(p) }
func (p PackageSlice) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p PackageSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PackageSlice) Sort() {
	sort.Sort(p)
	for _, v := range p {
		ClassSlice(v.Classes).Sort()
		InterfaceSlice(v.Interfaces).Sort()
	}
}

func (p ClassSlice) Len() int           { return len(p) }
func (p ClassSlice) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p ClassSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ClassSlice) Sort() {
	sort.Sort(p)
	for _, v := range p {
		RelationSlice(v.Relations).Sort()
	}
}

func (p InterfaceSlice) Len() int           { return len(p) }
func (p InterfaceSlice) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p InterfaceSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p InterfaceSlice) Sort() {
	sort.Sort(p)
	for _, v := range p {
		RelationSlice(v.Relations).Sort()
	}
}

func (p RelationSlice) Len() int           { return len(p) }
func (p RelationSlice) Less(i, j int) bool { return p[i].Target < p[j].Target }
func (p RelationSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p RelationSlice) Sort()              { sort.Sort(p) }
