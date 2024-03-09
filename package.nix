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
    rev = "3753b62a6d22c2e0389f569d1a0c2130dca97334";
    hash = "sha256-ypwXDvvWEm6neOVd2BmwBKdZ/dlVClhvnF03HHeyM48=";
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
