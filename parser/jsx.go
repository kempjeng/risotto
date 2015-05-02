package parser

import (
	"github.com/mamaar/risotto/ast"
	"github.com/mamaar/risotto/token"
)

func (self *_parser) parseJSX() *ast.JSXBlock {
	jsx := &ast.JSXBlock{}

	jsx.OpeningElement = self.parseOpeningElement()
	if jsx.OpeningElement.SelfClosing {
		return jsx
	}
	jsx.ClosingElement = self.parseClosingElement(jsx.OpeningElement.Name.Name)

	return jsx
}

func (self *_parser) parseClosingElement(expectedName string) ast.JSXElement {
	closing := ast.JSXElement{}
	closing.LeftTag = self.expect(token.LESS)
	self.expect(token.SLASH)

	if self.token == token.IDENTIFIER {
		closing.Name = self.parseIdentifier()
		if closing.Name.Name != expectedName {
			self.error(self.idx, "Closing element does not match the name of any opening elements")
		}
	}

	closing.RightTag = self.expect(token.GREATER)
	return closing
}

func (self *_parser) parseOpeningElement() ast.JSXElement {
	open := ast.JSXElement{}
	open.LeftTag = self.expect(token.LESS)

	if self.token == token.IDENTIFIER {
		open.Name = self.parseIdentifier()
	}

	for self.token == token.IDENTIFIER {
		open.PropertyList = append(open.PropertyList, self.parseJSXProperty())
	}

	if self.token == token.SLASH {
		self.next()
		open.RightTag = self.expect(token.GREATER)
		open.SelfClosing = true
		return open
	}

	open.RightTag = self.expect(token.GREATER)
	return open
}

func (self *_parser) parseJSXProperty() ast.Property {
	p := ast.Property{}

	p.Key = self.parseIdentifier().Name
	self.expect(token.ASSIGN)
	p.Value = self.parsePrimaryExpression()

	return p
}
