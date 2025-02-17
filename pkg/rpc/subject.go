package rpc

import (
	"context"
	"slices"
	"strings"
	"unicode"

	"github.com/nats-io/nats.go"
	"github.com/synternet/data-layer-sdk/pkg/service"
	"github.com/synternet/data-layer-sdk/x/synternet/rpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Publisher interface from your NATS wrapper
type Publisher interface {
	Serve(handler service.ServiceHandler, suffixes ...string) (*nats.Subscription, error)
	RequestFrom(ctx context.Context, msg proto.Message, resp proto.Message, tokens ...string) (service.Message, error)
	SubscribeTo(handler service.MessageHandler, tokens ...string) (*nats.Subscription, error)
	PublishTo(msg proto.Message, tokens ...string) error
	PublishToRpc(msg proto.Message, replyTo string, tokens ...string) error
	RpcInbox(suffixes ...string) string
	Unmarshal(nmsg service.Message, msg protoreflect.ProtoMessage) (nats.Header, error)
	Subject(suffixes ...string) string
}

// normalizePascalCase normalizes a Pascal Case string by converting a leading run
// of uppercase letters to proper Pascal Case. For example:
//
//	"PKGFarAway" -> "PkgFarAway"
//	"PKG" -> "pkg" (if the entire string is uppercase)
func normalizePascalCase(s string) string {
	if s == "" {
		return s
	}
	// If the entire string is uppercase, simply lowercase it.
	if s == strings.ToUpper(s) {
		return strings.ToLower(s)
	}
	// Count the number of uppercase letters at the beginning.
	runLen := 0
	for _, r := range s {
		if unicode.IsUpper(r) {
			runLen++
		} else {
			break
		}
	}
	// If more than one character is uppercase at the start, normalize the prefix.
	if runLen > 1 {
		runLen -= 1
		// Convert the prefix so that the first letter remains uppercase,
		// and the rest become lowercase.
		prefix := s[:runLen]
		normalizedPrefix := string(prefix[0]) + strings.ToLower(prefix[1:])
		return normalizedPrefix + s[runLen:]
	}
	return s
}

// splitPascalCase splits a Pascal Case string into a slice of lower case strings.
func splitPascalCase(s string) []string {
	s = normalizePascalCase(s)

	var words []string
	start := 0

	for i, r := range s {
		// If we find an uppercase letter that is not the first character,
		// then the current word ends just before it.
		if i > 0 && unicode.IsUpper(r) {
			str := strings.ToLower(s[start:i])
			if str != "" {
				words = append(words, str)
			}
			start = i
		}
	}
	// Append the last word.
	str := strings.ToLower(s[start:])
	if str != "" {
		words = append(words, str)
	}
	return words
}

// deriveServiceTokens derives tokens from the service's full name.
// For example, "pkg.MyService" becomes tokens like ["pkg", "my", "service"].
func deriveServiceTokens(serviceName protoreflect.FullName) []string {
	if serviceName == "" {
		return []string{}
	}
	tokens := []string{"service"}
	name := splitPascalCase(strings.TrimSuffix(string(serviceName.Name()), "Service"))
	pkg := serviceName.Parent()
	pkgTokens := splitPascalCase(string(pkg.Name()))
	if len(pkgTokens) != 0 {
		tokens = append(tokens, pkgTokens...)
	}
	return append(tokens, name...)
}

// getCustomServicePrefix extracts a custom subject prefix from the service options, if set.
func getCustomServicePrefix(serviceDescriptor protoreflect.Descriptor) string {
	if serviceDescriptor == nil {
		return ""
	}
	opts := serviceDescriptor.Options()
	if opts == nil {
		return ""
	}
	ext := proto.GetExtension(opts, rpc.E_SubjectPrefix)
	if prefix, ok := ext.(string); ok {
		return prefix
	}
	return ""
}

// getCustomMethodSuffix extracts a custom subject suffix from the method options, if set.
func getCustomMethodSuffix(methodDescriptor protoreflect.Descriptor) string {
	if methodDescriptor == nil {
		return ""
	}
	opts := methodDescriptor.Options()
	if opts == nil {
		return ""
	}
	ext := proto.GetExtension(opts, rpc.E_SubjectSuffix)
	if suffix, ok := ext.(string); ok {
		return suffix
	}
	return ""
}

func disableSubscription(methodDescriptor protoreflect.Descriptor) bool {
	if methodDescriptor == nil {
		return false
	}
	opts := methodDescriptor.Options()
	if opts == nil {
		return false
	}
	ext := proto.GetExtension(opts, rpc.E_DisableInputs)
	if suffix, ok := ext.(bool); ok {
		return suffix
	}
	return false
}

// deriveSubject generates the final subject tokens using the following logic:
//  1. If a non-empty prefix is provided as an argument, use it (split by dots).
//  2. Otherwise, check if the service descriptor has a custom subject prefix (via extension).
//     If set, use that; if not, derive tokens from the service's full name.
//  3. For the method portion, check if the method descriptor has a custom subject suffix.
//     If set, use that; otherwise, derive tokens from the method name.
func deriveSubject(prefix string, serviceDescriptor, methodDescriptor protoreflect.Descriptor) []string {
	var tokens []string

	// Use the explicit prefix if provided.
	tokens = append(tokens, prefix)

	// Service prefix.
	if sp := getCustomServicePrefix(serviceDescriptor); sp != "" {
		tokens = append(tokens, sp)
	} else {
		tokens = append(tokens, deriveServiceTokens(serviceDescriptor.FullName())...)
	}

	// Method suffix
	if ms := getCustomMethodSuffix(methodDescriptor); ms != "" {
		tokens = append(tokens, ms)
	} else {
		tokens = append(tokens, splitPascalCase(string(methodDescriptor.Name()))...)
	}

	return slices.DeleteFunc(tokens, func(s string) bool { return s == "" })
}
