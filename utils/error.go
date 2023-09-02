package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"blockwatch.cc/tzgo/rpc"
)

type internalRPCError struct {
	ID          string                 `json:"id"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
}

func (e *internalRPCError) Error() string {
	return fmt.Sprintf("%s: %s", e.ID, e.Description)
}

func ExtractInternalRPCerror(errMsg string) error {
	// Find the index of "Error:"
	idx := strings.Index(errMsg, "Error:")
	if idx == -1 {
		return fmt.Errorf("could not find 'Error:' in message")
	}

	// Extract the part of the message after "Error:"
	errPart := strings.TrimSpace(errMsg[idx+len("Error:"):])
	dec := json.NewDecoder(strings.NewReader(errPart))

	var internalError internalRPCError
	if err := dec.Decode(&internalError); err != nil {
		return fmt.Errorf("failed to decode error JSON: %w", err)
	}

	return &internalError
}

type errorResponseBodyMessage struct {
	Msg string `json:"msg"`
}

func TryUnwrapRPCError(err error) error {
	if rpcError, ok := err.(rpc.RPCError); ok {
		body := rpcError.Body()
		if len(body) == 0 {
			return err
		}

		var responseBodyMessages []errorResponseBodyMessage
		if bodyParseError := json.Unmarshal(body, &responseBodyMessages); bodyParseError != nil {
			return errors.Join(err, bodyParseError)
		}

		result := err
		for _, v := range responseBodyMessages {
			result = errors.Join(result, ExtractInternalRPCerror(v.Msg))
		}

		return result
	}
	return err
}
