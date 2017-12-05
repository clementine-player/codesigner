package codesigner

import (
  "context"
  "fmt"
)

type CodeSigner struct{}

func (s *CodeSigner) SignPackage(ctx context.Context, req *SignPackageRequest) (*SignPackageReply, error) {
  return nil, fmt.Errorf("Not implemented")
}

func (s *CodeSigner) VerifyPackage(ctx context.Context, req *VerifyPackageRequest) (*VerifyPackageReply, error) {
  return nil, fmt.Errorf("Not implemented")
}
