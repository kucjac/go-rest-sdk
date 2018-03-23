package response

import (
	"github.com/kucjac/go-rest-sdk/refutils"
	"github.com/kucjac/go-rest-sdk/resterrors"
)

// DefaultBody is a structure that defines default response body
// It implements Responser interface
type DefaultBody struct {
	Content map[string]interface{} `json:"content,omitempty"`
	Errors  []*resterrors.Error    `json:"errors,omitempty"`
}

// AddContent adds the given 'content' to the default body content
// The 'content' is saved with key - lowercased model struct name
// pluralized if slice provided.
// I.e. Providing model type Foo struct{} - would be saved as 'foo'
// 		But slice []Foo or []*Foo would result in 'foos'
// With this method DefaultBody implements ContentAdder interface
func (d *DefaultBody) AddContent(content ...interface{}) {
	d.addContent(content...)
}

// WithContent adds the given 'content' to the DefaultBody Content field
// The content is added with the same rules as with AddContent() method
// The difference is that this method may be used as callback - after
// processing it returns itself.
func (d *DefaultBody) WithContent(content ...interface{}) Responser {
	d.addContent(content...)
	return d
}

// AddErrors adds the given 'errors' to the default body content.
// With this method DefaultBody implements ErrorAdder interface
func (d *DefaultBody) AddErrors(errors ...*resterrors.Error) {
	d.addErrors(errors...)
}

// WithErrors adds the given 'errors' to the default body content.
// This method acts like AddErrors() method, but in addition it may be used
// as a callback function - returning *DefaultBody itself
func (d *DefaultBody) WithErrors(errors ...*resterrors.Error) Responser {
	d.addErrors(errors...)
	return d
}

// New creates and returns new *DefaultBody of given type
// Implements Responser New() method
func (d *DefaultBody) New() Responser {
	return &DefaultBody{Content: make(map[string]interface{})}
}

// NewErrored implements Responser NewErrored() method
// In this implementation there is no difference between New and NewErrored method.
func (d *DefaultBody) NewErrored() Responser {
	return &DefaultBody{}
}

func (d *DefaultBody) addContent(contents ...interface{}) {
	for _, content := range contents {
		d.Content[refutils.ModelName(content)] = content
	}
}

func (d *DefaultBody) addErrors(errors ...*resterrors.Error) {
	d.Errors = append(d.Errors, errors...)
}
