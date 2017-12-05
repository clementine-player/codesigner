package codesigner

import (
  "context"
  "fmt"
  "io/ioutil"
  "os/exec"
)

type CodeSigner struct{}

func unlockKeychain(password string, keychainPath string) (string, error) {
  cmd := exec.Command("security", "unlock-keychain", "-p", password, "buildbot.keychain")
  out, err := cmd.CombinedOutput()
  return string(out), err
}

func (s *CodeSigner) SignPackage(ctx context.Context, req *SignPackageRequest) (*SignPackageReply, error) {
  unlock, err := unlockKeychain(req.GetPassword(), "buildbot.keychain")
  if err != nil {
    return nil, fmt.Errorf("Failed to unlock keychain: %v; %s", err, unlock)
  }
  temp, err := ioutil.TempFile("", "codesigner")
  if err != nil {
    return nil, fmt.Errorf("Failed to create temp file for signing: %v", err)
  }
  _, err = temp.Write(req.GetPackage())
  if err != nil {
    return nil, fmt.Errorf("Failed to write temp file: %v", err)
  }
  temp.Close()
  cmd := exec.Command("codesign", "-fv", "-s", req.GetDeveloperId(), temp.Name())
  out, err := cmd.CombinedOutput()
  if err != nil {
    return nil, fmt.Errorf("Failed to codesign: %s", out)
  }
  signed, err := ioutil.ReadFile(temp.Name())
  if err != nil {
    return nil, fmt.Errorf("Failed to read back signed data: %v", err)
  }
  return &SignPackageReply{
    SignedPackage: signed,
    CodesignOutput: string(out),
  }, nil
}

func (s *CodeSigner) VerifyPackage(ctx context.Context, req *VerifyPackageRequest) (*VerifyPackageReply, error) {
  return nil, fmt.Errorf("Not implemented")
}
