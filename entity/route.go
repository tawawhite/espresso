// Package entity provides the Espresso domain models. They are
// used to represent the entire website and its components.
package entity

// Route is a relative URL that holds multiple pages as well as
// an index page. It is the first part of a FQN.
type Route struct {
	Children  map[string]*Route // Children is a set of sub-routes of the current route.
	Pages     []Page            // Pages is a set of pages available under the route.
	IndexPage IndexPage         // IndexPage is the index page generated by Espresso or provided by the user.
}
