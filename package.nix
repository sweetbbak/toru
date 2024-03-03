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

  meta = with lib; {
    homepage = "https://github.com/sweetbbak/toru";
    description = "Bittorrent streaming CLI tool. Stream anime torrents, real-time with no waiting for downloads";
    license = licenses.mit;
    maintainers = with maintainers; [sweetbbak];
    mainProgram = "toru";
  };
}
