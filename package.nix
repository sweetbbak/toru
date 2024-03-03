{
  lib,
  fetchFromGitHub,
  buildGoModule,
}:
buildGoModule rec {
    pname = "toru";
    version = "0.1";

  src = fetchFromGitHub {
    owner = "sweetbbak";
    repo = "toru";
    rev = "c5f6134b42d6dc8fd3f52e598f6fc6b1eae74888";
    hash = "sha256-6jZnuE2xNF0lrYVjBXTBGQCkT0vxPzHlmbnRDLPcglE=";
  };

  vendorHash = "sha256-bmG2qBlyN6aYdWvplXClylwCk2pDfr30w0ztgUXP71g=";
  CGO_ENABLED = 0;
  ldflags = ["-s" "-w"];
  tags = ["torrent" "bittorrent" "anime"];

    buildPhase = ''
        go mod vendor
        go build -o toru ./cmd/toru
    '';

    installPhase = ''
        mkdir -p $out/bin
        mv toru $out/bin
    '';

    meta = with lib; {
        homepage = "https://github.com/sweetbbak/toru";
        description = "Stream anime from the command line with the power of torrents";
        license = licenses.mit;
        maintainers = with maintainers; [ sweetbbak ];
        mainProgram = "toru";
    };
}
