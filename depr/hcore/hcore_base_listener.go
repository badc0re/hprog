// Code generated from HCore.g4 by ANTLR 4.7.1. DO NOT EDIT.

package HCore // HCore
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseHCoreListener is a complete listener for a parse tree produced by HCoreParser.
type BaseHCoreListener struct{}

var _ HCoreListener = &BaseHCoreListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseHCoreListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseHCoreListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseHCoreListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseHCoreListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BaseHCoreListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseHCoreListener) ExitStart(ctx *StartContext) {}

// EnterOp is called when production Op is entered.
func (s *BaseHCoreListener) EnterOp(ctx *OpContext) {}

// ExitOp is called when production Op is exited.
func (s *BaseHCoreListener) ExitOp(ctx *OpContext) {}

// EnterOperation is called when production Operation is entered.
func (s *BaseHCoreListener) EnterOperation(ctx *OperationContext) {}

// ExitOperation is called when production Operation is exited.
func (s *BaseHCoreListener) ExitOperation(ctx *OperationContext) {}

// EnterNumber is called when production Number is entered.
func (s *BaseHCoreListener) EnterNumber(ctx *NumberContext) {}

// ExitNumber is called when production Number is exited.
func (s *BaseHCoreListener) ExitNumber(ctx *NumberContext) {}
