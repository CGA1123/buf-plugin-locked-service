package main

import (
	"context"

	"buf.build/go/bufplugin/check"
	"buf.build/go/bufplugin/check/checkutil"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	lockedServiceRuleID = "LOCKED_SERVICE"
)

var (
	lockedServiceRuleSpec = &check.RuleSpec{
		ID:          lockedServiceRuleID,
		Default:     true,
		Purpose:     "Checks that a locked service hasn't been changed.",
		CategoryIDs: []string{},
		Type:        check.RuleTypeBreaking,
		Handler:     checkutil.NewServicePairRuleHandler(checkLockedServiceDoesntChange, checkutil.WithoutImports()),
	}
	spec = &check.Spec{
		Rules: []*check.RuleSpec{
			lockedServiceRuleSpec,
		},
	}
)

func main() {
	check.Main(spec)
}

func checkLockedServiceDoesntChange(
	ctx context.Context,
	w check.ResponseWriter,
	r check.Request,
	current protoreflect.ServiceDescriptor,
	against protoreflect.ServiceDescriptor,
) error {
	lockedServices, ok := getLockedServices(r)
	if !ok {
		return nil
	}

	if !isLockedService(current, lockedServices) {
		return nil
	}

	previousMethods := methodDescriptors(against)
	for _, method := range methodDescriptors(current) {
		if previouslyDefined(method, previousMethods) {
			continue
		}

		w.AddAnnotation(
			check.WithMessagef(
				"Method %q on locked service %q is new.",
				method.Name(),
				current.FullName(),
			),
			check.WithDescriptor(method),
			check.WithAgainstDescriptor(against),
		)
	}

	return nil
}

func isLockedService(svc protoreflect.ServiceDescriptor, lockedServices []string) bool {
	for _, locked := range lockedServices {
		if string(svc.FullName()) != locked {
			continue
		}

		return true
	}

	return false
}

func getLockedServices(r check.Request) ([]string, bool) {
	raw, ok := r.Options().Get("locked_services")
	if !ok {
		return nil, false
	}

	list, ok := raw.([]string)
	if !ok {
		return nil, false
	}

	return list, true
}

func methodDescriptors(svc protoreflect.ServiceDescriptor) []protoreflect.MethodDescriptor {
	mtds := svc.Methods()

	descs := make([]protoreflect.MethodDescriptor, mtds.Len())
	for i := 0; i < mtds.Len(); i++ {
		descs[i] = mtds.Get(i)
	}

	return descs
}

func previouslyDefined(method protoreflect.MethodDescriptor, previousMethods []protoreflect.MethodDescriptor) bool {
	for _, previous := range previousMethods {
		if method.FullName() != previous.FullName() {
			continue
		}

		return true
	}

	return false
}
