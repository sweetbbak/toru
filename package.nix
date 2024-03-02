{ lib, fetchFromGitHub, buildGoModule }:

buildGoModule rec {
    pname = "toru";
    version = "0.1";

    src = fetchFromGitHub {
      owner = "sweetbbak";
      repo = "toru";
      rev = "9dc67d420208bb5f9debd260170d54035242c7ab";
      hash = "sha256-2Z5agQtF6p21rnAcjsRr+3QOJ0QGveKVH8e9LHpm3ZE=";
    };

    # vendorHash = lib.fakeHash;
    vendorHash = "sha256-alC4/2wTbjJYWGzTDTgQweOicN3xSqfnncok/j16+0E=";

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
