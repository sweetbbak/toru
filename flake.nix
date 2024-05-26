# {
#   inputs.nixpkgs.url = "github:NixOS/nixpkgs?ref=nixos-unstable";
#   outputs = {nixpkgs, ...}: {
#     packages.x86_64-linux = let
#       pkgs = nixpkgs.legacyPackages.x86_64-linux;
#     in {
#       default = pkgs.callPackage ./package.nix {};
#     };
#
#     packages.aarch64-darwin = let
#       pkgs = nixpkgs.legacyPackages.x86_64-linux;
#     in {
#       default = pkgs.callPackage ./package.nix {};
#     };
#   };
# }
{
  description = "Flake utils demo";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
        };
        # pkgs = nixpkgs.legacyPackages.${system};
        go = pkgs.go;
      in {
        devShell = pkgs.mkShell {
          buildInputs = [
            go
            pkgs.gopls
            pkgs.go-tools
          ];
          shellHook = ''
            export CGO_ENABLED=0
            # export GOOS=${system}
            # export GOARCH=${system}
          '';
        };

        # default = pkgs.callPackage ./nix/package.nix {};
        defaultPackage = pkgs.callPackage ./nix/package.nix {};
      }
    );
}
