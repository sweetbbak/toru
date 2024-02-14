let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-23.11";
  pkgs = import nixpkgs { config = {}; overlays = []; };
in

pkgs.mkShellNoCC {
  packages = with pkgs; [
    eza
    helix
    go_1_21
  ];

  nativeBuildInputs = with pkgs; [
    just
    go-outline
    gopls
    gopkgs
    go-tools
  ];
  # Env
  GOARCH = "amd64";
  GOOS = "linux";
}
