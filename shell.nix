let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-23.11";
  pkgs = import nixpkgs {
    config = {};
    overlays = [];
  };
in
  pkgs.mkShellNoCC {
    packages = with pkgs; [
      go_1_22
    ];

    nativeBuildInputs = with pkgs; [
      just
      go-outline
      gopls
      gopkgs
      go-tools
    ];

    buildInputs = with pkgs; [
      just
      gopls
      gopkgs
      go-tools
    ];

    # Env
    GOARCH = "amd64";
    GOOS = "linux";
  }
