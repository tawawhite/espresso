// Package build provides all functionality required for performing a
// build from reading content files to modelling the static site.
package build

import (
	"github.com/dominikbraun/espresso/model"
	"path/filepath"
	"strings"
)

// Site represents the actual website. It is a generic data model that
// holds all components and pages and can be rendered to a static site.
type Site struct {
	Nav    *model.Nav
	root   Route
	Footer *model.Footer
}

// Route represents a website Route. Each Route can have multiple pages
// associated with it, as well as multiple child routes. For example, a
// website route like /blog/my-category can be represented as:
//
//	"blog" {
//		Children:
//			"my-category" {
//				Pages: ...
//			}
//	}
//
// The root field of Site is considered as the root route that holds all
// sub-routes: "/blog" would be a child route of the site's root.
type Route struct {
	Pages    []*model.ArticlePage
	ListPage *model.ArticleListPage
	Children map[string]*Route
}

// newSite creates and initializes a new Site instance.
func newSite() *Site {
	s := Site{
		root: Route{
			Pages:    make([]*model.ArticlePage, 0),
			Children: make(map[string]*Route),
		},
	}
	return &s
}

// newRoute creates and initializes a new Route instance.
func newRoute() *Route {
	r := Route{
		Pages: make([]*model.ArticlePage, 0),
		ListPage: &model.ArticleListPage{
			Page:     model.Page{},
			Articles: make([]*model.Article, 0),
		},
		Children: make(map[string]*Route),
	}
	return &r
}

// registerPage registers a given page under the route (path) that is
// stored in page.Path. This path must not end with a trailing slash.
//
// If the route doesn't exist yet, all of its required child-routes will
// be created until the entire page path is depicted.
func (s *Site) registerPage(page *model.ArticlePage) {
	node := &s.root
	segments := strings.Split(page.Path, "/")

	for i, seg := range segments {
		// If the child route (identified by the segment) doesn't exist,
		// create a new route under the current segment key.
		if _, exists := node.Children[seg]; !exists {
			node.Children[seg] = newRoute()
			// Set the "absolute" path of the list page to the current route
			// by joining all segments up to the current segment.
			node.Children[seg].ListPage.Path = filepath.Join(segments[:i]...)
		}
		// Append the page to the current segment if it is the last one.
		if i == len(segments)-1 {
			node.Children[seg].Pages = append(node.Children[seg].Pages, page)
			break
		}
		// Walk down the tree to the next segment.
		node = node.Children[seg]
	}
}

// WalkRoutes walks all site routes recursively and invokes a function
// for each route. depth specifies the maximal depth that the route tree
// will be walked down. Use -1 to walk down to the lowest level.
func (s *Site) WalkRoutes(walkFn func(r *Route), depth int) {
	s.walkRoute(&s.root, walkFn, depth, 0)
}

// walkRoute is used internally by WalkRoutes and should not be called
// by other functions. It is the actual implementation of WalkRoutes.
func (s *Site) walkRoute(route *Route, walkFn func(r *Route), depth int, currentDepth int) {
	if depth != -1 && currentDepth == depth {
		return
	}
	currentDepth++

	for _, route := range route.Children {
		walkFn(route)
		s.walkRoute(route, walkFn, depth, currentDepth)
	}
}
