{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.ginkgo
    pkgs.credhub
    pkgs.bosh
    pkgs.jq
    pkgs.coreutils-prefixed
    pkgs.libfaketime
  ];
}
