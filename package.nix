{
  lib,
  fetchFromGitHub,
  buildGoModule,
}:
buildGoModule {
  pname = "toru";
  version = "0.2";

  src = fetchFromGitHub {
    owner = "sweetbbak";
    repo = "toru";
    rev = "cf1729a116927b5ffc28eb99360628a27097fbe6";
    hash = "sha256-Na4JJ4RvghY/cx1H6XwMoP6963hobZyvZqhoKsGPDxs=";
  };

  vendorHash = "sha256-04fHwblTspzecnTizUlFqLwtcnBsHrFRvLX1eSXztRI=";

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
