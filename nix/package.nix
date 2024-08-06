{
  lib,
  go,
  fetchFromGitHub,
  buildGoModule,
}:
buildGoModule.override {inherit go;} {
  pname = "toru";
  version = "0.3.2";

  src = fetchFromGitHub {
    owner = "sweetbbak";
    repo = "toru";
    rev = "2bb096506fb154670c036ee1d9938984af4b22d8";
    hash = "sha256-rfnZvo1O46EW794T0/JaV9Wb8JxnTrbo/aD7EgeHcFc=";
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
