package ast

import (
	"fmt"
	"github.com/elliotchance/c2go/util"
	"regexp"
)

type Position struct {
	File      string // The relative or absolute file path.
	Line      int    // Start line
	LineEnd   int    // End line
	Column    int    // Start column
	ColumnEnd int    // End column

	// This is the original string that was converted. This is used for
	// debugging. We could derive this value from the other properties to save
	// on a bit of memory, but let worry about that later.
	StringValue string
}

func NewPositionFromString(s string) Position {
	if s == "<invalid sloc>" || s == "" {
		return Position{}
	}

	re := regexp.MustCompile(`^col:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Column:      util.Atoi(groups[1]),
		}
	}

	re = regexp.MustCompile(`^col:(\d+), col:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Column:      util.Atoi(groups[1]),
			ColumnEnd:   util.Atoi(groups[2]),
		}
	}

	re = regexp.MustCompile(`^line:(\d+), line:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Line:        util.Atoi(groups[1]),
			LineEnd:     util.Atoi(groups[2]),
		}
	}

	re = regexp.MustCompile(`^col:(\d+), line:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Column:      util.Atoi(groups[1]),
			Line:        util.Atoi(groups[2]),
		}
	}

	re = regexp.MustCompile(`^line:(\d+):(\d+), line:(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Line:        util.Atoi(groups[1]),
			Column:      util.Atoi(groups[2]),
			LineEnd:     util.Atoi(groups[3]),
			ColumnEnd:   util.Atoi(groups[4]),
		}
	}

	re = regexp.MustCompile(`^col:(\d+), line:(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Column:      util.Atoi(groups[1]),
			LineEnd:     util.Atoi(groups[2]),
			ColumnEnd:   util.Atoi(groups[3]),
		}
	}

	re = regexp.MustCompile(`^line:(\d+):(\d+), col:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Line:        util.Atoi(groups[1]),
			Column:      util.Atoi(groups[2]),
			ColumnEnd:   util.Atoi(groups[3]),
		}
	}

	re = regexp.MustCompile(`^line:(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Line:        util.Atoi(groups[1]),
			Column:      util.Atoi(groups[2]),
		}
	}

	// This must be below all of the others.
	re = regexp.MustCompile(`^([^:]+):(\d+):(\d+), col:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			File:        groups[1],
			Line:        util.Atoi(groups[2]),
			Column:      util.Atoi(groups[3]),
			ColumnEnd:   util.Atoi(groups[4]),
		}
	}

	re = regexp.MustCompile(`^([^:]+):(\d+):(\d+), line:(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			File:        groups[1],
			Line:        util.Atoi(groups[2]),
			Column:      util.Atoi(groups[3]),
			LineEnd:     util.Atoi(groups[4]),
			ColumnEnd:   util.Atoi(groups[5]),
		}
	}

	re = regexp.MustCompile(`^([^:]+):(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			File:        groups[1],
			Line:        util.Atoi(groups[2]),
			Column:      util.Atoi(groups[3]),
		}
	}

	re = regexp.MustCompile(`^([^:]+):(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			File:        groups[1],
			Line:        util.Atoi(groups[2]),
			Column:      util.Atoi(groups[3]),
		}
	}

	re = regexp.MustCompile(`^col:(\d+), ([^:]+):(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Column:      util.Atoi(groups[1]),
		}
	}

	re = regexp.MustCompile(`^([^:]+):(\d+):(\d+), ([^:]+):(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			File:        groups[1],
			Line:        util.Atoi(groups[2]),
			Column:      util.Atoi(groups[3]),
			LineEnd:     util.Atoi(groups[5]),
			ColumnEnd:   util.Atoi(groups[6]),
		}
	}

	panic("unable to understand position '" + s + "'")
}

func mergePositions(p1, p2 Position) Position {
	if p2.File != "" {
		p1.File = p2.File
		p1.Line = 0
		p1.LineEnd = 0
		p1.Column = 0
		p1.ColumnEnd = 0
	}

	if p2.Line != 0 {
		p1.Line = p2.Line
		p1.LineEnd = 0
	}

	if p2.LineEnd != 0 {
		p1.LineEnd = p2.LineEnd
	}

	if p2.Column != 0 {
		p1.Column = p2.Column
		p1.ColumnEnd = 0
	}

	if p2.ColumnEnd != 0 {
		p1.ColumnEnd = p2.ColumnEnd
	}

	return p1
}

var pos Position

func FixPositions(nodes []Node) {
	pos = Position{File: "<built-in>"}
	fixPositions(nodes)
}

func fixPositions(nodes []Node) {
	for _, node := range nodes {
		if node != nil {
			pos = mergePositions(pos, node.Position())
			setPosition(node, pos)
			fixPositions(node.Children())
		}
	}
}

func setPosition(node Node, position Position) {
	switch n := node.(type) {
	case *AlignedAttr:
		n.Pos = position
	case *AlwaysInlineAttr:
		n.Pos = position
	case *ArraySubscriptExpr:
		n.Pos = position
	case *AsmLabelAttr:
		n.Pos = position
	case *AvailabilityAttr:
		n.Pos = position
	case *BinaryOperator:
		n.Pos = position
	case *BreakStmt:
		n.Pos = position
	case *CallExpr:
		n.Pos = position
	case *CaseStmt:
		n.Pos = position
	case *CharacterLiteral:
		n.Pos = position
	case *CompoundStmt:
		n.Pos = position
	case *ConditionalOperator:
		n.Pos = position
	case *ConstAttr:
		n.Pos = position
	case *ContinueStmt:
		n.Pos = position
	case *CompoundAssignOperator:
		n.Pos = position
	case *CStyleCastExpr:
		n.Pos = position
	case *DeclRefExpr:
		n.Pos = position
	case *DeclStmt:
		n.Pos = position
	case *DefaultStmt:
		n.Pos = position
	case *DeprecatedAttr:
		n.Pos = position
	case *DoStmt:
		n.Pos = position
	case *EnumConstantDecl:
		n.Pos = position
	case *EnumDecl:
		n.Pos = position
	case *FieldDecl:
		n.Pos = position
	case *FloatingLiteral:
		n.Pos = position
	case *FormatAttr:
		n.Pos = position
	case *FunctionDecl:
		n.Pos = position
	case *ForStmt:
		n.Pos = position
	case *GCCAsmStmt:
		n.Pos = position
	case *GotoStmt:
		n.Pos = position
	case *IfStmt:
		n.Pos = position
	case *ImplicitCastExpr:
		n.Pos = position
	case *ImplicitValueInitExpr:
		n.Pos = position
	case *IndirectFieldDecl:
		n.Pos = position
	case *InitListExpr:
		n.Pos = position
	case *IntegerLiteral:
		n.Pos = position
	case *LabelStmt:
		n.Pos = position
	case *MallocAttr:
		n.Pos = position
	case *MaxFieldAlignmentAttr:
		n.Pos = position
	case *MemberExpr:
		n.Pos = position
	case *ModeAttr:
		n.Pos = position
	case *NoInlineAttr:
		n.Pos = position
	case *NoThrowAttr:
		n.Pos = position
	case *NonNullAttr:
		n.Pos = position
	case *OffsetOfExpr:
		n.Pos = position
	case *PackedAttr:
		n.Pos = position
	case *ParenExpr:
		n.Pos = position
	case *ParmVarDecl:
		n.Pos = position
	case *PredefinedExpr:
		n.Pos = position
	case *PureAttr:
		n.Pos = position
	case *RecordDecl:
		n.Pos = position
	case *RestrictAttr:
		n.Pos = position
	case *ReturnStmt:
		n.Pos = position
	case *ReturnsTwiceAttr:
		n.Pos = position
	case *StringLiteral:
		n.Pos = position
	case *SwitchStmt:
		n.Pos = position
	case *TransparentUnionAttr:
		n.Pos = position
	case *TypedefDecl:
		n.Pos = position
	case *UnaryExprOrTypeTraitExpr:
		n.Pos = position
	case *UnaryOperator:
		n.Pos = position
	case *UnusedAttr:
		n.Pos = position
	case *VAArgExpr:
		n.Pos = position
	case *VarDecl:
		n.Pos = position
	case *WarnUnusedResultAttr:
		n.Pos = position
	case *WeakAttr:
		n.Pos = position
	case *WhileStmt:
		n.Pos = position
	case *TypedefType, *Typedef, *TranslationUnitDecl, *RecordType, *Record,
		*QualType, *PointerType, *ParenType, *IncompleteArrayType,
		*FunctionProtoType, *EnumType, *Enum, *ElaboratedType,
		*ConstantArrayType, *BuiltinType, *ArrayFiller, *Field:
		// These do not have positions so they can be ignored.
	default:
		panic(fmt.Sprintf("unknown node type: %+#v", node))
	}
}
