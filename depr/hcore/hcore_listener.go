// Code generated from HCore.g4 by ANTLR 4.7.1. DO NOT EDIT.

package HCore // HCore
import "github.com/antlr/antlr4/runtime/Go/antlr"

// HCoreListener is a complete listener for a parse tree produced by HCoreParser.
type HCoreListener interface {
	antlr.ParseTreeListener

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterOp is called when entering the Op production.
	EnterOp(c *OpContext)

	// EnterOperation is called when entering the Operation production.
	EnterOperation(c *OperationContext)

	// EnterNumber is called when entering the Number production.
	EnterNumber(c *NumberContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitOp is called when exiting the Op production.
	ExitOp(c *OpContext)

	// ExitOperation is called when exiting the Operation production.
	ExitOperation(c *OperationContext)

	// ExitNumber is called when exiting the Number production.
	ExitNumber(c *NumberContext)
}
