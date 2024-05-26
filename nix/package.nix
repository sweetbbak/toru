{
  lib,
  go,
  fetchFromGitHub,
  buildGoModule,
}:
buildGoModule.override {inherit go;} {
  pname = "toru";
  version = "0.3";

  src = fetchFromGitHub {
    owner = "sweetbbak";
    repo = "toru";
    rev = "bb0d81894f59bf4c849502bd9cdab42ba79fdffe";
    hash = "sha256-oIPXODgXO5HvDLGxUT43IibkznNbJIQA5DvC3QDcWRw=";
  };

  vendorHash = "sha256-04fHwblTspzecnTizUlFqLwtcnBsHrFRvLX1eSXztRI=";

  CGO_ENABLED = 0;
  ldflags = ["-s" "-w"];
  tags = ["torrent" "bittorrent" "anime"];

  buildPhase = ''
    go build -ldflags="-s -w" ./cmd/toru
  '';

  installPhase = ''
    mkdir -p $out/share/zsh/site-functions
    cp completion/_toru $out/share/zsh/site-functions/_toru

    mkdir -p $out/etc/bash_completion.d
    cp completion/_toru_bash $out/etc/bash_completion.d/_toru

    mkdir -p $out/share/man
    cp toru.1 $out/share/man

    mkdir -p $out/bin
    mv toru $out/bin
  '';

  meta = with lib; {
    homepage = "https://github.com/sweetbbak/toru";
    description = "Bittorrent streaming CLI tool. Stream anime torrents, real-time with no waiting for downloads";
    license = licenses.mit;
    maintainers = with maintainers; [sweetbbak];
    mainProgram = "toru";
  };
}
