{ lib, fetchFromGitHub, buildGoModule }:

buildGoModule rec {
    pname = "toru";
    version = "0.1";

    src = fetchFromGitHub {
      owner = "sweetbbak";
      repo = "toru";
      rev = "v0.1";
      hash = "sha256-GiptGi2OZy6eiDa/g8GwPVZpDd8QUxZuzaEy+9J+v9I=";
    };

    # vendorHash = lib.fakeHash;
    vendorHash = "sha256-wbXpWyOLB5CwnzvRuM8BvlO5BK+AY4uQdGkV151DwvA=";

    CGO_ENABLED = 1;
    ldflags = [ "-s" "-w" ];

    tags = [ "torrent" "bittorrent" "anime" ];
    proxyVendor = true;

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
