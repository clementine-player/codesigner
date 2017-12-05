package codesigner

import (
  "context"
  "flag"
  "fmt"
  "io/ioutil"
  "os/exec"
  "sync"
)

var defaultKeychainPath = flag.String("keychain", "buildbot.keychain", "Path to keychain containing developer IDs")
var password = flag.String("password", "", "Password for the keychain")

type CodeSigner struct{
  lock sync.Mutex
}

func unlockKeychain(password string, keychainPath string) (string, error) {
  cmd := exec.Command("security", "unlock-keychain", "-p", password, keychainPath)
  out, err := cmd.CombinedOutput()
  return string(out), err
}

func lockKeychain(keychainPath string) (string, error) {
  cmd := exec.Command("security", "lock-keychain", keychainPath)
  out, err := cmd.CombinedOutput()
  return string(out), err
}

func (s *CodeSigner) SignPackage(ctx context.Context, req *SignPackageRequest) (*SignPackageReply, error) {
  s.lock.Lock()
  defer s.lock.Unlock()

  unlock, err := unlockKeychain(*password, *defaultKeychainPath)
  if err != nil {
    return nil, fmt.Errorf("Failed to unlock keychain: %v; %s", err, unlock)
  }
  defer lockKeychain(*defaultKeychainPath)
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
  temp, err := ioutil.TempFile("", "codesigner")
  if err != nil {
    return nil, fmt.Errorf("Failed to create temp file for verifying: %v", err)
  }
  _, err = temp.Write(req.GetPackage())
  if err != nil {
    return nil, fmt.Errorf("Failed to write temp file: %v", err)
  }
  temp.Close()
  cmd := exec.Command("codesign", "-vvv", temp.Name())
  out, err := cmd.CombinedOutput()
  return &VerifyPackageReply{
    Ok: err == nil,
    CodesignOutput: string(out),
  }, nil
}
